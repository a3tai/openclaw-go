package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ChatSend sends a chat message and returns the initial ack.
// Streaming content arrives via "chat" events on the event handler.
func (c *Client) ChatSend(ctx context.Context, params protocol.ChatSendParams) (*protocol.ChatSendResult, error) {
	var result protocol.ChatSendResult
	if err := c.sendRPCTyped(ctx, string(protocol.MethodChatSend), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ChatHistory retrieves chat history for a session.
func (c *Client) ChatHistory(ctx context.Context, params protocol.ChatHistoryParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, string(protocol.MethodChatHistory), params)
}

// ChatAbort aborts a running chat session.
func (c *Client) ChatAbort(ctx context.Context, params protocol.ChatAbortParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodChatAbort), params)
}

// ChatInject injects a message into a chat session.
func (c *Client) ChatInject(ctx context.Context, params protocol.ChatInjectParams) error {
	return c.sendRPCVoid(ctx, string(protocol.MethodChatInject), params)
}
