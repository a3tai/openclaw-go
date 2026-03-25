# Streaming

OpenClaw delivers responses as a stream of delta events. `openclaw-go` exposes this via event callbacks.

## Event Types

| Type | Description |
|------|-------------|
| `delta` | Incremental response chunk — **append** to your buffer |
| `final` | Response complete — contains full text |
| `error` | Stream error |
| `tool_call` | Agent is invoking a tool |
| `tool_result` | Tool returned a result |

::: warning Delta events are incremental
Each `delta` event contains only the **new** chunk, not the full message so far. You must accumulate them yourself.
:::

## Accumulating Deltas

```go
var buf strings.Builder

client.OnEvent(func(evt *protocol.StreamEvent) {
    switch evt.Type {
    case "delta":
        buf.WriteString(evt.Delta) // append each chunk
    case "final":
        fullResponse := buf.String()
        fmt.Println(fullResponse)
        buf.Reset()
    }
})
```

## Cancellation

Pass a cancellable context to abort a stream mid-flight:

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    time.Sleep(2 * time.Second)
    cancel() // abort after 2s
}()

if err := client.Send(ctx, "Write me a novel"); err != nil {
    if errors.Is(err, context.Canceled) {
        fmt.Println("cancelled")
    }
}
```

## Backpressure

The gateway uses `dropIfSlow: true` for delta broadcasts. If your consumer is slow, the server may drop events. For latency-sensitive applications, process deltas on a separate goroutine and use a buffered channel:

```go
deltas := make(chan string, 64)

client.OnEvent(func(evt *protocol.StreamEvent) {
    if evt.Type == "delta" {
        select {
        case deltas <- evt.Delta:
        default:
            // buffer full — handle backpressure
        }
    }
})
```
