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

// DiagnosticsStability retrieves gateway stability diagnostics.
func (c *Client) DiagnosticsStability(ctx context.Context, params protocol.DiagnosticsStabilityParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDiagnosticsStability), params)
}
