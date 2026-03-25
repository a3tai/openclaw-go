# Getting Started

`openclaw-go` is the Go client library for [OpenClaw](https://openclaw.ai). It provides typed clients for every OpenClaw API surface.

## Prerequisites

- Go 1.25 or later
- A running OpenClaw gateway (local or cloud)

## Install

```sh
go get github.com/a3tai/openclaw-go
```

## Connect to a Gateway

The `gateway` package is the primary entry point for most use cases.

```go
import (
    "context"
    "github.com/a3tai/openclaw-go/gateway"
)

client := gateway.NewClient(
    gateway.WithToken("your-token"),
    gateway.WithURL("ws://localhost:18789"),
)

ctx := context.Background()
if err := client.Connect(ctx); err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## Send a Message and Stream a Response

```go
client.OnEvent(func(evt *protocol.StreamEvent) {
    switch evt.Type {
    case "delta":
        fmt.Print(evt.Delta)      // stream chunks as they arrive
    case "final":
        fmt.Println()             // response complete
    case "error":
        log.Println("error:", evt.Error)
    }
})

if err := client.Send(ctx, "What's the weather like?"); err != nil {
    log.Fatal(err)
}

client.Wait() // block until the response finishes
```

## Next Steps

- [Authentication](./authentication) — token types, scopes, and refresh
- [Streaming](./streaming) — delta events, backpressure, and cancellation
- [gateway package reference](../packages/gateway)
