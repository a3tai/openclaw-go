// Package protocol defines the OpenClaw Gateway WebSocket protocol types.
//
// Reference: https://docs.openclaw.ai/gateway/protocol
package protocol

import "encoding/json"

// ProtocolVersion is the current protocol version.
const ProtocolVersion = 3

// ---------------------------------------------------------------------------
// Server Constants
// ---------------------------------------------------------------------------

const (
	// MaxPayloadBytes is the maximum size of a single WebSocket frame payload.
	MaxPayloadBytes = 25 * 1024 * 1024 // 25 MiB

	// MaxBufferedBytes is the maximum total buffered bytes per connection.
	MaxBufferedBytes = 50 * 1024 * 1024 // 50 MiB

	// DefaultTickIntervalMs is the default keepalive tick interval.
	DefaultTickIntervalMs = 30_000

	// DefaultHandshakeTimeoutMs is the default handshake timeout.
	DefaultHandshakeTimeoutMs = 10_000

	// DedupeTTLMs is the TTL for deduplication entries.
	DedupeTTLMs = 300_000

	// DedupeMax is the maximum number of deduplication entries.
	DedupeMax = 1000

	// DefaultMaxChatHistoryMessagesBytes is the default max chat history size.
	DefaultMaxChatHistoryMessagesBytes = 6 * 1024 * 1024 // 6 MiB

	// HealthRefreshIntervalMs is the interval for health snapshot refreshes.
	HealthRefreshIntervalMs = 60_000

	// SessionLabelMaxLength is the maximum length of a session label.
	SessionLabelMaxLength = 64
)

// ---------------------------------------------------------------------------
// Error Codes
// ---------------------------------------------------------------------------

// Error codes returned in gateway RPC error responses.
const (
	ErrorCodeNotLinked      = "NOT_LINKED"
	ErrorCodeNotPaired      = "NOT_PAIRED"
	ErrorCodeAgentTimeout   = "AGENT_TIMEOUT"
	ErrorCodeInvalidRequest = "INVALID_REQUEST"
	ErrorCodeUnavailable    = "UNAVAILABLE"
)

// ---------------------------------------------------------------------------
// Client IDs
// ---------------------------------------------------------------------------

// Known client ID values for ClientInfo.ID.
const (
	ClientIDWebchatUI   = "webchat-ui"
	ClientIDControlUI   = "openclaw-control-ui"
	ClientIDWebchat     = "webchat"
	ClientIDCLI         = "cli"
	ClientIDGateway     = "gateway-client"
	ClientIDMacOS       = "openclaw-macos"
	ClientIDIOS         = "openclaw-ios"
	ClientIDAndroid     = "openclaw-android"
	ClientIDNodeHost    = "node-host"
	ClientIDTest        = "test"
	ClientIDFingerprint = "fingerprint"
	ClientIDProbe       = "openclaw-probe"
)

// ---------------------------------------------------------------------------
// Client Modes
// ---------------------------------------------------------------------------

// Client mode values for ClientInfo.Mode.
const (
	ClientModeWebchat = "webchat"
	ClientModeCLI     = "cli"
	ClientModeUI      = "ui"
	ClientModeBackend = "backend"
	ClientModeNode    = "node"
	ClientModeProbe   = "probe"
	ClientModeTest    = "test"
)

// ---------------------------------------------------------------------------
// Client Capabilities
// ---------------------------------------------------------------------------

// Client capability values for ConnectParams.Caps.
const (
	ClientCapToolEvents = "tool-events"
)

// ---------------------------------------------------------------------------
// Frame Types
// ---------------------------------------------------------------------------

// FrameType identifies the kind of WebSocket frame.
type FrameType string

// Frame type constants for use with RawFrame, Request, Response, and Event.
const (
	FrameTypeRequest        FrameType = "req"
	FrameTypeResponse       FrameType = "res"
	FrameTypeEvent          FrameType = "event"
	FrameTypeInvoke         FrameType = "invoke"
	FrameTypeInvokeResponse FrameType = "invoke-res"
)

// ---------------------------------------------------------------------------
// Framing
// ---------------------------------------------------------------------------

// Request is a client→gateway RPC request frame.
type Request struct {
	Type   FrameType       `json:"type"`             // always "req"
	ID     string          `json:"id"`               // unique request id
	Method string          `json:"method"`           // RPC method name
	Params json.RawMessage `json:"params,omitempty"` // method-specific params
}

// Response is a gateway→client RPC response frame.
type Response struct {
	Type    FrameType       `json:"type"`              // always "res"
	ID      string          `json:"id"`                // matches request id
	OK      bool            `json:"ok"`                // success flag
	Payload json.RawMessage `json:"payload,omitempty"` // success payload
	Error   *ErrorPayload   `json:"error,omitempty"`   // error details
}

// Event is a uni-directional notification frame (either direction).
type Event struct {
	Type         FrameType       `json:"type"`                   // always "event"
	EventName    string          `json:"event"`                  // event name
	Payload      json.RawMessage `json:"payload,omitempty"`      // event-specific data
	Seq          *int64          `json:"seq,omitempty"`          // optional sequence number
	StateVersion *StateVersion   `json:"stateVersion,omitempty"` // optional state version
}

// ErrorPayload carries structured error information (spec: ErrorShape).
type ErrorPayload struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	Details      any    `json:"details,omitempty"`
	Retryable    *bool  `json:"retryable,omitempty"`
	RetryAfterMs *int   `json:"retryAfterMs,omitempty"`
}

// RawFrame is used for initial deserialization to determine the frame type.
type RawFrame struct {
	Type  FrameType `json:"type"`
	Event string    `json:"event,omitempty"` // only for event frames
}

// ---------------------------------------------------------------------------
// Roles & Scopes
// ---------------------------------------------------------------------------

// Role identifies a connection's role in the gateway.
type Role string

// Role constants for ConnectParams.Role.
const (
	RoleOperator Role = "operator"
	RoleNode     Role = "node"
)

// Scope is a capability scope for operator connections.
type Scope string

// Scope constants for ConnectParams.Scopes.
const (
	ScopeOperatorRead      Scope = "operator.read"
	ScopeOperatorWrite     Scope = "operator.write"
	ScopeOperatorAdmin     Scope = "operator.admin"
	ScopeOperatorApprovals Scope = "operator.approvals"
	ScopeOperatorPairing   Scope = "operator.pairing"
)

// ---------------------------------------------------------------------------
// State Version
// ---------------------------------------------------------------------------

// StateVersion tracks the version of presence and health snapshots.
type StateVersion struct {
	Presence int `json:"presence"`
	Health   int `json:"health"`
}

// ---------------------------------------------------------------------------
// Connect handshake
// ---------------------------------------------------------------------------

// ConnectChallenge is the server-initiated challenge sent before the client
// sends its connect request.
type ConnectChallenge struct {
	Nonce string `json:"nonce"`
	Ts    int64  `json:"ts"`
}

// ClientInfo describes the connecting client software.
type ClientInfo struct {
	ID              string `json:"id"`
	DisplayName     string `json:"displayName,omitempty"`
	Version         string `json:"version"`
	Platform        string `json:"platform"`
	DeviceFamily    string `json:"deviceFamily,omitempty"`
	ModelIdentifier string `json:"modelIdentifier,omitempty"`
	Mode            string `json:"mode"`
	InstanceID      string `json:"instanceId,omitempty"`
}

// DeviceIdentity carries the device's identity for pairing and auth.
type DeviceIdentity struct {
	ID        string `json:"id"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	SignedAt  int64  `json:"signedAt"`
	Nonce     string `json:"nonce,omitempty"`
}

// AuthParams carries auth credentials for a connect request.
type AuthParams struct {
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty"`
}

// ConnectParams is the params payload for a "connect" request.
type ConnectParams struct {
	MinProtocol int             `json:"minProtocol"`
	MaxProtocol int             `json:"maxProtocol"`
	Client      ClientInfo      `json:"client"`
	Role        Role            `json:"role,omitempty"`
	Scopes      []Scope         `json:"scopes,omitempty"`
	Caps        []string        `json:"caps,omitempty"`
	Commands    []string        `json:"commands,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`
	PathEnv     string          `json:"pathEnv,omitempty"`
	Auth        AuthParams      `json:"auth,omitempty"`
	Locale      string          `json:"locale,omitempty"`
	UserAgent   string          `json:"userAgent,omitempty"`
	Device      *DeviceIdentity `json:"device,omitempty"`
}

// ---------------------------------------------------------------------------
// Hello-OK response
// ---------------------------------------------------------------------------

// HelloOK is the payload returned in a successful connect response.
type HelloOK struct {
	Type          string        `json:"type"` // "hello-ok"
	Protocol      int           `json:"protocol"`
	Server        HelloServer   `json:"server"`
	Features      HelloFeatures `json:"features"`
	Snapshot      Snapshot      `json:"snapshot"`
	CanvasHostURL string        `json:"canvasHostUrl,omitempty"`
	Auth          *HelloAuth    `json:"auth,omitempty"`
	Policy        HelloPolicy   `json:"policy"`
}

// HelloServer identifies the gateway server.
type HelloServer struct {
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
	Host    string `json:"host,omitempty"`
	ConnID  string `json:"connId"`
}

// HelloFeatures lists the RPC methods and events the server supports.
type HelloFeatures struct {
	Methods []string `json:"methods"`
	Events  []string `json:"events"`
}

// HelloPolicy contains operational parameters from the server.
type HelloPolicy struct {
	MaxPayload       int `json:"maxPayload"`
	MaxBufferedBytes int `json:"maxBufferedBytes"`
	TickIntervalMs   int `json:"tickIntervalMs"`
}

// HelloAuth is returned when a device token is issued at connect time.
type HelloAuth struct {
	DeviceToken string   `json:"deviceToken"`
	Role        string   `json:"role"`
	Scopes      []string `json:"scopes"`
	IssuedAtMs  *int64   `json:"issuedAtMs,omitempty"`
}

// ---------------------------------------------------------------------------
// Snapshot
// ---------------------------------------------------------------------------

// Snapshot is the initial server state sent in hello-ok.
type Snapshot struct {
	Presence        []PresenceEntry  `json:"presence"`
	Health          json.RawMessage  `json:"health"`
	StateVersion    StateVersion     `json:"stateVersion"`
	UptimeMs        int64            `json:"uptimeMs"`
	ConfigPath      string           `json:"configPath,omitempty"`
	StateDir        string           `json:"stateDir,omitempty"`
	SessionDefaults *SessionDefaults `json:"sessionDefaults,omitempty"`
	AuthMode        string           `json:"authMode,omitempty"`
}

// SessionDefaults are the default session settings from the server.
type SessionDefaults struct {
	DefaultAgentID string `json:"defaultAgentId"`
	MainKey        string `json:"mainKey"`
	MainSessionKey string `json:"mainSessionKey"`
	Scope          string `json:"scope,omitempty"`
}

// ---------------------------------------------------------------------------
// Presence
// ---------------------------------------------------------------------------

// PresenceEntry is a single entry from system-presence.
type PresenceEntry struct {
	Host             string   `json:"host,omitempty"`
	IP               string   `json:"ip,omitempty"`
	Version          string   `json:"version,omitempty"`
	Platform         string   `json:"platform,omitempty"`
	DeviceFamily     string   `json:"deviceFamily,omitempty"`
	ModelIdentifier  string   `json:"modelIdentifier,omitempty"`
	Mode             string   `json:"mode,omitempty"`
	LastInputSeconds *int     `json:"lastInputSeconds,omitempty"`
	Reason           string   `json:"reason,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Text             string   `json:"text,omitempty"`
	Ts               int64    `json:"ts"`
	DeviceID         string   `json:"deviceId,omitempty"`
	Roles            []string `json:"roles,omitempty"`
	Scopes           []string `json:"scopes,omitempty"`
	InstanceID       string   `json:"instanceId,omitempty"`
}

// ---------------------------------------------------------------------------
// Exec approvals
// ---------------------------------------------------------------------------

// ExecApprovalRequestParams are the params for "exec.approval.request".
type ExecApprovalRequestParams struct {
	ID           string  `json:"id,omitempty"`
	Command      string  `json:"command"`
	Cwd          *string `json:"cwd,omitempty"`
	Host         *string `json:"host,omitempty"`
	Security     *string `json:"security,omitempty"`
	Ask          *string `json:"ask,omitempty"`
	AgentID      *string `json:"agentId,omitempty"`
	ResolvedPath *string `json:"resolvedPath,omitempty"`
	SessionKey   *string `json:"sessionKey,omitempty"`
	TimeoutMs    *int    `json:"timeoutMs,omitempty"`
	TwoPhase     *bool   `json:"twoPhase,omitempty"`
}

// ExecApprovalResolveParams are the params for "exec.approval.resolve".
type ExecApprovalResolveParams struct {
	ID       string `json:"id"`
	Decision string `json:"decision"`
}

// ExecApprovalRequested is the payload of an "exec.approval.requested" event.
type ExecApprovalRequested struct {
	ID           string  `json:"id,omitempty"`
	Command      string  `json:"command,omitempty"`
	Cwd          *string `json:"cwd,omitempty"`
	Host         *string `json:"host,omitempty"`
	Security     *string `json:"security,omitempty"`
	Ask          *string `json:"ask,omitempty"`
	AgentID      *string `json:"agentId,omitempty"`
	ResolvedPath *string `json:"resolvedPath,omitempty"`
	SessionKey   *string `json:"sessionKey,omitempty"`
	TimeoutMs    *int    `json:"timeoutMs,omitempty"`
	TwoPhase     *bool   `json:"twoPhase,omitempty"`
}

// ExecApprovalsGetParams are the params for "exec.approvals.get" (empty).
type ExecApprovalsGetParams struct{}

// ExecApprovalsSetParams are the params for "exec.approvals.set".
type ExecApprovalsSetParams struct {
	File     ExecApprovalsFile `json:"file"`
	BaseHash string            `json:"baseHash,omitempty"`
}

// ExecApprovalsNodeGetParams are the params for "exec.approvals.node.get".
type ExecApprovalsNodeGetParams struct {
	NodeID string `json:"nodeId"`
}

// ExecApprovalsNodeSetParams are the params for "exec.approvals.node.set".
type ExecApprovalsNodeSetParams struct {
	NodeID   string            `json:"nodeId"`
	File     ExecApprovalsFile `json:"file"`
	BaseHash string            `json:"baseHash,omitempty"`
}

// ExecApprovalsFile is the exec approvals configuration file.
type ExecApprovalsFile struct {
	Version  int                           `json:"version"`
	Socket   *ExecApprovalsSocket          `json:"socket,omitempty"`
	Defaults *ExecApprovalsDefaults        `json:"defaults,omitempty"`
	Agents   map[string]ExecApprovalsAgent `json:"agents,omitempty"`
}

// ExecApprovalsSocket is the socket configuration in exec approvals.
type ExecApprovalsSocket struct {
	Path  string `json:"path,omitempty"`
	Token string `json:"token,omitempty"`
}

// ExecApprovalsDefaults are the default exec approval settings.
type ExecApprovalsDefaults struct {
	Security        string `json:"security,omitempty"`
	Ask             string `json:"ask,omitempty"`
	AskFallback     string `json:"askFallback,omitempty"`
	AutoAllowSkills *bool  `json:"autoAllowSkills,omitempty"`
}

// ExecApprovalsAgent is per-agent exec approval settings.
type ExecApprovalsAgent struct {
	Security        string                        `json:"security,omitempty"`
	Ask             string                        `json:"ask,omitempty"`
	AskFallback     string                        `json:"askFallback,omitempty"`
	AutoAllowSkills *bool                         `json:"autoAllowSkills,omitempty"`
	Allowlist       []ExecApprovalsAllowlistEntry `json:"allowlist,omitempty"`
}

// ExecApprovalsAllowlistEntry is a single entry in the exec approvals allowlist.
type ExecApprovalsAllowlistEntry struct {
	ID               string `json:"id,omitempty"`
	Pattern          string `json:"pattern"`
	LastUsedAt       *int64 `json:"lastUsedAt,omitempty"`
	LastUsedCommand  string `json:"lastUsedCommand,omitempty"`
	LastResolvedPath string `json:"lastResolvedPath,omitempty"`
}

// ExecApprovalsSnapshot is the full exec approvals state.
type ExecApprovalsSnapshot struct {
	Path   string            `json:"path"`
	Exists bool              `json:"exists"`
	Hash   string            `json:"hash"`
	File   ExecApprovalsFile `json:"file"`
}

// ---------------------------------------------------------------------------
// Exec lifecycle (node events)
// ---------------------------------------------------------------------------

// ExecFinished is the payload of an "exec.finished" event from a node.
type ExecFinished struct {
	SessionKey string `json:"sessionKey"`
	RunID      string `json:"runId,omitempty"`
	Command    string `json:"command,omitempty"`
	ExitCode   *int   `json:"exitCode,omitempty"`
	TimedOut   *bool  `json:"timedOut,omitempty"`
	Success    *bool  `json:"success,omitempty"`
	Output     string `json:"output,omitempty"`
}

// ExecDenied is the payload of an "exec.denied" event from a node.
type ExecDenied struct {
	SessionKey string `json:"sessionKey"`
	RunID      string `json:"runId,omitempty"`
	Command    string `json:"command,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

// ---------------------------------------------------------------------------
// Node invoke (gateway→node)
// ---------------------------------------------------------------------------

// Invoke is a gateway→node command invocation.
type Invoke struct {
	Type    string          `json:"type"`    // "invoke"
	ID      string          `json:"id"`      // request id
	Command string          `json:"command"` // e.g. "camera.snap"
	Params  json.RawMessage `json:"params,omitempty"`
}

// InvokeResponse is a node→gateway response to an invoke.
type InvokeResponse struct {
	Type    string          `json:"type"`              // "invoke-res"
	ID      string          `json:"id"`                // matches invoke id
	OK      bool            `json:"ok"`                // success flag
	Payload json.RawMessage `json:"payload,omitempty"` // result
	Error   *ErrorPayload   `json:"error,omitempty"`   // error details
}

// ---------------------------------------------------------------------------
// Chat types
// ---------------------------------------------------------------------------

// ChatSendParams are the params for "chat.send".
type ChatSendParams struct {
	SessionKey     string          `json:"sessionKey"`
	Message        string          `json:"message"`
	Thinking       string          `json:"thinking,omitempty"`
	Deliver        *bool           `json:"deliver,omitempty"`
	Attachments    json.RawMessage `json:"attachments,omitempty"`
	TimeoutMs      *int            `json:"timeoutMs,omitempty"`
	IdempotencyKey string          `json:"idempotencyKey"`
}

// ChatHistoryParams are the params for "chat.history".
type ChatHistoryParams struct {
	SessionKey string `json:"sessionKey"`
	Limit      *int   `json:"limit,omitempty"`
}

// ChatAbortParams are the params for "chat.abort".
type ChatAbortParams struct {
	SessionKey string `json:"sessionKey"`
	RunID      string `json:"runId,omitempty"`
}

// ChatInjectParams are the params for "chat.inject".
type ChatInjectParams struct {
	SessionKey string `json:"sessionKey"`
	Message    string `json:"message"`
	Label      string `json:"label,omitempty"`
}

// ChatEvent is the payload of a "chat" event.
type ChatEvent struct {
	RunID        string          `json:"runId"`
	SessionKey   string          `json:"sessionKey"`
	Seq          int             `json:"seq"`
	State        string          `json:"state"` // "delta", "final", "aborted", "error"
	Message      json.RawMessage `json:"message,omitempty"`
	ErrorMessage string          `json:"errorMessage,omitempty"`
	Usage        json.RawMessage `json:"usage,omitempty"`
	StopReason   string          `json:"stopReason,omitempty"`
}

// ---------------------------------------------------------------------------
// Agent types
// ---------------------------------------------------------------------------

// InputProvenance describes the provenance of an agent input.
type InputProvenance struct {
	Kind             string `json:"kind"` // "external_user", "inter_session", "internal_system"
	SourceSessionKey string `json:"sourceSessionKey,omitempty"`
	SourceChannel    string `json:"sourceChannel,omitempty"`
	SourceTool       string `json:"sourceTool,omitempty"`
}

// AgentParams are the params for "agent".
type AgentParams struct {
	Message           string           `json:"message"`
	AgentID           string           `json:"agentId,omitempty"`
	To                string           `json:"to,omitempty"`
	ReplyTo           string           `json:"replyTo,omitempty"`
	SessionID         string           `json:"sessionId,omitempty"`
	SessionKey        string           `json:"sessionKey,omitempty"`
	Thinking          string           `json:"thinking,omitempty"`
	Deliver           *bool            `json:"deliver,omitempty"`
	Attachments       json.RawMessage  `json:"attachments,omitempty"`
	Channel           string           `json:"channel,omitempty"`
	ReplyChannel      string           `json:"replyChannel,omitempty"`
	AccountID         string           `json:"accountId,omitempty"`
	ReplyAccountID    string           `json:"replyAccountId,omitempty"`
	ThreadID          string           `json:"threadId,omitempty"`
	GroupID           string           `json:"groupId,omitempty"`
	GroupChannel      string           `json:"groupChannel,omitempty"`
	GroupSpace        string           `json:"groupSpace,omitempty"`
	Timeout           *int             `json:"timeout,omitempty"`
	Lane              string           `json:"lane,omitempty"`
	ExtraSystemPrompt string           `json:"extraSystemPrompt,omitempty"`
	InputProvenance   *InputProvenance `json:"inputProvenance,omitempty"`
	IdempotencyKey    string           `json:"idempotencyKey"`
	Label             string           `json:"label,omitempty"`
	SpawnedBy         string           `json:"spawnedBy,omitempty"`
}

// AgentEvent is the payload of an "agent" event.
type AgentEvent struct {
	RunID  string         `json:"runId"`
	Seq    int            `json:"seq"`
	Stream string         `json:"stream"`
	Ts     int64          `json:"ts"`
	Data   map[string]any `json:"data"`
}

// AgentIdentityParams are the params for "agent.identity.get".
type AgentIdentityParams struct {
	AgentID    string `json:"agentId,omitempty"`
	SessionKey string `json:"sessionKey,omitempty"`
}

// AgentIdentityResult is the result of "agent.identity.get".
type AgentIdentityResult struct {
	AgentID string `json:"agentId"`
	Name    string `json:"name,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
	Emoji   string `json:"emoji,omitempty"`
}

// AgentWaitParams are the params for "agent.wait".
type AgentWaitParams struct {
	RunID     string `json:"runId"`
	TimeoutMs *int   `json:"timeoutMs,omitempty"`
}

// ---------------------------------------------------------------------------
// Session types
// ---------------------------------------------------------------------------

// SessionsListParams are the params for "sessions.list".
type SessionsListParams struct {
	Limit                *int   `json:"limit,omitempty"`
	ActiveMinutes        *int   `json:"activeMinutes,omitempty"`
	IncludeGlobal        *bool  `json:"includeGlobal,omitempty"`
	IncludeUnknown       *bool  `json:"includeUnknown,omitempty"`
	IncludeDerivedTitles *bool  `json:"includeDerivedTitles,omitempty"`
	IncludeLastMessage   *bool  `json:"includeLastMessage,omitempty"`
	Label                string `json:"label,omitempty"`
	SpawnedBy            string `json:"spawnedBy,omitempty"`
	AgentID              string `json:"agentId,omitempty"`
	Search               string `json:"search,omitempty"`
}

// SessionsPreviewParams are the params for "sessions.preview".
type SessionsPreviewParams struct {
	Keys     []string `json:"keys"`
	Limit    *int     `json:"limit,omitempty"`
	MaxChars *int     `json:"maxChars,omitempty"`
}

// SessionsResolveParams are the params for "sessions.resolve".
type SessionsResolveParams struct {
	Key            string `json:"key,omitempty"`
	SessionID      string `json:"sessionId,omitempty"`
	Label          string `json:"label,omitempty"`
	AgentID        string `json:"agentId,omitempty"`
	SpawnedBy      string `json:"spawnedBy,omitempty"`
	IncludeGlobal  *bool  `json:"includeGlobal,omitempty"`
	IncludeUnknown *bool  `json:"includeUnknown,omitempty"`
}

// SessionsPatchParams are the params for "sessions.patch".
type SessionsPatchParams struct {
	Key             string  `json:"key"`
	Label           *string `json:"label,omitempty"`
	ThinkingLevel   *string `json:"thinkingLevel,omitempty"`
	VerboseLevel    *string `json:"verboseLevel,omitempty"`
	ReasoningLevel  *string `json:"reasoningLevel,omitempty"`
	ResponseUsage   *string `json:"responseUsage,omitempty"`
	ElevatedLevel   *string `json:"elevatedLevel,omitempty"`
	ExecHost        *string `json:"execHost,omitempty"`
	ExecSecurity    *string `json:"execSecurity,omitempty"`
	ExecAsk         *string `json:"execAsk,omitempty"`
	ExecNode        *string `json:"execNode,omitempty"`
	Model           *string `json:"model,omitempty"`
	SpawnedBy       *string `json:"spawnedBy,omitempty"`
	SpawnDepth      *int    `json:"spawnDepth,omitempty"`
	SendPolicy      *string `json:"sendPolicy,omitempty"`
	GroupActivation *string `json:"groupActivation,omitempty"`
}

// SessionsResetParams are the params for "sessions.reset".
type SessionsResetParams struct {
	Key    string `json:"key"`
	Reason string `json:"reason,omitempty"` // "new" or "reset"
}

// SessionsDeleteParams are the params for "sessions.delete".
type SessionsDeleteParams struct {
	Key              string `json:"key"`
	DeleteTranscript *bool  `json:"deleteTranscript,omitempty"`
}

// SessionsCompactParams are the params for "sessions.compact".
type SessionsCompactParams struct {
	Key      string `json:"key"`
	MaxLines *int   `json:"maxLines,omitempty"`
}

// SessionsGetParams are the params for "sessions.get".
type SessionsGetParams struct {
	Key        string `json:"key,omitempty"`
	SessionKey string `json:"sessionKey,omitempty"`
	Limit      *int   `json:"limit,omitempty"` // max messages to return (default 200, server-clamped)
}

// SessionsUsageParams are the params for "sessions.usage".
type SessionsUsageParams struct {
	Key                  string `json:"key,omitempty"`
	StartDate            string `json:"startDate,omitempty"`
	EndDate              string `json:"endDate,omitempty"`
	Limit                *int   `json:"limit,omitempty"`
	IncludeContextWeight *bool  `json:"includeContextWeight,omitempty"`
}

// ---------------------------------------------------------------------------
// Node types
// ---------------------------------------------------------------------------

// NodePairRequestParams are the params for "node.pair.request".
type NodePairRequestParams struct {
	NodeID          string   `json:"nodeId"`
	DisplayName     string   `json:"displayName,omitempty"`
	Platform        string   `json:"platform,omitempty"`
	Version         string   `json:"version,omitempty"`
	CoreVersion     string   `json:"coreVersion,omitempty"`
	UIVersion       string   `json:"uiVersion,omitempty"`
	DeviceFamily    string   `json:"deviceFamily,omitempty"`
	ModelIdentifier string   `json:"modelIdentifier,omitempty"`
	Caps            []string `json:"caps,omitempty"`
	Commands        []string `json:"commands,omitempty"`
	RemoteIP        string   `json:"remoteIp,omitempty"`
	Silent          *bool    `json:"silent,omitempty"`
}

// NodePairApproveParams are the params for "node.pair.approve".
type NodePairApproveParams struct {
	RequestID string `json:"requestId"`
}

// NodePairRejectParams are the params for "node.pair.reject".
type NodePairRejectParams struct {
	RequestID string `json:"requestId"`
}

// NodePairVerifyParams are the params for "node.pair.verify".
type NodePairVerifyParams struct {
	NodeID string `json:"nodeId"`
	Token  string `json:"token"`
}

// NodeRenameParams are the params for "node.rename".
type NodeRenameParams struct {
	NodeID      string `json:"nodeId"`
	DisplayName string `json:"displayName"`
}

// NodeDescribeParams are the params for "node.describe".
type NodeDescribeParams struct {
	NodeID string `json:"nodeId"`
}

// NodeInvokeParams are the params for "node.invoke".
type NodeInvokeParams struct {
	NodeID         string          `json:"nodeId"`
	Command        string          `json:"command"`
	Params         json.RawMessage `json:"params,omitempty"`
	TimeoutMs      *int            `json:"timeoutMs,omitempty"`
	IdempotencyKey string          `json:"idempotencyKey"`
}

// NodeInvokeResultParams are the params for "node.invoke.result".
type NodeInvokeResultParams struct {
	ID          string                 `json:"id"`
	NodeID      string                 `json:"nodeId"`
	OK          bool                   `json:"ok"`
	Payload     json.RawMessage        `json:"payload,omitempty"`
	PayloadJSON string                 `json:"payloadJSON,omitempty"`
	Error       *NodeInvokeResultError `json:"error,omitempty"`
}

// NodeInvokeResultError is the error structure in node invoke results.
type NodeInvokeResultError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// NodeEventParams are the params for "node.event".
type NodeEventParams struct {
	Event       string          `json:"event"`
	Payload     json.RawMessage `json:"payload,omitempty"`
	PayloadJSON string          `json:"payloadJSON,omitempty"`
}

// NodeInvokeRequestEvent is the payload of a "node.invoke.request" event.
type NodeInvokeRequestEvent struct {
	ID             string `json:"id"`
	NodeID         string `json:"nodeId"`
	Command        string `json:"command"`
	ParamsJSON     string `json:"paramsJSON,omitempty"`
	TimeoutMs      *int   `json:"timeoutMs,omitempty"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
}

// ---------------------------------------------------------------------------
// Device pairing types
// ---------------------------------------------------------------------------

// DevicePairApproveParams are the params for "device.pair.approve".
type DevicePairApproveParams struct {
	RequestID string `json:"requestId"`
}

// DevicePairRejectParams are the params for "device.pair.reject".
type DevicePairRejectParams struct {
	RequestID string `json:"requestId"`
}

// DevicePairRemoveParams are the params for "device.pair.remove".
type DevicePairRemoveParams struct {
	DeviceID string `json:"deviceId"`
}

// DeviceTokenRotateParams are the params for "device.token.rotate".
type DeviceTokenRotateParams struct {
	DeviceID string   `json:"deviceId"`
	Role     string   `json:"role"`
	Scopes   []string `json:"scopes,omitempty"`
}

// DeviceTokenRevokeParams are the params for "device.token.revoke".
type DeviceTokenRevokeParams struct {
	DeviceID string `json:"deviceId"`
	Role     string `json:"role"`
}

// DevicePairRequestedEvent is the payload of a "device.pair.requested" event.
type DevicePairRequestedEvent struct {
	RequestID   string   `json:"requestId"`
	DeviceID    string   `json:"deviceId"`
	PublicKey   string   `json:"publicKey"`
	DisplayName string   `json:"displayName,omitempty"`
	Platform    string   `json:"platform,omitempty"`
	ClientID    string   `json:"clientId,omitempty"`
	ClientMode  string   `json:"clientMode,omitempty"`
	Role        string   `json:"role,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	RemoteIP    string   `json:"remoteIp,omitempty"`
	Silent      *bool    `json:"silent,omitempty"`
	IsRepair    *bool    `json:"isRepair,omitempty"`
	Ts          int64    `json:"ts"`
}

// DevicePairResolvedEvent is the payload of a "device.pair.resolved" event.
type DevicePairResolvedEvent struct {
	RequestID string `json:"requestId"`
	DeviceID  string `json:"deviceId"`
	Decision  string `json:"decision"`
	Ts        int64  `json:"ts"`
}

// ---------------------------------------------------------------------------
// Config types
// ---------------------------------------------------------------------------

// ConfigGetParams are the params for "config.get" (empty).
type ConfigGetParams struct{}

// ConfigSetParams are the params for "config.set".
type ConfigSetParams struct {
	Raw      string `json:"raw"`
	BaseHash string `json:"baseHash,omitempty"`
}

// ConfigApplyParams are the params for "config.apply".
type ConfigApplyParams struct {
	Raw            string `json:"raw"`
	BaseHash       string `json:"baseHash,omitempty"`
	SessionKey     string `json:"sessionKey,omitempty"`
	Note           string `json:"note,omitempty"`
	RestartDelayMs *int   `json:"restartDelayMs,omitempty"`
}

// ConfigPatchParams are the params for "config.patch".
type ConfigPatchParams struct {
	Raw            string `json:"raw"`
	BaseHash       string `json:"baseHash,omitempty"`
	SessionKey     string `json:"sessionKey,omitempty"`
	Note           string `json:"note,omitempty"`
	RestartDelayMs *int   `json:"restartDelayMs,omitempty"`
}

// ConfigSchemaResponse is the response for "config.schema".
type ConfigSchemaResponse struct {
	Schema      json.RawMessage         `json:"schema"`
	UIHints     map[string]ConfigUIHint `json:"uiHints"`
	Version     string                  `json:"version"`
	GeneratedAt string                  `json:"generatedAt"`
}

// ConfigUIHint is a UI hint for a config field.
type ConfigUIHint struct {
	Label        string          `json:"label,omitempty"`
	Help         string          `json:"help,omitempty"`
	Group        string          `json:"group,omitempty"`
	Order        *int            `json:"order,omitempty"`
	Advanced     *bool           `json:"advanced,omitempty"`
	Sensitive    *bool           `json:"sensitive,omitempty"`
	Placeholder  string          `json:"placeholder,omitempty"`
	ItemTemplate json.RawMessage `json:"itemTemplate,omitempty"`
}

// ---------------------------------------------------------------------------
// Agents CRUD types
// ---------------------------------------------------------------------------

// AgentIdentity is the visual identity of an agent.
type AgentIdentity struct {
	Name      string `json:"name,omitempty"`
	Theme     string `json:"theme,omitempty"`
	Emoji     string `json:"emoji,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// AgentSummary is a summary of an agent.
type AgentSummary struct {
	ID       string         `json:"id"`
	Name     string         `json:"name,omitempty"`
	Identity *AgentIdentity `json:"identity,omitempty"`
}

// AgentsListResult is the result of "agents.list".
type AgentsListResult struct {
	DefaultID string         `json:"defaultId"`
	MainKey   string         `json:"mainKey"`
	Scope     string         `json:"scope"` // "per-sender" or "global"
	Agents    []AgentSummary `json:"agents"`
}

// AgentsCreateParams are the params for "agents.create".
type AgentsCreateParams struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
	Emoji     string `json:"emoji,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}

// AgentsCreateResult is the result of "agents.create".
type AgentsCreateResult struct {
	OK        bool   `json:"ok"`
	AgentID   string `json:"agentId"`
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
}

// AgentsUpdateParams are the params for "agents.update".
type AgentsUpdateParams struct {
	AgentID   string `json:"agentId"`
	Name      string `json:"name,omitempty"`
	Workspace string `json:"workspace,omitempty"`
	Model     string `json:"model,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}

// AgentsUpdateResult is the result of "agents.update".
type AgentsUpdateResult struct {
	OK      bool   `json:"ok"`
	AgentID string `json:"agentId"`
}

// AgentsDeleteParams are the params for "agents.delete".
type AgentsDeleteParams struct {
	AgentID     string `json:"agentId"`
	DeleteFiles *bool  `json:"deleteFiles,omitempty"`
}

// AgentsDeleteResult is the result of "agents.delete".
type AgentsDeleteResult struct {
	OK              bool   `json:"ok"`
	AgentID         string `json:"agentId"`
	RemovedBindings int    `json:"removedBindings"`
}

// AgentsFileEntry describes a file associated with an agent.
type AgentsFileEntry struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Missing     bool   `json:"missing"`
	Size        *int   `json:"size,omitempty"`
	UpdatedAtMs *int64 `json:"updatedAtMs,omitempty"`
	Content     string `json:"content,omitempty"`
}

// AgentsFilesListParams are the params for "agents.files.list".
type AgentsFilesListParams struct {
	AgentID string `json:"agentId"`
}

// AgentsFilesListResult is the result of "agents.files.list".
type AgentsFilesListResult struct {
	AgentID   string            `json:"agentId"`
	Workspace string            `json:"workspace"`
	Files     []AgentsFileEntry `json:"files"`
}

// AgentsFilesGetParams are the params for "agents.files.get".
type AgentsFilesGetParams struct {
	AgentID string `json:"agentId"`
	Name    string `json:"name"`
}

// AgentsFilesGetResult is the result of "agents.files.get".
type AgentsFilesGetResult struct {
	AgentID   string          `json:"agentId"`
	Workspace string          `json:"workspace"`
	File      AgentsFileEntry `json:"file"`
}

// AgentsFilesSetParams are the params for "agents.files.set".
type AgentsFilesSetParams struct {
	AgentID string `json:"agentId"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// AgentsFilesSetResult is the result of "agents.files.set".
type AgentsFilesSetResult struct {
	OK        bool            `json:"ok"`
	AgentID   string          `json:"agentId"`
	Workspace string          `json:"workspace"`
	File      AgentsFileEntry `json:"file"`
}

// ---------------------------------------------------------------------------
// Models types
// ---------------------------------------------------------------------------

// ModelChoice is a single model option.
type ModelChoice struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Provider      string `json:"provider"`
	ContextWindow *int   `json:"contextWindow,omitempty"`
	Reasoning     *bool  `json:"reasoning,omitempty"`
}

// ModelsListResult is the result of "models.list".
type ModelsListResult struct {
	Models []ModelChoice `json:"models"`
}

// ---------------------------------------------------------------------------
// Logs types
// ---------------------------------------------------------------------------

// LogsTailParams are the params for "logs.tail".
type LogsTailParams struct {
	Cursor   *int `json:"cursor,omitempty"`
	Limit    *int `json:"limit,omitempty"`
	MaxBytes *int `json:"maxBytes,omitempty"`
}

// LogsTailResult is the result of "logs.tail".
type LogsTailResult struct {
	File      string   `json:"file"`
	Cursor    int      `json:"cursor"`
	Size      int      `json:"size"`
	Lines     []string `json:"lines"`
	Truncated *bool    `json:"truncated,omitempty"`
	Reset     *bool    `json:"reset,omitempty"`
}

// ---------------------------------------------------------------------------
// Cron types
// ---------------------------------------------------------------------------

// CronSchedule is a schedule definition (union: at, every, cron).
type CronSchedule struct {
	Kind      string `json:"kind"` // "at", "every", "cron"
	At        string `json:"at,omitempty"`
	EveryMs   *int   `json:"everyMs,omitempty"`
	AnchorMs  *int64 `json:"anchorMs,omitempty"`
	Expr      string `json:"expr,omitempty"`
	Tz        string `json:"tz,omitempty"`
	StaggerMs *int   `json:"staggerMs,omitempty"`
}

// CronPayload is the payload for a cron job (union: systemEvent, agentTurn).
type CronPayload struct {
	Kind                       string `json:"kind"` // "systemEvent" or "agentTurn"
	Text                       string `json:"text,omitempty"`
	Message                    string `json:"message,omitempty"`
	Model                      string `json:"model,omitempty"`
	Thinking                   string `json:"thinking,omitempty"`
	TimeoutSeconds             *int   `json:"timeoutSeconds,omitempty"`
	AllowUnsafeExternalContent *bool  `json:"allowUnsafeExternalContent,omitempty"`
	Deliver                    *bool  `json:"deliver,omitempty"`
	Channel                    string `json:"channel,omitempty"`
	To                         string `json:"to,omitempty"`
	BestEffortDeliver          *bool  `json:"bestEffortDeliver,omitempty"`
}

// CronDelivery is the delivery configuration for a cron job.
type CronDelivery struct {
	Mode       string `json:"mode"` // "none", "announce", "webhook"
	Channel    string `json:"channel,omitempty"`
	BestEffort *bool  `json:"bestEffort,omitempty"`
	To         string `json:"to,omitempty"`
}

// CronJobState is the runtime state of a cron job.
type CronJobState struct {
	NextRunAtMs       *int64 `json:"nextRunAtMs,omitempty"`
	RunningAtMs       *int64 `json:"runningAtMs,omitempty"`
	LastRunAtMs       *int64 `json:"lastRunAtMs,omitempty"`
	LastStatus        string `json:"lastStatus,omitempty"` // "ok", "error", "skipped"
	LastError         string `json:"lastError,omitempty"`
	LastDurationMs    *int64 `json:"lastDurationMs,omitempty"`
	ConsecutiveErrors *int   `json:"consecutiveErrors,omitempty"`
}

// CronJob is a full cron job definition.
type CronJob struct {
	ID             string        `json:"id"`
	AgentID        string        `json:"agentId,omitempty"`
	SessionKey     string        `json:"sessionKey,omitempty"`
	Name           string        `json:"name"`
	Description    string        `json:"description,omitempty"`
	Enabled        bool          `json:"enabled"`
	DeleteAfterRun *bool         `json:"deleteAfterRun,omitempty"`
	CreatedAtMs    int64         `json:"createdAtMs"`
	UpdatedAtMs    int64         `json:"updatedAtMs"`
	Schedule       CronSchedule  `json:"schedule"`
	SessionTarget  string        `json:"sessionTarget"` // "main" or "isolated"
	WakeMode       string        `json:"wakeMode"`      // "next-heartbeat" or "now"
	Payload        CronPayload   `json:"payload"`
	Delivery       *CronDelivery `json:"delivery,omitempty"`
	State          CronJobState  `json:"state"`
}

// CronListParams are the params for "cron.list".
type CronListParams struct {
	IncludeDisabled *bool `json:"includeDisabled,omitempty"`
}

// CronAddParams are the params for "cron.add".
type CronAddParams struct {
	Name           string        `json:"name"`
	AgentID        *string       `json:"agentId,omitempty"`
	SessionKey     *string       `json:"sessionKey,omitempty"`
	Description    string        `json:"description,omitempty"`
	Enabled        *bool         `json:"enabled,omitempty"`
	DeleteAfterRun *bool         `json:"deleteAfterRun,omitempty"`
	Schedule       CronSchedule  `json:"schedule"`
	SessionTarget  string        `json:"sessionTarget"` // "main" or "isolated"
	WakeMode       string        `json:"wakeMode"`      // "next-heartbeat" or "now"
	Payload        CronPayload   `json:"payload"`
	Delivery       *CronDelivery `json:"delivery,omitempty"`
}

// CronJobPatch is a partial update for a cron job.
type CronJobPatch struct {
	Name           string        `json:"name,omitempty"`
	AgentID        *string       `json:"agentId,omitempty"`
	SessionKey     *string       `json:"sessionKey,omitempty"`
	Description    string        `json:"description,omitempty"`
	Enabled        *bool         `json:"enabled,omitempty"`
	DeleteAfterRun *bool         `json:"deleteAfterRun,omitempty"`
	Schedule       *CronSchedule `json:"schedule,omitempty"`
	SessionTarget  string        `json:"sessionTarget,omitempty"`
	WakeMode       string        `json:"wakeMode,omitempty"`
	Payload        *CronPayload  `json:"payload,omitempty"`
	Delivery       *CronDelivery `json:"delivery,omitempty"`
	State          *CronJobState `json:"state,omitempty"`
}

// CronUpdateParams are the params for "cron.update".
type CronUpdateParams struct {
	ID    string       `json:"id,omitempty"`
	JobID string       `json:"jobId,omitempty"`
	Patch CronJobPatch `json:"patch"`
}

// CronRemoveParams are the params for "cron.remove".
type CronRemoveParams struct {
	ID    string `json:"id,omitempty"`
	JobID string `json:"jobId,omitempty"`
}

// CronRunParams are the params for "cron.run".
type CronRunParams struct {
	ID    string `json:"id,omitempty"`
	JobID string `json:"jobId,omitempty"`
	Mode  string `json:"mode,omitempty"` // "due" or "force"
}

// CronRunsParams are the params for "cron.runs".
type CronRunsParams struct {
	ID    string `json:"id,omitempty"`
	JobID string `json:"jobId,omitempty"`
	Limit *int   `json:"limit,omitempty"`
}

// CronRunLogEntry is a single entry in the cron run log.
type CronRunLogEntry struct {
	Ts          int64  `json:"ts"`
	JobID       string `json:"jobId"`
	Action      string `json:"action"` // "finished"
	Status      string `json:"status,omitempty"`
	Error       string `json:"error,omitempty"`
	Summary     string `json:"summary,omitempty"`
	SessionID   string `json:"sessionId,omitempty"`
	SessionKey  string `json:"sessionKey,omitempty"`
	RunAtMs     *int64 `json:"runAtMs,omitempty"`
	DurationMs  *int64 `json:"durationMs,omitempty"`
	NextRunAtMs *int64 `json:"nextRunAtMs,omitempty"`
}

// ---------------------------------------------------------------------------
// Channels / Talk types
// ---------------------------------------------------------------------------

// TalkModeParams are the params for "talk.mode".
type TalkModeParams struct {
	Enabled bool   `json:"enabled"`
	Phase   string `json:"phase,omitempty"`
}

// TalkConfigParams are the params for "talk.config".
type TalkConfigParams struct {
	IncludeSecrets *bool `json:"includeSecrets,omitempty"`
}

// TalkConfigResult is the result of "talk.config".
type TalkConfigResult struct {
	Config TalkConfigData `json:"config"`
}

// TalkConfigData holds the talk configuration sections.
type TalkConfigData struct {
	Talk    *TalkConfigTalk    `json:"talk,omitempty"`
	Session *TalkConfigSession `json:"session,omitempty"`
	UI      *TalkConfigUI      `json:"ui,omitempty"`
}

// TalkConfigTalk is the talk section of talk config.
type TalkConfigTalk struct {
	VoiceID           string            `json:"voiceId,omitempty"`
	VoiceAliases      map[string]string `json:"voiceAliases,omitempty"`
	ModelID           string            `json:"modelId,omitempty"`
	OutputFormat      string            `json:"outputFormat,omitempty"`
	APIKey            string            `json:"apiKey,omitempty"`
	InterruptOnSpeech *bool             `json:"interruptOnSpeech,omitempty"`
}

// TalkConfigSession is the session section of talk config.
type TalkConfigSession struct {
	MainKey string `json:"mainKey,omitempty"`
}

// TalkConfigUI is the UI section of talk config.
type TalkConfigUI struct {
	SeamColor string `json:"seamColor,omitempty"`
}

// ChannelsStatusParams are the params for "channels.status".
type ChannelsStatusParams struct {
	Probe     *bool `json:"probe,omitempty"`
	TimeoutMs *int  `json:"timeoutMs,omitempty"`
}

// ChannelsStatusResult is the result of "channels.status".
type ChannelsStatusResult struct {
	Ts                      int64                               `json:"ts"`
	ChannelOrder            []string                            `json:"channelOrder"`
	ChannelLabels           map[string]string                   `json:"channelLabels"`
	ChannelDetailLabels     map[string]string                   `json:"channelDetailLabels,omitempty"`
	ChannelSystemImages     map[string]string                   `json:"channelSystemImages,omitempty"`
	ChannelMeta             []ChannelUIMeta                     `json:"channelMeta,omitempty"`
	Channels                map[string]json.RawMessage          `json:"channels"`
	ChannelAccounts         map[string][]ChannelAccountSnapshot `json:"channelAccounts"`
	ChannelDefaultAccountID map[string]string                   `json:"channelDefaultAccountId"`
}

// ChannelUIMeta is UI metadata for a channel.
type ChannelUIMeta struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	DetailLabel string `json:"detailLabel"`
	SystemImage string `json:"systemImage,omitempty"`
}

// ChannelAccountSnapshot is a snapshot of a channel account.
type ChannelAccountSnapshot struct {
	AccountID              string          `json:"accountId"`
	Name                   string          `json:"name,omitempty"`
	Enabled                *bool           `json:"enabled,omitempty"`
	Configured             *bool           `json:"configured,omitempty"`
	Linked                 *bool           `json:"linked,omitempty"`
	Running                *bool           `json:"running,omitempty"`
	Connected              *bool           `json:"connected,omitempty"`
	ReconnectAttempts      *int            `json:"reconnectAttempts,omitempty"`
	LastConnectedAt        *int64          `json:"lastConnectedAt,omitempty"`
	LastError              string          `json:"lastError,omitempty"`
	LastStartAt            *int64          `json:"lastStartAt,omitempty"`
	LastStopAt             *int64          `json:"lastStopAt,omitempty"`
	LastInboundAt          *int64          `json:"lastInboundAt,omitempty"`
	LastOutboundAt         *int64          `json:"lastOutboundAt,omitempty"`
	LastProbeAt            *int64          `json:"lastProbeAt,omitempty"`
	Mode                   string          `json:"mode,omitempty"`
	DMPolicy               string          `json:"dmPolicy,omitempty"`
	AllowFrom              []string        `json:"allowFrom,omitempty"`
	TokenSource            string          `json:"tokenSource,omitempty"`
	BotTokenSource         string          `json:"botTokenSource,omitempty"`
	AppTokenSource         string          `json:"appTokenSource,omitempty"`
	BaseURL                string          `json:"baseUrl,omitempty"`
	AllowUnmentionedGroups *bool           `json:"allowUnmentionedGroups,omitempty"`
	CLIPath                *string         `json:"cliPath,omitempty"`
	DBPath                 *string         `json:"dbPath,omitempty"`
	Port                   *int            `json:"port,omitempty"`
	Probe                  json.RawMessage `json:"probe,omitempty"`
	Audit                  json.RawMessage `json:"audit,omitempty"`
	Application            json.RawMessage `json:"application,omitempty"`
}

// ChannelsLogoutParams are the params for "channels.logout".
type ChannelsLogoutParams struct {
	Channel   string `json:"channel"`
	AccountID string `json:"accountId,omitempty"`
}

// ---------------------------------------------------------------------------
// Skills types
// ---------------------------------------------------------------------------

// SkillsStatusParams are the params for "skills.status".
type SkillsStatusParams struct {
	AgentID string `json:"agentId,omitempty"`
}

// SkillsBinsResult is the result of "skills.bins".
type SkillsBinsResult struct {
	Bins []string `json:"bins"`
}

// SkillsInstallParams are the params for "skills.install".
type SkillsInstallParams struct {
	Name      string `json:"name"`
	InstallID string `json:"installId"`
	TimeoutMs *int   `json:"timeoutMs,omitempty"`
}

// SkillsUpdateParams are the params for "skills.update".
type SkillsUpdateParams struct {
	SkillKey string            `json:"skillKey"`
	Enabled  *bool             `json:"enabled,omitempty"`
	APIKey   string            `json:"apiKey,omitempty"`
	Env      map[string]string `json:"env,omitempty"`
}

// ---------------------------------------------------------------------------
// Wizard types
// ---------------------------------------------------------------------------

// WizardStartParams are the params for "wizard.start".
type WizardStartParams struct {
	Mode      string `json:"mode,omitempty"` // "local" or "remote"
	Workspace string `json:"workspace,omitempty"`
}

// WizardStartResult is the result of "wizard.start".
type WizardStartResult struct {
	SessionID string      `json:"sessionId"`
	Done      bool        `json:"done"`
	Step      *WizardStep `json:"step,omitempty"`
	Status    string      `json:"status,omitempty"` // "running", "done", "cancelled", "error"
	Error     string      `json:"error,omitempty"`
}

// WizardAnswer is the answer to a wizard step.
type WizardAnswer struct {
	StepID string          `json:"stepId"`
	Value  json.RawMessage `json:"value,omitempty"`
}

// WizardNextParams are the params for "wizard.next".
type WizardNextParams struct {
	SessionID string        `json:"sessionId"`
	Answer    *WizardAnswer `json:"answer,omitempty"`
}

// WizardNextResult is the result of "wizard.next".
type WizardNextResult struct {
	Done   bool        `json:"done"`
	Step   *WizardStep `json:"step,omitempty"`
	Status string      `json:"status,omitempty"` // "running", "done", "cancelled", "error"
	Error  string      `json:"error,omitempty"`
}

// WizardCancelParams are the params for "wizard.cancel".
type WizardCancelParams struct {
	SessionID string `json:"sessionId"`
}

// WizardStatusParams are the params for "wizard.status".
type WizardStatusParams struct {
	SessionID string `json:"sessionId"`
}

// WizardStatusResult is the result of "wizard.status".
type WizardStatusResult struct {
	Status string `json:"status"` // "running", "done", "cancelled", "error"
	Error  string `json:"error,omitempty"`
}

// WizardStep describes a single wizard step.
type WizardStep struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"` // "note", "select", "text", "confirm", "multiselect", "progress", "action"
	Title        string             `json:"title,omitempty"`
	Message      string             `json:"message,omitempty"`
	Options      []WizardStepOption `json:"options,omitempty"`
	InitialValue json.RawMessage    `json:"initialValue,omitempty"`
	Placeholder  string             `json:"placeholder,omitempty"`
	Sensitive    *bool              `json:"sensitive,omitempty"`
	Executor     string             `json:"executor,omitempty"` // "gateway" or "client"
}

// WizardStepOption is an option for a wizard step.
type WizardStepOption struct {
	Value json.RawMessage `json:"value"`
	Label string          `json:"label"`
	Hint  string          `json:"hint,omitempty"`
}

// ---------------------------------------------------------------------------
// Push types
// ---------------------------------------------------------------------------

// PushTestParams are the params for "push.test".
type PushTestParams struct {
	NodeID      string `json:"nodeId"`
	Title       string `json:"title,omitempty"`
	Body        string `json:"body,omitempty"`
	Environment string `json:"environment,omitempty"` // "sandbox" or "production"
}

// PushTestResult is the result of "push.test".
type PushTestResult struct {
	OK          bool   `json:"ok"`
	Status      int    `json:"status"`
	ApnsID      string `json:"apnsId,omitempty"`
	Reason      string `json:"reason,omitempty"`
	TokenSuffix string `json:"tokenSuffix"`
	Topic       string `json:"topic"`
	Environment string `json:"environment"` // "sandbox" or "production"
}

// ---------------------------------------------------------------------------
// Send / Poll / Wake types
// ---------------------------------------------------------------------------

// SendParams are the params for "send".
type SendParams struct {
	To             string   `json:"to"`
	Message        string   `json:"message,omitempty"`
	MediaURL       string   `json:"mediaUrl,omitempty"`
	MediaURLs      []string `json:"mediaUrls,omitempty"`
	GifPlayback    *bool    `json:"gifPlayback,omitempty"`
	Channel        string   `json:"channel,omitempty"`
	AccountID      string   `json:"accountId,omitempty"`
	ThreadID       string   `json:"threadId,omitempty"`
	SessionKey     string   `json:"sessionKey,omitempty"`
	IdempotencyKey string   `json:"idempotencyKey"`
}

// PollParams are the params for "poll".
type PollParams struct {
	To              string   `json:"to"`
	Question        string   `json:"question"`
	Options         []string `json:"options"`
	MaxSelections   *int     `json:"maxSelections,omitempty"`
	DurationSeconds *int     `json:"durationSeconds,omitempty"`
	DurationHours   *int     `json:"durationHours,omitempty"`
	Silent          *bool    `json:"silent,omitempty"`
	IsAnonymous     *bool    `json:"isAnonymous,omitempty"`
	ThreadID        string   `json:"threadId,omitempty"`
	Channel         string   `json:"channel,omitempty"`
	AccountID       string   `json:"accountId,omitempty"`
	IdempotencyKey  string   `json:"idempotencyKey"`
}

// WakeParams are the params for "wake".
type WakeParams struct {
	Mode string `json:"mode"` // "now" or "next-heartbeat"
	Text string `json:"text"`
}

// ---------------------------------------------------------------------------
// Update / Misc event types
// ---------------------------------------------------------------------------

// UpdateRunParams are the params for "update.run".
type UpdateRunParams struct {
	SessionKey     string `json:"sessionKey,omitempty"`
	Note           string `json:"note,omitempty"`
	RestartDelayMs *int   `json:"restartDelayMs,omitempty"`
	TimeoutMs      *int   `json:"timeoutMs,omitempty"`
}

// TickEvent is the payload of a "tick" event.
type TickEvent struct {
	Ts int64 `json:"ts"`
}

// ShutdownEvent is the payload of a "shutdown" event.
type ShutdownEvent struct {
	Reason            string `json:"reason"`
	RestartExpectedMs *int64 `json:"restartExpectedMs,omitempty"`
}

// ---------------------------------------------------------------------------
// Presence event
// ---------------------------------------------------------------------------

// PresenceEvent is the payload of a "presence" event.
type PresenceEvent struct {
	Presence []SystemPresence `json:"presence"`
}

// SystemPresence describes a connected system's presence.
type SystemPresence struct {
	Text             string   `json:"text"`
	Ts               int64    `json:"ts"`
	Host             string   `json:"host,omitempty"`
	IP               string   `json:"ip,omitempty"`
	Version          string   `json:"version,omitempty"`
	Platform         string   `json:"platform,omitempty"`
	DeviceFamily     string   `json:"deviceFamily,omitempty"`
	ModelIdentifier  string   `json:"modelIdentifier,omitempty"`
	LastInputSeconds *float64 `json:"lastInputSeconds,omitempty"`
	Mode             string   `json:"mode,omitempty"`
	Reason           string   `json:"reason,omitempty"`
	DeviceID         string   `json:"deviceId,omitempty"`
	Roles            []string `json:"roles,omitempty"`
	Scopes           []string `json:"scopes,omitempty"`
	InstanceID       string   `json:"instanceId,omitempty"`
}

// ---------------------------------------------------------------------------
// Health event
// ---------------------------------------------------------------------------

// HealthEvent is the payload of a "health" event.
type HealthEvent struct {
	OK               bool                            `json:"ok"`
	Ts               int64                           `json:"ts"`
	DurationMs       int64                           `json:"durationMs"`
	Channels         map[string]ChannelHealthSummary `json:"channels"`
	ChannelOrder     []string                        `json:"channelOrder"`
	ChannelLabels    map[string]string               `json:"channelLabels"`
	HeartbeatSeconds int                             `json:"heartbeatSeconds"`
	DefaultAgentID   string                          `json:"defaultAgentId"`
	Agents           []AgentHealthSummary            `json:"agents"`
	Sessions         HealthSessionsSummary           `json:"sessions"`
}

// ChannelHealthSummary describes a channel's health status.
type ChannelHealthSummary struct {
	AccountID   string          `json:"accountId,omitempty"`
	Configured  *bool           `json:"configured,omitempty"`
	Linked      *bool           `json:"linked,omitempty"`
	AuthAgeMs   *int64          `json:"authAgeMs,omitempty"`
	Probe       json.RawMessage `json:"probe,omitempty"`
	LastProbeAt *int64          `json:"lastProbeAt,omitempty"`
	Accounts    json.RawMessage `json:"accounts,omitempty"`
}

// AgentHealthSummary describes an agent's health summary.
type AgentHealthSummary struct {
	AgentID   string                `json:"agentId"`
	Name      string                `json:"name,omitempty"`
	IsDefault bool                  `json:"isDefault"`
	Heartbeat json.RawMessage       `json:"heartbeat"`
	Sessions  HealthSessionsSummary `json:"sessions"`
}

// HealthSessionsSummary describes session information in a health summary.
type HealthSessionsSummary struct {
	Path   string                `json:"path"`
	Count  int                   `json:"count"`
	Recent []HealthRecentSession `json:"recent"`
}

// HealthRecentSession describes a recent session in the health summary.
type HealthRecentSession struct {
	Key       string `json:"key"`
	UpdatedAt *int64 `json:"updatedAt,omitempty"`
	Age       *int64 `json:"age,omitempty"`
}

// ---------------------------------------------------------------------------
// Heartbeat event
// ---------------------------------------------------------------------------

// HeartbeatEvent is the payload of a "heartbeat" event.
type HeartbeatEvent struct {
	Ts            int64  `json:"ts"`
	Status        string `json:"status"` // "sent", "ok-empty", "ok-token", "skipped", "failed"
	To            string `json:"to,omitempty"`
	AccountID     string `json:"accountId,omitempty"`
	Preview       string `json:"preview,omitempty"`
	DurationMs    *int64 `json:"durationMs,omitempty"`
	HasMedia      *bool  `json:"hasMedia,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Channel       string `json:"channel,omitempty"`
	Silent        *bool  `json:"silent,omitempty"`
	IndicatorType string `json:"indicatorType,omitempty"` // "ok", "alert", "error"
}

// ---------------------------------------------------------------------------
// Voicewake changed event
// ---------------------------------------------------------------------------

// VoicewakeChangedEvent is the payload of a "voicewake.changed" event.
type VoicewakeChangedEvent struct {
	Triggers []string `json:"triggers"`
}

// ---------------------------------------------------------------------------
// Cron event
// ---------------------------------------------------------------------------

// CronEvent is the payload of a "cron" event.
type CronEvent struct {
	JobID       string            `json:"jobId"`
	Action      string            `json:"action"` // "added", "updated", "removed", "started", "finished"
	RunAtMs     *int64            `json:"runAtMs,omitempty"`
	DurationMs  *int64            `json:"durationMs,omitempty"`
	Status      string            `json:"status,omitempty"` // "ok", "error", "skipped"
	Error       string            `json:"error,omitempty"`
	Summary     string            `json:"summary,omitempty"`
	SessionID   string            `json:"sessionId,omitempty"`
	SessionKey  string            `json:"sessionKey,omitempty"`
	NextRunAtMs *int64            `json:"nextRunAtMs,omitempty"`
	Model       string            `json:"model,omitempty"`
	Provider    string            `json:"provider,omitempty"`
	Usage       *CronUsageSummary `json:"usage,omitempty"`
}

// CronUsageSummary describes token usage for a cron run.
type CronUsageSummary struct {
	InputTokens      *int `json:"input_tokens,omitempty"`
	OutputTokens     *int `json:"output_tokens,omitempty"`
	TotalTokens      *int `json:"total_tokens,omitempty"`
	CacheReadTokens  *int `json:"cache_read_tokens,omitempty"`
	CacheWriteTokens *int `json:"cache_write_tokens,omitempty"`
}

// ---------------------------------------------------------------------------
// Node pair events
// ---------------------------------------------------------------------------

// NodePairRequestedEvent is the payload of a "node.pair.requested" event.
type NodePairRequestedEvent struct {
	RequestID       string   `json:"requestId"`
	NodeID          string   `json:"nodeId"`
	DisplayName     string   `json:"displayName,omitempty"`
	Platform        string   `json:"platform,omitempty"`
	Version         string   `json:"version,omitempty"`
	CoreVersion     string   `json:"coreVersion,omitempty"`
	UIVersion       string   `json:"uiVersion,omitempty"`
	DeviceFamily    string   `json:"deviceFamily,omitempty"`
	ModelIdentifier string   `json:"modelIdentifier,omitempty"`
	Caps            []string `json:"caps,omitempty"`
	Commands        []string `json:"commands,omitempty"`
	RemoteIP        string   `json:"remoteIp,omitempty"`
	Silent          *bool    `json:"silent,omitempty"`
	Ts              int64    `json:"ts"`
}

// NodePairResolvedEvent is the payload of a "node.pair.resolved" event.
type NodePairResolvedEvent struct {
	RequestID string `json:"requestId"`
	NodeID    string `json:"nodeId"`
	Decision  string `json:"decision"`
	Ts        int64  `json:"ts"`
}

// ---------------------------------------------------------------------------
// Web login types
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// TTS types
// ---------------------------------------------------------------------------

// TTSStatusResult is the result of "tts.status".
type TTSStatusResult struct {
	Enabled           bool     `json:"enabled"`
	Auto              string   `json:"auto"`
	Provider          string   `json:"provider"`
	FallbackProvider  *string  `json:"fallbackProvider"`
	FallbackProviders []string `json:"fallbackProviders"`
	PrefsPath         string   `json:"prefsPath"`
	HasOpenAIKey      bool     `json:"hasOpenAIKey"`
	HasElevenLabsKey  bool     `json:"hasElevenLabsKey"`
	EdgeEnabled       bool     `json:"edgeEnabled"`
}

// TTSProviderInfo describes a TTS provider.
type TTSProviderInfo struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Configured bool     `json:"configured"`
	Models     []string `json:"models"`
	Voices     []string `json:"voices,omitempty"`
}

// TTSProvidersResult is the result of "tts.providers".
type TTSProvidersResult struct {
	Providers []TTSProviderInfo `json:"providers"`
	Active    string            `json:"active"`
}

// TTSEnableResult is the result of "tts.enable".
type TTSEnableResult struct {
	Enabled bool `json:"enabled"`
}

// TTSDisableResult is the result of "tts.disable".
type TTSDisableResult struct {
	Enabled bool `json:"enabled"`
}

// TTSConvertParams are the params for "tts.convert".
type TTSConvertParams struct {
	Text    string `json:"text"`
	Channel string `json:"channel,omitempty"`
}

// TTSConvertResult is the result of "tts.convert".
type TTSConvertResult struct {
	AudioPath       string `json:"audioPath"`
	Provider        string `json:"provider,omitempty"`
	OutputFormat    string `json:"outputFormat,omitempty"`
	VoiceCompatible *bool  `json:"voiceCompatible,omitempty"`
}

// TTSSetProviderParams are the params for "tts.setProvider".
type TTSSetProviderParams struct {
	Provider string `json:"provider"`
}

// TTSSetProviderResult is the result of "tts.setProvider".
type TTSSetProviderResult struct {
	Provider string `json:"provider"`
}

// ---------------------------------------------------------------------------
// exec.approval.waitDecision types
// ---------------------------------------------------------------------------

// ExecApprovalWaitDecisionParams are the params for "exec.approval.waitDecision".
type ExecApprovalWaitDecisionParams struct {
	ID string `json:"id"`
}

// ExecApprovalWaitDecisionResult is the result of "exec.approval.waitDecision".
type ExecApprovalWaitDecisionResult struct {
	ID          string  `json:"id"`
	Decision    *string `json:"decision"` // "allow-once", "allow-always", "deny", or null
	CreatedAtMs *int64  `json:"createdAtMs,omitempty"`
	ExpiresAtMs *int64  `json:"expiresAtMs,omitempty"`
}

// ExecApprovalRequestResult is the result of "exec.approval.request".
type ExecApprovalRequestResult struct {
	ID          string  `json:"id"`
	Status      string  `json:"status,omitempty"` // "accepted" (for twoPhase first response)
	Decision    *string `json:"decision,omitempty"`
	CreatedAtMs int64   `json:"createdAtMs"`
	ExpiresAtMs int64   `json:"expiresAtMs"`
}

// ExecApprovalResolveResult is the result of "exec.approval.resolve".
type ExecApprovalResolveResult struct {
	OK bool `json:"ok"`
}

// ExecApprovalResolvedEvent is the payload of an "exec.approval.resolved" event.
type ExecApprovalResolvedEvent struct {
	ID         string `json:"id"`
	Decision   string `json:"decision"`
	ResolvedBy string `json:"resolvedBy,omitempty"`
	Ts         int64  `json:"ts"`
}

// ---------------------------------------------------------------------------
// Web login types
// ---------------------------------------------------------------------------

// WebLoginStartParams are the params for starting web login.
type WebLoginStartParams struct {
	Force     *bool  `json:"force,omitempty"`
	TimeoutMs *int   `json:"timeoutMs,omitempty"`
	Verbose   *bool  `json:"verbose,omitempty"`
	AccountID string `json:"accountId,omitempty"`
}

// WebLoginWaitParams are the params for waiting on web login.
type WebLoginWaitParams struct {
	TimeoutMs *int   `json:"timeoutMs,omitempty"`
	AccountID string `json:"accountId,omitempty"`
}
