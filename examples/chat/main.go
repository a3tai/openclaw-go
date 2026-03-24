// Command chat demonstrates an interactive chat session via the Gateway.
//
// It connects to the gateway, sends a chat message, and displays the
// streaming agent response events.
//
// Usage:
//
//	go run ./examples/chat
//
// Requires the mock server: go run ./examples/server
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Chat Example ===")
	fmt.Println()

	// Connect to the gateway as an operator.
	client := gateway.NewClient(
		gateway.WithToken("example-token"),
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(protocol.ScopeOperatorRead, protocol.ScopeOperatorWrite),
		gateway.WithOnEvent(func(ev protocol.Event) {
			// Handle chat events as they arrive.
			if ev.EventName == protocol.EventChat {
				var data map[string]any
				if json.Unmarshal(ev.Payload, &data) == nil {
					state, _ := data["state"].(string)
					switch state {
					case "delta":
						msg, _ := data["message"].(string)
						fmt.Print(msg)
					case "final":
						fmt.Println()
						fmt.Println("  [chat] Agent finished")
					case "error":
						errMsg, _ := data["errorMessage"].(string)
						fmt.Printf("  [chat] Error: %s\n", errMsg)
					}
				}
			}
		}),
	)
	defer client.Close()

	fmt.Println("Connecting to gateway...")
	if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
		log.Fatalf("Connect: %v", err)
	}

	hello := client.Hello()
	fmt.Printf("Connected (protocol=%d, server=%s)\n\n",
		hello.Protocol, hello.Server.Version)

	// Send a chat message.
	fmt.Println("Sending chat message: 'What is OpenClaw?'")
	fmt.Println("---")
	result, err := client.ChatSend(ctx, protocol.ChatSendParams{
		SessionKey:     "main",
		Message:        "What is OpenClaw?",
		IdempotencyKey: fmt.Sprintf("chat-%d", time.Now().UnixNano()),
	})
	if err != nil {
		log.Fatalf("ChatSend: %v", err)
	}

	// The result is the initial chat event (or the final one for non-streaming).
	data, _ := json.MarshalIndent(result, "  ", "  ")
	fmt.Printf("---\nChatSend result:\n  %s\n", data)

	// Fetch chat history.
	fmt.Println("\nFetching chat history...")
	limit := 10
	history, err := client.ChatHistory(ctx, protocol.ChatHistoryParams{
		SessionKey: "main",
		Limit:      &limit,
	})
	if err != nil {
		log.Fatalf("ChatHistory: %v", err)
	}
	fmt.Printf("History response: %d bytes\n", len(history))

	fmt.Println("\n=== Done ===")
}
