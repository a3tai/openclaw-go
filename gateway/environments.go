package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// EnvironmentsList lists available execution environments.
func (c *Client) EnvironmentsList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodEnvironmentsList), struct{}{})
}

// EnvironmentsStatus retrieves the status of execution environments.
func (c *Client) EnvironmentsStatus(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodEnvironmentsStatus), struct{}{})
}
