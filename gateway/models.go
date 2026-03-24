package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// ModelsList lists available models.
func (c *Client) ModelsList(ctx context.Context) (*protocol.ModelsListResult, error) {
	var result protocol.ModelsListResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodModelsList), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
