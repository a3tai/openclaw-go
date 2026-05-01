package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// PushWebSubscribe registers a web push subscription endpoint.
func (c *Client) PushWebSubscribe(ctx context.Context, params protocol.PushWebSubscribeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodPushWebSubscribe), params)
}

// PushWebUnsubscribe removes a web push subscription by endpoint.
func (c *Client) PushWebUnsubscribe(ctx context.Context, params protocol.PushWebUnsubscribeParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodPushWebUnsubscribe), params)
}

// PushWebTest sends a test web push notification to all registered subscriptions.
func (c *Client) PushWebTest(ctx context.Context, params protocol.PushWebTestParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodPushWebTest), params)
}

// AssistantMediaGet retrieves a media attachment for the assistant.
func (c *Client) AssistantMediaGet(ctx context.Context, params any) (json.RawMessage, error) {
	if params == nil {
		params = struct{}{}
	}
	return c.sendRPC(ctx, string(protocol.MethodAssistantMediaGet), params)
}
