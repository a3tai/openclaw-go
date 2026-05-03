package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// DiagnosticsStability retrieves gateway stability diagnostics.
func (c *Client) DiagnosticsStability(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDiagnosticsStability), struct{}{})
}

// UpdateStatus retrieves the current update status.
func (c *Client) UpdateStatus(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodUpdateStatus), struct{}{})
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
func (c *Client) VoiceWakeRoutingGet(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodVoiceWakeRoutingGet), struct{}{})
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
