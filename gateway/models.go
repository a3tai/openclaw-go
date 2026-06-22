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

// ModelsAuthStatus returns the authentication health status for all configured model providers.
func (c *Client) ModelsAuthStatus(ctx context.Context) (*protocol.ModelsAuthStatusResult, error) {
	var result protocol.ModelsAuthStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodModelsAuthStatus), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
