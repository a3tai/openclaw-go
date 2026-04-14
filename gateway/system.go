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

// DoctorMemoryBackfillDreamDiary triggers a backfill of the agent's dream diary.
func (c *Client) DoctorMemoryBackfillDreamDiary(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDoctorMemoryBackfillDreamDiary), struct{}{})
}

// DoctorMemoryResetDreamDiary resets the agent's dream diary.
func (c *Client) DoctorMemoryResetDreamDiary(ctx context.Context) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodDoctorMemoryResetDreamDiary), struct{}{})
}

// DoctorMemoryResetGroundedShortTerm resets the grounded short-term memory store.
func (c *Client) DoctorMemoryResetGroundedShortTerm(ctx context.Context) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodDoctorMemoryResetGroundedShortTerm), struct{}{})
}
