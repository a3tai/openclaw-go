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
