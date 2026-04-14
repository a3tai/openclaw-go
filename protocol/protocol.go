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

// EventName identifies a gateway protocol event name.
type EventName string

// Known event name constants.
const (
	// EventConnectChallenge is the initial server challenge sent before connect.
	EventConnectChallenge EventName = "connect.challenge"
	// EventTick is the keepalive tick event from the gateway.
	EventTick EventName = "tick"
	// EventShutdown indicates gateway shutdown or maintenance state.
	EventShutdown EventName = "shutdown"
	// EventPresence carries presence snapshot/update information.
	EventPresence EventName = "presence"
	// EventHealth carries health snapshot/update information.
	EventHealth EventName = "health"
	// EventHeartbeat is the heartbeat stream event from agent/runtime.
	EventHeartbeat EventName = "heartbeat"
	// EventChat carries chat stream updates (delta/final/error).
	EventChat EventName = "chat"
	// EventAgent carries agent stream updates.
	EventAgent EventName = "agent"
	// EventCron carries cron execution lifecycle events.
	EventCron EventName = "cron"
	// EventVoicewakeChanged notifies that voicewake config changed.
	EventVoicewakeChanged EventName = "voicewake.changed"
	// EventNodeInvokeRequest is a gateway->node invoke request event.
	EventNodeInvokeRequest EventName = "node.invoke.request"
	// EventNodePairRequested is emitted when node pairing is requested.
	EventNodePairRequested EventName = "node.pair.requested"
	// EventNodePairResolved is emitted when node pairing is approved/rejected.
	EventNodePairResolved EventName = "node.pair.resolved"
	// EventDevicePairRequest is emitted when device pairing is requested.
	EventDevicePairRequest EventName = "device.pair.requested"
	// EventDevicePairResolve is emitted when device pairing is resolved.
	EventDevicePairResolve EventName = "device.pair.resolved"
	// EventExecApproval is emitted when command execution needs approval.
	EventExecApproval EventName = "exec.approval.requested"
	// EventExecFinished is emitted when command execution finishes.
	EventExecFinished EventName = "exec.finished"
	// EventExecDenied is emitted when command execution is denied.
	EventExecDenied EventName = "exec.denied"
	// EventPluginApprovalRequested is emitted when a plugin action needs approval.
	EventPluginApprovalRequested EventName = "plugin.approval.requested"
	// EventPluginApprovalResolved is emitted when a plugin approval is resolved.
	EventPluginApprovalResolved EventName = "plugin.approval.resolved"
)

// MethodName identifies a gateway RPC method name.
type MethodName string

// Known RPC method name constants.
const (
	// Agent identity and lifecycle.
	MethodAgentIdentityGet MethodName = "agent.identity.get"
	MethodAgentWait        MethodName = "agent.wait"

	// Agent management and attached files.
	MethodAgentsCreate    MethodName = "agents.create"
	MethodAgentsDelete    MethodName = "agents.delete"
	MethodAgentsFilesGet  MethodName = "agents.files.get"
	MethodAgentsFilesList MethodName = "agents.files.list"
	MethodAgentsFilesSet  MethodName = "agents.files.set"
	MethodAgentsList      MethodName = "agents.list"
	MethodAgentsUpdate    MethodName = "agents.update"

	// Browser and channel integration.
	// Deprecated: browser.request is not in upstream BASE_METHODS; treat as a deprecation candidate.
	MethodBrowserRequest MethodName = "browser.request"
	MethodChannelsLogout MethodName = "channels.logout"
	MethodChannelsStatus MethodName = "channels.status"
	MethodMessageAction  MethodName = "message.action"

	// Chat request/stream controls.
	MethodChatAbort   MethodName = "chat.abort"
	MethodChatHistory MethodName = "chat.history"
	MethodChatInject  MethodName = "chat.inject"
	MethodChatSend    MethodName = "chat.send"

	// Configuration and schema operations.
	MethodConfigApply        MethodName = "config.apply"
	MethodConfigGet          MethodName = "config.get"
	MethodConfigPatch        MethodName = "config.patch"
	MethodConfigSchema       MethodName = "config.schema"
	MethodConfigSchemaLookup MethodName = "config.schema.lookup"
	MethodConfigSet          MethodName = "config.set"

	// Command catalog.
	MethodCommandsList MethodName = "commands.list"

	// Cron management and execution history.
	MethodCronAdd    MethodName = "cron.add"
	MethodCronList   MethodName = "cron.list"
	MethodCronRemove MethodName = "cron.remove"
	MethodCronRun    MethodName = "cron.run"
	MethodCronRuns   MethodName = "cron.runs"
	MethodCronStatus MethodName = "cron.status"
	MethodCronUpdate MethodName = "cron.update"

	// Device pairing and token management.
	MethodDevicePairApprove MethodName = "device.pair.approve"
	MethodDevicePairList    MethodName = "device.pair.list"
	MethodDevicePairReject  MethodName = "device.pair.reject"
	MethodDevicePairRemove  MethodName = "device.pair.remove"
	MethodDeviceTokenRevoke MethodName = "device.token.revoke"
	MethodDeviceTokenRotate MethodName = "device.token.rotate"

	// Commands.
	MethodCommandsList MethodName = "commands.list"

	// Diagnostics and execution approvals.
	MethodDoctorMemoryDreamDiary MethodName = "doctor.memory.dreamDiary"
	MethodDoctorMemoryStatus     MethodName = "doctor.memory.status"
	MethodExecApprovalGet        MethodName = "exec.approval.get"
	MethodExecApprovalList       MethodName = "exec.approval.list"
	MethodExecApprovalRequest    MethodName = "exec.approval.request"
	MethodExecApprovalResolve    MethodName = "exec.approval.resolve"
	MethodExecApprovalsGet       MethodName = "exec.approvals.get"
	MethodExecApprovalsNodeGet   MethodName = "exec.approvals.node.get"
	MethodExecApprovalsNodeSet   MethodName = "exec.approvals.node.set"
	MethodExecApprovalsSet       MethodName = "exec.approvals.set"

	// Message actions.
	MethodMessageAction MethodName = "message.action"

	// Plugin approvals.
	MethodPluginApprovalList         MethodName = "plugin.approval.list"
	MethodPluginApprovalRequest      MethodName = "plugin.approval.request"
	MethodPluginApprovalWaitDecision MethodName = "plugin.approval.waitDecision"
	MethodPluginApprovalResolve      MethodName = "plugin.approval.resolve"

	// Gateway status and logs.
	MethodGatewayIdentityGet MethodName = "gateway.identity.get"
	MethodLastHeartbeat      MethodName = "last-heartbeat"
	MethodLogsTail           MethodName = "logs.tail"
	MethodModelsList         MethodName = "models.list"

	// Node RPC and node queue/pairing operations.
	MethodNodeCanvasCapabilityRefresh MethodName = "node.canvas.capability.refresh"
	MethodNodeDescribe                MethodName = "node.describe"
	MethodNodeEvent                   MethodName = "node.event"
	MethodNodeInvoke                  MethodName = "node.invoke"
	MethodNodeInvokeResult            MethodName = "node.invoke.result"
	MethodNodeList                    MethodName = "node.list"
	MethodNodePairApprove             MethodName = "node.pair.approve"
	MethodNodePairList                MethodName = "node.pair.list"
	MethodNodePairReject              MethodName = "node.pair.reject"
	MethodNodePairRequest             MethodName = "node.pair.request"
	MethodNodePairVerify              MethodName = "node.pair.verify"
	MethodNodePendingAck              MethodName = "node.pending.ack"
	MethodNodePendingDrain            MethodName = "node.pending.drain"
	MethodNodePendingEnqueue          MethodName = "node.pending.enqueue"
	MethodNodePendingPull             MethodName = "node.pending.pull"
	MethodNodeRename                  MethodName = "node.rename"

	// Utility and secrets.
	MethodPushTest       MethodName = "push.test"
	MethodSecretsReload  MethodName = "secrets.reload"
	MethodSecretsResolve MethodName = "secrets.resolve"

	// Session lifecycle, messaging, compaction, and usage analytics.
	MethodSessionsAbort               MethodName = "sessions.abort"
	MethodSessionsCompact             MethodName = "sessions.compact"
	MethodSessionsCompactionBranch    MethodName = "sessions.compaction.branch"
	MethodSessionsCompactionGet       MethodName = "sessions.compaction.get"
	MethodSessionsCompactionList      MethodName = "sessions.compaction.list"
	MethodSessionsCompactionRestore   MethodName = "sessions.compaction.restore"
	MethodSessionsCreate              MethodName = "sessions.create"
	MethodSessionsDelete              MethodName = "sessions.delete"
	MethodSessionsGet                 MethodName = "sessions.get"
	MethodSessionsList                MethodName = "sessions.list"
	MethodSessionsMessagesSubscribe   MethodName = "sessions.messages.subscribe"
	MethodSessionsMessagesUnsubscribe MethodName = "sessions.messages.unsubscribe"
	MethodSessionsPatch               MethodName = "sessions.patch"
	MethodSessionsPreview             MethodName = "sessions.preview"
	MethodSessionsReset               MethodName = "sessions.reset"
	MethodSessionsResolve             MethodName = "sessions.resolve"
	MethodSessionsSend                MethodName = "sessions.send"
	MethodSessionsSteer               MethodName = "sessions.steer"
	MethodSessionsSubscribe           MethodName = "sessions.subscribe"
	MethodSessionsUnsubscribe         MethodName = "sessions.unsubscribe"
	MethodSessionsUsage               MethodName = "sessions.usage"
	MethodSessionsUsageLogs           MethodName = "sessions.usage.logs"
	MethodSessionsUsageTimeseries     MethodName = "sessions.usage.timeseries"

	// Runtime controls and skills.
	MethodSetHeartbeats MethodName = "set-heartbeats"
	MethodSkillsBins    MethodName = "skills.bins"
	MethodSkillsDetail  MethodName = "skills.detail"
	MethodSkillsInstall MethodName = "skills.install"
	MethodSkillsSearch  MethodName = "skills.search"
	MethodSkillsStatus  MethodName = "skills.status"
	MethodSkillsUpdate  MethodName = "skills.update"

	// System presence/events and voice interactions.
	MethodSystemEvent    MethodName = "system-event"
	MethodSystemPresence MethodName = "system-presence"
	MethodTalkConfig     MethodName = "talk.config"
	MethodTalkMode       MethodName = "talk.mode"
	MethodTalkSpeak      MethodName = "talk.speak"

	// Tool and TTS interfaces.
	MethodToolsCatalog   MethodName = "tools.catalog"
	MethodToolsEffective MethodName = "tools.effective"
	MethodTTSConvert     MethodName = "tts.convert"
	MethodTTSDisable     MethodName = "tts.disable"
	MethodTTSEnable      MethodName = "tts.enable"
	MethodTTSProviders   MethodName = "tts.providers"
	MethodTTSStatus      MethodName = "tts.status"

	// Update, usage, and voice wake configuration.
	MethodUpdateRun    MethodName = "update.run"
	MethodUsageCost    MethodName = "usage.cost"
	MethodUsageStatus  MethodName = "usage.status"
	MethodVoiceWakeGet MethodName = "voicewake.get"
	MethodVoiceWakeSet MethodName = "voicewake.set"

	// Web auth and interactive setup wizard.
	MethodWebLoginStart MethodName = "web.login.start"
	MethodWebLoginWait  MethodName = "web.login.wait"
	MethodWizardCancel  MethodName = "wizard.cancel"
	MethodWizardNext    MethodName = "wizard.next"
	MethodWizardStart   MethodName = "wizard.start"
	MethodWizardStatus  MethodName = "wizard.status"
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
	EventName    EventName       `json:"event"`                  // event name
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
	Event EventName `json:"event,omitempty"` // only for event frames
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

// ExecApprovalGetParams are the params for "exec.approval.get".
type ExecApprovalGetParams struct {
	// ID is the approval request ID returned by exec.approval.request.
	ID string `json:"id"`
}

// ExecApprovalGetResult is the result of "exec.approval.get".
// Host, NodeID, AgentID, and CommandPreview may be nil when not applicable
// to the approval context.
type ExecApprovalGetResult struct {
	ID               string   `json:"id"`
	CommandText      string   `json:"commandText"`
	CommandPreview   *string  `json:"commandPreview,omitempty"`
	AllowedDecisions []string `json:"allowedDecisions"`
	Host             *string  `json:"host,omitempty"`
	NodeID           *string  `json:"nodeId,omitempty"`
	AgentID          *string  `json:"agentId,omitempty"`
	ExpiresAtMs      int64    `json:"expiresAtMs"`
}

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
// ArgPattern optionally constrains which argument patterns are allowed for the entry.
type ExecApprovalsAllowlistEntry struct {
	ID               string `json:"id,omitempty"`
	Pattern          string `json:"pattern"`
	ArgPattern       string `json:"argPattern,omitempty"`
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
	MaxChars   *int   `json:"maxChars,omitempty"`
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

// ChatSendResult is the RPC response from "chat.send".
// This is the initial ack — streaming content arrives via "chat" events.
type ChatSendResult struct {
	RunID          string `json:"runId"`
	Status         string `json:"status"`                         // "started"
	MessageSeq     *int   `json:"messageSeq,omitempty"`           // sequence number if applicable
	InterruptedRun *bool  `json:"interruptedActiveRun,omitempty"` // true if a prior run was interrupted
}

// ChatEvent is the payload of a "chat" event (received via the event stream).
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

// SessionsCreateParams are the params for "sessions.create".
type SessionsCreateParams struct {
	Key              string `json:"key,omitempty"`
	AgentID          string `json:"agentId,omitempty"`
	Label            string `json:"label,omitempty"`
	Model            string `json:"model,omitempty"`
	ParentSessionKey string `json:"parentSessionKey,omitempty"`
	Task             string `json:"task,omitempty"`
	Message          string `json:"message,omitempty"`
}

// SessionsSendParams are the params for "sessions.send" and "sessions.steer".
type SessionsSendParams struct {
	Key            string          `json:"key"`
	Message        string          `json:"message"`
	Thinking       string          `json:"thinking,omitempty"`
	Attachments    json.RawMessage `json:"attachments,omitempty"`
	TimeoutMs      *int            `json:"timeoutMs,omitempty"`
	IdempotencyKey string          `json:"idempotencyKey,omitempty"`
}

// SessionsAbortParams are the params for "sessions.abort".
type SessionsAbortParams struct {
	Key   string `json:"key"`
	RunID string `json:"runId,omitempty"`
}

// SessionsMessagesSubscribeParams are the params for "sessions.messages.subscribe".
type SessionsMessagesSubscribeParams struct {
	Key string `json:"key"`
}

// SessionsMessagesUnsubscribeParams are the params for "sessions.messages.unsubscribe".
type SessionsMessagesUnsubscribeParams struct {
	Key string `json:"key"`
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

// SessionsUsageTimeseriesParams are the params for "sessions.usage.timeseries".
type SessionsUsageTimeseriesParams struct {
	Key string `json:"key"`
}

// SessionsUsageLogsParams are the params for "sessions.usage.logs".
type SessionsUsageLogsParams struct {
	Key   string `json:"key"`
	Limit *int   `json:"limit,omitempty"`
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

// NodePendingEnqueueParams are the params for "node.pending.enqueue".
// Added in openclaw v2026.3.11.
type NodePendingEnqueueParams struct {
	NodeID      string `json:"nodeId"`
	Type        string `json:"type"` // "status.request" or "location.request"
	Priority    string `json:"priority,omitempty"`
	ExpiresInMs *int   `json:"expiresInMs,omitempty"`
	Wake        *bool  `json:"wake,omitempty"`
}

// NodePendingDrainParams are the params for "node.pending.drain".
// Added in openclaw v2026.3.11.
type NodePendingDrainParams struct {
	MaxItems *int `json:"maxItems,omitempty"`
}

// NodePendingAckParams are the params for "node.pending.ack".
type NodePendingAckParams struct {
	IDs []string `json:"ids"`
}

// NodePendingWorkItem is a queued node.pending work item.
type NodePendingWorkItem struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Priority    string         `json:"priority"`
	CreatedAtMs int64          `json:"createdAtMs"`
	ExpiresAtMs *int64         `json:"expiresAtMs,omitempty"`
	Payload     map[string]any `json:"payload,omitempty"`
}

// NodePendingEnqueueResult is the result of "node.pending.enqueue".
type NodePendingEnqueueResult struct {
	NodeID        string              `json:"nodeId"`
	Revision      int                 `json:"revision"`
	Queued        NodePendingWorkItem `json:"queued"`
	WakeTriggered bool                `json:"wakeTriggered"`
}

// NodePendingDrainResult is the result of "node.pending.drain".
type NodePendingDrainResult struct {
	NodeID   string                `json:"nodeId"`
	Revision int                   `json:"revision"`
	Items    []NodePendingWorkItem `json:"items"`
	HasMore  bool                  `json:"hasMore"`
}

// NodePendingAction is a pending action returned by "node.pending.pull".
type NodePendingAction struct {
	ID           string  `json:"id"`
	Command      string  `json:"command"`
	ParamsJSON   *string `json:"paramsJSON"`
	EnqueuedAtMs int64   `json:"enqueuedAtMs"`
}

// NodePendingPullResult is the result of "node.pending.pull".
type NodePendingPullResult struct {
	NodeID  string              `json:"nodeId"`
	Actions []NodePendingAction `json:"actions"`
}

// NodePendingAckResult is the result of "node.pending.ack".
type NodePendingAckResult struct {
	NodeID         string   `json:"nodeId"`
	AckedIDs       []string `json:"ackedIds"`
	RemainingCount int      `json:"remainingCount"`
}

// NodeCanvasCapabilityRefreshResult is the result of "node.canvas.capability.refresh".
type NodeCanvasCapabilityRefreshResult struct {
	CanvasCapability            string `json:"canvasCapability"`
	CanvasCapabilityExpiresAtMs int64  `json:"canvasCapabilityExpiresAtMs"`
	CanvasHostURL               string `json:"canvasHostUrl"`
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

// ConfigSchemaLookupParams are the params for "config.schema.lookup".
type ConfigSchemaLookupParams struct {
	Path string `json:"path"`
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

// ConfigSchemaLookupChild describes one child entry in a schema lookup result.
type ConfigSchemaLookupChild struct {
	Key         string          `json:"key"`
	Path        string          `json:"path"`
	Type        json.RawMessage `json:"type,omitempty"`
	Required    bool            `json:"required"`
	HasChildren bool            `json:"hasChildren"`
	Hint        *ConfigUIHint   `json:"hint,omitempty"`
	HintPath    string          `json:"hintPath,omitempty"`
}

// ConfigSchemaLookupResult is the result for "config.schema.lookup".
type ConfigSchemaLookupResult struct {
	Path     string                    `json:"path"`
	Schema   json.RawMessage           `json:"schema"`
	Hint     *ConfigUIHint             `json:"hint,omitempty"`
	HintPath string                    `json:"hintPath,omitempty"`
	Children []ConfigSchemaLookupChild `json:"children"`
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

// AgentSummaryModel is the optional model override for an agent.
type AgentSummaryModel struct {
	Primary   string   `json:"primary,omitempty"`
	Fallbacks []string `json:"fallbacks,omitempty"`
}

// AgentSummary is a summary of an agent.
type AgentSummary struct {
	ID        string             `json:"id"`
	Name      string             `json:"name,omitempty"`
	Identity  *AgentIdentity     `json:"identity,omitempty"`
	Workspace string             `json:"workspace,omitempty"`
	Model     *AgentSummaryModel `json:"model,omitempty"`
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
	LightContext               *bool  `json:"lightContext,omitempty"`
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
	// Deprecated: use Enabled instead.
	IncludeDisabled *bool  `json:"includeDisabled,omitempty"`
	Limit           *int   `json:"limit,omitempty"`
	Offset          *int   `json:"offset,omitempty"`
	Query           string `json:"query,omitempty"`
	Enabled         string `json:"enabled,omitempty"` // "all", "enabled", "disabled"
	SortBy          string `json:"sortBy,omitempty"`  // "nextRunAtMs", "updatedAtMs", "name"
	SortDir         string `json:"sortDir,omitempty"` // "asc", "desc"
}

// CronListResult is the paginated result of "cron.list".
type CronListResult struct {
	Jobs       []CronJob `json:"jobs"`
	Total      int       `json:"total"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	HasMore    bool      `json:"hasMore"`
	NextOffset *int      `json:"nextOffset"`
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
// Use ID for the server-assigned job UUID. JobID is a wire-compat alias.
// When both Statuses and Status are set, the server uses Statuses.
type CronRunsParams struct {
	Scope            string   `json:"scope,omitempty"` // "job" or "all"
	ID               string   `json:"id,omitempty"`
	JobID            string   `json:"jobId,omitempty"` // alias for ID
	Limit            *int     `json:"limit,omitempty"`
	Offset           *int     `json:"offset,omitempty"`
	Statuses         []string `json:"statuses,omitempty"`         // ["ok","error","skipped"]; takes precedence over Status
	Status           string   `json:"status,omitempty"`           // convenience alias: "all","ok","error","skipped"
	DeliveryStatuses []string `json:"deliveryStatuses,omitempty"` // takes precedence over DeliveryStatus
	DeliveryStatus   string   `json:"deliveryStatus,omitempty"`   // convenience alias
	Query            string   `json:"query,omitempty"`
	SortDir          string   `json:"sortDir,omitempty"` // "asc", "desc"
}

// CronRunLogEntry is a single entry in the cron run log.
type CronRunLogEntry struct {
	Ts             int64            `json:"ts"`
	JobID          string           `json:"jobId"`
	Action         string           `json:"action"` // "finished"
	Status         string           `json:"status,omitempty"`
	Error          string           `json:"error,omitempty"`
	Summary        string           `json:"summary,omitempty"`
	Delivered      *bool            `json:"delivered,omitempty"`
	DeliveryStatus string           `json:"deliveryStatus,omitempty"`
	DeliveryError  string           `json:"deliveryError,omitempty"`
	SessionID      string           `json:"sessionId,omitempty"`
	SessionKey     string           `json:"sessionKey,omitempty"`
	RunAtMs        *int64           `json:"runAtMs,omitempty"`
	DurationMs     *int64           `json:"durationMs,omitempty"`
	NextRunAtMs    *int64           `json:"nextRunAtMs,omitempty"`
	Model          string           `json:"model,omitempty"`
	Provider       string           `json:"provider,omitempty"`
	Usage          *CronRunLogUsage `json:"usage,omitempty"`
	JobName        string           `json:"jobName,omitempty"`
}

// CronRunLogUsage contains token usage for a cron run.
type CronRunLogUsage struct {
	InputTokens      *int `json:"input_tokens,omitempty"`
	OutputTokens     *int `json:"output_tokens,omitempty"`
	TotalTokens      *int `json:"total_tokens,omitempty"`
	CacheReadTokens  *int `json:"cache_read_tokens,omitempty"`
	CacheWriteTokens *int `json:"cache_write_tokens,omitempty"`
}

// CronRunsResult is the paginated result of "cron.runs".
type CronRunsResult struct {
	Entries []CronRunLogEntry `json:"entries"`
	Total   int               `json:"total"`
	Offset  int               `json:"offset"`
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

// TalkSpeakParams are the params for "talk.speak".
// Added in openclaw v2026.3.22.
type TalkSpeakParams struct {
	Text         string   `json:"text"`
	VoiceID      string   `json:"voiceId,omitempty"`
	ModelID      string   `json:"modelId,omitempty"`
	OutputFormat string   `json:"outputFormat,omitempty"`
	Speed        *float64 `json:"speed,omitempty"`
	Stability    *float64 `json:"stability,omitempty"`
	Similarity   *float64 `json:"similarity,omitempty"`
	Style        *float64 `json:"style,omitempty"`
	SpeakerBoost *bool    `json:"speakerBoost,omitempty"`
	Seed         *int     `json:"seed,omitempty"`
	Normalize    string   `json:"normalize,omitempty"`
	Language     string   `json:"language,omitempty"`
}

// TalkSpeakResult is the result of "talk.speak".
type TalkSpeakResult struct {
	AudioBase64     string `json:"audioBase64"`
	Provider        string `json:"provider"`
	OutputFormat    string `json:"outputFormat,omitempty"`
	VoiceCompatible *bool  `json:"voiceCompatible,omitempty"`
	MimeType        string `json:"mimeType,omitempty"`
	FileExtension   string `json:"fileExtension,omitempty"`
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
	Name                          string `json:"name"`
	InstallID                     string `json:"installId"`
	DangerouslyForceUnsafeInstall *bool  `json:"dangerouslyForceUnsafeInstall,omitempty"`
	TimeoutMs                     *int   `json:"timeoutMs,omitempty"`
}

// SkillsUpdateParams are the params for "skills.update".
type SkillsUpdateParams struct {
	SkillKey string            `json:"skillKey"`
	Enabled  *bool             `json:"enabled,omitempty"`
	APIKey   string            `json:"apiKey,omitempty"`
	Env      map[string]string `json:"env,omitempty"`
}

// SkillsSearchParams are the params for "skills.search".
// Both Query and Limit are optional. Limit is capped at 100 by the server.
type SkillsSearchParams struct {
	Query string `json:"query,omitempty"`
	Limit *int   `json:"limit,omitempty"`
}

// SkillsSearchItem is a single result from "skills.search".
type SkillsSearchItem struct {
	Score       float64 `json:"score"`
	Slug        string  `json:"slug"`
	DisplayName string  `json:"displayName"`
	Summary     string  `json:"summary,omitempty"`
	Version     string  `json:"version,omitempty"`
	UpdatedAt   *int64  `json:"updatedAt,omitempty"`
}

// SkillsSearchResult is the result of "skills.search".
type SkillsSearchResult struct {
	Results []SkillsSearchItem `json:"results"`
}

// SkillsDetailParams are the params for "skills.detail".
type SkillsDetailParams struct {
	Slug string `json:"slug"`
}

// SkillsDetailSkill is the skill object in "skills.detail" result.
type SkillsDetailSkill struct {
	Slug        string            `json:"slug"`
	DisplayName string            `json:"displayName"`
	Summary     string            `json:"summary,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	CreatedAt   int64             `json:"createdAt"`
	UpdatedAt   int64             `json:"updatedAt"`
}

// SkillsDetailVersion is the latest version info in "skills.detail" result.
type SkillsDetailVersion struct {
	Version   string `json:"version"`
	CreatedAt int64  `json:"createdAt"`
	Changelog string `json:"changelog,omitempty"`
}

// SkillsDetailMetadata is the metadata in "skills.detail" result.
type SkillsDetailMetadata struct {
	OS      []string `json:"os,omitempty"`
	Systems []string `json:"systems,omitempty"`
}

// SkillsDetailOwner is the owner info in "skills.detail" result.
type SkillsDetailOwner struct {
	Handle      *string `json:"handle,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	Image       *string `json:"image,omitempty"`
}

// SkillsDetailResult is the result of "skills.detail".
// Skill is nil when no skill matches the requested slug.
// LatestVersion, Metadata, and Owner may be nil when not available.
type SkillsDetailResult struct {
	Skill         *SkillsDetailSkill    `json:"skill"`
	LatestVersion *SkillsDetailVersion  `json:"latestVersion,omitempty"`
	Metadata      *SkillsDetailMetadata `json:"metadata,omitempty"`
	Owner         *SkillsDetailOwner    `json:"owner,omitempty"`
}

// ToolsCatalogParams are the params for "tools.catalog".
type ToolsCatalogParams struct {
	AgentID        string `json:"agentId,omitempty"`
	IncludePlugins *bool  `json:"includePlugins,omitempty"`
}

// ToolsEffectiveParams are the params for "tools.effective".
// SessionKey is required; AgentID is optional and must match the session's agent when provided.
type ToolsEffectiveParams struct {
	SessionKey string `json:"sessionKey"`
	AgentID    string `json:"agentId,omitempty"`
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

// MessageActionToolContext identifies the current channel context for plugin-dispatched actions.
type MessageActionToolContext struct {
	CurrentChannelID       *string `json:"currentChannelId,omitempty"`
	CurrentChannelProvider *string `json:"currentChannelProvider,omitempty"`
	CurrentThreadTS        *string `json:"currentThreadTs,omitempty"`
	CurrentMessageID       any     `json:"currentMessageId,omitempty"`
}

// MessageActionParams are the params for "message.action".
type MessageActionParams struct {
	Channel           string                    `json:"channel"`
	Action            string                    `json:"action"`
	Params            map[string]any            `json:"params"`
	AccountID         *string                   `json:"accountId,omitempty"`
	RequesterSenderID *string                   `json:"requesterSenderId,omitempty"`
	SenderIsOwner     *bool                     `json:"senderIsOwner,omitempty"`
	SessionKey        *string                   `json:"sessionKey,omitempty"`
	SessionID         *string                   `json:"sessionId,omitempty"`
	AgentID           *string                   `json:"agentId,omitempty"`
	ToolContext       *MessageActionToolContext `json:"toolContext,omitempty"`
	IdempotencyKey    string                    `json:"idempotencyKey"`
}

// CommandArgChoice is a fixed value allowed for a command argument.
type CommandArgChoice struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// CommandArg describes a single command argument.
type CommandArg struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	Required    *bool              `json:"required,omitempty"`
	Choices     []CommandArgChoice `json:"choices,omitempty"`
	Dynamic     *bool              `json:"dynamic,omitempty"`
}

// CommandEntry describes an available command.
type CommandEntry struct {
	Name        string       `json:"name"`
	NativeName  *string      `json:"nativeName,omitempty"`
	TextAliases []string     `json:"textAliases,omitempty"`
	Description string       `json:"description"`
	Category    *string      `json:"category,omitempty"`
	Source      string       `json:"source"`
	Scope       string       `json:"scope"`
	AcceptsArgs bool         `json:"acceptsArgs"`
	Args        []CommandArg `json:"args,omitempty"`
}

// CommandsListParams are the params for "commands.list".
type CommandsListParams struct {
	AgentID     *string `json:"agentId,omitempty"`
	Provider    *string `json:"provider,omitempty"`
	Scope       *string `json:"scope,omitempty"`
	IncludeArgs *bool   `json:"includeArgs,omitempty"`
}

// CommandsListResult is the result of "commands.list".
type CommandsListResult struct {
	Commands []CommandEntry `json:"commands"`
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

// SecretsReloadParams are the params for "secrets.reload".
type SecretsReloadParams struct{}

// SecretsResolveParams are the params for "secrets.resolve".
type SecretsResolveParams struct {
	CommandName string   `json:"commandName"`
	TargetIDs   []string `json:"targetIds"`
}

// SecretsResolveAssignment is one resolved assignment for "secrets.resolve".
type SecretsResolveAssignment struct {
	Path         string          `json:"path,omitempty"`
	PathSegments []string        `json:"pathSegments"`
	Value        json.RawMessage `json:"value"`
}

// SecretsResolveResult is the result for "secrets.resolve".
type SecretsResolveResult struct {
	OK               *bool                      `json:"ok,omitempty"`
	Assignments      []SecretsResolveAssignment `json:"assignments,omitempty"`
	Diagnostics      []string                   `json:"diagnostics,omitempty"`
	InactiveRefPaths []string                   `json:"inactiveRefPaths,omitempty"`
}

// GatewayIdentityResult is the result for "gateway.identity.get".
type GatewayIdentityResult struct {
	DeviceID  string `json:"deviceId"`
	PublicKey string `json:"publicKey"`
}

// DoctorMemoryEmbeddingStatus describes embedding provider availability.
type DoctorMemoryEmbeddingStatus struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// DoctorMemoryStatusResult is the result for "doctor.memory.status".
type DoctorMemoryStatusResult struct {
	AgentID   string                      `json:"agentId"`
	Provider  string                      `json:"provider,omitempty"`
	Embedding DoctorMemoryEmbeddingStatus `json:"embedding"`
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

// ---------------------------------------------------------------------------
// plugin.approval types
// ---------------------------------------------------------------------------

// PluginApprovalRequestParams are the params for "plugin.approval.request".
type PluginApprovalRequestParams struct {
	PluginID            *string `json:"pluginId,omitempty"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Severity            *string `json:"severity,omitempty"` // "info", "warning", "critical"
	ToolName            *string `json:"toolName,omitempty"`
	ToolCallID          *string `json:"toolCallId,omitempty"`
	AgentID             *string `json:"agentId,omitempty"`
	SessionKey          *string `json:"sessionKey,omitempty"`
	TurnSourceChannel   *string `json:"turnSourceChannel,omitempty"`
	TurnSourceTo        *string `json:"turnSourceTo,omitempty"`
	TurnSourceAccountID *string `json:"turnSourceAccountId,omitempty"`
	TurnSourceThreadID  any     `json:"turnSourceThreadId,omitempty"` // string or number
	TimeoutMs           *int    `json:"timeoutMs,omitempty"`
	TwoPhase            *bool   `json:"twoPhase,omitempty"`
}

// PluginApprovalRequestResult is the result of "plugin.approval.request".
type PluginApprovalRequestResult struct {
	ID          string  `json:"id"`
	Status      string  `json:"status,omitempty"` // "accepted" (twoPhase first response)
	Decision    *string `json:"decision"`
	CreatedAtMs int64   `json:"createdAtMs"`
	ExpiresAtMs int64   `json:"expiresAtMs"`
}

// PluginApprovalWaitDecisionParams are the params for "plugin.approval.waitDecision".
type PluginApprovalWaitDecisionParams struct {
	ID string `json:"id"`
}

// PluginApprovalWaitDecisionResult is the result of "plugin.approval.waitDecision".
type PluginApprovalWaitDecisionResult struct {
	ID          string  `json:"id"`
	Decision    *string `json:"decision"` // "allow-once", "allow-always", "deny", or null
	CreatedAtMs *int64  `json:"createdAtMs,omitempty"`
	ExpiresAtMs *int64  `json:"expiresAtMs,omitempty"`
}

// PluginApprovalResolveParams are the params for "plugin.approval.resolve".
type PluginApprovalResolveParams struct {
	ID       string `json:"id"`
	Decision string `json:"decision"` // "allow-once", "allow-always", "deny"
}

// PluginApprovalResolveResult is the result of "plugin.approval.resolve".
type PluginApprovalResolveResult struct {
	OK bool `json:"ok"`
}

// PluginApprovalRequestedEvent is the payload of a "plugin.approval.requested" event.
type PluginApprovalRequestedEvent struct {
	ID          string                      `json:"id"`
	Request     PluginApprovalRequestParams `json:"request"`
	CreatedAtMs int64                       `json:"createdAtMs"`
	ExpiresAtMs int64                       `json:"expiresAtMs"`
}

// PluginApprovalResolvedEvent is the payload of a "plugin.approval.resolved" event.
type PluginApprovalResolvedEvent struct {
	ID         string                       `json:"id"`
	Decision   string                       `json:"decision"`
	ResolvedBy *string                      `json:"resolvedBy,omitempty"`
	Ts         int64                        `json:"ts"`
	Request    *PluginApprovalRequestParams `json:"request,omitempty"`
}

// ---------------------------------------------------------------------------
// doctor.memory.dreamDiary types
// ---------------------------------------------------------------------------

// DoctorMemoryDreamDiaryResult is the result for "doctor.memory.dreamDiary".
// Found indicates whether a dream diary file was located in the agent workspace.
// Content is only set when Found is true.
type DoctorMemoryDreamDiaryResult struct {
	AgentID     string  `json:"agentId"`
	Found       bool    `json:"found"`
	Path        string  `json:"path"`
	Content     *string `json:"content,omitempty"`
	UpdatedAtMs *int64  `json:"updatedAtMs,omitempty"`
}

// ---------------------------------------------------------------------------
// exec.approval.list types
// ---------------------------------------------------------------------------

// ExecApprovalListEntry is one pending exec approval in an "exec.approval.list" result.
type ExecApprovalListEntry struct {
	ID          string                    `json:"id"`
	Request     ExecApprovalRequestParams `json:"request"`
	CreatedAtMs int64                     `json:"createdAtMs"`
	ExpiresAtMs int64                     `json:"expiresAtMs"`
}

// ---------------------------------------------------------------------------
// plugin.approval.list types
// ---------------------------------------------------------------------------

// PluginApprovalListEntry is one pending plugin approval in a "plugin.approval.list" result.
type PluginApprovalListEntry struct {
	ID          string                      `json:"id"`
	Request     PluginApprovalRequestParams `json:"request"`
	CreatedAtMs int64                       `json:"createdAtMs"`
	ExpiresAtMs int64                       `json:"expiresAtMs"`
}
