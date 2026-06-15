package gateway

import (
	"crypto/tls"
	"encoding/json"
	"time"

	"github.com/a3tai/openclaw-go/identity"
	"github.com/a3tai/openclaw-go/protocol"
)

// Option configures a Client.
type Option func(*options)

type options struct {
	// Auth
	token          string
	password       string
	bootstrapToken string

	// Client identity
	clientInfo   protocol.ClientInfo
	role         protocol.Role
	scopes       []protocol.Scope
	caps         []string
	commands     []string
	permissions  map[string]bool
	device       *protocol.DeviceIdentity
	deviceSigner func(p identity.SigningParams) *protocol.DeviceIdentity
	locale       string
	userAgent    string

	// TLS
	tlsConfig *tls.Config

	// Timeouts
	connectTimeout time.Duration
	tickInterval   time.Duration // overridden by server policy

	// Event handler
	onEvent func(protocol.Event)

	// Invoke handler (node mode)
	onInvoke func(protocol.Invoke) protocol.InvokeResponse

	// Testing hooks
	marshalJSON    func(any) ([]byte, error)
	marshalRequest func(id, method string, params any) ([]byte, error)
}

func defaultOptions() options {
	return options{
		clientInfo: protocol.ClientInfo{
			ID:       protocol.ClientIDGateway,
			Version:  "0.1.0",
			Platform: "go",
			Mode:     protocol.ClientModeBackend,
		},
		role:           protocol.RoleOperator,
		scopes:         []protocol.Scope{protocol.ScopeOperatorRead, protocol.ScopeOperatorWrite},
		locale:         "en-US",
		userAgent:      "openclaw-go/0.1.0",
		connectTimeout: time.Duration(protocol.DefaultHandshakeTimeoutMs) * time.Millisecond,
		tickInterval:   time.Duration(protocol.DefaultTickIntervalMs) * time.Millisecond,
		marshalJSON:    json.Marshal,
		marshalRequest: protocol.MarshalRequest,
	}
}

// WithToken sets the bearer token for authentication.
func WithToken(token string) Option {
	return func(o *options) { o.token = token }
}

// WithPassword sets the password for authentication.
func WithPassword(password string) Option {
	return func(o *options) { o.password = password }
}

// WithBootstrapToken sets a short-lived setup-code bootstrap token (from
// `openclaw qr`) presented in the connect auth alongside a device identity.
// The gateway redeems it to establish and approve the device on first
// connect, then issues a durable device token in hello-ok. A device
// identity (WithIdentity) is required for the gateway to accept it.
func WithBootstrapToken(token string) Option {
	return func(o *options) { o.bootstrapToken = token }
}

// WithClientInfo overrides the default client identity.
func WithClientInfo(info protocol.ClientInfo) Option {
	return func(o *options) { o.clientInfo = info }
}

// WithRole sets the connection role (operator or node).
func WithRole(role protocol.Role) Option {
	return func(o *options) { o.role = role }
}

// WithScopes sets the operator scopes.
func WithScopes(scopes ...protocol.Scope) Option {
	return func(o *options) { o.scopes = scopes }
}

// WithCaps sets the node capability categories.
func WithCaps(caps ...string) Option {
	return func(o *options) { o.caps = caps }
}

// WithCommands sets the node command allowlist.
func WithCommands(commands ...string) Option {
	return func(o *options) { o.commands = commands }
}

// WithPermissions sets the node permission toggles.
func WithPermissions(permissions map[string]bool) Option {
	return func(o *options) { o.permissions = permissions }
}

// WithDevice sets the device identity for pairing.
func WithDevice(device protocol.DeviceIdentity) Option {
	return func(o *options) { o.device = &device }
}

// WithLocale sets the locale string.
func WithLocale(locale string) Option {
	return func(o *options) { o.locale = locale }
}

// WithUserAgent sets the user-agent string.
func WithUserAgent(ua string) Option {
	return func(o *options) { o.userAgent = ua }
}

// WithTLSConfig sets a custom TLS configuration for the WebSocket connection.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(o *options) { o.tlsConfig = cfg }
}

// WithConnectTimeout sets the timeout for the initial connect handshake.
func WithConnectTimeout(d time.Duration) Option {
	return func(o *options) { o.connectTimeout = d }
}

// WithOnEvent registers a callback for incoming Event frames.
func WithOnEvent(fn func(protocol.Event)) Option {
	return func(o *options) { o.onEvent = fn }
}

// WithOnInvoke registers a handler for incoming Invoke frames (node mode).
func WithOnInvoke(fn func(protocol.Invoke) protocol.InvokeResponse) Option {
	return func(o *options) { o.onInvoke = fn }
}
