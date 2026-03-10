package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// SessionsGet retrieves messages for a session by key.
// Added in openclaw v2026.3.7.
func (c *Client) SessionsGet(ctx context.Context, params protocol.SessionsGetParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "sessions.get", params)
}

// SessionsList lists sessions matching the given criteria.
func (c *Client) SessionsList(ctx context.Context, params protocol.SessionsListParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "sessions.list", params)
}

// SessionsPreview retrieves previews for the given session keys.
func (c *Client) SessionsPreview(ctx context.Context, params protocol.SessionsPreviewParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "sessions.preview", params)
}

// SessionsResolve resolves a session by key, ID, or other criteria.
func (c *Client) SessionsResolve(ctx context.Context, params protocol.SessionsResolveParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "sessions.resolve", params)
}

// SessionsPatch patches session settings.
func (c *Client) SessionsPatch(ctx context.Context, params protocol.SessionsPatchParams) error {
	return c.sendRPCVoid(ctx, "sessions.patch", params)
}

// SessionsReset resets a session.
func (c *Client) SessionsReset(ctx context.Context, params protocol.SessionsResetParams) error {
	return c.sendRPCVoid(ctx, "sessions.reset", params)
}

// SessionsDelete deletes a session.
func (c *Client) SessionsDelete(ctx context.Context, params protocol.SessionsDeleteParams) error {
	return c.sendRPCVoid(ctx, "sessions.delete", params)
}

// SessionsCompact compacts a session's history.
func (c *Client) SessionsCompact(ctx context.Context, params protocol.SessionsCompactParams) error {
	return c.sendRPCVoid(ctx, "sessions.compact", params)
}

// SessionsUsage retrieves usage data for a session.
func (c *Client) SessionsUsage(ctx context.Context, params protocol.SessionsUsageParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "sessions.usage", params)
}
