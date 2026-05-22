package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ChannelsStatus retrieves the status of all channels.
func (c *Client) ChannelsStatus(ctx context.Context, params protocol.ChannelsStatusParams) (*protocol.ChannelsStatusResult, error) {
	var result protocol.ChannelsStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodChannelsStatus), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ChannelsLogout logs out of a channel account.
func (c *Client) ChannelsLogout(ctx context.Context, params protocol.ChannelsLogoutParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodChannelsLogout), params)
}

// TalkConfig retrieves the talk (voice) configuration.
func (c *Client) TalkConfig(ctx context.Context, params protocol.TalkConfigParams) (*protocol.TalkConfigResult, error) {
	var result protocol.TalkConfigResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodTalkConfig), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TalkMode sets the talk mode (voice).
func (c *Client) TalkMode(ctx context.Context, params protocol.TalkModeParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodTalkMode), params)
}

// TalkSpeak synthesizes speech audio through the configured talk provider.
func (c *Client) TalkSpeak(ctx context.Context, params protocol.TalkSpeakParams) (*protocol.TalkSpeakResult, error) {
	var result protocol.TalkSpeakResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodTalkSpeak), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// WebLoginStart starts an interactive web login flow for a channel provider.
func (c *Client) WebLoginStart(ctx context.Context, params protocol.WebLoginStartParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodWebLoginStart), params)
}

// WebLoginWait waits for completion of an interactive web login flow.
func (c *Client) WebLoginWait(ctx context.Context, params protocol.WebLoginWaitParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodWebLoginWait), params)
}

// ChannelsStart starts a channel account.
func (c *Client) ChannelsStart(ctx context.Context, params protocol.ChannelsStartParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodChannelsStart), params)
}

// ChannelsStop stops a channel account.
func (c *Client) ChannelsStop(ctx context.Context, params protocol.ChannelsStopParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodChannelsStop), params)
}

// TalkCatalog retrieves the talk catalog.
func (c *Client) TalkCatalog(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkCatalog), struct{}{})
}

// TalkClientCreate creates a talk client.
func (c *Client) TalkClientCreate(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkClientCreate), params)
}

// TalkClientToolCall submits a tool call result from a talk client.
func (c *Client) TalkClientToolCall(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkClientToolCall), params)
}

// TalkSessionCreate creates a real-time talk session.
func (c *Client) TalkSessionCreate(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionCreate), params)
}

// TalkSessionJoin joins an existing talk session.
func (c *Client) TalkSessionJoin(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionJoin), params)
}

// TalkSessionAppendAudio appends audio to a talk session.
func (c *Client) TalkSessionAppendAudio(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionAppendAudio), params)
}

// TalkSessionStartTurn starts a turn in a talk session.
func (c *Client) TalkSessionStartTurn(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionStartTurn), params)
}

// TalkSessionEndTurn ends the current turn in a talk session.
func (c *Client) TalkSessionEndTurn(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionEndTurn), params)
}

// TalkSessionCancelTurn cancels the current turn in a talk session.
func (c *Client) TalkSessionCancelTurn(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionCancelTurn), params)
}

// TalkSessionCancelOutput cancels queued output in a talk session.
func (c *Client) TalkSessionCancelOutput(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionCancelOutput), params)
}

// TalkSessionSubmitToolResult submits a tool result into a talk session.
func (c *Client) TalkSessionSubmitToolResult(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTalkSessionSubmitToolResult), params)
}

// TalkSessionClose closes a talk session.
func (c *Client) TalkSessionClose(ctx context.Context, params any) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodTalkSessionClose), params)
}
