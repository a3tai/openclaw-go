---
title: toolsinvoke
---

# toolsinvoke

```
import "github.com/a3tai/openclaw-go/toolsinvoke"
```

Package `toolsinvoke` implements an HTTP client for the OpenClaw Tools Invoke endpoint (`POST /tools/invoke`). This allows calling gateway-registered tools directly over HTTP without a WebSocket connection.

## Client Setup

```go
client := &toolsinvoke.Client{
    BaseURL:        "http://localhost:18789",
    Token:          "my-token",
    MessageChannel: "cli",        // optional: x-openclaw-message-channel header
    AccountID:      "acct-123",   // optional: x-openclaw-account-id header
    HTTPClient:     &http.Client{}, // optional
}
```

## Invoking Tools

```go
resp, err := client.Invoke(ctx, toolsinvoke.Request{
    Tool:       "sessions_list",
    Action:     "json",
    Args:       map[string]any{"filter": "active"},
    SessionKey: "main",
    DryRun:     false,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("ok=%v result=%s\n", resp.OK, resp.Result)
```

## Types

### Request

```go
type Request struct {
    Tool       string         `json:"tool"`
    Action     string         `json:"action,omitempty"`
    Args       map[string]any `json:"args,omitempty"`
    SessionKey string         `json:"sessionKey,omitempty"`
    DryRun     bool           `json:"dryRun,omitempty"`
}
```

### Response

```go
type Response struct {
    OK     bool            `json:"ok"`
    Result json.RawMessage `json:"result,omitempty"`
    Error  *ErrorDetail    `json:"error,omitempty"`
}

type ErrorDetail struct {
    Type    string `json:"type"`
    Message string `json:"message"`
}
```

## Error Handling

The package distinguishes between transport-level errors and tool-level errors:

```go
resp, err := client.Invoke(ctx, req)
if err != nil {
    var httpErr *toolsinvoke.HTTPError
    var invokeErr *toolsinvoke.InvokeError

    if errors.As(err, &httpErr) {
        // Transport error (401, 405, 429)
        fmt.Printf("HTTP %d: %s (retry-after: %s)\n",
            httpErr.StatusCode, httpErr.Body, httpErr.RetryAfter)
    } else if errors.As(err, &invokeErr) {
        // Tool-level error (JSON body with ok:false)
        // Note: resp is still populated with the error details
        fmt.Printf("invoke error: %s: %s\n", invokeErr.Type, invokeErr.Message)
    }
}
```

- **`HTTPError`**: Returned for 401 (unauthorized), 405 (method not allowed), 429 (rate limited), and non-JSON responses. Includes `RetryAfter` from the response header.
- **`InvokeError`**: Returned when the endpoint returns a JSON body with `ok: false`. The `Response` is still populated so you can inspect `resp.Error` for details.
