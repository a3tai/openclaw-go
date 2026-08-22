package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// AssistantMediaGet retrieves the current assistant media.
func (c *Client) AssistantMediaGet(ctx context.Context) (*protocol.AssistantMediaGetResult, error) {
	var result protocol.AssistantMediaGetResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAssistantMediaGet), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
