---
title: discovery
---

# discovery

```
import "github.com/a3tai/openclaw-go/discovery"
```

Package `discovery` discovers OpenClaw Gateway instances on the local network using mDNS/DNS-SD (Bonjour). The gateway advertises itself as `_openclaw-gw._tcp`.

## Platform Support

| Platform | Implementation | Tool |
|----------|---------------|------|
| macOS | Native | `dns-sd` CLI |
| Linux | Native | `avahi-browse` CLI |
| Windows | Stub | Returns "not supported" error |

## Scanning for Gateways

```go
browser := discovery.NewBrowser()

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

beacons, err := browser.Browse(ctx)
if err != nil {
    log.Fatal(err)
}

for _, b := range beacons {
    fmt.Printf("Found: %s\n", b.DisplayName)
    fmt.Printf("  Host: %s:%d\n", b.Host, b.Port)
    fmt.Printf("  WebSocket: %s\n", b.WebSocketURL())
    fmt.Printf("  HTTP: %s\n", b.HTTPURL())
    fmt.Printf("  TLS: %v\n", b.GatewayTLS)
    if b.TailnetDNS != "" {
        fmt.Printf("  Tailnet: %s\n", b.TailnetDNS)
    }
}
```

## Beacon Fields

The `Beacon` struct contains all information from the mDNS TXT records:

| Field | Description |
|-------|-------------|
| `InstanceName` | mDNS instance name |
| `Domain` | DNS-SD domain (e.g., `"local."`) |
| `Host` | Resolved hostname from SRV record |
| `Port` | Resolved port from SRV record |
| `DisplayName` | Human-readable gateway name |
| `LanHost` | `<hostname>.local` address |
| `TailnetDNS` | MagicDNS hostname (Tailscale) |
| `GatewayPort` | Gateway WebSocket/HTTP port |
| `GatewayTLS` | Whether TLS is enabled |
| `GatewayTLSFingerprint` | SHA-256 certificate fingerprint |
| `SSHPort` | SSH port on the gateway host |
| `CLIPath` | Absolute path to the `openclaw` CLI |
| `Role` | Beacon role (always `"gateway"`) |
| `Transport` | Transport type (always `"gateway"`) |
| `CanvasPort` | Canvas host port if enabled |
| `TXT` | Raw TXT record key-value map |

## URL Helpers

```go
// WebSocket URL (ws:// or wss:// based on GatewayTLS)
wsURL := beacon.WebSocketURL()

// HTTP URL (http:// or https:// based on GatewayTLS)
httpURL := beacon.HTTPURL()
```

Host resolution priority: SRV host > TailnetDNS > LanHost > `127.0.0.1`.

Port resolution priority: SRV port > GatewayPort > `18789` (default).

## Constants

```go
discovery.ServiceType        // "_openclaw-gw._tcp"
discovery.DefaultGatewayPort // 18789
```
