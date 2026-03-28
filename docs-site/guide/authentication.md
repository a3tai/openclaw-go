# Authentication

OpenClaw uses device identity authentication for gateway connections. A shared token bootstraps the initial connection, but a paired Ed25519 device keypair is required for scoped access on real servers.

## Why Device Identity?

The gateway server clears self-declared scopes for clients that don't present a cryptographically signed device identity. Without it, the connection succeeds but every scoped RPC call fails with `missing scope`. This prevents clients from escalating privileges by declaring scopes they don't have.

## Setup

### 1. Generate a Device Keypair

```go
import "github.com/a3tai/openclaw-go/identity"

store, err := identity.NewStore("~/.openclaw-go/identity")
if err != nil {
    log.Fatal(err)
}
id, err := store.LoadOrGenerate()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Device ID: %s\n", id.DeviceID)
```

This creates an Ed25519 keypair in `~/.openclaw-go/identity/keypair.json` (permissions `0600`). The device ID is a SHA-256 hash of the public key.

### 2. Connect with Identity

```go
deviceToken := store.LoadDeviceToken() // empty on first run

client := gateway.NewClient(
    gateway.WithToken("your-shared-gateway-token"),
    gateway.WithIdentity(id, deviceToken),
    gateway.WithRole(protocol.RoleOperator),
    gateway.WithScopes(
        protocol.ScopeOperatorRead,
        protocol.ScopeOperatorWrite,
    ),
)
```

### 3. Approve the Device

On first connection, the server responds with `NOT_PAIRED: pairing required`. Approve the device from a client that's already paired:

```sh
# From the CLI or Control UI
openclaw device pair approve --device <device-id>
```

After approval, reconnect. The server issues a device token in the `hello-ok` response. The `identity.Store` persists it for subsequent connections.

## Roles

| Role | Constant | Description |
|------|----------|-------------|
| Operator | `protocol.RoleOperator` | Human sessions — uses scopes for permission control |
| Node | `protocol.RoleNode` | Capability host — bypasses scope checks but still requires device identity |

```go
// Operator (default)
gateway.WithRole(protocol.RoleOperator)

// Node
gateway.WithRole(protocol.RoleNode)
```

## Scopes

Operator connections declare scopes that control which RPC methods are accessible:

```go
gateway.WithScopes(
    protocol.ScopeOperatorRead,
    protocol.ScopeOperatorWrite,
    protocol.ScopeOperatorAdmin,
)
```

| Scope | Constant | Grants |
|-------|----------|--------|
| `operator.read` | `ScopeOperatorRead` | Read sessions, agents, config, nodes, cron, models, logs, health, presence |
| `operator.write` | `ScopeOperatorWrite` | Send messages, chat, invoke nodes, approval requests, talk, system events |
| `operator.admin` | `ScopeOperatorAdmin` | Config changes, session management, cron CRUD, agent CRUD, node management, TTS, secrets, wizard |
| `operator.approvals` | `ScopeOperatorApprovals` | Resolve exec/plugin approval decisions |
| `operator.pairing` | `ScopeOperatorPairing` | Device and node pairing operations |

::: warning Common Mistake
Many RPC methods (cron.add, agents.create, config.set, sessions.patch) require `operator.admin`. If you get `missing scope: operator.admin` errors, add `ScopeOperatorAdmin` to your `WithScopes` call.
:::

## Password Authentication

For gateways configured with password auth instead of token auth:

```go
gateway.WithPassword("your-password")
```

## Node Mode

Nodes provide execution capabilities (shell, filesystem, tools) to the gateway. Connect as a node with capabilities and an invoke handler:

```go
client := gateway.NewClient(
    gateway.WithToken("node-token"),
    gateway.WithIdentity(id, deviceToken),
    gateway.WithRole(protocol.RoleNode),
    gateway.WithCaps("exec", "fs"),
    gateway.WithCommands("shell", "read_file", "write_file"),
    gateway.WithOnInvoke(func(inv protocol.Invoke) protocol.InvokeResponse {
        // Handle commands from the gateway
        return protocol.InvokeResponse{OK: true, Payload: result}
    }),
)
```
