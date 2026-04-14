package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// SessionsGet retrieves messages for a session by key.
// Added in openclaw v2026.3.7.
func (c *Client) SessionsGet(ctx context.Context, params protocol.SessionsGetParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsGet), params)
}

// SessionsList lists sessions matching the given criteria.
func (c *Client) SessionsList(ctx context.Context, params protocol.SessionsListParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsList), params)
}

// SessionsPreview retrieves previews for the given session keys.
func (c *Client) SessionsPreview(ctx context.Context, params protocol.SessionsPreviewParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsPreview), params)
}

// SessionsResolve resolves a session by key, ID, or other criteria.
func (c *Client) SessionsResolve(ctx context.Context, params protocol.SessionsResolveParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsResolve), params)
}

// SessionsCreate creates a session and optionally sends an initial message.
func (c *Client) SessionsCreate(ctx context.Context, params protocol.SessionsCreateParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsCreate), params)
}

// SessionsSend sends a message into a session.
func (c *Client) SessionsSend(ctx context.Context, params protocol.SessionsSendParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsSend), params)
}

// SessionsSteer sends an interrupting message into a session.
func (c *Client) SessionsSteer(ctx context.Context, params protocol.SessionsSendParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsSteer), params)
}

// SessionsAbort aborts an active session run.
func (c *Client) SessionsAbort(ctx context.Context, params protocol.SessionsAbortParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsAbort), params)
}

// SessionsSubscribe subscribes the connection to global sessions changed events.
func (c *Client) SessionsSubscribe(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsSubscribe), struct{}{})
}

// SessionsUnsubscribe unsubscribes the connection from global sessions changed events.
func (c *Client) SessionsUnsubscribe(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsUnsubscribe), struct{}{})
}

// SessionsMessagesSubscribe subscribes to message updates for one session.
func (c *Client) SessionsMessagesSubscribe(ctx context.Context, params protocol.SessionsMessagesSubscribeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsMessagesSubscribe), params)
}

// SessionsMessagesUnsubscribe unsubscribes from message updates for one session.
func (c *Client) SessionsMessagesUnsubscribe(ctx context.Context, params protocol.SessionsMessagesUnsubscribeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsMessagesUnsubscribe), params)
}

// SessionsPatch patches session settings.
func (c *Client) SessionsPatch(ctx context.Context, params protocol.SessionsPatchParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSessionsPatch), params)
}

// SessionsReset resets a session.
func (c *Client) SessionsReset(ctx context.Context, params protocol.SessionsResetParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSessionsReset), params)
}

// SessionsDelete deletes a session.
func (c *Client) SessionsDelete(ctx context.Context, params protocol.SessionsDeleteParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSessionsDelete), params)
}

// SessionsCompact compacts a session's history.
func (c *Client) SessionsCompact(ctx context.Context, params protocol.SessionsCompactParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSessionsCompact), params)
}

// SessionsUsage retrieves usage data for a session.
func (c *Client) SessionsUsage(ctx context.Context, params protocol.SessionsUsageParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsUsage), params)
}

// SessionsUsageTimeseries retrieves a token/cost timeline for one session.
func (c *Client) SessionsUsageTimeseries(ctx context.Context, params protocol.SessionsUsageTimeseriesParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsUsageTimeseries), params)
}

// SessionsUsageLogs retrieves usage log entries for one session.
func (c *Client) SessionsUsageLogs(ctx context.Context, params protocol.SessionsUsageLogsParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSessionsUsageLogs), params)
}

// SessionsCompactionList lists compaction checkpoints for a session.
func (c *Client) SessionsCompactionList(ctx context.Context, params protocol.SessionsCompactionListParams) (*protocol.SessionsCompactionListResult, error) {
	var result protocol.SessionsCompactionListResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSessionsCompactionList), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SessionsCompactionGet retrieves a specific compaction checkpoint by ID.
func (c *Client) SessionsCompactionGet(ctx context.Context, params protocol.SessionsCompactionGetParams) (*protocol.SessionsCompactionGetResult, error) {
	var result protocol.SessionsCompactionGetResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSessionsCompactionGet), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SessionsCompactionBranch creates a new session branched from a compaction checkpoint.
func (c *Client) SessionsCompactionBranch(ctx context.Context, params protocol.SessionsCompactionBranchParams) (*protocol.SessionsCompactionBranchResult, error) {
	var result protocol.SessionsCompactionBranchResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSessionsCompactionBranch), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SessionsCompactionRestore restores a session to the state at a compaction checkpoint.
func (c *Client) SessionsCompactionRestore(ctx context.Context, params protocol.SessionsCompactionRestoreParams) (*protocol.SessionsCompactionRestoreResult, error) {
	var result protocol.SessionsCompactionRestoreResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSessionsCompactionRestore), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
