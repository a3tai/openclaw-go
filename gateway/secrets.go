package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// SecretsReload reloads secrets from configured secret sources.
func (c *Client) SecretsReload(ctx context.Context) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodSecretsReload), protocol.SecretsReloadParams{})
}

// SecretsResolve resolves secret assignments for command targets.
func (c *Client) SecretsResolve(ctx context.Context, params protocol.SecretsResolveParams) (*protocol.SecretsResolveResult, error) {
	var result protocol.SecretsResolveResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodSecretsResolve), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
