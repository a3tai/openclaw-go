# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v2026.4.12] - 2026-04-14

Sync to upstream openclaw v2026.4.12. 1200 upstream files changed; 8 new RPC methods added; 29 methods touched across changed server-method files.

### Added

- `CommandsList(ctx, CommandsListParams) (json.RawMessage, error)` — list available commands for an agent (`commands.list`). Params: `AgentID`, `Provider`, `Scope`, `IncludeArgs`.
- `ExecApprovalList(ctx) (json.RawMessage, error)` — list pending exec approval requests (`exec.approval.list`). Requires `operator.approvals` scope.
- `MessageAction(ctx, MessageActionParams) (json.RawMessage, error)` — perform a message action through a channel (`message.action`). Required params: `Channel`, `Action`, `Params`, `IdempotencyKey`.
- `PluginApprovalList(ctx) (json.RawMessage, error)` — list pending plugin approval requests (`plugin.approval.list`). Requires `operator.approvals` scope.
- `SessionsCompactionList(ctx, SessionsCompactionListParams)` — list compaction checkpoints for a session (`sessions.compaction.list`).
- `SessionsCompactionGet(ctx, SessionsCompactionGetParams)` — retrieve a single compaction checkpoint (`sessions.compaction.get`).
- `SessionsCompactionBranch(ctx, SessionsCompactionBranchParams)` — create a new session branched from a compaction checkpoint (`sessions.compaction.branch`).
- `SessionsCompactionRestore(ctx, SessionsCompactionRestoreParams)` — restore a session to a compaction checkpoint (`sessions.compaction.restore`).
- `CommandsListParams` protocol type with `AgentID`, `Provider`, `Scope`, `IncludeArgs` fields
- `MessageActionParams` protocol type with `Channel`, `Action`, `Params`, `IdempotencyKey`, and optional fields
- `SessionsCompactionListParams`, `SessionsCompactionGetParams`, `SessionsCompactionBranchParams`, `SessionsCompactionRestoreParams` protocol types
- `MethodCommandsList`, `MethodExecApprovalList`, `MethodMessageAction`, `MethodPluginApprovalList`, `MethodSessionsCompactionBranch`, `MethodSessionsCompactionGet`, `MethodSessionsCompactionList`, `MethodSessionsCompactionRestore` method name constants

## [Unreleased]

### Breaking Changes

- `ChatSend` now returns `(*ChatSendResult, error)` instead of `(*ChatEvent, error)`. The RPC response is an ack `{runId, status}`, not a chat event. Streaming content arrives via the event handler. Callers must update: `result.Status` instead of `result.State`, and `ChatEvent` is now only for event stream payloads.
- `Presence` now returns `([]PresenceEntry, error)` instead of `(map[string]PresenceEntry, error)`. The upstream server returns a flat array, not a keyed map. Callers must iterate the slice instead of indexing by key.
- `CronList` now returns `(*CronListResult, error)` instead of `([]CronJob, error)`. The server returns a paginated object `{jobs, total, offset, limit, hasMore, nextOffset}`. Access jobs via `result.Jobs`.
- `CronRuns` now returns `(*CronRunsResult, error)` instead of `([]CronRunLogEntry, error)`. The server returns a paginated object `{entries, total, offset}`. Access entries via `result.Entries`.

### Added

- `ChatSendResult` type for the `chat.send` RPC ack response
- `CronListResult` and `CronRunsResult` paginated result types
- `CronRunLogUsage` type for per-run token usage tracking
- Extended `CronListParams` with `Limit`, `Offset`, `Query`, `Enabled`, `SortBy`, `SortDir`
- Extended `CronRunsParams` with `Scope`, `Offset`, `Statuses`, `DeliveryStatuses`, `Query`, `SortDir`
- Extended `CronRunLogEntry` with `Delivered`, `DeliveryStatus`, `Model`, `Provider`, `Usage`, `JobName`
- `PluginApprovalRequest`, `PluginApprovalResolve`, `PluginApprovalWaitDecision` methods and types
- `examples/internal/gwconn` shared helper for device identity auth in examples
- `docs/gateway-methods.md` comprehensive method reference (122 methods)
- `WithIdentity(id, deviceToken)` documented in gateway options

### Fixed

- All gateway examples now use device identity auth (`WithIdentity`) for real server compatibility. Without it, the upstream server clears self-declared scopes.
- All gateway examples accept CLI args (`<token> <host> [identity-dir]`) instead of hardcoded values
- Cron example uses `systemEvent` payload kind (required for main agent)
- Cron example captures server-assigned job UUID for subsequent operations
- CI E2E tests pass required args to example binaries

### Deprecated

- `CronListParams.IncludeDisabled` — use `Enabled` field instead (`"all"`, `"enabled"`, `"disabled"`)

## [v2026.4.2] - 2026-04-03

Sync to upstream openclaw v2026.4.2. Two new gateway RPC methods for ClawHub skill discovery.

### Added

- `SkillsSearch(ctx, SkillsSearchParams) (*SkillsSearchResult, error)` — search ClawHub for skills by query string and optional limit
- `SkillsDetail(ctx, SkillsDetailParams) (*SkillsDetailResult, error)` — fetch full detail for a skill by slug, including latest version, metadata, and owner
- `SkillsSearchParams`, `SkillsSearchItem`, `SkillsSearchResult` protocol types
- `SkillsDetailParams`, `SkillsDetailResult`, `SkillsDetailSkill`, `SkillsDetailVersion`, `SkillsDetailMetadata`, `SkillsDetailOwner` protocol types
- `MethodSkillsSearch` and `MethodSkillsDetail` method name constants

## [v2026.4.1] - 2026-04-01

Sync to upstream openclaw v2026.4.1. No new gateway RPC methods in this release; 23 methods touched across changed server-method files. Protocol type parity updated for cron and chat-history schema changes.

### Breaking Changes

- `CronPayload.Deliver`, `CronPayload.Channel`, `CronPayload.To`, `CronPayload.BestEffortDeliver` removed — the upstream schema dropped these fields from the `agentTurn` payload kind (server enforces `additionalProperties: false`). Remove these fields from any `CronPayload{Kind: "agentTurn", ...}` literals.

### Added

- `ChatHistoryParams.MaxChars` optional int field — limits history response size in characters (upstream `maxChars`, max 500,000)
- `CronPayload.LightContext` optional bool field — upstream parity for `agentTurn` payload

## [v2026.3.31] - 2026-04-01

Sync to upstream openclaw v2026.3.31. No new gateway RPC methods in this release; 39 methods touched across 16 changed server-method files. Protocol type parity updated for schema changes.

### Added

- `AgentSummaryModel` protocol type — optional model override (`primary`, `fallbacks`) exposed on `AgentSummary`
- `AgentSummary.Workspace` and `AgentSummary.Model` fields (upstream `agents-models-skills.ts` schema update)
- `SkillsInstallParams.DangerouslyForceUnsafeInstall` optional boolean field (upstream skills schema update)

## [v2026.3.28] - 2026-03-29

Sync to upstream openclaw v2026.3.28. No new gateway RPC methods in this release; protocol parity maintained. The `plugin.approval.*` methods added in the v2026.3.24 sync cycle cover the new upstream plugin approval hook surface.

### Changed

- No wire-protocol changes required. This release consists of provider integrations (xAI Responses API, MiniMax image-01), plugin SDK hooks (`requireApproval` in `before_tool_call`), channel fixes, and numerous bug fixes in the upstream gateway. The Go SDK has full parity for this version.

## [v2026.3.24] - 2026-03-26

Sync to upstream openclaw v2026.3.24. 1539 upstream files changed; 1 RPC method added (`tools.effective`); 30 methods touched.

### Added

- `ToolsEffective(ctx, ToolsEffectiveParams)` — runtime-effective tool inventory for a session (`tools.effective` RPC). `SessionKey` is required; `AgentID` is optional and must match the session's agent when provided. Requires `operator.read` scope.
- `ToolsEffectiveParams` protocol type with `SessionKey` and `AgentID` fields
- `SessionsGet(ctx, SessionsGetParams)` — retrieve session messages by key (`sessions.get` RPC, openclaw v2026.3.7)
- `SessionsGetParams` protocol type with `Key`, `SessionKey`, and `Limit` fields
- `NodePendingEnqueue(ctx, NodePendingEnqueueParams)` and `NodePendingDrain(ctx, NodePendingDrainParams)` for pending node work queue operations (`node.pending.enqueue`, `node.pending.drain`; openclaw v2026.3.11)
- `NodePendingEnqueueParams`, `NodePendingDrainParams`, `NodePendingWorkItem`, `NodePendingEnqueueResult`, and `NodePendingDrainResult` protocol types
- `TalkSpeak(ctx, TalkSpeakParams)` typed method and `TalkSpeakParams`/`TalkSpeakResult` protocol types for voice synthesis via gateway (`talk.speak`; openclaw v2026.3.22)

## [1.0.0] - 2026-02-20

Initial public release of `openclaw-go`, a Go client library for the
[OpenClaw](https://openclaw.ai) AI gateway.

Single external dependency: [gorilla/websocket v1.5.3](https://github.com/gorilla/websocket).
All library packages target 100% statement coverage and pass the Go race detector.

### Added

#### `protocol` — Gateway Wire Types

- Complete set of Go types for the OpenClaw Gateway WebSocket protocol (v3)
- Frame types: `Request`, `Response`, `Event`, `Invoke`, `InvokeResponse`
- Handshake types: `ConnectChallenge`, `ConnectParams`, `HelloOK`, `Snapshot`, `HelloPolicy`
- All role and scope constants: `RoleOperator`, `RoleNode`, `ScopeOperatorRead`, `ScopeOperatorWrite`, `ScopeOperatorAdmin`, `ScopeOperatorApprovals`, `ScopeOperatorPairing`
- Chat types: `ChatSendParams`, `ChatHistoryParams`, `ChatAbortParams`, `ChatInjectParams`, `ChatEvent`
- Agent types: `AgentParams`, `AgentWaitParams`, `AgentIdentityParams`, `AgentEvent`
- Session types: `SessionsListParams`, `SessionsPreviewParams`, `SessionsPatchParams`, `SessionsResetParams`, `SessionsDeleteParams`, `SessionsUsageParams`, `SessionsCompactParams`, `SessionsResolveParams`
- Agents CRUD types: `AgentsCreateParams/Result`, `AgentsUpdateParams/Result`, `AgentsDeleteParams/Result`, `AgentsFilesListResult`, `AgentsFilesGetParams/Result`, `AgentsFilesSetParams/Result`
- Exec approval types: `ExecApprovalRequestParams/Result`, `ExecApprovalResolveParams/Result`, `ExecApprovalWaitDecisionParams/Result`, `ExecApprovalsSnapshot`, `ExecApprovalsFile`
- Config types: `ConfigGetParams`, `ConfigSetParams`, `ConfigPatchParams`, `ConfigApplyParams`, `ConfigSchemaResponse`
- Node types: `NodePairRequestParams`, `NodePairApproveParams`, `NodePairRejectParams`, `NodePairVerifyParams`, `NodeDescribeParams`, `NodeInvokeParams`, `NodeInvokeResultParams`, `NodeEventParams`
- Device pairing types: `DevicePairApproveParams`, `DevicePairRejectParams`, `DevicePairRemoveParams`, `DeviceTokenRotateParams`, `DeviceTokenRevokeParams`
- Cron types: `CronAddParams`, `CronUpdateParams`, `CronRemoveParams`, `CronRunParams`, `CronRunsParams`, `CronJob`, `CronRunLogEntry`, `CronSchedule`, `CronPayload`
- TTS types: `TTSConvertParams/Result`, `TTSSetProviderParams/Result`, `TTSStatusResult`, `TTSProvidersResult`, `TTSEnableResult`, `TTSDisableResult`
- Channel and talk types: `ChannelsStatusParams/Result`, `TalkModeParams`, `TalkConfigParams/Result`
- Skills types: `SkillsStatusParams`, `SkillsBinsResult`, `SkillsInstallParams`, `SkillsUpdateParams`
- Wizard types: `WizardStartParams/Result`, `WizardNextParams/Result`, `WizardCancelParams`, `WizardStatusParams/Result`, `WizardStep`
- Misc types: `SendParams`, `WakeParams`, `PollParams`, `PushTestParams/Result`, `PresenceEntry`, `ModelsListResult`, `LogsTailParams/Result`
- Serialization helpers: `MarshalEvent`, `MarshalResponse`, `MarshalErrorResponse`, `ParseFrame`
- Server constants: `ProtocolVersion` (3), `MaxPayloadBytes`, `MaxBufferedBytes`, `DefaultTickIntervalMs`, `DedupeTTLMs`, `SessionLabelMaxLength`
- Client ID and mode constants for all first-party clients

#### `gateway` — WebSocket Gateway Client

- `Client` struct with full connection lifecycle: challenge→connect→hello-ok handshake
- Background read loop with frame dispatch (requests, events, invocations)
- Keepalive tick loop using the gateway's `TickIntervalMs` policy
- Pending request/response correlation with atomic request ID counter
- Functional options: `WithToken`, `WithPassword`, `WithRole`, `WithScopes`, `WithCaps`, `WithCommands`, `WithPermissions`, `WithDevice`, `WithLocale`, `WithUserAgent`, `WithTLSConfig`, `WithConnectTimeout`, `WithOnEvent`, `WithOnInvoke`, `WithClientInfo`
- `Connect(ctx, wsURL)` — connects and completes the handshake
- `Close()` — graceful shutdown; `Done()` channel signals completion
- `Hello()` — returns the `HelloOK` payload from the server
- `Send(ctx, method, params)` — low-level typed RPC send
- `SendEvent(eventName, payload)` — send a uni-directional event frame
- **96+ typed RPC methods**, organized by domain:
  - Chat: `ChatSend`, `ChatHistory`, `ChatAbort`, `ChatInject`
  - Agent: `Agent`, `AgentIdentity`, `AgentWait`
  - Sessions: `SessionsList`, `SessionsPreview`, `SessionsResolve`, `SessionsPatch`, `SessionsReset`, `SessionsDelete`, `SessionsCompact`, `SessionsUsage`
  - Agents CRUD: `AgentsList`, `AgentsCreate`, `AgentsUpdate`, `AgentsDelete`, `AgentsFilesList`, `AgentsFilesGet`, `AgentsFilesSet`
  - Config: `ConfigGet`, `ConfigSet`, `ConfigPatch`, `ConfigApply`, `ConfigSchema`
  - Exec approvals: `ExecApprovalRequest`, `ExecApprovalResolve`, `ExecApprovalWaitDecision`, `ExecApprovalsGet`, `ExecApprovalsSet`, `ExecApprovalsNodeGet`, `ExecApprovalsNodeSet`
  - Nodes: `NodeList`, `NodeDescribe`, `NodeInvoke`, `NodeInvokeResult`, `NodeEvent`, `NodeRename`
  - Node pairing: `NodePairRequest`, `NodePairList`, `NodePairApprove`, `NodePairReject`, `NodePairVerify`
  - Device pairing: `DevicePairList`, `DevicePairApprove`, `DevicePairReject`, `DevicePairRemove`, `DeviceTokenRotate`, `DeviceTokenRevoke`
  - Cron: `CronList`, `CronStatus`, `CronAdd`, `CronUpdate`, `CronRemove`, `CronRun`, `CronRuns`
  - TTS: `TTSStatus`, `TTSProviders`, `TTSEnable`, `TTSDisable`, `TTSConvert`, `TTSSetProvider`
  - Channels: `ChannelsStatus`, `ChannelsLogout`, `TalkConfig`, `TalkMode`
  - Models: `ModelsList`
  - Logs: `LogsTail`
  - Skills: `SkillsStatus`, `SkillsBins`, `SkillsInstall`, `SkillsUpdate`
  - Wizard: `WizardStart`, `WizardNext`, `WizardCancel`, `WizardStatus`
  - Presence: `Presence`
  - Misc: `Health`, `Status`, `SendMessage`, `Wake`, `LastHeartbeat`, `SetHeartbeats`, `SystemEvent`, `UpdateRun`, `PushTest`, `BrowserRequest`, `VoiceWakeGet`, `VoiceWakeSet`, `UsageStatus`, `UsageCost`, `Poll`
- Concurrent-safe: mutex-protected writes, atomic request counters, synchronized pending-response map

#### `chatcompletions` — OpenAI-compatible Chat Completions Client

- `Client` struct with `BaseURL`, `Token`, `AgentID`, `SessionKey`, `HTTPClient` fields
- `Create(ctx, Request) (*Response, error)` — non-streaming chat completion
- `CreateStream(ctx, Request) (*Stream, error)` — server-sent events streaming completion
- `Stream.Recv() (*Chunk, error)` — read next SSE chunk; returns `io.EOF` at end
- `Stream.Close()` — close the underlying HTTP response body
- Request type with `Model`, `Messages`, `MaxTokens`, `Temperature`, `Stream` fields
- Response type with `Choices`, `Usage` (prompt/completion/total tokens), `Model`, `ID`
- Streaming chunk type with `Choices[].Delta.Content`
- `Authorization: Bearer` token header support
- `x-openclaw-agent-id` and `x-openclaw-session-key` custom header support

#### `openresponses` — OpenAI Responses API Client

- `Client` struct with `BaseURL`, `Token`, `AgentID`, `SessionKey`, `HTTPClient` fields
- `Create(ctx, Request) (*Response, error)` — non-streaming responses request
- `CreateStream(ctx, Request) (*Stream, error)` — SSE streaming with typed events
- `Stream.Recv() (*StreamEvent, error)` — read next typed SSE event; returns `io.EOF` at end
- `Stream.Close()` — close the underlying HTTP response body
- `InputFromItems([]InputItem) Input` — helper to build structured input
- `MessageItem(role, text string) InputItem` — helper to build a text message item
- Structured input items: text messages, images, audio, resource links, function call outputs
- Tool definitions: `ToolDefinition` with `FunctionTool` (name, description, parameters schema)
- Response type with `Output` items, `Usage` (input/output/total tokens), `Status`, `ID`
- Typed SSE event types: `response.created`, `response.output_text.delta`, `response.output_text.done`, `response.function_call_arguments.delta`, `response.completed`, and more
- `Authorization: Bearer` token, `x-openclaw-agent-id`, `x-openclaw-session-key` header support

#### `toolsinvoke` — Tools Invoke HTTP Client

- `Client` struct with `BaseURL`, `Token`, `MessageChannel`, `AccountID`, `HTTPClient` fields
- `Invoke(ctx, Request) (*Response, error)` — calls `POST /tools/invoke`
- Request type with `Tool`, `Action`, `Args`, `SessionKey`, `DryRun` fields
- Response type with `OK`, `Result` (`json.RawMessage`), `Error` (`*ErrorDetail`)
- `ErrorDetail` with `Type` and `Message` fields
- `Authorization: Bearer` token, `x-openclaw-message-channel`, `x-openclaw-account-id` header support

#### `discovery` — mDNS/DNS-SD Gateway Discovery

- `Browser` struct for discovering OpenClaw Gateway instances on the local network
- `NewBrowser() *Browser` — create a new browser
- `Browser.Browse(ctx) ([]Beacon, error)` — scan for gateways, returns after context deadline
- `Beacon` struct with all mDNS TXT record fields: `Host`, `Port`, `DisplayName`, `LanHost`, `TailnetDNS`, `GatewayPort`, `GatewayTLS`, `GatewayTLSFingerprint`, `SSHPort`, `CLIPath`, `Role`, `CanvasPort`
- `Beacon.WebSocketURL() string` — derives `ws://` or `wss://` URL with host priority: SRV host → TailnetDNS → LanHost → `127.0.0.1`
- `Beacon.HTTPURL() string` — derives `http://` or `https://` base URL
- macOS implementation via `dns-sd -B` / `dns-sd -L` (Bonjour)
- Linux implementation via `avahi-browse`
- Windows: returns an informative `not supported` error
- `ServiceType` constant: `_openclaw-gw._tcp`

#### `acp` — Agent Client Protocol Server

- `Server` struct for handling ACP (Agent Client Protocol) over stdio
- `NewServer(handler Handler, r io.Reader, w io.Writer) *Server` — create a new ACP server
- `Server.Serve(ctx) error` — run the JSON-RPC 2.0 / NDJSON dispatch loop
- `Server.Notify(method string, params any) error` — send a server-initiated notification
- `Handler` interface with all ACP methods:
  - `Initialize`, `Authenticate`
  - `NewSession`, `LoadSession`, `ListSessions`, `ForkSession`, `ResumeSession`
  - `Prompt`, `Cancel`
  - `SetSessionMode`, `SetSessionModel`, `SetSessionConfigOption`
- Full ACP type coverage: `InitializeRequest/Response`, `AuthenticateRequest/Response`, `NewSessionRequest/Response`, `LoadSessionRequest/Response`, `ListSessionsRequest/Response`, `ForkSessionRequest/Response`, `ResumeSessionRequest/Response`, `PromptRequest/Response`, `CancelNotification`, `SetSessionModeRequest/Response`, `SetSessionModelRequest/Response`, `SetSessionConfigOptionRequest/Response`
- `AgentCapabilities` with `LoadSession`, `PromptCapabilities` (image, embedded context), `SessionCapabilities` (list, fork, resume)
- Session update notification support via `Notify`
- `ProtocolVersion` constant: `1`
- Standard and extended JSON-RPC 2.0 error codes

#### Examples

Twelve runnable examples in `examples/`:

- `server` — mock OpenClaw Gateway for local development (WebSocket + HTTP)
- `client` — demonstrates all three APIs: WebSocket, Chat Completions, Tools Invoke
- `chat` — interactive chat session with streaming event handling
- `openresponses` — OpenAI Responses API with tool definitions and SSE streaming
- `agents` — agent CRUD: list, create, update, files, delete
- `sessions` — session management: list, preview, patch, usage, reset
- `approvals` — exec approval flow: listen for events, approve/reject, admin config
- `pairing` — node and device pairing workflows
- `config` — gateway configuration: get, schema, patch, apply
- `cron` — cron job management: list, add, run, history, remove
- `node` — connect as a node: declare capabilities, handle invocations, send events
- `discovery` — scan the LAN for OpenClaw gateways via mDNS
- `acp` — ACP agent server implementation over stdio

[1.0.0]: https://github.com/a3tai/openclaw-go/releases/tag/v1.0.0
