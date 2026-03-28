// Command cron demonstrates cron job management via the Gateway.
//
// It connects and exercises:
//   - List cron jobs
//   - Add a new cron job
//   - Run a cron job manually
//   - View cron run history
//   - Remove a cron job
//
// Usage:
//
//	go run ./examples/cron
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

func boolPtr(v bool) *bool    { return &v }
func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Cron Example ===")
	fmt.Println()

	client := gateway.NewClient(
		gateway.WithToken("example-token"),
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
		protocol.ScopeOperatorRead,
		protocol.ScopeOperatorWrite,
		protocol.ScopeOperatorAdmin, // required for cron.add, cron.run, cron.remove
	),
	)
	defer client.Close()

	fmt.Println("Connecting...")
	if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
		log.Fatalf("Connect: %v", err)
	}
	fmt.Println("Connected")
	fmt.Println()

	// List cron jobs.
	fmt.Println("--- List Cron Jobs ---")
	jobs, err := client.CronList(ctx, protocol.CronListParams{
		IncludeDisabled: boolPtr(true),
	})
	if err != nil {
		fmt.Printf("CronList: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(jobs, "  ", "  ")
		fmt.Printf("Jobs:\n  %s\n", data)
	}

	// Add a cron job.
	fmt.Println("\n--- Add Cron Job ---")
	addResult, err := client.CronAdd(ctx, protocol.CronAddParams{
		Name:       "daily-summary",
		AgentID:    strPtr("main"),
		SessionKey: strPtr("cron-daily"),
		Enabled:    boolPtr(true),
		Schedule: protocol.CronSchedule{
			Kind: "cron",
			Expr: "0 9 * * *",
		},
		SessionTarget: "main",
		WakeMode:      "now",
		Payload: protocol.CronPayload{
			Kind:    "agentTurn",
			Message: "Generate a daily summary of recent activity.",
		},
	})
	if err != nil {
		fmt.Printf("CronAdd: %v\n", err)
	} else {
		fmt.Printf("Added: %s\n", formatJSON(addResult))
	}

	// Get cron status.
	fmt.Println("\n--- Cron Status ---")
	status, err := client.CronStatus(ctx)
	if err != nil {
		fmt.Printf("CronStatus: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", formatJSON(status))
	}

	// View run history.
	fmt.Println("\n--- Cron Runs ---")
	runs, err := client.CronRuns(ctx, protocol.CronRunsParams{
		JobID: "daily-summary",
		Limit: intPtr(5),
	})
	if err != nil {
		fmt.Printf("CronRuns: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(runs, "  ", "  ")
		fmt.Printf("Runs:\n  %s\n", data)
	}

	// Run a job manually.
	fmt.Println("\n--- Manual Run ---")
	err = client.CronRun(ctx, protocol.CronRunParams{
		JobID: "daily-summary",
		Mode:  "force",
	})
	if err != nil {
		fmt.Printf("CronRun: %v\n", err)
	} else {
		fmt.Println("Job triggered")
	}

	// Remove the job.
	fmt.Println("\n--- Remove Cron Job ---")
	err = client.CronRemove(ctx, protocol.CronRemoveParams{
		JobID: "daily-summary",
	})
	if err != nil {
		fmt.Printf("CronRemove: %v\n", err)
	} else {
		fmt.Println("Job removed")
	}

	fmt.Println("\n=== Done ===")
}

func formatJSON(data json.RawMessage) string {
	var v any
	json.Unmarshal(data, &v)
	out, _ := json.MarshalIndent(v, "  ", "  ")
	return string(out)
}
