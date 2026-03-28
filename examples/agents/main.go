// Command agents demonstrates agent CRUD operations via the Gateway.
//
// It connects and exercises:
//   - List agents
//   - Create an agent
//   - Update agent configuration
//   - Manage agent files
//   - Delete an agent
//
// Usage:
//
//	go run ./examples/agents
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

func boolPtr(v bool) *bool { return &v }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Agents CRUD Example ===")
	fmt.Println()

	client := gateway.NewClient(
		gateway.WithToken("example-token"),
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
		protocol.ScopeOperatorRead,
		protocol.ScopeOperatorWrite,
		protocol.ScopeOperatorAdmin, // required for agents.create, agents.update, agents.delete
	),
	)
	defer client.Close()

	fmt.Println("Connecting...")
	if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
		log.Fatalf("Connect: %v", err)
	}
	fmt.Println("Connected")
	fmt.Println()

	// List agents.
	fmt.Println("--- List Agents ---")
	listResult, err := client.AgentsList(ctx)
	if err != nil {
		fmt.Printf("AgentsList: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(listResult, "  ", "  ")
		fmt.Printf("Result:\n  %s\n", data)
	}

	// Create an agent.
	fmt.Println("\n--- Create Agent ---")
	createResult, err := client.AgentsCreate(ctx, protocol.AgentsCreateParams{
		Name:  "my-assistant",
		Emoji: "robot",
	})
	if err != nil {
		fmt.Printf("AgentsCreate: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(createResult, "  ", "  ")
		fmt.Printf("Created:\n  %s\n", data)
	}

	// List agent files.
	fmt.Println("\n--- List Agent Files ---")
	filesResult, err := client.AgentsFilesList(ctx, protocol.AgentsFilesListParams{
		AgentID: "my-assistant",
	})
	if err != nil {
		fmt.Printf("AgentsFilesList: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(filesResult, "  ", "  ")
		fmt.Printf("Files:\n  %s\n", data)
	}

	// Set an agent file.
	fmt.Println("\n--- Set Agent File ---")
	setResult, err := client.AgentsFilesSet(ctx, protocol.AgentsFilesSetParams{
		AgentID: "my-assistant",
		Name:    "system.md",
		Content: "You are a helpful assistant.",
	})
	if err != nil {
		fmt.Printf("AgentsFilesSet: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(setResult, "  ", "  ")
		fmt.Printf("Set result:\n  %s\n", data)
	}

	// Update agent.
	fmt.Println("\n--- Update Agent ---")
	err = client.AgentsUpdate(ctx, protocol.AgentsUpdateParams{
		AgentID: "my-assistant",
		Name:    "My Updated Assistant",
		Model:   "gpt-4",
	})
	if err != nil {
		fmt.Printf("AgentsUpdate: %v\n", err)
	} else {
		fmt.Println("Agent updated")
	}

	// Delete agent.
	fmt.Println("\n--- Delete Agent ---")
	deleteResult, err := client.AgentsDelete(ctx, protocol.AgentsDeleteParams{
		AgentID:     "my-assistant",
		DeleteFiles: boolPtr(true),
	})
	if err != nil {
		fmt.Printf("AgentsDelete: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(deleteResult, "  ", "  ")
		fmt.Printf("Deleted:\n  %s\n", data)
	}

	fmt.Println("\n=== Done ===")
}
