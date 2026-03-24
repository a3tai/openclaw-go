package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// Presence fetches the current presence entries from the gateway.
func (c *Client) Presence(ctx context.Context) (map[string]protocol.PresenceEntry, error) {
	var entries map[string]protocol.PresenceEntry
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSystemPresence), nil, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
