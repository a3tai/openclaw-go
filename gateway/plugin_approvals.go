package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// PluginApprovalRequest submits a new plugin approval request.
// The gateway broadcasts a "plugin.approval.requested" event to connected approval clients
// and waits for a decision. If no approval clients are connected, the decision is nil.
// Requires the operator.approvals scope.
func (c *Client) PluginApprovalRequest(ctx context.Context, params protocol.PluginApprovalRequestParams) (*protocol.PluginApprovalRequestResult, error) {
	var result protocol.PluginApprovalRequestResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodPluginApprovalRequest), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PluginApprovalWaitDecision waits for a decision on a pending plugin approval.
// Requires the operator.approvals scope.
func (c *Client) PluginApprovalWaitDecision(ctx context.Context, params protocol.PluginApprovalWaitDecisionParams) (*protocol.PluginApprovalWaitDecisionResult, error) {
	var result protocol.PluginApprovalWaitDecisionResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodPluginApprovalWaitDecision), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PluginApprovalResolve resolves a pending plugin approval request.
// decision must be one of: "allow-once", "allow-always", "deny".
// Requires the operator.approvals scope.
func (c *Client) PluginApprovalResolve(ctx context.Context, params protocol.PluginApprovalResolveParams) (*protocol.PluginApprovalResolveResult, error) {
	var result protocol.PluginApprovalResolveResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodPluginApprovalResolve), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
