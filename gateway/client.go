// Package gateway implements an OpenClaw Gateway WebSocket client.
//
// Reference: https://docs.openclaw.ai/gateway/protocol
package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/a3tai/openclaw-go/identity"
	"github.com/a3tai/openclaw-go/protocol"
	"github.com/gorilla/websocket"
)

// Client is a WebSocket client for the OpenClaw Gateway protocol.
type Client struct {
	opts options

	conn   *websocket.Conn
	connMu sync.Mutex

	// pending tracks outstanding request IDs to their response channels.
	pending   map[string]chan *protocol.Response
	pendingMu sync.Mutex

	// hello stores the server's hello-ok payload after a successful connect.
	hello *protocol.HelloOK

	// seq is an atomic counter for generating unique request IDs.
	seq atomic.Int64

	// done is closed when the client is shut down.
	done     chan struct{}
	doneOnce sync.Once

	// readLoopDone is closed when the read loop exits.
	readLoopDone chan struct{}

	// tickStop stops the keepalive ticker.
	tickStop     chan struct{}
	tickStopOnce sync.Once
}

// NewClient creates a new Gateway client but does not connect.
// Call Connect to initiate the WebSocket handshake.
func NewClient(opts ...Option) *Client {
	o := defaultOptions()
	for _, fn := range opts {
		fn(&o)
	}
	return &Client{
		opts:         o,
		pending:      make(map[string]chan *protocol.Response),
		done:         make(chan struct{}),
		readLoopDone: make(chan struct{}),
		tickStop:     make(chan struct{}),
	}
}

// Hello returns the server's hello-ok payload. Returns nil if not connected.
func (c *Client) Hello() *protocol.HelloOK {
	return c.hello
}

// Connect dials the gateway at the given WebSocket URL and performs the
// protocol handshake (challenge → connect → hello-ok).
func (c *Client) Connect(ctx context.Context, wsURL string) error {
	dialer := websocket.Dialer{
		HandshakeTimeout: c.opts.connectTimeout,
		TLSClientConfig:  c.opts.tlsConfig,
	}

	conn, _, err := dialer.DialContext(ctx, wsURL, http.Header{})
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()

	// 1. Read the connect.challenge event.
	challenge, err := c.readChallenge(ctx)
	if err != nil {
		conn.Close()
		return fmt.Errorf("read challenge: %w", err)
	}

	// 2. Send the connect request.
	reqID := c.nextID()
	params := c.buildConnectParams(challenge)

	data, err := c.opts.marshalRequest(reqID, "connect", params)
	if err != nil {
		conn.Close()
		return fmt.Errorf("marshal connect: %w", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		conn.Close()
		return fmt.Errorf("send connect: %w", err)
	}

	// 3. Read the connect response (hello-ok).
	if err := c.readHelloOK(ctx, reqID); err != nil {
		conn.Close()
		return fmt.Errorf("hello: %w", err)
	}

	// Start background loops.
	go c.readLoop()
	go c.tickLoop()

	return nil
}

// Send sends a request and waits for the matching response.
func (c *Client) Send(ctx context.Context, method string, params any) (*protocol.Response, error) {
	id := c.nextID()
	data, err := protocol.MarshalRequest(id, method, params)
	if err != nil {
		return nil, err
	}

	ch := make(chan *protocol.Response, 1)
	c.pendingMu.Lock()
	c.pending[id] = ch
	c.pendingMu.Unlock()
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	c.connMu.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	c.connMu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.done:
		return nil, errors.New("client closed")
	case resp := <-ch:
		return resp, nil
	}
}

// SendEvent sends a one-way event frame.
func (c *Client) SendEvent(eventName string, payload any) error {
	data, err := protocol.MarshalEvent(eventName, payload)
	if err != nil {
		return err
	}
	c.connMu.Lock()
	defer c.connMu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Close gracefully shuts down the client.
func (c *Client) Close() error {
	select {
	case <-c.done:
		return nil // already closed
	default:
	}
	c.doneOnce.Do(func() { close(c.done) })
	c.tickStopOnce.Do(func() { close(c.tickStop) })

	c.connMu.Lock()
	conn := c.conn
	c.connMu.Unlock()

	if conn != nil {
		// Send a close frame and then close.
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		return conn.Close()
	}
	return nil
}

// Done returns a channel that is closed when the client shuts down.
func (c *Client) Done() <-chan struct{} {
	return c.done
}

// ---------------------------------------------------------------------------
// Internal
// ---------------------------------------------------------------------------

func (c *Client) nextID() string {
	return fmt.Sprintf("go-%d", c.seq.Add(1))
}

func (c *Client) readChallenge(ctx context.Context) (*protocol.ConnectChallenge, error) {
	_ = c.conn.SetReadDeadline(time.Now().Add(c.opts.connectTimeout))
	defer func() { _ = c.conn.SetReadDeadline(time.Time{}) }()

	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var ev protocol.Event
	if err := json.Unmarshal(msg, &ev); err != nil {
		return nil, fmt.Errorf("unmarshal challenge event: %w", err)
	}
	if ev.Type != protocol.FrameTypeEvent || ev.EventName != "connect.challenge" {
		return nil, fmt.Errorf("expected connect.challenge event, got type=%s event=%s", ev.Type, ev.EventName)
	}
	var ch protocol.ConnectChallenge
	if err := json.Unmarshal(ev.Payload, &ch); err != nil {
		return nil, fmt.Errorf("unmarshal challenge payload: %w", err)
	}
	return &ch, nil
}

func (c *Client) buildConnectParams(challenge *protocol.ConnectChallenge) protocol.ConnectParams {
	params := protocol.ConnectParams{
		MinProtocol: protocol.ProtocolVersion,
		MaxProtocol: protocol.ProtocolVersion,
		Client:      c.opts.clientInfo,
		Role:        c.opts.role,
		Scopes:      c.opts.scopes,
		Caps:        c.opts.caps,
		Commands:    c.opts.commands,
		Permissions: c.opts.permissions,
		Locale:      c.opts.locale,
		UserAgent:   c.opts.userAgent,
	}

	// Auth: token or password.
	if c.opts.token != "" {
		params.Auth = protocol.AuthParams{Token: c.opts.token}
	} else if c.opts.password != "" {
		params.Auth = protocol.AuthParams{Password: c.opts.password}
	}

	// Device identity (includes challenge nonce for signing).
	if c.opts.deviceSigner != nil {
		nonce := ""
		if challenge != nil {
			nonce = challenge.Nonce
		}
		// Build the signing params from the connect request fields so the
		// signer can construct the v2 pipe-delimited auth payload.
		scopeStrs := make([]string, len(c.opts.scopes))
		for i, s := range c.opts.scopes {
			scopeStrs[i] = string(s)
		}
		sp := identity.SigningParams{
			ClientID:   c.opts.clientInfo.ID,
			ClientMode: c.opts.clientInfo.Mode,
			Role:       string(c.opts.role),
			Scopes:     scopeStrs,
			Token:      c.opts.token,
			Nonce:      nonce,
		}
		params.Device = c.opts.deviceSigner(sp)
	} else if c.opts.device != nil {
		dev := *c.opts.device
		if challenge != nil {
			dev.Nonce = challenge.Nonce
		}
		params.Device = &dev
	}

	return params
}

func (c *Client) readHelloOK(ctx context.Context, reqID string) error {
	_ = c.conn.SetReadDeadline(time.Now().Add(c.opts.connectTimeout))
	defer func() { _ = c.conn.SetReadDeadline(time.Time{}) }()

	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}
	var resp protocol.Response
	if err := json.Unmarshal(msg, &resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	if resp.ID != reqID {
		return fmt.Errorf("response id mismatch: got %q, want %q", resp.ID, reqID)
	}
	if !resp.OK {
		if resp.Error != nil {
			return fmt.Errorf("connect rejected: %s: %s", resp.Error.Code, resp.Error.Message)
		}
		return errors.New("connect rejected (no error details)")
	}

	var hello protocol.HelloOK
	if err := json.Unmarshal(resp.Payload, &hello); err != nil {
		return fmt.Errorf("unmarshal hello-ok: %w", err)
	}
	c.hello = &hello

	// Apply server-provided tick interval.
	if hello.Policy.TickIntervalMs > 0 {
		c.opts.tickInterval = time.Duration(hello.Policy.TickIntervalMs) * time.Millisecond
	}

	return nil
}

func (c *Client) readLoop() {
	defer close(c.readLoopDone)
	for {
		select {
		case <-c.done:
			return
		default:
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			// Connection closed or error — shut down.
			c.doneOnce.Do(func() { close(c.done) })
			return
		}

		frame, err := protocol.ParseFrame(msg)
		if err != nil {
			continue
		}

		switch frame.Type {
		case protocol.FrameTypeResponse:
			var resp protocol.Response
			if json.Unmarshal(msg, &resp) == nil {
				c.pendingMu.Lock()
				ch, ok := c.pending[resp.ID]
				c.pendingMu.Unlock()
				if ok {
					ch <- &resp
				}
			}

		case protocol.FrameTypeEvent:
			var ev protocol.Event
			if json.Unmarshal(msg, &ev) == nil && c.opts.onEvent != nil {
				// Dispatch asynchronously so the readLoop never blocks on
				// a slow event handler (e.g. manager's handleGatewayEvent
				// acquiring locks or writing to a buffered channel). Blocking
				// here would stall all WebSocket reads, causing stream
				// deltas and finals to queue up behind slow handlers.
				go c.opts.onEvent(ev)
			}

		case protocol.FrameTypeRequest:
			// Shouldn't normally happen, but handle gracefully.

		case protocol.FrameTypeInvoke:
			// Gateway→Node invoke
			var inv protocol.Invoke
			if json.Unmarshal(msg, &inv) == nil && c.opts.onInvoke != nil {
				go func() {
					res := c.opts.onInvoke(inv)
					res.Type = string(protocol.FrameTypeInvokeResponse)
					res.ID = inv.ID
					data, err := c.opts.marshalJSON(res)
					if err != nil {
						return
					}
					select {
					case <-c.done:
						return
					default:
					}
					c.connMu.Lock()
					_ = c.conn.WriteMessage(websocket.TextMessage, data)
					c.connMu.Unlock()
				}()
			}
		}
	}
}

func (c *Client) tickLoop() {
	ticker := time.NewTicker(c.opts.tickInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.tickStop:
			return
		case <-c.done:
			return
		case <-ticker.C:
			c.connMu.Lock()
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			c.connMu.Unlock()
			if err != nil {
				return
			}
		}
	}
}
