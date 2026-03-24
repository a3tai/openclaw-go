package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// CronList lists cron jobs.
func (c *Client) CronList(ctx context.Context, params protocol.CronListParams) ([]protocol.CronJob, error) {
	var result []protocol.CronJob
	if err := c.sendRPCTyped(ctx, string(protocol.MethodCronList), params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CronStatus retrieves the cron system status.
func (c *Client) CronStatus(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodCronStatus), struct{}{})
}

// CronAdd adds a new cron job.
func (c *Client) CronAdd(ctx context.Context, params protocol.CronAddParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodCronAdd), params)
}

// CronUpdate updates an existing cron job.
func (c *Client) CronUpdate(ctx context.Context, params protocol.CronUpdateParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodCronUpdate), params)
}

// CronRemove removes a cron job.
func (c *Client) CronRemove(ctx context.Context, params protocol.CronRemoveParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodCronRemove), params)
}

// CronRun manually runs a cron job.
func (c *Client) CronRun(ctx context.Context, params protocol.CronRunParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodCronRun), params)
}

// CronRuns retrieves the run history for a cron job.
func (c *Client) CronRuns(ctx context.Context, params protocol.CronRunsParams) ([]protocol.CronRunLogEntry, error) {
	var result []protocol.CronRunLogEntry
	if err := c.sendRPCTyped(ctx, string(protocol.MethodCronRuns), params, &result); err != nil {
		return nil, err
	}
	return result, nil
}
