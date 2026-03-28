// Command client demonstrates all three OpenClaw client APIs:
//
//  1. WebSocket Gateway protocol (connect, presence, events, approvals)
//  2. OpenAI-compatible Chat Completions (non-streaming and streaming)
//  3. Tools Invoke HTTP API
//
// Usage:
//
//	go run ./examples/client <token> <host> [identity-dir]
//
// The host URL is used for both WebSocket (gateway) and HTTP (chat/tools) APIs.
// For WebSocket, the scheme is converted to ws/wss automatically.
// For HTTP, the scheme is converted to http/https automatically.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/a3tai/openclaw-go/chatcompletions"
	"github.com/a3tai/openclaw-go/examples/internal/gwconn"
	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
	"github.com/a3tai/openclaw-go/toolsinvoke"
)

func main() {
	cfg := gwconn.ParseArgs("Usage: client <token> <host> [identity-dir]")
	cfg.PrintIdentityInfo()

	// Derive an HTTP URL from the WebSocket URL for the REST APIs.
	httpURL := cfg.WSURL
	httpURL = strings.Replace(httpURL, "wss://", "https://", 1)
	httpURL = strings.Replace(httpURL, "ws://", "http://", 1)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println()
	fmt.Println("=== OpenClaw Go Client Example ===")
	fmt.Println()

	// 1. Gateway WebSocket API
	fmt.Println("--- 1. WebSocket Gateway ---")
	demonstrateGateway(ctx, cfg)
	fmt.Println()

	// 2. Chat Completions API
	fmt.Println("--- 2. Chat Completions ---")
	demonstrateChatCompletions(ctx, httpURL, cfg.Token)
	fmt.Println()

	// 3. Tools Invoke API
	fmt.Println("--- 3. Tools Invoke ---")
	demonstrateToolsInvoke(ctx, httpURL, cfg.Token)
	fmt.Println()

	fmt.Println("=== Done ===")
}

func demonstrateGateway(ctx context.Context, cfg *gwconn.Config) {
	client := cfg.NewClient(
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
			protocol.ScopeOperatorRead,
			protocol.ScopeOperatorWrite,
			protocol.ScopeOperatorApprovals,
		),
		gateway.WithLocale("en-US"),
		gateway.WithUserAgent("openclaw-go-example/1.0"),
		gateway.WithConnectTimeout(5*time.Second),
		gateway.WithOnEvent(func(ev protocol.Event) {
			fmt.Printf("  [event] %s\n", ev.EventName)
		}),
	)
	defer client.Close()

	cfg.Connect(ctx, client)

	hello := client.Hello()
	fmt.Printf("  Protocol: %d, TickInterval: %dms\n",
		hello.Protocol, hello.Policy.TickIntervalMs)

	// Fetch presence.
	fmt.Println("  Fetching presence...")
	presence, err := client.Presence(ctx)
	if err != nil {
		fmt.Printf("  Presence: %v\n", err)
	} else {
		for _, entry := range presence {
			fmt.Printf("  Presence: %s -> roles=%v\n", entry.DeviceID, entry.Roles)
		}
	}

	// Resolve an exec approval.
	fmt.Println("  Resolving exec approval...")
	_, err = client.ExecApprovalResolve(ctx, protocol.ExecApprovalResolveParams{
		ID:       "approval-example",
		Decision: "approved",
	})
	if err != nil {
		fmt.Printf("  ExecApprovalResolve: %v\n", err)
	} else {
		fmt.Println("  Approval resolved successfully")
	}

	// Send a custom event.
	fmt.Println("  Sending event...")
	err = client.SendEvent(string(protocol.EventExecFinished), protocol.ExecFinished{
		SessionKey: "main",
		RunID:      "run-example",
	})
	if err != nil {
		fmt.Printf("  SendEvent: %v\n", err)
	} else {
		fmt.Println("  Event sent successfully")
	}
}

func demonstrateChatCompletions(ctx context.Context, httpURL, token string) {
	client := &chatcompletions.Client{
		BaseURL:    httpURL,
		Token:      token,
		AgentID:    "main",
		SessionKey: "example-session",
	}

	// Non-streaming completion.
	fmt.Println("  Creating non-streaming completion...")
	resp, err := client.Create(ctx, chatcompletions.Request{
		Model: "openclaw:main",
		Messages: []chatcompletions.Message{
			{Role: "user", Content: "Hello, OpenClaw!"},
		},
	})
	if err != nil {
		fmt.Printf("  Create: %v\n", err)
		return
	}
	fmt.Printf("  Response: %s\n", resp.Choices[0].Message.Content)
	if resp.Usage != nil {
		fmt.Printf("  Usage: prompt=%d, completion=%d, total=%d\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	// Streaming completion.
	fmt.Println("  Creating streaming completion...")
	stream, err := client.CreateStream(ctx, chatcompletions.Request{
		Model: "openclaw:main",
		Messages: []chatcompletions.Message{
			{Role: "user", Content: "Tell me about Go"},
		},
	})
	if err != nil {
		fmt.Printf("  CreateStream: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Print("  Stream: ")
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("  Recv: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}

func demonstrateToolsInvoke(ctx context.Context, httpURL, token string) {
	client := &toolsinvoke.Client{
		BaseURL:        httpURL,
		Token:          token,
		MessageChannel: "cli",
	}

	// Invoke sessions_list tool.
	fmt.Println("  Invoking sessions_list tool...")
	resp, err := client.Invoke(ctx, toolsinvoke.Request{
		Tool:   "sessions_list",
		Action: "json",
	})
	if err != nil {
		fmt.Printf("  Invoke: %v\n", err)
		return
	}
	fmt.Printf("  OK: %v, Result: %s\n", resp.OK, string(resp.Result))

	// Invoke a nonexistent tool to show error handling.
	fmt.Println("  Invoking nonexistent tool (expecting error)...")
	resp, err = client.Invoke(ctx, toolsinvoke.Request{
		Tool: "nonexistent_tool",
	})
	if err != nil {
		fmt.Printf("  Expected error: %v\n", err)
		if resp != nil && resp.Error != nil {
			fmt.Printf("  Error detail: type=%s message=%s\n", resp.Error.Type, resp.Error.Message)
		}
	}

	// Invoke with args and dry run.
	fmt.Println("  Invoking with args and dry run...")
	resp, err = client.Invoke(ctx, toolsinvoke.Request{
		Tool:       "sessions_list",
		Action:     "json",
		Args:       map[string]any{"filter": "active"},
		SessionKey: "main",
		DryRun:     true,
	})
	if err != nil {
		fmt.Printf("  Invoke: %v\n", err)
		return
	}

	var result json.RawMessage
	if resp.Result != nil {
		result = resp.Result
	}
	fmt.Printf("  OK: %v, Result: %s\n", resp.OK, string(result))
}
