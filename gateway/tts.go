package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// TTSStatus retrieves TTS status.
func (c *Client) TTSStatus(ctx context.Context) (*protocol.TTSStatusResult, error) {
	var result protocol.TTSStatusResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodTTSStatus), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TTSPersonas retrieves TTS persona presets.
func (c *Client) TTSPersonas(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTTSPersonas), struct{}{})
}

// TTSProviders retrieves TTS provider information.
func (c *Client) TTSProviders(ctx context.Context) (*protocol.TTSProvidersResult, error) {
	var result protocol.TTSProvidersResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodTTSProviders), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TTSEnable enables TTS.
func (c *Client) TTSEnable(ctx context.Context) (*protocol.TTSEnableResult, error) {
	var result protocol.TTSEnableResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodTTSEnable), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TTSDisable disables TTS.
func (c *Client) TTSDisable(ctx context.Context) (*protocol.TTSDisableResult, error) {
	var result protocol.TTSDisableResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodTTSDisable), struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TTSConvert converts text to speech.
func (c *Client) TTSConvert(ctx context.Context, params protocol.TTSConvertParams) (*protocol.TTSConvertResult, error) {
	var result protocol.TTSConvertResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodTTSConvert), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TTSSetProvider sets the active TTS provider.
func (c *Client) TTSSetProvider(ctx context.Context, params protocol.TTSSetProviderParams) (*protocol.TTSSetProviderResult, error) {
	var result protocol.TTSSetProviderResult
	if err := c.sendRPCTyped(ctx, "tts.setProvider", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
