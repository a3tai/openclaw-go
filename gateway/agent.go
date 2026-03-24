package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// Agent triggers an agent turn and returns the raw response.
func (c *Client) Agent(ctx context.Context, params protocol.AgentParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "agent", params)
}

// AgentIdentity retrieves the identity of an agent.
func (c *Client) AgentIdentity(ctx context.Context, params protocol.AgentIdentityParams) (*protocol.AgentIdentityResult, error) {
	var result protocol.AgentIdentityResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodAgentIdentityGet), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AgentWait waits for an agent run to complete.
func (c *Client) AgentWait(ctx context.Context, params protocol.AgentWaitParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodAgentWait), params)
}
