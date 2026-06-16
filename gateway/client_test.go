package gateway

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/a3tai/openclaw-go/protocol"
	"github.com/gorilla/websocket"
)

// mockGateway is a test helper that implements the OpenClaw gateway handshake.
type mockGateway struct {
	t        *testing.T
	upgrader websocket.Upgrader
	conns    []*websocket.Conn
	mu       sync.Mutex

	// ready is closed when the mock has entered its read loop.
	ready chan struct{}

	// onConnect is called with the connect params after the handshake.
	onConnect func(protocol.ConnectParams)

	// onRequest is called for any non-connect request.
	onRequest func(*websocket.Conn, protocol.Request)

	// rejectConnect causes the gateway to reject the connect request.
	rejectConnect bool

	// rejectWithError causes the rejection to include an error payload.
	rejectWithError *protocol.ErrorPayload

	// sendBadChallenge causes the gateway to send a non-challenge first frame.
	sendBadChallenge bool

	// sendGarbage causes the gateway to send invalid JSON as the first frame.
	sendGarbage bool

	// sendBadHelloPayload causes the gateway to send unparseable hello-ok payload.
	sendBadHelloPayload bool

	// sendMismatchedID causes the response to have a wrong ID.
	sendMismatchedID bool
}

// waitReady blocks until the mock has entered its read loop.
func (m *mockGateway) waitReady(t *testing.T) {
	t.Helper()
	select {
	case <-m.ready:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for mock gateway ready")
	}
}

func newMockGateway(t *testing.T) *mockGateway {
	return &mockGateway{
		t:        t,
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		ready:    make(chan struct{}, 1),
	}
}

func (m *mockGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.t.Logf("upgrade error: %v", err)
		return
	}
	m.mu.Lock()
	m.conns = append(m.conns, conn)
	m.mu.Unlock()

	// Send garbage instead of challenge.
	if m.sendGarbage {
		conn.WriteMessage(websocket.TextMessage, []byte("not json"))
		return
	}

	// Send a non-challenge event.
	if m.sendBadChallenge {
		evData, _ := protocol.MarshalEvent("some.other.event", map[string]string{"foo": "bar"})
		conn.WriteMessage(websocket.TextMessage, evData)
		return
	}

	// 1. Send connect.challenge.
	challenge := protocol.ConnectChallenge{Nonce: "test-nonce", Ts: time.Now().UnixMilli()}
	evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
	conn.WriteMessage(websocket.TextMessage, evData)

	// 2. Read the connect request.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return
	}
	var req protocol.Request
	if err := json.Unmarshal(msg, &req); err != nil {
		m.t.Logf("unmarshal req: %v", err)
		return
	}
	if req.Method != "connect" {
		m.t.Logf("expected connect, got %s", req.Method)
		return
	}

	var params protocol.ConnectParams
	json.Unmarshal(req.Params, &params)
	if m.onConnect != nil {
		m.onConnect(params)
	}

	// Reject connect.
	if m.rejectConnect {
		if m.rejectWithError != nil {
			respData, _ := protocol.MarshalErrorResponse(req.ID, *m.rejectWithError)
			conn.WriteMessage(websocket.TextMessage, respData)
		} else {
			// ok:false, no error detail.
			resp := protocol.Response{
				Type: protocol.FrameTypeResponse,
				ID:   req.ID,
				OK:   false,
			}
			data, _ := json.Marshal(resp)
			conn.WriteMessage(websocket.TextMessage, data)
		}
		return
	}

	// Mismatched ID.
	if m.sendMismatchedID {
		hello := protocol.HelloOK{
			Type:     "hello-ok",
			Protocol: protocol.ProtocolVersion,
			Policy:   protocol.HelloPolicy{TickIntervalMs: 15000},
		}
		respData, _ := protocol.MarshalResponse("wrong-id", hello)
		conn.WriteMessage(websocket.TextMessage, respData)
		return
	}

	// Bad hello-ok payload: send raw bytes with a payload that is valid JSON
	// but wrong type for HelloOK (a string instead of object).
	if m.sendBadHelloPayload {
		raw := fmt.Sprintf(`{"type":"res","id":%q,"ok":true,"payload":"not an object"}`, req.ID)
		conn.WriteMessage(websocket.TextMessage, []byte(raw))
		return
	}

	// 3. Send hello-ok response.
	hello := protocol.HelloOK{
		Type:     "hello-ok",
		Protocol: protocol.ProtocolVersion,
		Policy:   protocol.HelloPolicy{TickIntervalMs: 15000},
	}
	respData, _ := protocol.MarshalResponse(req.ID, hello)
	conn.WriteMessage(websocket.TextMessage, respData)

	// Signal that the mock is ready (done writing, entering read loop).
	select {
	case m.ready <- struct{}{}:
	default:
	}

	// 4. Serve subsequent requests.
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var r protocol.Request
		if json.Unmarshal(msg, &r) == nil && r.Type == protocol.FrameTypeRequest {
			if m.onRequest != nil {
				m.onRequest(conn, r)
			}
		}
	}
}

func startMockGateway(t *testing.T) (*mockGateway, string, func()) {
	mg := newMockGateway(t)
	srv := httptest.NewServer(mg)
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	return mg, wsURL, srv.Close
}

// --- Basic connect ---

func TestConnect(t *testing.T) {
	var gotParams protocol.ConnectParams
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onConnect = func(p protocol.ConnectParams) {
		gotParams = p
	}

	client := NewClient(
		WithToken("test-token"),
		WithConnectTimeout(5*time.Second),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	hello := client.Hello()
	if hello == nil {
		t.Fatal("hello is nil")
	}
	if hello.Protocol != protocol.ProtocolVersion {
		t.Errorf("protocol = %d, want %d", hello.Protocol, protocol.ProtocolVersion)
	}
	if gotParams.MinProtocol != protocol.ProtocolVersion {
		t.Errorf("minProtocol = %d, want %d", gotParams.MinProtocol, protocol.ProtocolVersion)
	}
	if gotParams.Role != protocol.RoleOperator {
		t.Errorf("role = %q, want %q", gotParams.Role, protocol.RoleOperator)
	}
	if gotParams.Auth.Token != "test-token" {
		t.Errorf("auth.token = %q, want %q", gotParams.Auth.Token, "test-token")
	}
}

// --- Connect with password auth ---

func TestConnectWithPassword(t *testing.T) {
	var gotParams protocol.ConnectParams
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onConnect = func(p protocol.ConnectParams) {
		gotParams = p
	}

	client := NewClient(
		WithPassword("my-password"),
		WithConnectTimeout(5*time.Second),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if gotParams.Auth.Password != "my-password" {
		t.Errorf("auth.password = %q, want 'my-password'", gotParams.Auth.Password)
	}
}

func TestConnectWithBootstrapToken(t *testing.T) {
	var gotParams protocol.ConnectParams
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onConnect = func(p protocol.ConnectParams) {
		gotParams = p
	}

	client := NewClient(
		WithBootstrapToken("setup-code-bootstrap"),
		WithDevice(protocol.DeviceIdentity{ID: "dev-1", PublicKey: "pk", Signature: "sig", SignedAt: 1}),
		WithConnectTimeout(5*time.Second),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Errorf("Connect: %v", err)
		return
	}
	if gotParams.Auth.BootstrapToken != "setup-code-bootstrap" {
		t.Errorf("auth.bootstrapToken = %q, want 'setup-code-bootstrap'", gotParams.Auth.BootstrapToken)
	}
	// A shared token may ride alongside the bootstrap token, but here none
	// was set, so it must stay empty rather than be clobbered.
	if gotParams.Auth.Token != "" {
		t.Errorf("auth.token = %q, want empty", gotParams.Auth.Token)
	}
}

// --- Connect with device identity ---

func TestConnectWithDevice(t *testing.T) {
	var gotParams protocol.ConnectParams
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onConnect = func(p protocol.ConnectParams) {
		gotParams = p
	}

	client := NewClient(
		WithToken("tok"),
		WithDevice(protocol.DeviceIdentity{
			ID:        "dev-42",
			PublicKey: "pk",
			Signature: "sig",
			SignedAt:  12345,
		}),
		WithConnectTimeout(5*time.Second),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if gotParams.Device == nil {
		t.Fatal("device is nil")
	}
	if gotParams.Device.ID != "dev-42" {
		t.Errorf("device.id = %q, want 'dev-42'", gotParams.Device.ID)
	}
	// Nonce should be set from the challenge.
	if gotParams.Device.Nonce != "test-nonce" {
		t.Errorf("device.nonce = %q, want 'test-nonce'", gotParams.Device.Nonce)
	}
}

// --- All option setters ---

func TestAllOptions(t *testing.T) {
	var gotParams protocol.ConnectParams
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onConnect = func(p protocol.ConnectParams) {
		gotParams = p
	}

	client := NewClient(
		WithToken("tok"),
		WithScopes(protocol.ScopeOperatorAdmin, protocol.ScopeOperatorApprovals),
		WithLocale("de-DE"),
		WithUserAgent("test-agent/2.0"),
		WithTLSConfig(&tls.Config{InsecureSkipVerify: true}), // not used for non-TLS, but exercises the setter
		WithConnectTimeout(5*time.Second),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if len(gotParams.Scopes) != 2 || gotParams.Scopes[0] != protocol.ScopeOperatorAdmin {
		t.Errorf("scopes = %v", gotParams.Scopes)
	}
	if gotParams.Locale != "de-DE" {
		t.Errorf("locale = %q, want 'de-DE'", gotParams.Locale)
	}
	if gotParams.UserAgent != "test-agent/2.0" {
		t.Errorf("userAgent = %q, want 'test-agent/2.0'", gotParams.UserAgent)
	}
}

// --- Connect failure: bad URL ---

func TestConnectDialError(t *testing.T) {
	client := NewClient(WithConnectTimeout(1 * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Connect(ctx, "ws://127.0.0.1:1") // nothing listening
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "dial") {
		t.Errorf("error = %q, want to contain 'dial'", err.Error())
	}
}

// --- Connect failure: bad challenge (garbage JSON) ---

func TestConnectBadChallengeGarbage(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()
	mg.sendGarbage = true

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "challenge") {
		t.Errorf("error = %q, want to contain 'challenge'", err.Error())
	}
}

// --- Connect failure: wrong event type instead of challenge ---

func TestConnectBadChallengeWrongEvent(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()
	mg.sendBadChallenge = true

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "expected connect.challenge") {
		t.Errorf("error = %q, want to contain 'expected connect.challenge'", err.Error())
	}
}

// --- Connect rejected with error payload ---

func TestConnectRejectedWithError(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()
	mg.rejectConnect = true
	mg.rejectWithError = &protocol.ErrorPayload{Code: "AUTH_FAILED", Message: "invalid token"}

	client := NewClient(WithToken("bad-tok"), WithConnectTimeout(5*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "AUTH_FAILED") {
		t.Errorf("error = %q, want to contain 'AUTH_FAILED'", err.Error())
	}
}

// --- Connect rejected without error detail ---

func TestConnectRejectedNoDetail(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()
	mg.rejectConnect = true

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "no error details") {
		t.Errorf("error = %q, want to contain 'no error details'", err.Error())
	}
}

// --- Connect: mismatched response ID ---

func TestConnectMismatchedID(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()
	mg.sendMismatchedID = true

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Errorf("error = %q, want to contain 'mismatch'", err.Error())
	}
}

// --- Connect: bad hello-ok payload ---

func TestConnectBadHelloPayload(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()
	mg.sendBadHelloPayload = true

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "hello") {
		t.Errorf("error = %q, want to contain 'hello'", err.Error())
	}
}

// --- Send and Presence ---

func TestSendRequest(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "system-presence" {
			// Real server returns a flat array of presence entries.
			entries := []protocol.PresenceEntry{
				{DeviceID: "dev-1", Roles: []string{"operator"}, Ts: 1},
			}
			respData, _ := protocol.MarshalResponse(req.ID, entries)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	presence, err := client.Presence(ctx)
	if err != nil {
		t.Fatalf("Presence: %v", err)
	}
	if len(presence) != 1 {
		t.Fatalf("presence len = %d, want 1", len(presence))
	}
	if presence[0].DeviceID != "dev-1" {
		t.Errorf("deviceId = %q, want %q", presence[0].DeviceID, "dev-1")
	}
}

// --- Presence error with error payload ---

func TestPresenceErrorPayload(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "system-presence" {
			respData, _ := protocol.MarshalErrorResponse(req.ID, protocol.ErrorPayload{
				Code:    "FORBIDDEN",
				Message: "not allowed",
			})
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	_, err := client.Presence(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "FORBIDDEN") {
		t.Errorf("error = %q, want to contain 'FORBIDDEN'", err.Error())
	}
}

// --- Presence error without error payload (ok:false, no error detail) ---

func TestPresenceErrorNoDetail(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "system-presence" {
			resp := protocol.Response{
				Type: protocol.FrameTypeResponse,
				ID:   req.ID,
				OK:   false,
			}
			data, _ := json.Marshal(resp)
			conn.WriteMessage(websocket.TextMessage, data)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	_, err := client.Presence(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "request failed") {
		t.Errorf("error = %q, want to contain 'request failed'", err.Error())
	}
}

// --- Presence: unmarshal error (bad payload) ---

func TestPresenceUnmarshalError(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "system-presence" {
			// Send raw JSON with an invalid payload that will fail
			// json.Unmarshal into map[string]PresenceEntry.
			// Use a valid JSON response but payload that is wrong type (string instead of object).
			raw := fmt.Sprintf(`{"type":"res","id":%q,"ok":true,"payload":"not a map"}`, req.ID)
			conn.WriteMessage(websocket.TextMessage, []byte(raw))
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	_, err := client.Presence(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("error = %q, want to contain 'unmarshal'", err.Error())
	}
}

// --- Exec approval resolve ---

func TestExecApprovalResolve(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "exec.approval.resolve" {
			respData, _ := protocol.MarshalResponse(req.ID, map[string]string{"status": "ok"})
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	_, err := client.ExecApprovalResolve(ctx, protocol.ExecApprovalResolveParams{
		ID:       "approval-1",
		Decision: "approved",
	})
	if err != nil {
		t.Fatalf("ExecApprovalResolve: %v", err)
	}
}

// --- Exec approval resolve: error with payload ---

func TestExecApprovalResolveError(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "exec.approval.resolve" {
			respData, _ := protocol.MarshalErrorResponse(req.ID, protocol.ErrorPayload{
				Code: "FORBIDDEN", Message: "not allowed",
			})
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	_, err := client.ExecApprovalResolve(ctx, protocol.ExecApprovalResolveParams{
		ID: "approval-1", Decision: "approved",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- readLoop: unparseable frame (continues) ---

func TestReadLoopBadFrame(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	received := make(chan struct{})
	client := NewClient(
		WithToken("tok"),
		WithConnectTimeout(5*time.Second),
		WithOnEvent(func(ev protocol.Event) {
			if ev.EventName == "good-event" {
				close(received)
			}
		}),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	// Send garbage, then a valid event. The readLoop should skip the garbage.
	conn.WriteMessage(websocket.TextMessage, []byte("not json"))
	evData, _ := protocol.MarshalEvent("good-event", map[string]string{})
	conn.WriteMessage(websocket.TextMessage, evData)

	select {
	case <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("did not receive event after bad frame")
	}
}

// --- readLoop: response for unknown request ID (no pending) ---

func TestReadLoopUnmatchedResponse(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	received := make(chan struct{})
	client := NewClient(
		WithToken("tok"),
		WithConnectTimeout(5*time.Second),
		WithOnEvent(func(ev protocol.Event) {
			if ev.EventName == "after-unmatched" {
				close(received)
			}
		}),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	// Send a response with an ID that nobody is waiting for.
	respData, _ := protocol.MarshalResponse("unknown-id", map[string]string{})
	conn.WriteMessage(websocket.TextMessage, respData)

	// Then send an event to confirm the loop is still running.
	evData, _ := protocol.MarshalEvent("after-unmatched", map[string]string{})
	conn.WriteMessage(websocket.TextMessage, evData)

	select {
	case <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("readLoop died after unmatched response")
	}
}

// --- Hello before connect returns nil ---

func TestHelloBeforeConnect(t *testing.T) {
	client := NewClient(WithToken("tok"))
	if client.Hello() != nil {
		t.Error("Hello() should be nil before Connect")
	}
}

// --- Hello with zero tick interval (doesn't override default) ---

func TestHelloZeroTickInterval(t *testing.T) {
	mg := newMockGateway(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		challenge := protocol.ConnectChallenge{Nonce: "n", Ts: 1}
		evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
		conn.WriteMessage(websocket.TextMessage, evData)

		_, msg, _ := conn.ReadMessage()
		var req protocol.Request
		json.Unmarshal(msg, &req)

		hello := protocol.HelloOK{
			Type:     "hello-ok",
			Protocol: protocol.ProtocolVersion,
			Policy:   protocol.HelloPolicy{TickIntervalMs: 0},
		}
		respData, _ := protocol.MarshalResponse(req.ID, hello)
		conn.WriteMessage(websocket.TextMessage, respData)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	if client.Hello().Policy.TickIntervalMs != 0 {
		t.Errorf("tickIntervalMs = %d, want 0", client.Hello().Policy.TickIntervalMs)
	}
}

// --- tickLoop fires pings ---

func TestTickLoopFires(t *testing.T) {
	mg := newMockGateway(t)
	pingReceived := make(chan struct{}, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// Set a ping handler on the server side to detect pings.
		conn.SetPingHandler(func(appData string) error {
			select {
			case pingReceived <- struct{}{}:
			default:
			}
			return conn.WriteMessage(websocket.PongMessage, []byte(appData))
		})

		challenge := protocol.ConnectChallenge{Nonce: "n", Ts: 1}
		evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
		conn.WriteMessage(websocket.TextMessage, evData)

		_, msg, _ := conn.ReadMessage()
		var req protocol.Request
		json.Unmarshal(msg, &req)

		// Send hello-ok with a very short tick interval (50ms).
		hello := protocol.HelloOK{
			Type:     "hello-ok",
			Protocol: protocol.ProtocolVersion,
			Policy:   protocol.HelloPolicy{TickIntervalMs: 50},
		}
		respData, _ := protocol.MarshalResponse(req.ID, hello)
		conn.WriteMessage(websocket.TextMessage, respData)

		// Keep reading to handle pings.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	select {
	case <-pingReceived:
		// Tick fired successfully.
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for tick ping")
	}
}

// --- readLoop handles bad challenge payload in connect ---

func TestReadChallengeUnmarshalPayloadError(t *testing.T) {
	mg := newMockGateway(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Send a connect.challenge event with invalid payload.
		raw := `{"type":"event","event":"connect.challenge","payload":"not an object"}`
		conn.WriteMessage(websocket.TextMessage, []byte(raw))
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	client := NewClient(WithToken("tok"), WithConnectTimeout(2*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unmarshal challenge payload") {
		t.Errorf("error = %q, want to contain 'unmarshal challenge payload'", err.Error())
	}
}

// --- readHelloOK receives invalid JSON response ---

func TestReadHelloOKUnmarshalError(t *testing.T) {
	mg := newMockGateway(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		challenge := protocol.ConnectChallenge{Nonce: "n", Ts: 1}
		evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
		conn.WriteMessage(websocket.TextMessage, evData)

		// Read connect request.
		conn.ReadMessage()

		// Send invalid JSON as response.
		conn.WriteMessage(websocket.TextMessage, []byte("not json"))
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	client := NewClient(WithToken("tok"), WithConnectTimeout(2*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unmarshal response") {
		t.Errorf("error = %q, want to contain 'unmarshal response'", err.Error())
	}
}

// --- readLoop: event without onEvent handler (no panic) ---

func TestReadLoopEventNoHandler(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	// No onEvent handler.
	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	// Send event — should not panic even with no handler.
	evData, _ := protocol.MarshalEvent("some.event", map[string]string{})
	conn.WriteMessage(websocket.TextMessage, evData)

	// Give time for the readLoop to process.
	time.Sleep(100 * time.Millisecond)
}

// --- readLoop: event handler is async (slow handler doesn't block) ---

func TestReadLoopEventAsync(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	slowBlock := make(chan struct{})
	fastSeen := make(chan struct{})

	client := NewClient(
		WithToken("tok"),
		WithConnectTimeout(5*time.Second),
		WithOnEvent(func(ev protocol.Event) {
			switch ev.EventName {
			case "slow":
				<-slowBlock
			case "fast":
				select {
				case <-fastSeen:
				default:
					close(fastSeen)
				}
			}
		}),
	)
	defer client.Close()
	defer close(slowBlock)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	slowData, _ := protocol.MarshalEvent("slow", map[string]string{})
	fastData, _ := protocol.MarshalEvent("fast", map[string]string{})
	conn.WriteMessage(websocket.TextMessage, slowData)
	conn.WriteMessage(websocket.TextMessage, fastData)

	select {
	case <-fastSeen:
	case <-time.After(2 * time.Second):
		t.Fatal("fast event not handled while slow handler blocked")
	}
}

// --- readLoop: invoke without onInvoke handler (no panic) ---

func TestReadLoopInvokeNoHandler(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	// No onInvoke handler.
	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	inv := protocol.Invoke{Type: "invoke", ID: "inv-1", Command: "test"}
	data, _ := json.Marshal(inv)
	conn.WriteMessage(websocket.TextMessage, data)

	time.Sleep(100 * time.Millisecond)
}

// --- readLoop: FrameTypeRequest is handled (dead path, but covered) ---

func TestReadLoopRequestFrame(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	received := make(chan struct{})
	client := NewClient(
		WithToken("tok"),
		WithConnectTimeout(5*time.Second),
		WithOnEvent(func(ev protocol.Event) {
			if ev.EventName == "after-req" {
				close(received)
			}
		}),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	// Send a frame with type "req" — goes to the FrameTypeRequest case.
	conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"req","id":"x","method":"test","params":{}}`))

	// Then send an event to verify readLoop is still alive.
	evData, _ := protocol.MarshalEvent("after-req", map[string]string{})
	conn.WriteMessage(websocket.TextMessage, evData)

	select {
	case <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("readLoop died after req frame")
	}
}

// --- tickLoop: connection closed while ticking ---

func TestTickLoopConnectionClosed(t *testing.T) {
	mg := newMockGateway(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		challenge := protocol.ConnectChallenge{Nonce: "n", Ts: 1}
		evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
		conn.WriteMessage(websocket.TextMessage, evData)

		_, msg, _ := conn.ReadMessage()
		var req protocol.Request
		json.Unmarshal(msg, &req)

		hello := protocol.HelloOK{
			Type:     "hello-ok",
			Protocol: protocol.ProtocolVersion,
			Policy:   protocol.HelloPolicy{TickIntervalMs: 50},
		}
		respData, _ := protocol.MarshalResponse(req.ID, hello)
		conn.WriteMessage(websocket.TextMessage, respData)

		// Close server side immediately to trigger tick write error.
		time.Sleep(30 * time.Millisecond)
		conn.Close()
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Wait for done to be signaled (from readLoop or tickLoop).
	select {
	case <-client.Done():
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for done after tick write error")
	}
}

// --- Send: marshal error ---

func TestSendMarshalError(t *testing.T) {
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// func values cannot be marshaled.
	_, err := client.Send(ctx, "test", func() {})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Send: write error (connection closed before write) ---

func TestSendWriteError(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Close the server-side connection to cause a write error.
	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()
	conn.Close()

	// Wait for readLoop to detect the close.
	time.Sleep(100 * time.Millisecond)

	_, err := client.Send(ctx, "test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Presence: Send returns error (connection closed) ---

func TestPresenceSendError(t *testing.T) {
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Close client first, so Send will return "client closed".
	client.Close()

	_, err := client.Presence(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ExecApprovalResolve: Send returns error (connection closed) ---

func TestExecApprovalResolveSendError(t *testing.T) {
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Close client first.
	client.Close()

	_, err := client.ExecApprovalResolve(ctx, protocol.ExecApprovalResolveParams{
		ID: "approval-1", Decision: "approved",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Connect: marshalRequest error ---

func TestConnectMarshalRequestError(t *testing.T) {
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	// Inject error into the per-client marshalRequest.
	client.opts.marshalRequest = func(id, method string, params any) ([]byte, error) {
		return nil, fmt.Errorf("injected marshal error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "marshal connect") {
		t.Errorf("error = %q, want to contain 'marshal connect'", err.Error())
	}
}

// --- Connect: WriteMessage error after challenge ---

func TestConnectWriteError(t *testing.T) {
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	// Override marshalRequest to succeed but close the connection before returning,
	// so the subsequent WriteMessage fails.
	client.opts.marshalRequest = func(id, method string, params any) ([]byte, error) {
		data, err := protocol.MarshalRequest(id, method, params)
		if err != nil {
			return nil, err
		}
		// Close the underlying connection to force WriteMessage to fail.
		client.connMu.Lock()
		client.conn.Close()
		client.connMu.Unlock()
		return data, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "send connect") {
		t.Errorf("error = %q, want to contain 'send connect'", err.Error())
	}
}

// --- readChallenge: server closes connection before sending anything ---

func TestReadChallengeReadError(t *testing.T) {
	mg := newMockGateway(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Close immediately — readChallenge's ReadMessage will fail.
		conn.Close()
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	client := NewClient(WithToken("tok"), WithConnectTimeout(2*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "read challenge") {
		t.Errorf("error = %q, want to contain 'read challenge'", err.Error())
	}
}

// --- readHelloOK: server closes connection after challenge+connect ---

func TestReadHelloOKReadError(t *testing.T) {
	mg := newMockGateway(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		challenge := protocol.ConnectChallenge{Nonce: "n", Ts: 1}
		evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
		conn.WriteMessage(websocket.TextMessage, evData)

		// Read the connect request, then close.
		conn.ReadMessage()
		conn.Close()
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	client := NewClient(WithToken("tok"), WithConnectTimeout(2*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Connect(ctx, wsURL)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "hello") {
		t.Errorf("error = %q, want to contain 'hello'", err.Error())
	}
}

// --- Send: WriteMessage error (close underlying conn, but keep done open) ---

func TestSendWriteErrorDirect(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		// Don't respond.
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Close the underlying websocket directly (bypass client.Close so done stays open).
	client.connMu.Lock()
	client.conn.Close()
	client.connMu.Unlock()

	_, err := client.Send(ctx, "test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "write") {
		t.Errorf("error = %q, want to contain 'write'", err.Error())
	}
}

// --- Invoke handler: json.Marshal error ---

func TestInvokeHandlerMarshalError(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	marshalCalled := make(chan struct{}, 1)
	invokeDone := make(chan struct{})

	client := NewClient(
		WithToken("tok"),
		WithRole(protocol.RoleNode),
		WithConnectTimeout(5*time.Second),
		WithOnInvoke(func(inv protocol.Invoke) protocol.InvokeResponse {
			defer func() { close(invokeDone) }()
			return protocol.InvokeResponse{OK: true}
		}),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Inject marshal error into the per-client marshalJSON (race-safe: set before invoke arrives).
	client.opts.marshalJSON = func(v any) ([]byte, error) {
		select {
		case marshalCalled <- struct{}{}:
		default:
		}
		return nil, fmt.Errorf("injected marshal error")
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	inv := protocol.Invoke{Type: "invoke", ID: "inv-1", Command: "test"}
	data, _ := json.Marshal(inv)
	conn.WriteMessage(websocket.TextMessage, data)

	// Wait for the marshal call and invoke handler to complete.
	select {
	case <-marshalCalled:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for marshal call")
	}

	select {
	case <-invokeDone:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for invoke handler")
	}
}

// --- Invoke handler: client closed during invoke ---

func TestInvokeHandlerClientClosed(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	invokeDone := make(chan struct{})
	var client *Client

	client = NewClient(
		WithToken("tok"),
		WithRole(protocol.RoleNode),
		WithConnectTimeout(5*time.Second),
		WithOnInvoke(func(inv protocol.Invoke) protocol.InvokeResponse {
			// Wait for client to close before returning, so the done check triggers.
			<-client.Done()
			defer func() { close(invokeDone) }()
			return protocol.InvokeResponse{OK: true, Payload: json.RawMessage(`{}`)}
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mg.waitReady(t)
	mg.mu.Lock()
	conn := mg.conns[len(mg.conns)-1]
	mg.mu.Unlock()

	inv := protocol.Invoke{Type: "invoke", ID: "inv-1", Command: "test"}
	data, _ := json.Marshal(inv)
	conn.WriteMessage(websocket.TextMessage, data)

	// Give time for the invoke handler to start.
	time.Sleep(50 * time.Millisecond)

	// Close the client, which closes c.done.
	client.Close()

	select {
	case <-invokeDone:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for invoke handler")
	}
}

// --- tickLoop: write error ---
// This test exercises the tick write error path by creating a client,
// connecting it, then directly calling tickLoop on a copy of the client
// with a closed connection and an open done channel.

func TestTickLoopWriteError(t *testing.T) {
	// Create a connected client so we can steal its conn.
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Create a separate client instance for the isolated tickLoop test.
	tc := &Client{
		opts: options{
			tickInterval: 10 * time.Millisecond,
		},
		conn:         client.conn,
		done:         make(chan struct{}),
		readLoopDone: make(chan struct{}),
		tickStop:     make(chan struct{}),
	}

	// Close the underlying TCP connection so the tick write will fail.
	tc.conn.UnderlyingConn().Close()

	// Run tickLoop — it should return quickly due to write error.
	done := make(chan struct{})
	go func() {
		tc.tickLoop()
		close(done)
	}()

	select {
	case <-done:
		// tickLoop exited due to write error — success.
	case <-time.After(3 * time.Second):
		t.Fatal("tickLoop did not exit after write error")
	}

	client.Close()
}

// --- SendEvent marshal error ---

func TestSendEventMarshalError(t *testing.T) {
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// func values cannot be marshaled to JSON.
	err := client.SendEvent("test", func() {})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- SendEvent: successful write ---

func TestSendEventSuccess(t *testing.T) {
	_, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	if err := client.SendEvent("test.event", map[string]string{"key": "val"}); err != nil {
		t.Fatalf("SendEvent: %v", err)
	}
}

// --- Send: context cancelled while waiting for response ---

func TestSendContextCancelled(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	// Don't reply to requests — let them hang.
	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Use a short-lived context for Send.
	sendCtx, sendCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer sendCancel()

	_, err := client.Send(sendCtx, "test.method", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("error = %v, want context.DeadlineExceeded", err)
	}
}

// --- Close when conn is nil (never connected) ---

func TestCloseWithoutConnect(t *testing.T) {
	client := NewClient(WithToken("tok"))
	if err := client.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	// Second close should be a no-op (already closed).
	if err := client.Close(); err != nil {
		t.Fatalf("Close again: %v", err)
	}
}

// --- Invoke handler: successful response write ---

func TestInvokeHandlerSuccess(t *testing.T) {
	mg := newMockGateway(t)
	invokeResult := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := mg.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// 1. Send challenge
		challenge := protocol.ConnectChallenge{Nonce: "n", Ts: 1}
		evData, _ := protocol.MarshalEvent("connect.challenge", challenge)
		conn.WriteMessage(websocket.TextMessage, evData)

		// 2. Read connect request
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var req protocol.Request
		json.Unmarshal(msg, &req)

		// 3. Send hello-ok
		hello := protocol.HelloOK{
			Type:     "hello-ok",
			Protocol: protocol.ProtocolVersion,
			Policy:   protocol.HelloPolicy{TickIntervalMs: 30000},
		}
		respData, _ := protocol.MarshalResponse(req.ID, hello)
		conn.WriteMessage(websocket.TextMessage, respData)

		// 4. Send an invoke frame
		inv := protocol.Invoke{Type: "invoke", ID: "inv-1", Command: "test.cmd"}
		invData, _ := json.Marshal(inv)
		conn.WriteMessage(websocket.TextMessage, invData)

		// 5. Read messages until we get the invoke-res
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var frame struct {
				Type string `json:"type"`
			}
			if json.Unmarshal(msg, &frame) == nil && frame.Type == "invoke-res" {
				close(invokeResult)
				return
			}
		}
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")

	client := NewClient(
		WithToken("tok"),
		WithRole(protocol.RoleNode),
		WithConnectTimeout(5*time.Second),
		WithOnInvoke(func(inv protocol.Invoke) protocol.InvokeResponse {
			return protocol.InvokeResponse{OK: true, Payload: json.RawMessage(`{"result":"done"}`)}
		}),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	select {
	case <-invokeResult:
		// Invoke response was written back successfully.
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for invoke response")
	}
}

// --- Option setters: WithClientInfo, WithCaps, WithCommands, WithPermissions ---

func TestOptionSetters(t *testing.T) {
	var gotParams protocol.ConnectParams
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onConnect = func(p protocol.ConnectParams) {
		gotParams = p
	}

	client := NewClient(
		WithToken("tok"),
		WithClientInfo(protocol.ClientInfo{
			ID:       protocol.ClientIDCLI,
			Version:  "1.0.0",
			Platform: "test",
			Mode:     protocol.ClientModeCLI,
		}),
		WithCaps("exec", "files"),
		WithCommands("ls", "cat"),
		WithPermissions(map[string]bool{"exec": true, "net": false}),
		WithConnectTimeout(5*time.Second),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	if gotParams.Client.ID != protocol.ClientIDCLI {
		t.Errorf("client.id = %q, want %q", gotParams.Client.ID, protocol.ClientIDCLI)
	}
	if len(gotParams.Caps) != 2 || gotParams.Caps[0] != "exec" {
		t.Errorf("caps = %v", gotParams.Caps)
	}
	if len(gotParams.Commands) != 2 || gotParams.Commands[0] != "ls" {
		t.Errorf("commands = %v", gotParams.Commands)
	}
	if len(gotParams.Permissions) != 2 || !gotParams.Permissions["exec"] {
		t.Errorf("permissions = %v", gotParams.Permissions)
	}
}
