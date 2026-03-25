# Authentication

OpenClaw uses token-based authentication for all gateway connections.

## Token Types

| Type | Use case |
|------|----------|
| `operator` | Direct human sessions — full permissions |
| `node` | Paired device/agent connections |

## Providing a Token

Pass your token when constructing the client:

```go
client := gateway.NewClient(
    gateway.WithToken("my-token"),
    gateway.WithURL("ws://localhost:18789"),
)
```

## Scopes

Operator tokens can carry scopes that restrict what the client can do:

```go
client := gateway.NewClient(
    gateway.WithToken("my-token"),
    gateway.WithScopes([]string{
        "operator.read",
        "operator.write",
    }),
)
```

Available scopes:

| Scope | Description |
|-------|-------------|
| `operator.read` | Read sessions, agents, and events |
| `operator.write` | Send messages and trigger actions |
| `operator.admin` | Full administrative access |
| `operator.approvals` | Approve or deny exec requests |
| `operator.pairing` | Pair new devices |
| `operator.talk.secrets` | Access secret values in responses |

## Node Mode

To connect as a paired node rather than an operator:

```go
client := gateway.NewClient(
    gateway.WithToken("node-token"),
    gateway.WithNodeMode(true),
)
```

Node connections use a separate handshake and have different permission semantics.
