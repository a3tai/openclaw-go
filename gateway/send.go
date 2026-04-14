package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// SendMessage sends a message via the gateway.
func (c *Client) SendMessage(ctx context.Context, params protocol.SendParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "send", params)
}

// Wake wakes up the gateway.
func (c *Client) Wake(ctx context.Context, params protocol.WakeParams) error {
	return c.sendRPCVoid(ctx, "wake", params)
}

// LastHeartbeat retrieves the last heartbeat timestamp.
func (c *Client) LastHeartbeat(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodLastHeartbeat), nil)
}

// SetHeartbeats enables or disables heartbeat events.
func (c *Client) SetHeartbeats(ctx context.Context, enabled bool) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSetHeartbeats), map[string]bool{"enabled": enabled})
}

// SystemEvent sends a system event.
func (c *Client) SystemEvent(ctx context.Context, params any) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSystemEvent), params)
}

// MessageAction performs a channel-specific action on a message.
func (c *Client) MessageAction(ctx context.Context, params protocol.MessageActionParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodMessageAction), params)
}
