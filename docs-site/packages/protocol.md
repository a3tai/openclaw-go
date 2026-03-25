---
title: protocol
---

# protocol

```
import "github.com/a3tai/openclaw-go/protocol"
```

Package `protocol` defines all wire types, constants, and serialization helpers for the OpenClaw Gateway WebSocket protocol (version 3). It contains no networking code -- only data structures and JSON marshal/unmarshal functions.

Every other package in this module depends on `protocol`.

## Frame Types

The gateway protocol uses typed JSON frames over WebSocket text messages:

| Frame Type | Constant | Direction | Purpose |
|------------|----------|-----------|---------|
| `req` | `FrameTypeRequest` | client -> gateway | RPC request |
| `res` | `FrameTypeResponse` | gateway -> client | RPC response |
| `event` | `FrameTypeEvent` | bidirectional | One-way notification |
| `invoke` | `FrameTypeInvoke` | gateway -> node | Invoke a node capability |
| `invoke-res` | `FrameTypeInvokeResponse` | node -> gateway | Invoke result |

```go
type Request struct {
    Type   FrameType       `json:"type"`
    ID     string          `json:"id"`
    Method string          `json:"method"`
    Params json.RawMessage `json:"params,omitempty"`
}

type Response struct {
    Type    FrameType       `json:"type"`
    ID      string          `json:"id"`
    OK      bool            `json:"ok"`
    Payload json.RawMessage `json:"payload,omitempty"`
    Error   *ErrorPayload   `json:"error,omitempty"`
}

type Event struct {
    Type         FrameType       `json:"type"`
    EventName    string          `json:"event"`
    Payload      json.RawMessage `json:"payload,omitempty"`
    Seq          *int64          `json:"seq,omitempty"`
    StateVersion *StateVersion   `json:"stateVersion,omitempty"`
}
```

## Serialization Helpers

```go
// Serialize frames
data, err := protocol.MarshalRequest(id, method, params)
data, err := protocol.MarshalResponse(id, payload)
data, err := protocol.MarshalErrorResponse(id, errPayload)
data, err := protocol.MarshalEvent(eventName, payload)

// Deserialize frames
frame, err := protocol.ParseFrame(data)      // determine type
req, err   := protocol.UnmarshalRequest(data)
resp, err  := protocol.UnmarshalResponse(data)
ev, err    := protocol.UnmarshalEvent(data)
```

## Roles and Scopes

```go
// Roles
protocol.RoleOperator  // "operator" - UI/CLI clients
protocol.RoleNode      // "node"     - capability hosts

// Operator scopes
protocol.ScopeOperatorRead      // "operator.read"
protocol.ScopeOperatorWrite     // "operator.write"
protocol.ScopeOperatorAdmin     // "operator.admin"
protocol.ScopeOperatorApprovals // "operator.approvals"
protocol.ScopeOperatorPairing   // "operator.pairing"
```

## Connect Handshake

The connection lifecycle is: challenge -> connect request -> hello-ok response.

```go
type ConnectChallenge struct {
    Nonce string `json:"nonce"`
    Ts    int64  `json:"ts"`
}

type ConnectParams struct {
    MinProtocol int            `json:"minProtocol"`
    MaxProtocol int            `json:"maxProtocol"`
    Auth        AuthParams     `json:"auth"`
    Client      ClientInfo     `json:"client"`
    Role        Role           `json:"role"`
    Scopes      []Scope        `json:"scopes,omitempty"`
    // ... plus Caps, Commands, Permissions, Device, Locale, UserAgent, PathEnv
}

type HelloOK struct {
    Protocol int            `json:"protocol"`
    Server   HelloServer    `json:"server"`
    Features HelloFeatures  `json:"features"`
    Snapshot *Snapshot      `json:"snapshot,omitempty"`
    Policy   HelloPolicy    `json:"policy"`
}
```

## RPC Parameter/Result Types

The package defines typed Go structs for every gateway RPC method. Major categories:

| Category | Example Types |
|----------|---------------|
| Chat | `ChatSendParams`, `ChatHistoryParams`, `ChatAbortParams`, `ChatEvent` |
| Agent | `AgentParams`, `AgentEvent`, `AgentIdentityParams` |
| Sessions | `SessionsListParams`, `SessionsPatchParams`, `SessionsUsageParams` |
| Agents CRUD | `AgentsCreateParams`, `AgentsUpdateParams`, `AgentsFilesGetParams` |
| Config | `ConfigGetParams`, `ConfigSetParams`, `ConfigApplyParams`, `ConfigSchemaResponse` |
| Nodes | `NodePairRequestParams`, `NodeInvokeParams`, `NodeEventParams` |
| Device Pairing | `DevicePairApproveParams`, `DeviceTokenRotateParams` |
| Exec Approvals | `ExecApprovalRequestParams`, `ExecApprovalResolveParams`, `ExecApprovalWaitDecisionParams` |
| Cron | `CronAddParams`, `CronUpdateParams`, `CronRunParams`, `CronJob` |
| TTS | `TTSConvertParams`, `TTSSetProviderParams`, `TTSStatusResult` |
| Channels/Talk | `TalkModeParams`, `TalkConfigParams`, `ChannelsStatusResult` |
| Skills | `SkillsInstallParams`, `SkillsUpdateParams`, `SkillsBinsResult` |
| Wizard | `WizardStartParams`, `WizardNextParams`, `WizardStep` |
| Models | `ModelChoice`, `ModelsListResult` |
| Logs | `LogsTailParams`, `LogsTailResult` |
| Push | `PushTestParams`, `PushTestResult` |
| Health/Presence | `PresenceEntry`, `HealthEvent`, `HeartbeatEvent` |

## Event Payload Types

```go
// Real-time events with fully-typed payloads
PresenceEvent           // client presence changes
HealthEvent             // system health snapshots
HeartbeatEvent          // periodic heartbeat with system stats
TickEvent               // keepalive tick
ShutdownEvent           // gateway shutdown
CronEvent               // cron job execution events
VoicewakeChangedEvent   // voice wake word state changes
ExecApprovalResolvedEvent // approval decision events
```

## Constants

```go
protocol.ProtocolVersion              // 3
protocol.MaxPayloadBytes              // 25 MiB
protocol.MaxBufferedBytes             // 50 MiB
protocol.DefaultTickIntervalMs        // 30,000
protocol.DefaultHandshakeTimeoutMs    // 10,000
protocol.SessionLabelMaxLength        // 64

// Error codes
protocol.ErrorCodeNotLinked
protocol.ErrorCodeNotPaired
protocol.ErrorCodeAgentTimeout
protocol.ErrorCodeInvalidRequest
protocol.ErrorCodeUnavailable

// Client IDs
protocol.ClientIDCLI
protocol.ClientIDGateway
protocol.ClientIDMacOS
// ... and more
```
