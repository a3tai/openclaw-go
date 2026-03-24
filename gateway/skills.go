package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// SkillsStatus retrieves the status of installed skills.
func (c *Client) SkillsStatus(ctx context.Context, params protocol.SkillsStatusParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsStatus), params)
}

// SkillsBins retrieves available skill binaries.
func (c *Client) SkillsBins(ctx context.Context) (*protocol.SkillsBinsResult, error) {
	var result protocol.SkillsBinsResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSkillsBins), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SkillsInstall installs a skill.
func (c *Client) SkillsInstall(ctx context.Context, params protocol.SkillsInstallParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsInstall), params)
}

// SkillsUpdate updates a skill's configuration.
func (c *Client) SkillsUpdate(ctx context.Context, params protocol.SkillsUpdateParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSkillsUpdate), params)
}
