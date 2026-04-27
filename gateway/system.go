package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// GatewayIdentityGet retrieves the gateway device identity and public key.
func (c *Client) GatewayIdentityGet(ctx context.Context) (*protocol.GatewayIdentityResult, error) {
	var result protocol.GatewayIdentityResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodGatewayIdentityGet), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DiagnosticsStability retrieves gateway stability diagnostics.
func (c *Client) DiagnosticsStability(ctx context.Context) (*protocol.DiagnosticsStabilityResult, error) {
	var result protocol.DiagnosticsStabilityResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDiagnosticsStability), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryStatus returns memory embedding/provider health status.
func (c *Client) DoctorMemoryStatus(ctx context.Context) (*protocol.DoctorMemoryStatusResult, error) {
	var result protocol.DoctorMemoryStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryStatus), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryDreamDiary retrieves the dream diary memory entries.
func (c *Client) DoctorMemoryDreamDiary(ctx context.Context) (*protocol.DoctorMemoryDreamDiaryResult, error) {
	var result protocol.DoctorMemoryDreamDiaryResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryDreamDiary), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryBackfillDreamDiary backfills the dream diary from existing session history.
func (c *Client) DoctorMemoryBackfillDreamDiary(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDoctorMemoryBackfillDreamDiary), struct{}{})
}

// DoctorMemoryResetDreamDiary resets the dream diary memory.
func (c *Client) DoctorMemoryResetDreamDiary(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDoctorMemoryResetDreamDiary), struct{}{})
}

// DoctorMemoryResetGroundedShortTerm resets the grounded short-term memory.
func (c *Client) DoctorMemoryResetGroundedShortTerm(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDoctorMemoryResetGroundedShortTerm), struct{}{})
}

// DoctorMemoryRepairDreamingArtifacts repairs dreaming artifacts in memory.
func (c *Client) DoctorMemoryRepairDreamingArtifacts(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDoctorMemoryRepairDreamingArtifacts), struct{}{})
}

// DoctorMemoryDedupeDreamDiary deduplicates dream diary memory entries.
func (c *Client) DoctorMemoryDedupeDreamDiary(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDoctorMemoryDedupeDreamDiary), struct{}{})
}
