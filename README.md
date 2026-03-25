# openclaw-go

[![CI](https://github.com/a3tai/openclaw-go/actions/workflows/ci.yml/badge.svg)](https://github.com/a3tai/openclaw-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/a3tai/openclaw-go/branch/main/graph/badge.svg)](https://codecov.io/gh/a3tai/openclaw-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/a3tai/openclaw-go)](https://goreportcard.com/report/github.com/a3tai/openclaw-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/a3tai/openclaw-go.svg)](https://pkg.go.dev/github.com/a3tai/openclaw-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/a3tai/openclaw-go)](go.mod)

**The** Go client library for [OpenClaw](https://openclaw.ai) — the open gateway for AI agents. Maintained by [A3T](https://a3t.ai).

```
go get github.com/a3tai/openclaw-go
```

Provides typed clients for the Gateway WebSocket protocol, OpenAI-compatible HTTP APIs, local network discovery, and the Agent Client Protocol (ACP).

## Install

```
go get github.com/a3tai/openclaw-go
```

Requires Go 1.25+. The only external dependency is [gorilla/websocket](https://github.com/gorilla/websocket).

## Packages

| Package | Import | Description |
|---------|--------|-------------|
| [`protocol`](docs/protocol.md) | `github.com/a3tai/openclaw-go/protocol` | Wire types, constants, and serialization for the Gateway WebSocket protocol (v3) |
| [`gateway`](docs/gateway.md) | `github.com/a3tai/openclaw-go/gateway` | WebSocket client with full handshake, 96+ typed RPC methods, event/invoke callbacks |
| [`chatcompletions`](docs/chatcompletions.md) | `github.com/a3tai/openclaw-go/chatcompletions` | OpenAI-compatible Chat Completions HTTP client (streaming + non-streaming) |
| [`openresponses`](docs/openresponses.md) | `github.com/a3tai/openclaw-go/openresponses` | OpenAI Responses API HTTP client with typed SSE events and tool calling |
| [`toolsinvoke`](docs/toolsinvoke.md) | `github.com/a3tai/openclaw-go/toolsinvoke` | Tools Invoke HTTP client (`POST /tools/invoke`) |
| [`discovery`](docs/discovery.md) | `github.com/a3tai/openclaw-go/discovery` | mDNS/DNS-SD local network gateway discovery |
| [`acp`](docs/acp.md) | `github.com/a3tai/openclaw-go/acp` | Agent Client Protocol (JSON-RPC 2.0 over NDJSON) server for IDE integration |

## Quick Start

### Gateway WebSocket Client

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := gateway.NewClient(
		gateway.WithToken("my-token"),
		gateway.WithOnEvent(func(ev protocol.Event) {
			fmt.Printf("event: %s\n", ev.EventName)
		}),
	)
	defer client.Close()

	if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
		log.Fatal(err)
	}

	result, err := client.ChatSend(ctx, protocol.ChatSendParams{
		SessionKey: "main",
		Message:    "Hello from Go!",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("chat response: %+v\n", result)
}
```

### Chat Completions (OpenAI-compatible)

```go
client := &chatcompletions.Client{
	BaseURL: "http://localhost:18789",
	Token:   "my-token",
}

resp, _ := client.Create(ctx, chatcompletions.Request{
	Model:    "openclaw:main",
	Messages: []chatcompletions.Message{
		{Role: "user", Content: "Hello!"},
	},
})
fmt.Println(resp.Choices[0].Message.Content)

// Streaming
stream, _ := client.CreateStream(ctx, chatcompletions.Request{
	Model:    "openclaw:main",
	Messages: []chatcompletions.Message{
		{Role: "user", Content: "Tell me about Go"},
	},
})
defer stream.Close()

for {
	chunk, err := stream.Recv()
	if err == io.EOF { break }
	fmt.Print(chunk.Choices[0].Delta.Content)
}
```

### Tools Invoke

```go
client := &toolsinvoke.Client{
	BaseURL: "http://localhost:18789",
	Token:   "my-token",
}

resp, _ := client.Invoke(ctx, toolsinvoke.Request{
	Tool:   "sessions_list",
	Action: "json",
})
fmt.Printf("result: %s\n", resp.Result)
```

### Network Discovery

```go
browser := discovery.NewBrowser()
beacons, _ := browser.Browse(ctx)
for _, b := range beacons {
	fmt.Printf("%s -> %s\n", b.DisplayName, b.WebSocketURL())
}
```

## Examples

Runnable examples are in [`examples/`](examples/). Start the mock server first:

```
go run ./examples/server
```

Then run any example:

```
go run ./examples/chat
go run ./examples/client
go run ./examples/openresponses
go run ./examples/agents
go run ./examples/sessions
go run ./examples/approvals
go run ./examples/pairing
go run ./examples/config
go run ./examples/cron
go run ./examples/node
go run ./examples/discovery
go run ./examples/acp
```

| Example | What it demonstrates |
|---------|---------------------|
| `server` | Mock gateway for local development |
| `client` | All three APIs: WebSocket, Chat Completions, Tools Invoke |
| `chat` | Interactive chat with streaming events |
| `openresponses` | OpenAI Responses API with tool definitions |
| `agents` | Agent CRUD: list, create, update, files, delete |
| `sessions` | Session management: list, preview, patch, usage, reset |
| `approvals` | Exec approval flow: listen, approve/reject, admin config |
| `pairing` | Node and device pairing workflows |
| `config` | Gateway configuration: get, schema, patch, apply |
| `cron` | Cron jobs: list, add, run, history, remove |
| `node` | Connect as a node: declare capabilities, handle invocations |
| `discovery` | Scan the LAN for gateways via mDNS |
| `acp` | ACP agent server over stdio |

## Testing

```
go test ./... -race
go vet ./...
```

All library packages target 100% statement coverage (except `discovery/` at 98.2% due to platform exec wrappers).

## Documentation

See the [`docs/`](docs/) directory for per-package API guides:

- [Protocol Types](docs/protocol.md)
- [Gateway Client](docs/gateway.md)
- [Chat Completions](docs/chatcompletions.md)
- [Open Responses](docs/openresponses.md)
- [Tools Invoke](docs/toolsinvoke.md)
- [Discovery](docs/discovery.md)
- [Agent Client Protocol (ACP)](docs/acp.md)

## About

This project is independently maintained by Rude Company LLC (d/b/a
[a3t.ai](https://a3t.ai)) and is not officially affiliated with the OpenClaw
project.

Many thanks to [Peter Steinberger](https://github.com/steipete) for creating
[OpenClaw](https://github.com/nicepkg/openclaw) and giving it to the world.
His vision for an open, local-first AI gateway has made projects like this
possible. With all love and respect.

**Author:** Steve Rude <steve@a3t.ai>

### Links

- [OpenClaw](https://openclaw.ai) -- the open gateway for AI agents
- [OpenClaw GitHub](https://github.com/nicepkg/openclaw)
- [a3t.ai](https://a3t.ai) -- website for Rude Company LLC, maintainers of this Go client library

## License

MIT -- see [LICENSE](LICENSE) for details.
