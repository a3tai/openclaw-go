package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ConfigGet retrieves the current gateway configuration.
func (c *Client) ConfigGet(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodConfigGet), protocol.ConfigGetParams{})
}

// ConfigSet replaces the gateway configuration.
func (c *Client) ConfigSet(ctx context.Context, params protocol.ConfigSetParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodConfigSet), params)
}

// ConfigApply applies a configuration change with optional restart.
func (c *Client) ConfigApply(ctx context.Context, params protocol.ConfigApplyParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodConfigApply), params)
}

// ConfigPatch patches the gateway configuration.
func (c *Client) ConfigPatch(ctx context.Context, params protocol.ConfigPatchParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodConfigPatch), params)
}

// ConfigSchema retrieves the configuration JSON schema.
func (c *Client) ConfigSchema(ctx context.Context) (*protocol.ConfigSchemaResponse, error) {
	var result protocol.ConfigSchemaResponse
	if err := c.sendRPCTyped(ctx, string(protocol.MethodConfigSchema), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ConfigSchemaLookup retrieves schema details for a specific config path.
func (c *Client) ConfigSchemaLookup(ctx context.Context, params protocol.ConfigSchemaLookupParams) (*protocol.ConfigSchemaLookupResult, error) {
	var result protocol.ConfigSchemaLookupResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodConfigSchemaLookup), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
