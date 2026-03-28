// Package gwconn provides shared gateway connection setup for examples.
//
// It handles CLI argument parsing, device identity management, and gateway
// client construction so each example can focus on its domain logic.
//
// Every gateway example needs device identity authentication for scoped
// operations on real servers. Without a device identity, the server clears
// self-declared scopes (they can't be cryptographically bound), causing
// "missing scope" errors on every RPC call.
package gwconn

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/identity"
	"github.com/a3tai/openclaw-go/protocol"
)

// Config holds parsed CLI args and resolved connection parameters.
type Config struct {
	Token       string
	WSURL       string
	IdentityDir string
	Store       *identity.Store
	Identity    *identity.Identity
	DeviceToken string
}

// ParseArgs parses standard gateway example CLI arguments:
//
//	<program> <token> <host> [identity-dir]
//
// host is a URL like https://mygateway.example.com or ws://localhost:18789/ws.
// identity-dir defaults to ~/.openclaw-go/identity.
func ParseArgs(programUsage string) *Config {
	if len(os.Args) < 3 {
		fmt.Println(programUsage)
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  token         Gateway shared token")
		fmt.Println("  host          Gateway URL (e.g. https://mygateway.example.com)")
		fmt.Println("  identity-dir  Device keypair directory (default: ~/.openclaw-go/identity)")
		os.Exit(1)
	}

	token := os.Args[1]
	if token == "" {
		fmt.Println("Error: token cannot be empty")
		os.Exit(1)
	}

	host, err := url.Parse(os.Args[2])
	if err != nil {
		fmt.Printf("Error parsing host: %s\n", err)
		os.Exit(1)
	}

	// Normalize the URL scheme for WebSocket.
	if host.Hostname() == "" {
		host.Host = "localhost:18789"
	}
	if host.Scheme == "" || host.Scheme == "http" {
		host.Scheme = "ws"
	}
	if host.Scheme == "https" || host.Scheme == "wss" {
		host.Scheme = "wss"
		if host.Port() == "80" {
			host.Host = fmt.Sprintf("%s:%s", host.Hostname(), "443")
		}
	}

	// Resolve identity directory.
	identityDir := ""
	if len(os.Args) > 3 {
		identityDir = os.Args[3]
	}
	if identityDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("UserHomeDir: %v", err)
		}
		identityDir = filepath.Join(home, ".openclaw-go", "identity")
	}

	// Load or generate the device keypair.
	store, err := identity.NewStore(identityDir)
	if err != nil {
		log.Fatalf("identity store: %v", err)
	}
	id, err := store.LoadOrGenerate()
	if err != nil {
		log.Fatalf("identity load: %v", err)
	}

	deviceToken := store.LoadDeviceToken()

	return &Config{
		Token:       token,
		WSURL:       host.String(),
		IdentityDir: identityDir,
		Store:       store,
		Identity:    id,
		DeviceToken: deviceToken,
	}
}

// ConnectToken returns the token to use for the connect request.
// If a device token exists (from prior pairing), it takes precedence
// over the shared gateway token.
func (c *Config) ConnectToken() string {
	if c.DeviceToken != "" {
		return c.DeviceToken
	}
	return c.Token
}

// PrintIdentityInfo logs the device identity status.
func (c *Config) PrintIdentityInfo() {
	fmt.Printf("Device ID: %s\n", c.Identity.DeviceID)
	if c.DeviceToken != "" {
		fmt.Println("Auth: device token (paired)")
	} else {
		fmt.Println("Auth: shared token (device may need pairing)")
	}
}

// NewClient creates a gateway client configured with device identity auth.
// The caller should pass additional options (scopes, event handlers, etc.)
// which are appended after the base auth options.
func (c *Config) NewClient(extraOpts ...gateway.Option) *gateway.Client {
	baseOpts := []gateway.Option{
		gateway.WithToken(c.ConnectToken()),
		gateway.WithIdentity(c.Identity, c.DeviceToken),
	}
	return gateway.NewClient(append(baseOpts, extraOpts...)...)
}

// SaveDeviceToken checks the hello-ok response for a newly issued device
// token and persists it for subsequent connections.
func (c *Config) SaveDeviceToken(client *gateway.Client) {
	hello := client.Hello()
	if hello == nil || hello.Auth == nil || hello.Auth.DeviceToken == "" {
		return
	}
	if err := c.Store.SaveDeviceToken(hello.Auth.DeviceToken); err != nil {
		log.Printf("Warning: failed to save device token: %v", err)
	} else {
		fmt.Println("Device token saved")
	}
}

// Connect dials the gateway and performs the full handshake. It also
// saves any device token issued in the hello-ok response.
func (c *Config) Connect(ctx context.Context, client *gateway.Client) {
	fmt.Printf("Connecting to %s...\n", c.WSURL)
	if err := client.Connect(ctx, c.WSURL); err != nil {
		log.Fatalf("Connect: %v", err)
	}
	fmt.Println("Connected")
	c.SaveDeviceToken(client)
}

// Scopes is a convenience re-export so examples don't need to import protocol
// just for scope constants.
var (
	ScopeRead      = protocol.ScopeOperatorRead
	ScopeWrite     = protocol.ScopeOperatorWrite
	ScopeAdmin     = protocol.ScopeOperatorAdmin
	ScopeApprovals = protocol.ScopeOperatorApprovals
	ScopePairing   = protocol.ScopeOperatorPairing
)
