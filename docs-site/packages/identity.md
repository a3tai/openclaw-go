# identity

```
import "github.com/a3tai/openclaw-go/identity"
```

Package `identity` manages Ed25519 device keypairs for OpenClaw gateway authentication. Real gateway servers require a signed device identity to grant scoped access — without it, the server clears self-declared scopes.

## Why Device Identity?

The OpenClaw gateway protocol uses a challenge-response handshake. During connect, the client signs a payload containing the challenge nonce, device ID, client info, role, scopes, token, and timestamp with its Ed25519 private key. The server verifies this signature to cryptographically bind the declared scopes to the device. Without a device identity, the server accepts the connection but strips all scopes.

## Store

The `Store` manages the keypair and device token files on disk.

```go
store, err := identity.NewStore("~/.openclaw-go/identity")
if err != nil {
    log.Fatal(err)
}
```

Files are stored with strict permissions:
- Directory: `0700` (owner only)
- `keypair.json`: `0600` (private key material)
- `device-token`: `0600` (bearer credential)

## Load or Generate

```go
id, err := store.LoadOrGenerate()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Device ID: %s\n", id.DeviceID)
```

On first call, this generates a new Ed25519 keypair and persists it. On subsequent calls, it loads the existing keypair. The device ID is a hex-encoded SHA-256 hash of the public key.

## Device Token

After a device is paired with the gateway, the server issues a device token in the `hello-ok` response. The store manages this token:

```go
// Load existing device token (empty string if not yet paired)
deviceToken := store.LoadDeviceToken()

// Save a device token received from the server
err := store.SaveDeviceToken(token)
```

## Using with the Gateway Client

Pass the identity and device token when constructing the gateway client:

```go
client := gateway.NewClient(
    gateway.WithToken(token),             // shared token or device token
    gateway.WithIdentity(id, deviceToken), // Ed25519 device identity
    gateway.WithRole(protocol.RoleOperator),
    gateway.WithScopes(protocol.ScopeOperatorRead, protocol.ScopeOperatorWrite),
)
```

The `WithIdentity` option configures the client to sign the connect handshake with the device's private key. The signed payload follows the v2 format:

```
v2|deviceId|clientId|clientMode|role|scopes|signedAtMs|token|nonce
```

## Identity Type

```go
type Identity struct {
    DeviceID  string           // hex SHA-256 of public key
    PublicKey ed25519.PublicKey // 32-byte Ed25519 public key
    // private key is unexported
}
```

The private key is not exported. Signing operations are performed via `BuildDeviceIdentity`, which is called internally by `gateway.WithIdentity`.

## Pairing Flow

1. **Generate keypair**: `store.LoadOrGenerate()` creates the Ed25519 keypair
2. **Connect**: Client sends signed device identity in the connect handshake
3. **Server rejects**: Returns `NOT_PAIRED: pairing required` with the device ID
4. **Approve**: Admin approves the device via CLI or Control UI
5. **Reconnect**: Server accepts and issues a device token in `hello-ok`
6. **Persist**: `store.SaveDeviceToken(token)` saves for future use
7. **Subsequent connections**: Use the device token for instant auth
