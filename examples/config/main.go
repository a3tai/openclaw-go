// Command config demonstrates gateway configuration management.
//
// It connects and exercises:
//   - Get current config
//   - Get config schema
//   - Patch config with partial changes
//   - Apply config with restart
//
// Usage:
//
//	go run ./examples/config
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

func intPtr(v int) *int { return &v }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Config Example ===")
	fmt.Println()

	client := gateway.NewClient(
		gateway.WithToken("example-token"),
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
		protocol.ScopeOperatorRead,
		protocol.ScopeOperatorWrite,
		protocol.ScopeOperatorAdmin, // required for config.set, config.apply
	),
	)
	defer client.Close()

	fmt.Println("Connecting...")
	if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
		log.Fatalf("Connect: %v", err)
	}
	fmt.Println("Connected")
	fmt.Println()

	// Get config.
	fmt.Println("--- Get Config ---")
	configData, err := client.ConfigGet(ctx)
	if err != nil {
		fmt.Printf("ConfigGet: %v\n", err)
	} else {
		fmt.Printf("Config: %s\n", formatJSON(configData))
	}

	// Get config schema.
	fmt.Println("\n--- Config Schema ---")
	schema, err := client.ConfigSchema(ctx)
	if err != nil {
		fmt.Printf("ConfigSchema: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(schema, "  ", "  ")
		fmt.Printf("Schema:\n  %s\n", data)
	}

	// Patch config (partial update). Raw is a YAML/JSON string.
	fmt.Println("\n--- Patch Config ---")
	err = client.ConfigPatch(ctx, protocol.ConfigPatchParams{
		Raw: `{"model":"gpt-4","exec":{"ask":"always"}}`,
	})
	if err != nil {
		fmt.Printf("ConfigPatch: %v\n", err)
	} else {
		fmt.Println("Config patched")
	}

	// Apply config (full replacement with restart).
	fmt.Println("\n--- Apply Config ---")
	err = client.ConfigApply(ctx, protocol.ConfigApplyParams{
		Raw:            `{"model":"gpt-4","exec":{"ask":"always","security":"strict"}}`,
		SessionKey:     "main",
		Note:           "Applied from Go example",
		RestartDelayMs: intPtr(1000),
	})
	if err != nil {
		fmt.Printf("ConfigApply: %v\n", err)
	} else {
		fmt.Println("Config applied")
	}

	fmt.Println("\n=== Done ===")
}

func formatJSON(data json.RawMessage) string {
	var v any
	json.Unmarshal(data, &v)
	out, _ := json.MarshalIndent(v, "  ", "  ")
	return string(out)
}
