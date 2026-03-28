package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// PluginApprovalRequest submits a plugin approval request and waits for a human
// decision. The call blocks until the request is approved, denied, or times out.
//
// Plugin approval lets an agent plugin ask for human permission before taking a
// consequential action (e.g. sending an email, deleting a file, making an API call).
//
// Requires the operator.approvals scope on the gateway connection.
func (c *Client) PluginApprovalRequest(ctx context.Context, params protocol.PluginApprovalRequestParams) (*protocol.PluginApprovalRequestResult, error) {
	var result protocol.PluginApprovalRequestResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodPluginApprovalRequest), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PluginApprovalResolve resolves a pending plugin approval request.
//
// Operators call this to approve or deny a request submitted via PluginApprovalRequest.
//
// Requires the operator.approvals scope on the gateway connection.
func (c *Client) PluginApprovalResolve(ctx context.Context, params protocol.PluginApprovalResolveParams) (*protocol.PluginApprovalResolveResult, error) {
	var result protocol.PluginApprovalResolveResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodPluginApprovalResolve), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
