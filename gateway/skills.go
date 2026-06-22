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

// SkillsSearch searches for skills on ClawHub by query string.
// Requires the operator.read scope.
// Query and Limit are both optional; Limit is capped at 100 by the server.
func (c *Client) SkillsSearch(ctx context.Context, params protocol.SkillsSearchParams) (*protocol.SkillsSearchResult, error) {
	var result protocol.SkillsSearchResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSkillsSearch), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SkillsDetail retrieves full metadata for a skill from ClawHub by slug.
// Requires the operator.read scope.
// The Skill field in the result is nil when no skill matches the slug.
func (c *Client) SkillsDetail(ctx context.Context, params protocol.SkillsDetailParams) (*protocol.SkillsDetailResult, error) {
	var result protocol.SkillsDetailResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSkillsDetail), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SkillsUploadBegin begins a chunked skill upload.
func (c *Client) SkillsUploadBegin(ctx context.Context, params protocol.SkillsUploadBeginParams) (*protocol.SkillsUploadBeginResult, error) {
	var result protocol.SkillsUploadBeginResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSkillsUploadBegin), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SkillsUploadChunk uploads a chunk of skill data.
func (c *Client) SkillsUploadChunk(ctx context.Context, params protocol.SkillsUploadChunkParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSkillsUploadChunk), params)
}

// SkillsUploadCommit finalizes a chunked skill upload and returns an install ID.
func (c *Client) SkillsUploadCommit(ctx context.Context, params protocol.SkillsUploadCommitParams) (*protocol.SkillsUploadCommitResult, error) {
	var result protocol.SkillsUploadCommitResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSkillsUploadCommit), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SkillsProposalsList lists skill proposals.
func (c *Client) SkillsProposalsList(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsList), params)
}

// SkillsProposalsCreate creates a skill proposal.
func (c *Client) SkillsProposalsCreate(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsCreate), params)
}

// SkillsProposalsInspect inspects a skill proposal.
func (c *Client) SkillsProposalsInspect(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsInspect), params)
}

// SkillsProposalsApply applies a skill proposal.
func (c *Client) SkillsProposalsApply(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsApply), params)
}

// SkillsProposalsReject rejects a skill proposal.
func (c *Client) SkillsProposalsReject(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsReject), params)
}

// SkillsProposalsRevise revises a skill proposal.
func (c *Client) SkillsProposalsRevise(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsRevise), params)
}

// SkillsProposalsUpdate updates a skill proposal.
func (c *Client) SkillsProposalsUpdate(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsUpdate), params)
}

// SkillsProposalsQuarantine quarantines a skill proposal.
func (c *Client) SkillsProposalsQuarantine(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodSkillsProposalsQuarantine), params)
}
