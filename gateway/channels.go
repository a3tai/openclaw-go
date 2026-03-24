package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// ChannelsStatus retrieves the status of all channels.
func (c *Client) ChannelsStatus(ctx context.Context, params protocol.ChannelsStatusParams) (*protocol.ChannelsStatusResult, error) {
	var result protocol.ChannelsStatusResult
	if err := c.sendRPCTyped(ctx, "channels.status", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ChannelsLogout logs out of a channel account.
func (c *Client) ChannelsLogout(ctx context.Context, params protocol.ChannelsLogoutParams) error {
	return c.sendRPCVoid(ctx, "channels.logout", params)
}

// TalkConfig retrieves the talk (voice) configuration.
func (c *Client) TalkConfig(ctx context.Context, params protocol.TalkConfigParams) (*protocol.TalkConfigResult, error) {
	var result protocol.TalkConfigResult
	if err := c.sendRPCTyped(ctx, "talk.config", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TalkMode sets the talk mode (voice).
func (c *Client) TalkMode(ctx context.Context, params protocol.TalkModeParams) error {
	return c.sendRPCVoid(ctx, "talk.mode", params)
}

// TalkSpeak synthesizes speech audio through the configured talk provider.
func (c *Client) TalkSpeak(ctx context.Context, params protocol.TalkSpeakParams) (*protocol.TalkSpeakResult, error) {
	var result protocol.TalkSpeakResult
	if err := c.sendRPCTyped(ctx, "talk.speak", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
