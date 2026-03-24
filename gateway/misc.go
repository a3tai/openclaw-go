package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// UpdateRun triggers a gateway update run.
func (c *Client) UpdateRun(ctx context.Context, params protocol.UpdateRunParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodUpdateRun), params)
}

// PushTest sends a test push notification.
func (c *Client) PushTest(ctx context.Context, params protocol.PushTestParams) (*protocol.PushTestResult, error) {
	var result protocol.PushTestResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodPushTest), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BrowserRequest makes a browser request via the gateway.
func (c *Client) BrowserRequest(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodBrowserRequest), params)
}

// VoiceWakeGet retrieves the voice wake configuration.
func (c *Client) VoiceWakeGet(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodVoiceWakeGet), struct{}{})
}

// VoiceWakeSet sets the voice wake configuration.
func (c *Client) VoiceWakeSet(ctx context.Context, params any) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodVoiceWakeSet), params)
}

// UsageStatus retrieves usage status.
func (c *Client) UsageStatus(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodUsageStatus), struct{}{})
}

// UsageCost retrieves usage cost information.
func (c *Client) UsageCost(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodUsageCost), params)
}

// Poll creates a poll.
func (c *Client) Poll(ctx context.Context, params protocol.PollParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "poll", params)
}
