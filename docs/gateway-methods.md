# Gateway Method Reference

Complete reference for all OpenClaw Gateway WebSocket RPC methods as of upstream v2026.3.24.

Each method is available as a typed Go function on `*gateway.Client`. Methods are grouped by domain. Return types marked `json.RawMessage` return raw JSON — callers unmarshal into their own types.

> **Device identity required.** Real gateway servers clear self-declared scopes for clients without a signed device identity. Use `gateway.WithIdentity(id, deviceToken)` — see the [gateway overview](gateway.md) for setup.

## Chat

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `ChatSend` | `chat.send` | `ChatSendParams` | `*ChatSendResult{RunID, Status}` | `operator.write` |
| `ChatHistory` | `chat.history` | `ChatHistoryParams` | `json.RawMessage` | `operator.read` |
| `ChatAbort` | `chat.abort` | `ChatAbortParams` | `error` | `operator.write` |
| `ChatInject` | `chat.inject` | `ChatInjectParams` | `error` | `operator.write` |

`ChatSend` returns an initial ack `{runId, status: "started"}`. Streaming content arrives via `"chat"` events on the event handler. The `ChatEvent` type is for those events, not the RPC response.

## Sessions

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `SessionsList` | `sessions.list` | `SessionsListParams` | `json.RawMessage` | `operator.read` |
| `SessionsGet` | `sessions.get` | `SessionsGetParams` | `json.RawMessage` | `operator.read` |
| `SessionsPreview` | `sessions.preview` | `SessionsPreviewParams` | `json.RawMessage` | `operator.read` |
| `SessionsResolve` | `sessions.resolve` | `SessionsResolveParams` | `json.RawMessage` | `operator.read` |
| `SessionsCreate` | `sessions.create` | `SessionsCreateParams` | `json.RawMessage` | `operator.write` |
| `SessionsSend` | `sessions.send` | `SessionsSendParams` | `json.RawMessage` | `operator.write` |
| `SessionsSteer` | `sessions.steer` | `SessionsSteerParams` | `json.RawMessage` | `operator.write` |
| `SessionsAbort` | `sessions.abort` | `SessionsAbortParams` | `json.RawMessage` | `operator.write` |
| `SessionsSubscribe` | `sessions.subscribe` | `SessionsSubscribeParams` | `json.RawMessage` | `operator.read` |
| `SessionsUnsubscribe` | `sessions.unsubscribe` | `SessionsUnsubscribeParams` | `json.RawMessage` | `operator.read` |
| `SessionsMessagesSubscribe` | `sessions.messages.subscribe` | `SessionsMessagesSubscribeParams` | `json.RawMessage` | `operator.read` |
| `SessionsMessagesUnsubscribe` | `sessions.messages.unsubscribe` | `SessionsMessagesUnsubscribeParams` | `json.RawMessage` | `operator.read` |
| `SessionsPatch` | `sessions.patch` | `SessionsPatchParams` | `error` | `operator.admin` |
| `SessionsReset` | `sessions.reset` | `SessionsResetParams` | `error` | `operator.admin` |
| `SessionsDelete` | `sessions.delete` | `SessionsDeleteParams` | `error` | `operator.admin` |
| `SessionsCompact` | `sessions.compact` | `SessionsCompactParams` | `error` | `operator.admin` |
| `SessionsUsage` | `sessions.usage` | `SessionsUsageParams` | `json.RawMessage` | `operator.read` |
| `SessionsUsageTimeseries` | `sessions.usage.timeseries` | `SessionsUsageTimeseriesParams` | `json.RawMessage` | `operator.read` |
| `SessionsUsageLogs` | `sessions.usage.logs` | `SessionsUsageLogsParams` | `json.RawMessage` | `operator.read` |

### Upstream response shapes

- `sessions.list` returns `{sessions: SessionRow[], total: number}`
- `sessions.get` returns `{messages: Message[]}`
- `sessions.preview` returns `{ts: number, previews: SessionPreviewEntry[]}`
- `sessions.create` returns `{ok, key, sessionId, entry, runStarted, runId?, status?}`
- `sessions.send` returns `{runId, status: "started", messageSeq?}`
- `sessions.abort` returns `{ok, abortedRunId, status}`
- `sessions.subscribe` returns `{subscribed: boolean}`
- `sessions.usage` returns `{updatedAt, startDate, endDate, sessions: [...], totals, aggregates}`

## Agents CRUD

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `AgentsList` | `agents.list` | (none) | `*AgentsListResult{Agents, DefaultID, MainKey, Scope}` | `operator.read` |
| `AgentsCreate` | `agents.create` | `AgentsCreateParams` | `*AgentsCreateResult{OK, AgentID, Name, Workspace}` | `operator.admin` |
| `AgentsUpdate` | `agents.update` | `AgentsUpdateParams` | `error` | `operator.admin` |
| `AgentsDelete` | `agents.delete` | `AgentsDeleteParams` | `*AgentsDeleteResult{OK, AgentID, RemovedBindings}` | `operator.admin` |
| `AgentsFilesList` | `agents.files.list` | `AgentsFilesListParams` | `*AgentsFilesListResult{AgentID, Workspace, Files}` | `operator.read` |
| `AgentsFilesGet` | `agents.files.get` | `AgentsFilesGetParams` | `*AgentsFilesGetResult{AgentID, Workspace, File}` | `operator.read` |
| `AgentsFilesSet` | `agents.files.set` | `AgentsFilesSetParams` | `*AgentsFilesSetResult{OK, AgentID, Workspace, File}` | `operator.admin` |

## Agent

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `Agent` | `agent` | `AgentParams` | `json.RawMessage` | `operator.read` |
| `AgentIdentity` | `agent.identity.get` | `AgentIdentityParams` | `*AgentIdentityResult{AgentID, Name, Avatar, Emoji}` | `operator.read` |
| `AgentWait` | `agent.wait` | `AgentWaitParams` | `json.RawMessage` | `operator.read` |

## Config

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `ConfigGet` | `config.get` | (none) | `json.RawMessage` | `operator.read` |
| `ConfigSet` | `config.set` | `ConfigSetParams` | `error` | `operator.admin` |
| `ConfigApply` | `config.apply` | `ConfigApplyParams` | `error` | `operator.admin` |
| `ConfigPatch` | `config.patch` | `ConfigPatchParams` | `error` | `operator.admin` |
| `ConfigSchema` | `config.schema` | (none) | `*ConfigSchemaResponse{Schema, UIHints, Version, GeneratedAt}` | `operator.read` |
| `ConfigSchemaLookup` | `config.schema.lookup` | `ConfigSchemaLookupParams` | `*ConfigSchemaLookupResult{Path, Schema, Hint, HintPath, Children}` | `operator.read` |

### Upstream response shapes

- `config.get` returns `{raw: string, path: string, config: RedactedConfig, valid: boolean, hash: string}`

## Cron

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `CronList` | `cron.list` | `CronListParams` | `*CronListResult{Jobs, Total, Offset, Limit, HasMore, NextOffset}` | `operator.read` |
| `CronStatus` | `cron.status` | (none) | `json.RawMessage` | `operator.read` |
| `CronAdd` | `cron.add` | `CronAddParams` | `json.RawMessage` (returns full CronJob) | `operator.admin` |
| `CronUpdate` | `cron.update` | `CronUpdateParams` | `error` | `operator.admin` |
| `CronRemove` | `cron.remove` | `CronRemoveParams` | `error` | `operator.admin` |
| `CronRun` | `cron.run` | `CronRunParams` | `error` | `operator.admin` |
| `CronRuns` | `cron.runs` | `CronRunsParams` | `*CronRunsResult{Entries, Total, Offset}` | `operator.read` |

`CronList` and `CronRuns` return paginated results. Use `Offset`/`Limit` params to page.

## Exec Approvals

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `ExecApprovalRequest` | `exec.approval.request` | `ExecApprovalRequestParams` | `*ExecApprovalRequestResult{ID, Status, Decision, CreatedAtMs, ExpiresAtMs}` | `operator.write` |
| `ExecApprovalResolve` | `exec.approval.resolve` | `ExecApprovalResolveParams` | `*ExecApprovalResolveResult{OK}` | `operator.approvals` |
| `ExecApprovalWaitDecision` | `exec.approval.waitDecision` | `ExecApprovalWaitDecisionParams` | `*ExecApprovalWaitDecisionResult{ID, Decision}` | `operator.write` |
| `ExecApprovalsGet` | `exec.approvals.get` | (none) | `*ExecApprovalsSnapshot{Path, Exists, Hash, File}` | `operator.read` |
| `ExecApprovalsSet` | `exec.approvals.set` | `ExecApprovalsSetParams` | `error` | `operator.admin` |
| `ExecApprovalsNodeGet` | `exec.approvals.node.get` | `ExecApprovalsNodeGetParams` | `*ExecApprovalsSnapshot` | `operator.read` |
| `ExecApprovalsNodeSet` | `exec.approvals.node.set` | `ExecApprovalsNodeSetParams` | `error` | `operator.admin` |

## Plugin Approvals

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `PluginApprovalRequest` | `plugin.approval.request` | `PluginApprovalRequestParams` | `*PluginApprovalRequestResult` | `operator.write` |
| `PluginApprovalResolve` | `plugin.approval.resolve` | `PluginApprovalResolveParams` | `*PluginApprovalResolveResult` | `operator.approvals` |

## Nodes

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `NodeList` | `node.list` | (none) | `json.RawMessage` | `operator.read` |
| `NodeDescribe` | `node.describe` | `NodeDescribeParams` | `json.RawMessage` | `operator.read` |
| `NodeInvoke` | `node.invoke` | `NodeInvokeParams` | `json.RawMessage` | `operator.write` |
| `NodeInvokeResult` | `node.invoke.result` | `NodeInvokeResultParams` | `error` | (node role) |
| `NodeEvent` | `node.event` | `NodeEventParams` | `error` | (node role) |
| `NodeRename` | `node.rename` | `NodeRenameParams` | `error` | `operator.admin` |
| `NodePendingEnqueue` | `node.pending.enqueue` | `NodePendingEnqueueParams` | `*NodePendingEnqueueResult{NodeID, Revision, Queued, WakeTriggered}` | `operator.admin` |
| `NodePendingDrain` | `node.pending.drain` | `NodePendingDrainParams` | `*NodePendingDrainResult{NodeID, Revision, Items, HasMore}` | (node role) |
| `NodePendingPull` | `node.pending.pull` | `NodePendingPullParams` | `*NodePendingPullResult{NodeID, Actions}` | (node role) |
| `NodePendingAck` | `node.pending.ack` | `NodePendingAckParams` | `*NodePendingAckResult{NodeID, AckedIDs, RemainingCount}` | (node role) |
| `NodeCanvasCapabilityRefresh` | `node.canvas.capability.refresh` | `NodeCanvasCapabilityRefreshParams` | `*NodeCanvasCapabilityRefreshResult` | (node role) |

### Upstream response shapes

- `node.list` returns `{ts: number, nodes: NodeEntry[]}`
- `node.describe` returns full node descriptor object
- `node.invoke` returns `{ok, nodeId, command, payload, payloadJSON}`

## Node Pairing

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `NodePairRequest` | `node.pair.request` | `NodePairRequestParams` | `json.RawMessage` | (node role) |
| `NodePairList` | `node.pair.list` | (none) | `json.RawMessage` | `operator.read` |
| `NodePairApprove` | `node.pair.approve` | `NodePairApproveParams` | `error` | `operator.admin` |
| `NodePairReject` | `node.pair.reject` | `NodePairRejectParams` | `error` | `operator.admin` |
| `NodePairVerify` | `node.pair.verify` | `NodePairVerifyParams` | `json.RawMessage` | (node role) |

## Device Pairing

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `DevicePairList` | `device.pair.list` | (none) | `json.RawMessage` | `operator.read` |
| `DevicePairApprove` | `device.pair.approve` | `DevicePairApproveParams` | `error` | `operator.admin` |
| `DevicePairReject` | `device.pair.reject` | `DevicePairRejectParams` | `error` | `operator.admin` |
| `DevicePairRemove` | `device.pair.remove` | `DevicePairRemoveParams` | `error` | `operator.admin` |
| `DeviceTokenRotate` | `device.token.rotate` | `DeviceTokenRotateParams` | `json.RawMessage` | `operator.admin` |
| `DeviceTokenRevoke` | `device.token.revoke` | `DeviceTokenRevokeParams` | `error` | `operator.admin` |

### Upstream response shapes

- `device.pair.list` returns `{pending: PendingEntry[], paired: RedactedDevice[]}`
- `device.token.rotate` returns `{deviceId, role, token, scopes, rotatedAtMs}`

## Channels & Talk

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `ChannelsStatus` | `channels.status` | `ChannelsStatusParams` | `*ChannelsStatusResult{Ts, Channels, ChannelAccounts, ...}` | `operator.read` |
| `ChannelsLogout` | `channels.logout` | `ChannelsLogoutParams` | `error` | `operator.admin` |
| `TalkConfig` | `talk.config` | `TalkConfigParams` | `*TalkConfigResult{Config{Talk, Session, UI}}` | `operator.read` |
| `TalkMode` | `talk.mode` | `TalkModeParams` | `error` | `operator.write` |
| `TalkSpeak` | `talk.speak` | `TalkSpeakParams` | `*TalkSpeakResult{AudioBase64, Provider, MimeType, ...}` | `operator.write` |
| `WebLoginStart` | `web.login.start` | `WebLoginStartParams` | `json.RawMessage` | `operator.write` |
| `WebLoginWait` | `web.login.wait` | `WebLoginWaitParams` | `json.RawMessage` | `operator.write` |

## TTS

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `TTSStatus` | `tts.status` | (none) | `*TTSStatusResult{Enabled, Provider, FallbackProviders, ...}` | `operator.read` |
| `TTSProviders` | `tts.providers` | (none) | `*TTSProvidersResult{Providers, Active}` | `operator.read` |
| `TTSEnable` | `tts.enable` | (none) | `*TTSEnableResult{Enabled}` | `operator.admin` |
| `TTSDisable` | `tts.disable` | (none) | `*TTSDisableResult{Enabled}` | `operator.admin` |
| `TTSConvert` | `tts.convert` | `TTSConvertParams` | `*TTSConvertResult{AudioPath, Provider, OutputFormat, VoiceCompatible}` | `operator.write` |
| `TTSSetProvider` | `tts.setProvider` | `TTSSetProviderParams` | `*TTSSetProviderResult{Provider}` | `operator.admin` |

## Models

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `ModelsList` | `models.list` | (none) | `*ModelsListResult{Models}` | `operator.read` |

## Logs

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `LogsTail` | `logs.tail` | `LogsTailParams` | `*LogsTailResult{File, Cursor, Size, Lines, Truncated, Reset}` | `operator.read` |

## Health & System

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `Health` | `health` | (none) | `json.RawMessage` | (none) |
| `Status` | `status` | (none) | `json.RawMessage` | `operator.read` |
| `Presence` | `system-presence` | (none) | `[]PresenceEntry` | `operator.read` |
| `GatewayIdentityGet` | `gateway.identity.get` | (none) | `*GatewayIdentityResult{DeviceID, PublicKey}` | `operator.read` |
| `DoctorMemoryStatus` | `doctor.memory.status` | `DoctorMemoryStatusParams` | `*DoctorMemoryStatusResult` | `operator.read` |

## Skills

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `SkillsStatus` | `skills.status` | `SkillsStatusParams` | `json.RawMessage` | `operator.read` |
| `SkillsBins` | `skills.bins` | (none) | `*SkillsBinsResult{Bins}` | `operator.read` |
| `SkillsInstall` | `skills.install` | `SkillsInstallParams` | `json.RawMessage` | `operator.admin` |
| `SkillsUpdate` | `skills.update` | `SkillsUpdateParams` | `error` | `operator.admin` |

## Tools

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `ToolsCatalog` | `tools.catalog` | `ToolsCatalogParams` | `json.RawMessage` | `operator.read` |
| `ToolsEffective` | `tools.effective` | `ToolsEffectiveParams` | `json.RawMessage` | `operator.read` |

### Upstream response shapes

- `tools.catalog` returns `{agentId, profiles: ProfileOption[], groups: ToolCatalogGroup[]}`
- `tools.effective` returns effective tool inventory keyed by agent/session context

## Secrets

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `SecretsReload` | `secrets.reload` | (none) | `error` | `operator.admin` |
| `SecretsResolve` | `secrets.resolve` | `SecretsResolveParams` | `*SecretsResolveResult{OK, Assignments, Diagnostics, InactiveRefPaths}` | `operator.admin` |

## Wizard

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `WizardStart` | `wizard.start` | `WizardStartParams` | `*WizardStartResult{SessionID, Done, Step, Status, Error}` | `operator.admin` |
| `WizardNext` | `wizard.next` | `WizardNextParams` | `*WizardNextResult{Done, Step, Status, Error}` | `operator.admin` |
| `WizardCancel` | `wizard.cancel` | `WizardCancelParams` | `error` | `operator.admin` |
| `WizardStatus` | `wizard.status` | (none) | `*WizardStatusResult{Status, Error}` | `operator.admin` |

## Messaging & Misc

| Method | Protocol | Params | Result | Scope |
|--------|----------|--------|--------|-------|
| `SendMessage` | `send` | `SendParams` | `json.RawMessage` | `operator.write` |
| `Wake` | `wake` | `WakeParams` | `error` | `operator.write` |
| `LastHeartbeat` | `last-heartbeat` | (none) | `json.RawMessage` | `operator.read` |
| `SetHeartbeats` | `set-heartbeats` | `SetHeartbeatsParams` | `error` | `operator.write` |
| `SystemEvent` | `system-event` | `SystemEventParams` | `error` | `operator.write` |
| `UpdateRun` | `update.run` | `UpdateRunParams` | `json.RawMessage` | `operator.admin` |
| `PushTest` | `push.test` | `PushTestParams` | `*PushTestResult{OK, Status, ApnsID, Reason, ...}` | `operator.admin` |
| `BrowserRequest` | `browser.request` | `BrowserRequestParams` | `json.RawMessage` | `operator.write` |
| `VoiceWakeGet` | `voicewake.get` | (none) | `json.RawMessage` | `operator.read` |
| `VoiceWakeSet` | `voicewake.set` | `VoiceWakeSetParams` | `error` | `operator.admin` |
| `UsageStatus` | `usage.status` | (none) | `json.RawMessage` | `operator.read` |
| `UsageCost` | `usage.cost` | `UsageCostParams` | `json.RawMessage` | `operator.read` |
| `Poll` | `poll` | `PollParams` | `json.RawMessage` | `operator.write` |

## Scope Reference

| Scope | Grants |
|-------|--------|
| `operator.read` | Read-only access to sessions, config, agents, nodes, cron, models, logs, health, presence |
| `operator.write` | Send messages, chat, invoke nodes, approval requests, talk, system events |
| `operator.admin` | Config changes, session management, cron management, agent CRUD, node management, TTS, secrets, wizard |
| `operator.approvals` | Resolve exec/plugin approval decisions |
| `operator.pairing` | Device and node pairing operations |

## Audit Status

Full audit completed against upstream `openclaw/openclaw` at `v2026.3.24`. All 122 typed methods verified against server handler response shapes. See [PR #56](https://github.com/a3tai/openclaw-go/pull/56) for the audit details and fixes applied.
