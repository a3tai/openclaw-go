// Command sessions demonstrates session management via the Gateway.
//
// It connects with device identity and exercises the session lifecycle:
//   - List sessions
//   - Preview sessions
//   - Patch session settings
//   - Reset and delete sessions
//
// Usage:
//
//	go run ./examples/sessions <token> <host> [identity-dir]
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/a3tai/openclaw-go/examples/internal/gwconn"
	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }

func main() {
	cfg := gwconn.ParseArgs("Usage: sessions <token> <host> [identity-dir]")
	cfg.PrintIdentityInfo()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println()
	fmt.Println("=== OpenClaw Sessions Example ===")
	fmt.Println()

	client := cfg.NewClient(
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
			protocol.ScopeOperatorRead,
			protocol.ScopeOperatorWrite,
			protocol.ScopeOperatorAdmin, // required for sessions.patch, sessions.reset
		),
	)
	defer client.Close()

	if err := cfg.Connect(ctx, client); err != nil { log.Fatal(err) }
	fmt.Println()

	// List sessions.
	fmt.Println("--- List Sessions ---")
	listResult, err := client.SessionsList(ctx, protocol.SessionsListParams{
		Limit:         intPtr(10),
		ActiveMinutes: intPtr(60),
	})
	if err != nil {
		fmt.Printf("SessionsList: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", formatJSON(listResult))
	}

	// Preview sessions.
	fmt.Println("\n--- Preview Sessions ---")
	previewResult, err := client.SessionsPreview(ctx, protocol.SessionsPreviewParams{
		Keys:     []string{"main"},
		Limit:    intPtr(5),
		MaxChars: intPtr(200),
	})
	if err != nil {
		fmt.Printf("SessionsPreview: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", formatJSON(previewResult))
	}

	// Patch a session.
	fmt.Println("\n--- Patch Session ---")
	err = client.SessionsPatch(ctx, protocol.SessionsPatchParams{
		Key:   "main",
		Label: strPtr("My Session"),
		Model: strPtr("gpt-4"),
	})
	if err != nil {
		fmt.Printf("SessionsPatch: %v\n", err)
	} else {
		fmt.Println("Session patched successfully")
	}

	// Get session usage.
	fmt.Println("\n--- Session Usage ---")
	usageResult, err := client.SessionsUsage(ctx, protocol.SessionsUsageParams{
		Key: "main",
	})
	if err != nil {
		fmt.Printf("SessionsUsage: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", formatJSON(usageResult))
	}

	// Reset a session.
	fmt.Println("\n--- Reset Session ---")
	err = client.SessionsReset(ctx, protocol.SessionsResetParams{
		Key:    "main",
		Reason: "example reset",
	})
	if err != nil {
		fmt.Printf("SessionsReset: %v\n", err)
	} else {
		fmt.Println("Session reset successfully")
	}

	fmt.Println("\n=== Done ===")
}

func formatJSON(data json.RawMessage) string {
	var v any
	json.Unmarshal(data, &v)
	out, _ := json.MarshalIndent(v, "  ", "  ")
	return string(out)
}
