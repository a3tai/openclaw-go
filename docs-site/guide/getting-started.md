# Getting Started

`openclaw-go` is the Go client library for [OpenClaw](https://openclaw.ai). It provides typed clients for every OpenClaw API surface.

## Prerequisites

- Go 1.24 or later
- A running OpenClaw gateway (local or cloud)

## Install

```sh
go get github.com/a3tai/openclaw-go@latest
```

## Connect to a Gateway

The `gateway` package is the primary entry point. Real servers require device identity authentication — without it, the gateway clears self-declared scopes and all scoped RPC calls will fail.

```go
import (
    "context"
    "log"

    "github.com/a3tai/openclaw-go/gateway"
    "github.com/a3tai/openclaw-go/identity"
    "github.com/a3tai/openclaw-go/protocol"
)

// Load or generate a device keypair.
store, _ := identity.NewStore("~/.openclaw-go/identity")
id, _ := store.LoadOrGenerate()
deviceToken := store.LoadDeviceToken()

client := gateway.NewClient(
    gateway.WithToken("your-gateway-token"),
    gateway.WithIdentity(id, deviceToken),
    gateway.WithRole(protocol.RoleOperator),
    gateway.WithScopes(protocol.ScopeOperatorRead, protocol.ScopeOperatorWrite),
)
defer client.Close()

ctx := context.Background()
if err := client.Connect(ctx, "wss://your-gateway.example.com"); err != nil {
    log.Fatal(err)
}

hello := client.Hello()
fmt.Printf("Connected — protocol=%d server=%s\n", hello.Protocol, hello.Server.Version)
```

::: tip Device Pairing
On first connection, the gateway returns `NOT_PAIRED`. Approve the device from the Control UI or CLI, then reconnect. The server issues a device token that the `identity.Store` saves automatically for subsequent connections.
:::

## Send a Message

Use `ChatSend` to send a message. The response is an ack — streaming content arrives via events:

```go
result, err := client.ChatSend(ctx, protocol.ChatSendParams{
    SessionKey:     "main",
    Message:        "What is OpenClaw?",
    IdempotencyKey: "msg-1",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Run started: %s (status: %s)\n", result.RunID, result.Status)
```

## Handle Streaming Events

Register an event handler when creating the client to receive streaming deltas:

```go
client := gateway.NewClient(
    // ... auth options ...
    gateway.WithOnEvent(func(ev protocol.Event) {
        if ev.EventName == "chat" {
            var chat protocol.ChatEvent
            json.Unmarshal(ev.Payload, &chat)
            switch chat.State {
            case "delta":
                fmt.Print(string(chat.Message))
            case "final":
                fmt.Println("\n--- done ---")
            case "error":
                fmt.Printf("Error: %s\n", chat.ErrorMessage)
            }
        }
    }),
)
```

## Local Development

For local development against the mock server:

```sh
# Start the mock server
go run ./examples/server

# Run an example
go run ./examples/chat my-token ws://localhost:18789/ws
```

## Next Steps

- [Authentication](./authentication) — device identity, token types, scopes
- [Streaming](./streaming) — event handling, chat events, cancellation
- [Gateway Method Reference](/packages/gateway-methods) — all 120+ RPC methods
- [gateway package](/packages/gateway) — full package documentation
