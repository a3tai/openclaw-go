// Command cron demonstrates cron job management via the Gateway.
//
// It connects with device identity authentication and exercises:
//   - List cron jobs
//   - Add a new cron job
//   - Run a cron job manually
//   - View cron run history
//   - Remove a cron job
//
// Usage:
//
//	go run ./examples/cron <token> <host> [identity-dir]
//
// The device must be paired with the gateway before scoped operations will
// work. On first run, the example generates a new Ed25519 keypair. Approve
// the device on the gateway (e.g. via the Control UI or CLI), then re-run.
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

func boolPtr(v bool) *bool    { return &v }
func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }

func main() {
	cfg := gwconn.ParseArgs("Usage: cron <token> <host> [identity-dir]")
	cfg.PrintIdentityInfo()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println()
	fmt.Println("=== OpenClaw Cron Example ===")
	fmt.Println()

	client := cfg.NewClient(
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
			protocol.ScopeOperatorRead,
			protocol.ScopeOperatorWrite,
			protocol.ScopeOperatorAdmin, // required for cron.add, cron.run, cron.remove
		),
	)
	defer client.Close()

	if err := cfg.Connect(ctx, client); err != nil { log.Fatal(err) }
	fmt.Println()

	// List cron jobs.
	fmt.Println("--- List Cron Jobs ---")
	jobs, err := client.CronList(ctx, protocol.CronListParams{
		IncludeDisabled: boolPtr(true),
	})
	if err != nil {
		fmt.Printf("CronList: %v\n", err)
	} else {
		fmt.Printf("Jobs: %d (total: %d, hasMore: %v)\n", len(jobs.Jobs), jobs.Total, jobs.HasMore)
		data, _ := json.MarshalIndent(jobs.Jobs, "  ", "  ")
		fmt.Printf("  %s\n", data)
	}

	// Add a cron job using systemEvent payload (required for main agent).
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
			Kind: "systemEvent",
			Text: "Generate a daily summary of recent activity.",
		},
	})

	// Extract the server-assigned job ID for subsequent operations.
	var jobID string
	if err != nil {
		fmt.Printf("CronAdd: %v\n", err)
	} else {
		fmt.Printf("Added: %s\n", formatJSON(addResult))
		var added struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(addResult, &added) == nil {
			jobID = added.ID
		}
	}

	// Get cron status.
	fmt.Println("\n--- Cron Status ---")
	status, err := client.CronStatus(ctx)
	if err != nil {
		fmt.Printf("CronStatus: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", formatJSON(status))
	}

	// View run history (use server-assigned ID, not the name).
	fmt.Println("\n--- Cron Runs ---")
	if jobID == "" {
		fmt.Println("Skipped (no job ID from CronAdd)")
	} else {
		runs, err := client.CronRuns(ctx, protocol.CronRunsParams{
			ID:    jobID,
			Limit: intPtr(5),
		})
		if err != nil {
			fmt.Printf("CronRuns: %v\n", err)
		} else {
			fmt.Printf("Runs: %d (total: %d)\n", len(runs.Entries), runs.Total)
			data, _ := json.MarshalIndent(runs.Entries, "  ", "  ")
			fmt.Printf("  %s\n", data)
		}
	}

	// Run a job manually (use server-assigned ID).
	fmt.Println("\n--- Manual Run ---")
	if jobID == "" {
		fmt.Println("Skipped (no job ID from CronAdd)")
	} else {
		err = client.CronRun(ctx, protocol.CronRunParams{
			ID:   jobID,
			Mode: "force",
		})
		if err != nil {
			fmt.Printf("CronRun: %v\n", err)
		} else {
			fmt.Println("Job triggered")
		}
	}

	// Remove the job (use server-assigned ID).
	fmt.Println("\n--- Remove Cron Job ---")
	if jobID == "" {
		fmt.Println("Skipped (no job ID from CronAdd)")
	} else {
		err = client.CronRemove(ctx, protocol.CronRemoveParams{
			ID: jobID,
		})
		if err != nil {
			fmt.Printf("CronRemove: %v\n", err)
		} else {
			fmt.Println("Job removed")
		}
	}

	fmt.Println("\n=== Done ===")
}

func formatJSON(data json.RawMessage) string {
	var v any
	json.Unmarshal(data, &v)
	out, _ := json.MarshalIndent(v, "  ", "  ")
	return string(out)
}
