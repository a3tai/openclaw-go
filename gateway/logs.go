package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// LogsTail retrieves the latest log lines.
func (c *Client) LogsTail(ctx context.Context, params protocol.LogsTailParams) (*protocol.LogsTailResult, error) {
	var result protocol.LogsTailResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodLogsTail), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
