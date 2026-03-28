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

	// Require a scheme to avoid silently connecting to the wrong host.
	if host.Scheme == "" || host.Hostname() == "" {
		fmt.Println("Error: host must include a scheme (e.g. wss://mygateway.example.com)")
		os.Exit(1)
	}

	// Normalize the URL scheme for WebSocket.
	switch host.Scheme {
	case "http":
		host.Scheme = "ws"
	case "https":
		host.Scheme = "wss"
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
func (c *Config) Connect(ctx context.Context, client *gateway.Client) error {
	fmt.Printf("Connecting to %s...\n", c.WSURL)
	if err := client.Connect(ctx, c.WSURL); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	fmt.Println("Connected")
	c.SaveDeviceToken(client)
	return nil
}
