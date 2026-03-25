---
title: chatcompletions
---

# chatcompletions

```
import "github.com/a3tai/openclaw-go/chatcompletions"
```

Package `chatcompletions` implements an HTTP client for the OpenClaw OpenAI-compatible Chat Completions endpoint (`POST /v1/chat/completions`). Supports both non-streaming and SSE streaming modes.

## Client Setup

```go
client := &chatcompletions.Client{
    BaseURL:    "http://localhost:18789",
    Token:      "my-token",
    AgentID:    "main",               // optional: x-openclaw-agent-id header
    SessionKey: "my-session",         // optional: x-openclaw-session-key header
    HTTPClient: &http.Client{},       // optional: custom HTTP client
}
```

## Non-Streaming

```go
resp, err := client.Create(ctx, chatcompletions.Request{
    Model: "openclaw:main",
    Messages: []chatcompletions.Message{
        {Role: "user", Content: "Hello!"},
    },
    Temperature: &temp,
    MaxTokens:   &maxTok,
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.Choices[0].Message.Content)
fmt.Printf("tokens: %d\n", resp.Usage.TotalTokens)
```

## Streaming

```go
stream, err := client.CreateStream(ctx, chatcompletions.Request{
    Model: "openclaw:main",
    Messages: []chatcompletions.Message{
        {Role: "user", Content: "Tell me about Go"},
    },
})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    chunk, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    if len(chunk.Choices) > 0 {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
```

## Types

### Request

```go
type Request struct {
    Model            string    `json:"model"`
    Messages         []Message `json:"messages"`
    Stream           bool      `json:"stream,omitempty"`
    Temperature      *float64  `json:"temperature,omitempty"`
    TopP             *float64  `json:"top_p,omitempty"`
    MaxTokens        *int      `json:"max_tokens,omitempty"`
    Stop             []string  `json:"stop,omitempty"`
    PresencePenalty  *float64  `json:"presence_penalty,omitempty"`
    FrequencyPenalty *float64  `json:"frequency_penalty,omitempty"`
    User             string    `json:"user,omitempty"`
    N                *int      `json:"n,omitempty"`
}
```

### Message

```go
type Message struct {
    Role       string     `json:"role"`
    Content    string     `json:"content,omitempty"`
    Name       string     `json:"name,omitempty"`
    ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
    ToolCallID string     `json:"tool_call_id,omitempty"`
}
```

### Response / StreamChunk

```go
type Response struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   *Usage   `json:"usage,omitempty"`
}

type StreamChunk struct {
    ID      string        `json:"id"`
    Object  string        `json:"object"`
    Created int64         `json:"created"`
    Model   string        `json:"model"`
    Choices []StreamDelta `json:"choices"`
}
```

## Error Handling

Non-200 responses return an `*HTTPError`:

```go
resp, err := client.Create(ctx, req)
if err != nil {
    var httpErr *chatcompletions.HTTPError
    if errors.As(err, &httpErr) {
        fmt.Printf("status=%d retry-after=%s\n", httpErr.StatusCode, httpErr.RetryAfter)
    }
}
```

The `RetryAfter` field is populated from the `Retry-After` response header when present (e.g., on 429 responses).

## OpenClaw Headers

The client automatically sets these headers when the corresponding fields are non-empty:

| Header | Client Field |
|--------|-------------|
| `Authorization: Bearer <token>` | `Token` |
| `x-openclaw-agent-id` | `AgentID` |
| `x-openclaw-session-key` | `SessionKey` |
