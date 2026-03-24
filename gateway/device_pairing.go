package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// DevicePairList lists device pairing entries.
func (c *Client) DevicePairList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDevicePairList), struct{}{})
}

// DevicePairApprove approves a device pairing request.
func (c *Client) DevicePairApprove(ctx context.Context, params protocol.DevicePairApproveParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodDevicePairApprove), params)
}

// DevicePairReject rejects a device pairing request.
func (c *Client) DevicePairReject(ctx context.Context, params protocol.DevicePairRejectParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodDevicePairReject), params)
}

// DevicePairRemove removes a paired device.
func (c *Client) DevicePairRemove(ctx context.Context, params protocol.DevicePairRemoveParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodDevicePairRemove), params)
}

// DeviceTokenRotate rotates a device token.
func (c *Client) DeviceTokenRotate(ctx context.Context, params protocol.DeviceTokenRotateParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodDeviceTokenRotate), params)
}

// DeviceTokenRevoke revokes a device token.
func (c *Client) DeviceTokenRevoke(ctx context.Context, params protocol.DeviceTokenRevokeParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodDeviceTokenRevoke), params)
}
