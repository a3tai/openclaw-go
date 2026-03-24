package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// WizardStart starts a new wizard session.
func (c *Client) WizardStart(ctx context.Context, params protocol.WizardStartParams) (*protocol.WizardStartResult, error) {
	var result protocol.WizardStartResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodWizardStart), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// WizardNext advances to the next wizard step.
func (c *Client) WizardNext(ctx context.Context, params protocol.WizardNextParams) (*protocol.WizardNextResult, error) {
	var result protocol.WizardNextResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodWizardNext), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// WizardCancel cancels a wizard session.
func (c *Client) WizardCancel(ctx context.Context, params protocol.WizardCancelParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodWizardCancel), params)
}

// WizardStatus retrieves the status of a wizard session.
func (c *Client) WizardStatus(ctx context.Context, params protocol.WizardStatusParams) (*protocol.WizardStatusResult, error) {
	var result protocol.WizardStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodWizardStatus), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
