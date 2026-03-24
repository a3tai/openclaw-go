package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ToolsCatalog returns the gateway tool catalog grouped by profile/source.
func (c *Client) ToolsCatalog(ctx context.Context, params protocol.ToolsCatalogParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodToolsCatalog), params)
}
