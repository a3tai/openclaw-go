---
layout: home

hero:
  name: openclaw-go
  text: Go client library for OpenClaw
  tagline: Connect agents, stream responses, discover gateways — typed and production-ready.
  actions:
    - theme: brand
      text: Get Started →
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/a3tai/openclaw-go

features:
  - icon: 🔌
    title: gateway
    details: WebSocket client with full handshake, 96+ typed RPC methods, streaming event callbacks, and node mode support.
    link: /packages/gateway
    linkText: Explore package
  - icon: 📡
    title: protocol
    details: Wire types, constants, and serialization for the OpenClaw Gateway WebSocket protocol (v3).
    link: /packages/protocol
    linkText: Explore package
  - icon: 💬
    title: chatcompletions
    details: OpenAI-compatible Chat Completions HTTP client with streaming SSE and non-streaming modes.
    link: /packages/chatcompletions
    linkText: Explore package
  - icon: 🔍
    title: discovery
    details: mDNS/DNS-SD local network discovery to locate OpenClaw gateways on your LAN.
    link: /packages/discovery
    linkText: Explore package
  - icon: 🤖
    title: acp
    details: Agent Client Protocol server — JSON-RPC 2.0 over NDJSON for IDE and tool integration.
    link: /packages/acp
    linkText: Explore package
  - icon: 🛠
    title: toolsinvoke
    details: Typed HTTP client for the Tools Invoke endpoint to run agent tools directly.
    link: /packages/toolsinvoke
    linkText: Explore package
---

## Install

```sh
go get github.com/a3tai/openclaw-go
```

Requires **Go 1.25+**. Single external dependency: [gorilla/websocket](https://github.com/gorilla/websocket).

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/a3tai/openclaw-go/gateway"
    "github.com/a3tai/openclaw-go/protocol"
)

func main() {
    client := gateway.NewClient(
        gateway.WithToken("my-token"),
        gateway.WithURL("ws://localhost:18789"),
    )

    if err := client.Connect(context.Background()); err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    client.OnEvent(func(evt *protocol.StreamEvent) {
        if evt.Type == "delta" {
            fmt.Print(evt.Delta)
        }
    })

    if err := client.Send(context.Background(), "Hello from Go!"); err != nil {
        log.Fatal(err)
    }

    client.Wait()
}
```

## About

`openclaw-go` is maintained by [A3T](https://a3t.ai) — the enterprise agent platform built on OpenClaw.

[![CI](https://github.com/a3tai/openclaw-go/actions/workflows/ci.yml/badge.svg)](https://github.com/a3tai/openclaw-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/a3tai/openclaw-go.svg)](https://pkg.go.dev/github.com/a3tai/openclaw-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/a3tai/openclaw-go/blob/main/LICENSE)
