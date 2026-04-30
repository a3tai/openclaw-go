package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// AssistantMediaGet retrieves assistant media.
func (c *Client) AssistantMediaGet(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodAssistantMediaGet), params)
}

// UpdateStatus retrieves the current update status including restart sentinel.
func (c *Client) UpdateStatus(ctx context.Context) (*protocol.UpdateStatusResult, error) {
	var result protocol.UpdateStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodUpdateStatus), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

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
//
// Deprecated: browser.request is not present in upstream OpenClaw's BASE_METHODS
// and is a deprecation candidate. Avoid using this in new code.
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

// VoiceWakeRoutingGet retrieves the voice wake routing configuration.
func (c *Client) VoiceWakeRoutingGet(ctx context.Context) (*protocol.VoiceWakeRoutingGetResult, error) {
	var result protocol.VoiceWakeRoutingGetResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodVoiceWakeRoutingGet), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// VoiceWakeRoutingSet sets the voice wake routing configuration.
func (c *Client) VoiceWakeRoutingSet(ctx context.Context, params protocol.VoiceWakeRoutingSetParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodVoiceWakeRoutingSet), params)
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
