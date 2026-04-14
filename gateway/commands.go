package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// CommandsList lists the commands available on the gateway.
func (c *Client) CommandsList(ctx context.Context, params protocol.CommandsListParams) (*protocol.CommandsListResult, error) {
	var result protocol.CommandsListResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodCommandsList), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
