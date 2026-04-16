package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// ExecApprovalResolve resolves a pending exec approval request.
// Requires the operator.approvals scope.
func (c *Client) ExecApprovalResolve(ctx context.Context, params protocol.ExecApprovalResolveParams) (*protocol.ExecApprovalResolveResult, error) {
	var result protocol.ExecApprovalResolveResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodExecApprovalResolve), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecApprovalGet retrieves details about a pending exec approval by ID.
// Requires the operator.approvals scope.
func (c *Client) ExecApprovalGet(ctx context.Context, params protocol.ExecApprovalGetParams) (*protocol.ExecApprovalGetResult, error) {
	var result protocol.ExecApprovalGetResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodExecApprovalGet), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecApprovalRequest submits a new exec approval request.
func (c *Client) ExecApprovalRequest(ctx context.Context, params protocol.ExecApprovalRequestParams) (*protocol.ExecApprovalRequestResult, error) {
	var result protocol.ExecApprovalRequestResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodExecApprovalRequest), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecApprovalList lists all pending exec approval requests.
// Requires the operator.approvals scope.
func (c *Client) ExecApprovalList(ctx context.Context) ([]protocol.ExecApprovalListEntry, error) {
	var result []protocol.ExecApprovalListEntry
	if err := c.sendRPCTyped(ctx, string(protocol.MethodExecApprovalList), struct{}{}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ExecApprovalWaitDecision waits for a decision on a pending exec approval.
// Requires the operator.approvals scope.
func (c *Client) ExecApprovalWaitDecision(ctx context.Context, params protocol.ExecApprovalWaitDecisionParams) (*protocol.ExecApprovalWaitDecisionResult, error) {
	var result protocol.ExecApprovalWaitDecisionResult
	if err := c.sendRPCTyped(ctx, "exec.approval.waitDecision", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecApprovalsGet retrieves the exec approvals configuration.
func (c *Client) ExecApprovalsGet(ctx context.Context) (*protocol.ExecApprovalsSnapshot, error) {
	var snap protocol.ExecApprovalsSnapshot
	if err := c.sendRPCTyped(ctx, string(protocol.MethodExecApprovalsGet), protocol.ExecApprovalsGetParams{}, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// ExecApprovalsSet updates the exec approvals configuration.
func (c *Client) ExecApprovalsSet(ctx context.Context, params protocol.ExecApprovalsSetParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodExecApprovalsSet), params)
}

// ExecApprovalsNodeGet retrieves exec approvals for a specific node.
func (c *Client) ExecApprovalsNodeGet(ctx context.Context, params protocol.ExecApprovalsNodeGetParams) (*protocol.ExecApprovalsSnapshot, error) {
	var snap protocol.ExecApprovalsSnapshot
	if err := c.sendRPCTyped(ctx, string(protocol.MethodExecApprovalsNodeGet), params, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// ExecApprovalsNodeSet updates exec approvals for a specific node.
func (c *Client) ExecApprovalsNodeSet(ctx context.Context, params protocol.ExecApprovalsNodeSetParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodExecApprovalsNodeSet), params)
}
