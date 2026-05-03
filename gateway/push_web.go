package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// PushWebSubscribe registers a Web Push subscription endpoint.
func (c *Client) PushWebSubscribe(ctx context.Context, params protocol.PushWebSubscribeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodPushWebSubscribe), params)
}

// PushWebUnsubscribe removes a Web Push subscription endpoint.
func (c *Client) PushWebUnsubscribe(ctx context.Context, params protocol.PushWebUnsubscribeParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodPushWebUnsubscribe), params)
}

// PushWebTest sends a test Web Push notification to the subscribed endpoint.
func (c *Client) PushWebTest(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodPushWebTest), struct{}{})
}

// AssistantMediaGet retrieves assistant-uploaded media by ID.
func (c *Client) AssistantMediaGet(ctx context.Context, params protocol.AssistantMediaGetParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodAssistantMediaGet), params)
}
