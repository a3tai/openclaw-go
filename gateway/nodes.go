package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// NodeList lists connected nodes.
func (c *Client) NodeList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.list", struct{}{})
}

// NodeDescribe describes a specific node.
func (c *Client) NodeDescribe(ctx context.Context, params protocol.NodeDescribeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.describe", params)
}

// NodeInvoke invokes a command on a node.
func (c *Client) NodeInvoke(ctx context.Context, params protocol.NodeInvokeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.invoke", params)
}

// NodeInvokeResult sends an invoke result from a node to the gateway.
func (c *Client) NodeInvokeResult(ctx context.Context, params protocol.NodeInvokeResultParams) error {
	return c.sendRPCVoid(ctx, "node.invoke.result", params)
}

// NodeEvent sends an event from a node to the gateway.
func (c *Client) NodeEvent(ctx context.Context, params protocol.NodeEventParams) error {
	return c.sendRPCVoid(ctx, "node.event", params)
}

// NodeRename renames a node.
func (c *Client) NodeRename(ctx context.Context, params protocol.NodeRenameParams) error {
	return c.sendRPCVoid(ctx, "node.rename", params)
}

// NodePendingEnqueue enqueues pending work for a node.
func (c *Client) NodePendingEnqueue(ctx context.Context, params protocol.NodePendingEnqueueParams) (*protocol.NodePendingEnqueueResult, error) {
	var result protocol.NodePendingEnqueueResult
	if err := c.sendRPCTyped(ctx, "node.pending.enqueue", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NodePendingDrain drains pending work for the connected node identity.
func (c *Client) NodePendingDrain(ctx context.Context, params protocol.NodePendingDrainParams) (*protocol.NodePendingDrainResult, error) {
	var result protocol.NodePendingDrainResult
	if err := c.sendRPCTyped(ctx, "node.pending.drain", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
