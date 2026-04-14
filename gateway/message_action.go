package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// MessageAction dispatches a channel message action through the gateway.
//
// Scope: operator.write
func (c *Client) MessageAction(ctx context.Context, params protocol.MessageActionParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodMessageAction), params)
}
