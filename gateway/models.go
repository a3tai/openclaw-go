package gateway

import (
	"context"
	"encoding/json"

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

// ModelsAuthStatus retrieves authentication status for the model provider.
func (c *Client) ModelsAuthStatus(ctx context.Context) (*protocol.ModelsAuthStatusResult, error) {
	var result protocol.ModelsAuthStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodModelsAuthStatus), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModelsAuthLogout logs out from the model provider.
func (c *Client) ModelsAuthLogout(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodModelsAuthLogout), struct{}{})
}
