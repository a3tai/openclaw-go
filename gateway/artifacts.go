package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ArtifactsList lists available artifacts.
func (c *Client) ArtifactsList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodArtifactsList), struct{}{})
}

// ArtifactsGet retrieves a single artifact by ID.
func (c *Client) ArtifactsGet(ctx context.Context, params protocol.ArtifactsGetParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodArtifactsGet), params)
}

// ArtifactsDownload downloads an artifact by ID.
func (c *Client) ArtifactsDownload(ctx context.Context, params protocol.ArtifactsDownloadParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodArtifactsDownload), params)
}
