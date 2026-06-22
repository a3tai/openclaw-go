package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// TasksList lists background tasks.
func (c *Client) TasksList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTasksList), struct{}{})
}

// TasksGet retrieves a single background task by ID.
func (c *Client) TasksGet(ctx context.Context, params protocol.TasksGetParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodTasksGet), params)
}

// TasksCancel cancels a running background task.
func (c *Client) TasksCancel(ctx context.Context, params protocol.TasksCancelParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodTasksCancel), params)
}
