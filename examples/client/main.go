// Command client demonstrates all three OpenClaw client APIs:
//
//  1. WebSocket Gateway protocol (connect, presence, events, approvals)
//  2. OpenAI-compatible Chat Completions (non-streaming and streaming)
//  3. Tools Invoke HTTP API
//
// It expects the mock server to be running (go run ./examples/server).
//
// Usage:
//
//	go run ./examples/client
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/a3tai/openclaw-go/chatcompletions"
	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
	"github.com/a3tai/openclaw-go/toolsinvoke"
)

const (
	wsURL   = "ws://localhost:18789/ws"
	httpURL = "http://localhost:18789"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Go Client Example ===")
	fmt.Println()

	// 1. Gateway WebSocket API
	fmt.Println("--- 1. WebSocket Gateway ---")
	demonstrateGateway(ctx)
	fmt.Println()

	// 2. Chat Completions API
	fmt.Println("--- 2. Chat Completions ---")
	demonstrateChatCompletions(ctx)
	fmt.Println()

	// 3. Tools Invoke API
	fmt.Println("--- 3. Tools Invoke ---")
	demonstrateToolsInvoke(ctx)
	fmt.Println()

	fmt.Println("=== Done ===")
}

func demonstrateGateway(ctx context.Context) {
	// Create a client with operator role and event handler.
	client := gateway.NewClient(
		gateway.WithToken("example-token"),
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

	// Connect to the gateway.
	fmt.Println("  Connecting to gateway...")
	if err := client.Connect(ctx, wsURL); err != nil {
		log.Fatalf("  Connect: %v", err)
	}

	hello := client.Hello()
	fmt.Printf("  Connected! Protocol: %d, TickInterval: %dms\n",
		hello.Protocol, hello.Policy.TickIntervalMs)

	// Fetch presence.
	fmt.Println("  Fetching presence...")
	presence, err := client.Presence(ctx)
	if err != nil {
		log.Fatalf("  Presence: %v", err)
	}
	for id, entry := range presence {
		fmt.Printf("  Presence: %s -> roles=%v\n", id, entry.Roles)
	}

	// Resolve an exec approval.
	fmt.Println("  Resolving exec approval...")
	_, err = client.ExecApprovalResolve(ctx, protocol.ExecApprovalResolveParams{
		ID:       "approval-example",
		Decision: "approved",
	})
	if err != nil {
		log.Fatalf("  ExecApprovalResolve: %v", err)
	}
	fmt.Println("  Approval resolved successfully")

	// Send a custom event.
	fmt.Println("  Sending event...")
	err = client.SendEvent(string(protocol.EventExecFinished), protocol.ExecFinished{
		SessionKey: "main",
		RunID:      "run-example",
	})
	if err != nil {
		log.Fatalf("  SendEvent: %v", err)
	}
	fmt.Println("  Event sent successfully")
}

func demonstrateChatCompletions(ctx context.Context) {
	client := &chatcompletions.Client{
		BaseURL:    httpURL,
		Token:      "example-token",
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
		log.Fatalf("  Create: %v", err)
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
		log.Fatalf("  CreateStream: %v", err)
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

func demonstrateToolsInvoke(ctx context.Context) {
	client := &toolsinvoke.Client{
		BaseURL:        httpURL,
		Token:          "example-token",
		MessageChannel: "cli",
	}

	// Invoke sessions_list tool.
	fmt.Println("  Invoking sessions_list tool...")
	resp, err := client.Invoke(ctx, toolsinvoke.Request{
		Tool:   "sessions_list",
		Action: "json",
	})
	if err != nil {
		log.Fatalf("  Invoke: %v", err)
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
		log.Fatalf("  Invoke: %v", err)
	}

	var result json.RawMessage
	if resp.Result != nil {
		result = resp.Result
	}
	fmt.Printf("  OK: %v, Result: %s\n", resp.OK, string(result))
}
