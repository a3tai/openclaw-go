// Command server runs a mock OpenClaw Gateway for local development.
//
// It accepts WebSocket connections on /ws and serves:
//   - Gateway protocol handshake (challenge → connect → hello-ok)
//   - system-presence requests
//   - exec.approval.resolve requests
//   - Gateway→Node invoke forwarding
//
// It also exposes the OpenAI-compatible HTTP API:
//   - POST /v1/chat/completions (non-streaming and streaming)
//   - POST /tools/invoke
//
// Usage:
//
//	go run ./examples/server
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/a3tai/openclaw-go/protocol"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWS)
	mux.HandleFunc("/v1/chat/completions", handleChatCompletions)
	mux.HandleFunc("/tools/invoke", handleToolsInvoke)

	addr := ":18789"
	log.Printf("Mock OpenClaw Gateway listening on %s", addr)
	log.Printf("  WebSocket: ws://localhost%s/ws", addr)
	log.Printf("  Chat API:  http://localhost%s/v1/chat/completions", addr)
	log.Printf("  Tools API: http://localhost%s/tools/invoke", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// ---------------------------------------------------------------------------
// WebSocket gateway
// ---------------------------------------------------------------------------

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("[ws] new connection from %s", r.RemoteAddr)

	// 1. Send connect.challenge event.
	challenge := protocol.ConnectChallenge{
		Nonce: fmt.Sprintf("nonce-%d", time.Now().UnixNano()),
		Ts:    time.Now().UnixMilli(),
	}
	evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
	if err := conn.WriteMessage(websocket.TextMessage, evData); err != nil {
		log.Printf("[ws] write challenge: %v", err)
		return
	}

	// 2. Read the connect request.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("[ws] read connect: %v", err)
		return
	}
	var req protocol.Request
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Printf("[ws] unmarshal connect: %v", err)
		return
	}

	var params protocol.ConnectParams
	json.Unmarshal(req.Params, &params)
	log.Printf("[ws] connect: role=%s, client=%s/%s, scopes=%v",
		params.Role, params.Client.ID, params.Client.Version, params.Scopes)

	// 3. Send hello-ok response.
	hello := protocol.HelloOK{
		Type:     "hello-ok",
		Protocol: protocol.ProtocolVersion,
		Server: protocol.HelloServer{
			Version: "mock-1.0.0",
			ConnID:  fmt.Sprintf("conn-%d", time.Now().UnixNano()),
		},
		Features: protocol.HelloFeatures{
			Methods: []string{"system-presence", "exec.approval.resolve", "chat.send", "chat.history"},
			Events:  []string{"tick", "exec.approval.requested", "chat"},
		},
		Snapshot: protocol.Snapshot{
			Presence:     []protocol.PresenceEntry{{Ts: time.Now().UnixMilli(), DeviceID: "mock-device", Roles: []string{"operator"}}},
			Health:       json.RawMessage(`{}`),
			StateVersion: protocol.StateVersion{Presence: 1, Health: 1},
			UptimeMs:     0,
			AuthMode:     "token",
		},
		Policy: protocol.HelloPolicy{
			MaxPayload:       protocol.MaxPayloadBytes,
			MaxBufferedBytes: protocol.MaxBufferedBytes,
			TickIntervalMs:   protocol.DefaultTickIntervalMs,
		},
	}
	respData, _ := protocol.MarshalResponse(req.ID, hello)
	if err := conn.WriteMessage(websocket.TextMessage, respData); err != nil {
		log.Printf("[ws] write hello-ok: %v", err)
		return
	}
	log.Printf("[ws] handshake complete")

	// 4. Serve requests in a loop.
	var mu sync.Mutex
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[ws] read: %v", err)
			return
		}

		frame, err := protocol.ParseFrame(msg)
		if err != nil {
			log.Printf("[ws] parse frame: %v", err)
			continue
		}

		switch frame.Type {
		case protocol.FrameTypeRequest:
			var r protocol.Request
			if err := json.Unmarshal(msg, &r); err != nil {
				continue
			}
			log.Printf("[ws] request: method=%s id=%s", r.Method, r.ID)

			switch r.Method {
			case "system-presence":
				entries := map[string]protocol.PresenceEntry{
					"mock-device": {
						DeviceID: "mock-device",
						Roles:    []string{"operator"},
						Scopes:   []string{"operator.read", "operator.write"},
						Ts:       time.Now().UnixMilli(),
					},
				}
				data, _ := protocol.MarshalResponse(r.ID, entries)
				mu.Lock()
				conn.WriteMessage(websocket.TextMessage, data)
				mu.Unlock()

			case "exec.approval.resolve":
				data, _ := protocol.MarshalResponse(r.ID, map[string]string{"status": "ok"})
				mu.Lock()
				conn.WriteMessage(websocket.TextMessage, data)
				mu.Unlock()

			case "chat.send":
				var p protocol.ChatSendParams
				_ = json.Unmarshal(r.Params, &p)

				resp := protocol.ChatEvent{
					RunID:      fmt.Sprintf("run-%d", time.Now().UnixNano()),
					SessionKey: p.SessionKey,
					Seq:        1,
					State:      "final",
					Message:    json.RawMessage(fmt.Sprintf("%q", "Mock response to: "+p.Message)),
				}
				data, _ := protocol.MarshalResponse(r.ID, resp)
				mu.Lock()
				conn.WriteMessage(websocket.TextMessage, data)
				mu.Unlock()

			case "chat.history":
				var p protocol.ChatHistoryParams
				_ = json.Unmarshal(r.Params, &p)

				history := map[string]any{
					"sessionKey": p.SessionKey,
					"messages": []map[string]string{
						{"role": "user", "content": "What is OpenClaw?"},
						{"role": "assistant", "content": "Mock response to: What is OpenClaw?"},
					},
				}
				data, _ := protocol.MarshalResponse(r.ID, history)
				mu.Lock()
				conn.WriteMessage(websocket.TextMessage, data)
				mu.Unlock()

			default:
				data, _ := protocol.MarshalErrorResponse(r.ID, protocol.ErrorPayload{
					Code:    "UNKNOWN_METHOD",
					Message: fmt.Sprintf("unknown method: %s", r.Method),
				})
				mu.Lock()
				conn.WriteMessage(websocket.TextMessage, data)
				mu.Unlock()
			}

		case protocol.FrameTypeEvent:
			var ev protocol.Event
			if json.Unmarshal(msg, &ev) == nil {
				log.Printf("[ws] event: %s", ev.EventName)
			}

		case protocol.FrameTypeInvokeResponse:
			var res protocol.InvokeResponse
			if json.Unmarshal(msg, &res) == nil {
				log.Printf("[ws] invoke-response: id=%s ok=%v", res.ID, res.OK)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Chat Completions API
// ---------------------------------------------------------------------------

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	var req struct {
		Model    string `json:"model"`
		Stream   bool   `json:"stream"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Printf("[chat] model=%s stream=%v messages=%d", req.Model, req.Stream, len(req.Messages))

	if req.Stream {
		handleChatStream(w, req.Messages[len(req.Messages)-1].Content)
		return
	}

	resp := map[string]any{
		"id":      "chatcmpl-mock-1",
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   req.Model,
		"choices": []map[string]any{{
			"index": 0,
			"message": map[string]string{
				"role":    "assistant",
				"content": fmt.Sprintf("Mock response to: %s", req.Messages[len(req.Messages)-1].Content),
			},
			"finish_reason": "stop",
		}},
		"usage": map[string]int{
			"prompt_tokens":     10,
			"completion_tokens": 5,
			"total_tokens":      15,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleChatStream(w http.ResponseWriter, userMsg string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	words := []string{"Mock", " response", " to:", " " + userMsg}
	for i, word := range words {
		chunk := map[string]any{
			"id":      "chatcmpl-mock-1",
			"object":  "chat.completion.chunk",
			"created": time.Now().Unix(),
			"model":   "mock-model",
			"choices": []map[string]any{{
				"index": 0,
				"delta": map[string]string{"content": word},
			}},
		}
		if i == len(words)-1 {
			chunk["choices"].([]map[string]any)[0]["finish_reason"] = "stop"
		}
		data, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Fprint(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// ---------------------------------------------------------------------------
// Tools Invoke API
// ---------------------------------------------------------------------------

func handleToolsInvoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	var req struct {
		Tool       string         `json:"tool"`
		Action     string         `json:"action"`
		Args       map[string]any `json:"args"`
		SessionKey string         `json:"sessionKey"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Printf("[tools] tool=%s action=%s", req.Tool, req.Action)

	switch req.Tool {
	case "sessions_list":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"ok":     true,
			"result": []map[string]string{{"key": "main", "status": "active"}},
		})
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"ok": false,
			"error": map[string]string{
				"type":    "NOT_FOUND",
				"message": fmt.Sprintf("tool %q not found", req.Tool),
			},
		})
	}
}
