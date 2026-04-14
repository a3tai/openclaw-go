package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// CommandsList retrieves the list of available commands.
//
// Scope: operator.read
func (c *Client) CommandsList(ctx context.Context, params protocol.CommandsListParams) (*protocol.CommandsListResult, error) {
	var result protocol.CommandsListResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodCommandsList), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
