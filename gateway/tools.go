package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ToolsCatalog returns the gateway tool catalog grouped by profile/source.
func (c *Client) ToolsCatalog(ctx context.Context, params protocol.ToolsCatalogParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodToolsCatalog), params)
}

// ToolsEffective returns the runtime-effective tool inventory for a session.
// SessionKey is required. AgentID is optional and must match the session's agent when provided.
// Requires operator.read scope.
func (c *Client) ToolsEffective(ctx context.Context, params protocol.ToolsEffectiveParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodToolsEffective), params)
}

// ToolsInvoke directly invokes a named tool.
func (c *Client) ToolsInvoke(ctx context.Context, params protocol.ToolsInvokeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodToolsInvoke), params)
}
