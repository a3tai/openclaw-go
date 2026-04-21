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

// ChannelsStart starts a channel account.
func (c *Client) ChannelsStart(ctx context.Context, params protocol.ChannelsStartParams) (*protocol.ChannelsStartResult, error) {
	var result protocol.ChannelsStartResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodChannelsStart), params, &result); err != nil {
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
