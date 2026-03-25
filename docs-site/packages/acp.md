---
title: acp
---

# acp

```
import "github.com/a3tai/openclaw-go/acp"
```

Package `acp` implements an Agent Client Protocol (ACP) server. ACP uses JSON-RPC 2.0 over NDJSON (newline-delimited JSON) on stdio, enabling IDE clients to communicate with AI agent processes.

The IDE spawns the agent process and communicates over stdin/stdout. The `Server` handles the read/dispatch loop and bidirectional request/response correlation.

## Architecture

```
IDE (client)  <--stdin/stdout-->  ACP Server (agent)
                  JSON-RPC 2.0
                  NDJSON format
```

## Implementing a Handler

Implement the `Handler` interface to process ACP methods:

```go
type MyHandler struct{}

func (h *MyHandler) Initialize(ctx context.Context, req acp.InitializeRequest) (*acp.InitializeResponse, error) {
    return &acp.InitializeResponse{
        Name:    "my-agent",
        Version: "1.0.0",
        Capabilities: acp.AgentCapabilities{
            Streaming: true,
        },
    }, nil
}

func (h *MyHandler) Authenticate(ctx context.Context, req acp.AuthenticateRequest) (*acp.AuthenticateResponse, error) {
    return &acp.AuthenticateResponse{}, nil
}

func (h *MyHandler) NewSession(ctx context.Context, req acp.NewSessionRequest) (*acp.NewSessionResponse, error) {
    return &acp.NewSessionResponse{
        SessionID: "session-1",
    }, nil
}

func (h *MyHandler) Prompt(ctx context.Context, req acp.PromptRequest) (*acp.PromptResponse, error) {
    // Process the prompt, send session updates...
    return &acp.PromptResponse{
        StopReason: acp.StopReasonEndTurn,
    }, nil
}

// ... implement all 12 Handler methods
```

## Handler Interface

```go
type Handler interface {
    Initialize(ctx, InitializeRequest) (*InitializeResponse, error)
    Authenticate(ctx, AuthenticateRequest) (*AuthenticateResponse, error)
    NewSession(ctx, NewSessionRequest) (*NewSessionResponse, error)
    LoadSession(ctx, LoadSessionRequest) (*LoadSessionResponse, error)
    ListSessions(ctx, ListSessionsRequest) (*ListSessionsResponse, error)
    ForkSession(ctx, ForkSessionRequest) (*ForkSessionResponse, error)
    ResumeSession(ctx, ResumeSessionRequest) (*ResumeSessionResponse, error)
    Prompt(ctx, PromptRequest) (*PromptResponse, error)
    Cancel(ctx, CancelNotification)
    SetSessionMode(ctx, SetSessionModeRequest) (*SetSessionModeResponse, error)
    SetSessionModel(ctx, SetSessionModelRequest) (*SetSessionModelResponse, error)
    SetSessionConfigOption(ctx, SetSessionConfigOptionRequest) (*SetSessionConfigOptionResponse, error)
}
```

## Running the Server

```go
func main() {
    handler := &MyHandler{}
    server := acp.NewServer(handler, os.Stdin, os.Stdout)

    ctx := context.Background()
    if err := server.Serve(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Sending Notifications and Requests

The server can send notifications and requests back to the IDE client:

```go
// Send a session update notification
err := server.SessionUpdate(acp.SessionNotification{
    SessionID: "session-1",
    Updates: []acp.SessionUpdate{
        {
            Type: "text",
            Text: &acp.ContentBlock{Type: "text", Text: "Processing..."},
        },
    },
})

// Send a generic notification
err := server.SendNotification("session/update", params)

// Send a request to the client and wait for a response (agent -> client)
resp, err := server.SendRequest(ctx, "fs/read_text_file", acp.ReadTextFileRequest{
    Path: "/path/to/file.go",
})
```

## Session Update Types

The protocol supports 11 session update variants:

| Update Type | Description |
|-------------|-------------|
| `text` | Text content block |
| `tool_call` | Tool call started/updated |
| `tool_call_update` | Tool call progress update |
| `plan` | Agent plan |
| `usage` | Token usage statistics |
| `modes` | Available session modes changed |
| `models` | Available models changed |
| `config_options` | Config options changed |
| `available_commands` | Available commands changed |
| `status` | Session status change |
| `error` | Error notification |

## Key Types

### Content Blocks (discriminated union)

```go
type ContentBlock struct {
    Type         string          // "text", "image", "audio", "resource_link", "resource"
    Text         string          // for "text"
    MimeType     string          // for binary types
    Data         string          // base64 for binary types
    ResourceURI  string          // for "resource_link"
}
```

### Tool Calls

```go
type ToolCall struct {
    ID        string
    Type      string
    Name      string
    Status    string            // "in_progress", "completed", "errored", "cancelled"
    Kind      string            // "bash", "file_edit", "file_read", etc.
    Content   []ToolCallContent // output content
    Locations []any             // file locations
}
```

### Permissions

```go
type PermissionOption struct {
    ID          string
    Kind        string  // "always", "allow", "deny"
    Title       string
    Description string
}

// Request permission from the user
resp, err := server.SendRequest(ctx, "request_permission", acp.RequestPermissionRequest{
    ToolCall: toolCallUpdate,
    Options:  []acp.PermissionOption{...},
})
```

## Error Codes

```go
acp.ErrCodeParseError            // -32700
acp.ErrCodeInvalidRequest        // -32600
acp.ErrCodeMethodNotFound        // -32601
acp.ErrCodeInvalidParams         // -32602
acp.ErrCodeInternal              // -32603
acp.ErrCodeRequestCancelled      // -32800
acp.ErrCodeAuthenticationNeeded  // -32001
acp.ErrCodeResourceNotFound      // -32002
```

## Constants

```go
acp.ProtocolVersion  // 1

// Stop reasons
acp.StopReasonEndTurn
acp.StopReasonMaxTokens
acp.StopReasonStopSequence
acp.StopReasonToolUse
acp.StopReasonError
acp.StopReasonCancel

// Tool call statuses
acp.ToolCallStatusInProgress
acp.ToolCallStatusCompleted
acp.ToolCallStatusErrored
acp.ToolCallStatusCancelled
```
