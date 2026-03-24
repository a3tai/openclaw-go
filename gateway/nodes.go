package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// NodeList lists connected nodes.
func (c *Client) NodeList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodNodeList), struct{}{})
}

// NodeDescribe describes a specific node.
func (c *Client) NodeDescribe(ctx context.Context, params protocol.NodeDescribeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodNodeDescribe), params)
}

// NodeInvoke invokes a command on a node.
func (c *Client) NodeInvoke(ctx context.Context, params protocol.NodeInvokeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodNodeInvoke), params)
}

// NodeInvokeResult sends an invoke result from a node to the gateway.
func (c *Client) NodeInvokeResult(ctx context.Context, params protocol.NodeInvokeResultParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodNodeInvokeResult), params)
}

// NodeEvent sends an event from a node to the gateway.
func (c *Client) NodeEvent(ctx context.Context, params protocol.NodeEventParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodNodeEvent), params)
}

// NodeRename renames a node.
func (c *Client) NodeRename(ctx context.Context, params protocol.NodeRenameParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodNodeRename), params)
}

// NodePendingEnqueue enqueues pending work for a node.
func (c *Client) NodePendingEnqueue(ctx context.Context, params protocol.NodePendingEnqueueParams) (*protocol.NodePendingEnqueueResult, error) {
	var result protocol.NodePendingEnqueueResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodNodePendingEnqueue), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NodePendingDrain drains pending work for the connected node identity.
func (c *Client) NodePendingDrain(ctx context.Context, params protocol.NodePendingDrainParams) (*protocol.NodePendingDrainResult, error) {
	var result protocol.NodePendingDrainResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodNodePendingDrain), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NodePendingPull pulls queued actions for the connected node identity.
func (c *Client) NodePendingPull(ctx context.Context) (*protocol.NodePendingPullResult, error) {
	var result protocol.NodePendingPullResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodNodePendingPull), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NodePendingAck acknowledges one or more queued pending action IDs.
func (c *Client) NodePendingAck(ctx context.Context, params protocol.NodePendingAckParams) (*protocol.NodePendingAckResult, error) {
	var result protocol.NodePendingAckResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodNodePendingAck), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NodeCanvasCapabilityRefresh mints a fresh node canvas capability token.
func (c *Client) NodeCanvasCapabilityRefresh(ctx context.Context) (*protocol.NodeCanvasCapabilityRefreshResult, error) {
	var result protocol.NodeCanvasCapabilityRefreshResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodNodeCanvasCapabilityRefresh), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
