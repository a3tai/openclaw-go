// Command chat demonstrates an interactive chat session via the Gateway.
//
// It connects to the gateway with device identity, sends a chat message,
// and displays the streaming agent response events.
//
// Usage:
//
//	go run ./examples/chat <token> <host> [identity-dir]
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/a3tai/openclaw-go/examples/internal/gwconn"
	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

func main() {
	cfg := gwconn.ParseArgs("Usage: chat <token> <host> [identity-dir]")
	cfg.PrintIdentityInfo()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println()
	fmt.Println("=== OpenClaw Chat Example ===")
	fmt.Println()

	client := cfg.NewClient(
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

	cfg.Connect(ctx, client)

	hello := client.Hello()
	fmt.Printf("Protocol: %d, Server: %s\n\n",
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
		fmt.Printf("ChatSend: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(result, "  ", "  ")
		fmt.Printf("---\nChatSend result:\n  %s\n", data)
	}

	// Fetch chat history.
	fmt.Println("\nFetching chat history...")
	limit := 10
	history, err := client.ChatHistory(ctx, protocol.ChatHistoryParams{
		SessionKey: "main",
		Limit:      &limit,
	})
	if err != nil {
		fmt.Printf("ChatHistory: %v\n", err)
	} else {
		fmt.Printf("History response: %d bytes\n", len(history))
	}

	fmt.Println("\n=== Done ===")
}
