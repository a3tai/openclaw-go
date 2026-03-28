package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// Presence fetches the current presence entries from the gateway.
// The server returns a flat array of presence entries.
func (c *Client) Presence(ctx context.Context) ([]protocol.PresenceEntry, error) {
	var entries []protocol.PresenceEntry
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSystemPresence), nil, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
