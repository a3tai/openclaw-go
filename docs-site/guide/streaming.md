# Streaming

OpenClaw delivers agent responses as a stream of events over the WebSocket connection. The `gateway` package exposes this via the `WithOnEvent` callback.

## Event Flow

When you call `ChatSend`, the gateway returns an immediate ack (`ChatSendResult` with `RunID` and `Status`). The actual response content arrives asynchronously as `"chat"` events on the event handler.

```
ChatSend() ──► {runId: "abc", status: "started"}   (RPC response)

Event stream:
  chat {state: "delta", message: "Hello"}           (first chunk)
  chat {state: "delta", message: " world"}          (next chunk)
  chat {state: "final", message: "Hello world"}     (complete)
```

## Handling Chat Events

Register an event handler when creating the client:

```go
client := gateway.NewClient(
    // ... auth options ...
    gateway.WithOnEvent(func(ev protocol.Event) {
        if ev.EventName == "chat" {
            var chat protocol.ChatEvent
            if err := json.Unmarshal(ev.Payload, &chat); err != nil {
                return
            }
            switch chat.State {
            case "delta":
                fmt.Print(string(chat.Message))
            case "final":
                fmt.Println()
            case "error":
                fmt.Printf("Error: %s\n", chat.ErrorMessage)
            case "aborted":
                fmt.Println("[aborted]")
            }
        }
    }),
)
```

## ChatEvent Fields

| Field | Type | Description |
|-------|------|-------------|
| `RunID` | `string` | Correlates with the `ChatSendResult.RunID` |
| `SessionKey` | `string` | Session this event belongs to |
| `Seq` | `int` | Event sequence number |
| `State` | `string` | `"delta"`, `"final"`, `"aborted"`, `"error"` |
| `Message` | `json.RawMessage` | The content chunk (delta) or full message (final) |
| `ErrorMessage` | `string` | Error description when `State == "error"` |
| `Usage` | `json.RawMessage` | Token usage (present on final) |
| `StopReason` | `string` | Why the response ended (e.g. `"end_turn"`) |

## Other Event Types

The gateway emits many event types beyond chat:

| Event Name | Description |
|------------|-------------|
| `chat` | Agent response chunks and completion |
| `presence` | Presence updates from connected clients |
| `health` | Gateway health status changes |
| `heartbeat` | Periodic heartbeat |
| `tick` | Keepalive tick |
| `cron` | Cron job execution events |
| `exec.approval.requested` | Approval request from agent |
| `plugin.approval.requested` | Plugin approval request |
| `shutdown` | Gateway shutting down |

## Cancellation

Use `ChatAbort` to cancel a running response:

```go
// Start a chat
result, _ := client.ChatSend(ctx, protocol.ChatSendParams{
    SessionKey: "main",
    Message:    "Write me a novel",
    IdempotencyKey: "novel-1",
})

// Cancel it after 2 seconds
time.Sleep(2 * time.Second)
err := client.ChatAbort(ctx, protocol.ChatAbortParams{
    SessionKey: "main",
})
```

You can also use context cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := client.ChatSend(ctx, protocol.ChatSendParams{
    SessionKey: "main",
    Message:    "Long task",
    IdempotencyKey: "long-1",
})
// ctx timeout will cancel the RPC wait
```

## Buffered Processing

For latency-sensitive applications, process events on a separate goroutine:

```go
type ChatDelta struct {
    RunID   string
    Content string
}
deltas := make(chan ChatDelta, 64)

client := gateway.NewClient(
    // ... auth options ...
    gateway.WithOnEvent(func(ev protocol.Event) {
        if ev.EventName == "chat" {
            var chat protocol.ChatEvent
            json.Unmarshal(ev.Payload, &chat)
            if chat.State == "delta" {
                select {
                case deltas <- ChatDelta{RunID: chat.RunID, Content: string(chat.Message)}:
                default:
                    // buffer full — oldest deltas may be lost
                }
            }
        }
    }),
)

// Consumer goroutine
go func() {
    for d := range deltas {
        fmt.Print(d.Content)
    }
}()
```
