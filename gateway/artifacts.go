package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ArtifactsList lists artifacts, optionally filtered by session key.
func (c *Client) ArtifactsList(ctx context.Context, params protocol.ArtifactsListParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodArtifactsList), params)
}

// ArtifactsGet retrieves a single artifact by ID.
func (c *Client) ArtifactsGet(ctx context.Context, params protocol.ArtifactsGetParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodArtifactsGet), params)
}

// ArtifactsDownload downloads the binary content of an artifact.
func (c *Client) ArtifactsDownload(ctx context.Context, params protocol.ArtifactsDownloadParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodArtifactsDownload), params)
}
