package gateway

import (
	"context"

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

// DoctorMemoryStatus returns memory embedding/provider health status.
func (c *Client) DoctorMemoryStatus(ctx context.Context) (*protocol.DoctorMemoryStatusResult, error) {
	var result protocol.DoctorMemoryStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryStatus), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryDreamDiary retrieves the agent's DREAMS.md dream diary file.
func (c *Client) DoctorMemoryDreamDiary(ctx context.Context) (*protocol.DoctorMemoryDreamDiaryResult, error) {
	var result protocol.DoctorMemoryDreamDiaryResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryDreamDiary), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryBackfillDreamDiary backfills the dream diary from daily memory files.
func (c *Client) DoctorMemoryBackfillDreamDiary(ctx context.Context) (*protocol.DoctorMemoryDreamActionResult, error) {
	var result protocol.DoctorMemoryDreamActionResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryBackfillDreamDiary), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryResetDreamDiary removes backfilled entries from the dream diary.
func (c *Client) DoctorMemoryResetDreamDiary(ctx context.Context) (*protocol.DoctorMemoryDreamActionResult, error) {
	var result protocol.DoctorMemoryDreamActionResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryResetDreamDiary), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryResetGroundedShortTerm removes grounded short-term memory candidates.
func (c *Client) DoctorMemoryResetGroundedShortTerm(ctx context.Context) (*protocol.DoctorMemoryDreamActionResult, error) {
	var result protocol.DoctorMemoryDreamActionResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryResetGroundedShortTerm), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryRepairDreamingArtifacts repairs corrupted dreaming artifacts.
func (c *Client) DoctorMemoryRepairDreamingArtifacts(ctx context.Context) (*protocol.DoctorMemoryDreamActionResult, error) {
	var result protocol.DoctorMemoryDreamActionResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryRepairDreamingArtifacts), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DoctorMemoryDedupeDreamDiary removes duplicate entries from the dream diary.
func (c *Client) DoctorMemoryDedupeDreamDiary(ctx context.Context) (*protocol.DoctorMemoryDreamActionResult, error) {
	var result protocol.DoctorMemoryDreamActionResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodDoctorMemoryDedupeDreamDiary), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
