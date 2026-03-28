---
title: openresponses
---

# openresponses

```
import "github.com/a3tai/openclaw-go/openresponses"
```

Package `openresponses` implements an HTTP client for the OpenClaw OpenAI Responses API-compatible endpoint (`POST /v1/responses`). This is the newer API (vs Chat Completions) supporting structured input items, typed SSE streaming events, and tool/function calling.

## Client Setup

```go
client := &openresponses.Client{
    BaseURL:    "http://localhost:18789",
    Token:      "my-token",
    AgentID:    "main",               // optional
    SessionKey: "my-session",         // optional
    HTTPClient: &http.Client{},       // optional
}
```

## Non-Streaming

```go
maxTokens := 200
resp, err := client.Create(ctx, openresponses.Request{
    Model: "openclaw:main",
    Input: openresponses.InputFromItems([]openresponses.InputItem{
        openresponses.MessageItem("user", "What is OpenClaw?"),
    }),
    MaxOutputTokens: &maxTokens,
})
if err != nil {
    log.Fatal(err)
}

for _, item := range resp.Output {
    if item.Type == "message" {
        for _, part := range item.Content {
            if part.Type == "output_text" {
                fmt.Println(part.Text)
            }
        }
    }
}
```

## Streaming

```go
stream, err := client.CreateStream(ctx, openresponses.Request{
    Model: "openclaw:main",
    Input: openresponses.InputFromString("Hello!"),
})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    ev, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }

    switch ev.EventType {
    case "response.created":
        fmt.Println("response started")
    case "response.output_text.delta":
        var delta openresponses.OutputTextDeltaEvent
        json.Unmarshal(ev.RawData, &delta)
        fmt.Print(delta.Delta)
    case "response.completed":
        fmt.Println("\ndone")
    }
}
```

## Input Construction

The `Input` field accepts either a plain string or structured items:

```go
// Simple string input
req := openresponses.Request{
    Model: "openclaw:main",
    Input: openresponses.InputFromString("Hello!"),
}

// Structured input items
req := openresponses.Request{
    Model: "openclaw:main",
    Input: openresponses.InputFromItems([]openresponses.InputItem{
        openresponses.MessageItem("system", "You are helpful."),
        openresponses.MessageItem("user", "What is Go?"),
    }),
}

// Multi-part content
req := openresponses.Request{
    Model: "openclaw:main",
    Input: openresponses.InputFromItems([]openresponses.InputItem{
        openresponses.MessageItemParts("user", []openresponses.ContentPart{
            {Type: "input_text", Text: "Describe this image"},
            {Type: "input_image", ImageURL: "https://example.com/img.png"},
        }),
    }),
}

// Function call output (multi-turn tool use)
req := openresponses.Request{
    Model: "openclaw:main",
    Input: openresponses.InputFromItems([]openresponses.InputItem{
        openresponses.MessageItem("user", "What's the weather?"),
        openresponses.FunctionCallItem("call_123", "get_weather", `{"location":"SF"}`),
        openresponses.FunctionCallOutputItem("call_123", `{"temp":72,"unit":"F"}`),
    }),
}
```

## Tool Definitions

```go
req := openresponses.Request{
    Model: "openclaw:main",
    Input: openresponses.InputFromString("What's the weather in SF?"),
    Tools: []openresponses.ToolDefinition{
        {
            Type: "function",
            Function: openresponses.FunctionTool{
                Name:        "get_weather",
                Description: "Get current weather for a location",
                Parameters: map[string]any{
                    "type": "object",
                    "properties": map[string]any{
                        "location": map[string]string{
                            "type": "string",
                            "description": "City name",
                        },
                    },
                    "required": []string{"location"},
                },
            },
        },
    },
    ToolChoice: "auto", // or "none", "required", or openresponses.ToolChoiceFunction{...}
}
```

## Stream Event Types

| Event Type | Description |
|------------|-------------|
| `response.created` | Response object created |
| `response.in_progress` | Processing started |
| `response.output_item.added` | New output item added |
| `response.content_part.added` | New content part added |
| `response.output_text.delta` | Incremental text chunk |
| `response.output_text.done` | Text generation complete |
| `response.function_call_arguments.delta` | Incremental function args |
| `response.function_call_arguments.done` | Function args complete |
| `response.output_item.done` | Output item complete |
| `response.content_part.done` | Content part complete |
| `response.completed` | Full response complete |
| `response.failed` | Response failed |

Each event's `RawData` field contains the full JSON payload. Unmarshal into the corresponding typed struct (e.g., `OutputTextDeltaEvent`, `ResponseEvent`, `OutputItemEvent`).

## Error Handling

```go
resp, err := client.Create(ctx, req)
if err != nil {
    var httpErr *openresponses.HTTPError
    if errors.As(err, &httpErr) {
        fmt.Printf("status=%d retry-after=%s\n", httpErr.StatusCode, httpErr.RetryAfter)
    }
}
```
