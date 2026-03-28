// Command approvals demonstrates the exec approval flow via the Gateway.
//
// When an agent wants to execute a command, the gateway can require operator
// approval. This example connects with device identity and:
//   - Listens for exec.approval.requested events
//   - Auto-approves safe commands, rejects dangerous ones
//   - Shows the exec approvals admin API
//
// Usage:
//
//	go run ./examples/approvals <token> <host> [identity-dir]
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/a3tai/openclaw-go/examples/internal/gwconn"
	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

// dangerousCommands is a simple blocklist for demonstration.
var dangerousCommands = []string{"rm -rf", "sudo", "chmod 777", "mkfs"}

func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }

func main() {
	cfg := gwconn.ParseArgs("Usage: approvals <token> <host> [identity-dir]")
	cfg.PrintIdentityInfo()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println()
	fmt.Println("=== OpenClaw Exec Approvals Example ===")
	fmt.Println()

	client := cfg.NewClient(
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
			protocol.ScopeOperatorRead,
			protocol.ScopeOperatorWrite,
			protocol.ScopeOperatorApprovals,
		),
		gateway.WithOnEvent(func(ev protocol.Event) {
			if ev.EventName == protocol.EventExecApproval {
				var req protocol.ExecApprovalRequested
				if json.Unmarshal(ev.Payload, &req) == nil {
					cwd := "<unknown>"
					if req.Cwd != nil {
						cwd = *req.Cwd
					}
					host := "<unknown>"
					if req.Host != nil {
						host = *req.Host
					}
					fmt.Printf("  [approval] Command: %s (cwd: %s, host: %s)\n",
						req.Command, cwd, host)

					// Check if the command is dangerous.
					decision := "approved"
					for _, d := range dangerousCommands {
						if strings.Contains(req.Command, d) {
							decision = "denied"
							break
						}
					}
					fmt.Printf("  [approval] Decision: %s\n", decision)
				}
			}
		}),
	)
	defer client.Close()

	if err := cfg.Connect(ctx, client); err != nil { log.Fatal(err) }
	fmt.Println()

	// --- Resolve approvals ---
	fmt.Println("--- Resolve Exec Approval ---")
	resolveResult, err := client.ExecApprovalResolve(ctx, protocol.ExecApprovalResolveParams{
		ID:       "approval-1",
		Decision: "approved",
	})
	if err != nil {
		fmt.Printf("ExecApprovalResolve: %v\n", err)
	} else {
		fmt.Printf("Approval resolved: ok=%v\n", resolveResult.OK)
	}

	// --- Admin: Get approvals config ---
	fmt.Println("\n--- Get Exec Approvals Config ---")
	snapshot, err := client.ExecApprovalsGet(ctx)
	if err != nil {
		fmt.Printf("ExecApprovalsGet: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(snapshot, "  ", "  ")
		fmt.Printf("Config:\n  %s\n", data)
	}

	// --- Admin: Set approvals config ---
	fmt.Println("\n--- Set Exec Approvals Config ---")
	err = client.ExecApprovalsSet(ctx, protocol.ExecApprovalsSetParams{
		File: protocol.ExecApprovalsFile{
			Version: 1,
			Defaults: &protocol.ExecApprovalsDefaults{
				Ask:      "always",
				Security: "strict",
			},
		},
	})
	if err != nil {
		fmt.Printf("ExecApprovalsSet: %v\n", err)
	} else {
		fmt.Println("Approvals config updated")
	}

	// --- Request an exec approval (node-side) ---
	fmt.Println("\n--- Request Exec Approval ---")
	reqResult, err := client.ExecApprovalRequest(ctx, protocol.ExecApprovalRequestParams{
		ID:        "approval-new",
		Command:   "ls -la /tmp",
		Cwd:       strPtr("/tmp"),
		Host:      strPtr("localhost"),
		Security:  strPtr("normal"),
		TimeoutMs: intPtr(30000),
	})
	if err != nil {
		fmt.Printf("ExecApprovalRequest: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(reqResult, "  ", "  ")
		fmt.Printf("Request result:\n  %s\n", data)
	}

	fmt.Println("\n=== Done ===")
}
