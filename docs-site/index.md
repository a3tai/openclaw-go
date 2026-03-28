---
layout: home

hero:
  name: openclaw-go
  text: The Go SDK for OpenClaw
  tagline: Connect agents, stream responses, discover gateways — typed and production-ready.
  image:
    src: /mascot-gpt.png
    alt: openclaw-go mascot
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/a3tai/openclaw-go

features:
  - icon: "\U0001F50C"
    title: gateway
    details: WebSocket client with full handshake, 120+ typed RPC methods, device identity auth, streaming events, and node mode.
    link: /packages/gateway
    linkText: Explore
  - icon: "\U0001F4E1"
    title: protocol
    details: Wire types, constants, and serialization for the OpenClaw Gateway WebSocket protocol (v3).
    link: /packages/protocol
    linkText: Explore
  - icon: "\U0001F4AC"
    title: chatcompletions
    details: OpenAI-compatible Chat Completions HTTP client with streaming SSE and non-streaming modes.
    link: /packages/chatcompletions
    linkText: Explore
  - icon: "\U0001F916"
    title: openresponses
    details: OpenAI Responses API client with streaming SSE, function calls, and file search.
    link: /packages/openresponses
    linkText: Explore
  - icon: "\U0001F50D"
    title: discovery
    details: mDNS/DNS-SD local network discovery to locate OpenClaw gateways on your LAN.
    link: /packages/discovery
    linkText: Explore
  - icon: "\U0001F6E0"
    title: toolsinvoke
    details: Typed HTTP client for the Tools Invoke endpoint to run agent tools directly.
    link: /packages/toolsinvoke
    linkText: Explore
  - icon: "\U0001F5DD"
    title: identity
    details: Ed25519 device keypair management for authenticated gateway connections.
    link: /packages/identity
    linkText: Explore
  - icon: "\U0001F4BB"
    title: acp
    details: Agent Client Protocol server — JSON-RPC 2.0 over NDJSON for IDE and tool integration.
    link: /packages/acp
    linkText: Explore
---

## Install

```sh
go get github.com/a3tai/openclaw-go@latest
```

Requires **Go 1.24+**. Single external dependency: [gorilla/websocket](https://github.com/gorilla/websocket).

## Quick Start

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/a3tai/openclaw-go/gateway"
    "github.com/a3tai/openclaw-go/identity"
    "github.com/a3tai/openclaw-go/protocol"
)

func main() {
    // Load or generate a device keypair for authenticated access.
    store, _ := identity.NewStore("~/.openclaw-go/identity")
    id, _ := store.LoadOrGenerate()
    deviceToken := store.LoadDeviceToken()

    client := gateway.NewClient(
        gateway.WithToken("my-token"),
        gateway.WithIdentity(id, deviceToken),
        gateway.WithRole(protocol.RoleOperator),
        gateway.WithScopes(protocol.ScopeOperatorRead, protocol.ScopeOperatorWrite),
        gateway.WithOnEvent(func(ev protocol.Event) {
            if ev.EventName == "chat" {
                var chat protocol.ChatEvent
                json.Unmarshal(ev.Payload, &chat)
                if chat.State == "delta" {
                    fmt.Print(string(chat.Message))
                }
            }
        }),
    )
    defer client.Close()

    ctx := context.Background()
    if err := client.Connect(ctx, "wss://my-gateway.example.com"); err != nil {
        log.Fatal(err)
    }

    // Send a message — streaming deltas arrive via the event handler above.
    result, err := client.ChatSend(ctx, protocol.ChatSendParams{
        SessionKey:     "main",
        Message:        "Hello from Go!",
        IdempotencyKey: "hello-1",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Run started: %s\n", result.RunID)
}
```

## About

`openclaw-go` is maintained by [A3T](https://a3t.ai) — the enterprise agent platform built on OpenClaw.

[![CI](https://github.com/a3tai/openclaw-go/actions/workflows/ci.yml/badge.svg)](https://github.com/a3tai/openclaw-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/a3tai/openclaw-go.svg)](https://pkg.go.dev/github.com/a3tai/openclaw-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/a3tai/openclaw-go/blob/main/LICENSE)
