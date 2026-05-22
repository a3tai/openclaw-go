package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// PluginsUIDescriptors retrieves UI descriptors for installed plugins.
func (c *Client) PluginsUIDescriptors(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodPluginsUIDescriptors), struct{}{})
}

// PluginsSessionAction invokes a plugin-defined session action.
func (c *Client) PluginsSessionAction(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodPluginsSessionAction), params)
}
