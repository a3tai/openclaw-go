package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// AgentsList lists all agents.
func (c *Client) AgentsList(ctx context.Context) (*protocol.AgentsListResult, error) {
	var result protocol.AgentsListResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAgentsList), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AgentsCreate creates a new agent.
func (c *Client) AgentsCreate(ctx context.Context, params protocol.AgentsCreateParams) (*protocol.AgentsCreateResult, error) {
	var result protocol.AgentsCreateResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAgentsCreate), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AgentsUpdate updates an agent.
func (c *Client) AgentsUpdate(ctx context.Context, params protocol.AgentsUpdateParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodAgentsUpdate), params)
}

// AgentsDelete deletes an agent.
func (c *Client) AgentsDelete(ctx context.Context, params protocol.AgentsDeleteParams) (*protocol.AgentsDeleteResult, error) {
	var result protocol.AgentsDeleteResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAgentsDelete), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AgentsFilesList lists files for an agent.
func (c *Client) AgentsFilesList(ctx context.Context, params protocol.AgentsFilesListParams) (*protocol.AgentsFilesListResult, error) {
	var result protocol.AgentsFilesListResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAgentsFilesList), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AgentsFilesGet retrieves a specific agent file.
func (c *Client) AgentsFilesGet(ctx context.Context, params protocol.AgentsFilesGetParams) (*protocol.AgentsFilesGetResult, error) {
	var result protocol.AgentsFilesGetResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAgentsFilesGet), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AgentsFilesSet creates or updates an agent file.
func (c *Client) AgentsFilesSet(ctx context.Context, params protocol.AgentsFilesSetParams) (*protocol.AgentsFilesSetResult, error) {
	var result protocol.AgentsFilesSetResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAgentsFilesSet), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
