package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// Health retrieves the gateway health status.
func (c *Client) Health(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "health", nil)
}

// Status retrieves the gateway status.
func (c *Client) Status(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "status", nil)
}

// DiagnosticsStability retrieves the gateway stability diagnostics snapshot.
func (c *Client) DiagnosticsStability(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDiagnosticsStability), params)
}
