package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// NodePairRequest requests pairing with the gateway.
func (c *Client) NodePairRequest(ctx context.Context, params protocol.NodePairRequestParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodNodePairRequest), params)
}

// NodePairList lists pending node pairing requests.
func (c *Client) NodePairList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodNodePairList), struct{}{})
}

// NodePairApprove approves a node pairing request.
func (c *Client) NodePairApprove(ctx context.Context, params protocol.NodePairApproveParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodNodePairApprove), params)
}

// NodePairReject rejects a node pairing request.
func (c *Client) NodePairReject(ctx context.Context, params protocol.NodePairRejectParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodNodePairReject), params)
}

// NodePairVerify verifies a node pairing.
func (c *Client) NodePairVerify(ctx context.Context, params protocol.NodePairVerifyParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodNodePairVerify), params)
}
