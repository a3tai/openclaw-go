package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// AssistantMediaGet retrieves assistant media availability and metadata.
func (c *Client) AssistantMediaGet(ctx context.Context, params protocol.AssistantMediaGetParams) (*protocol.AssistantMediaGetResult, error) {
	var result protocol.AssistantMediaGetResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAssistantMediaGet), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
