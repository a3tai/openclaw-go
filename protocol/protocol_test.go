package protocol

import (
	"encoding/json"
	"testing"
)

// unmarshallable is a type that cannot be marshaled to JSON.
type unmarshallable struct{}

func (unmarshallable) MarshalJSON() ([]byte, error) {
	return nil, &json.UnsupportedTypeError{}
}

func TestMarshalRequest(t *testing.T) {
	data, err := MarshalRequest("req-1", "connect", map[string]string{"key": "val"})
	if err != nil {
		t.Fatal(err)
	}
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		t.Fatal(err)
	}
	if req.Type != FrameTypeRequest {
		t.Errorf("type = %q, want %q", req.Type, FrameTypeRequest)
	}
	if req.ID != "req-1" {
		t.Errorf("id = %q, want %q", req.ID, "req-1")
	}
	if req.Method != "connect" {
		t.Errorf("method = %q, want %q", req.Method, "connect")
	}
}

func TestMarshalRequestError(t *testing.T) {
	_, err := MarshalRequest("req-1", "test", unmarshallable{})
	if err == nil {
		t.Fatal("expected error marshaling unmarshallable params")
	}
}

func TestMarshalResponse(t *testing.T) {
	data, err := MarshalResponse("req-1", map[string]string{"status": "ok"})
	if err != nil {
		t.Fatal(err)
	}
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.OK {
		t.Error("ok = false, want true")
	}
	if resp.ID != "req-1" {
		t.Errorf("id = %q, want %q", resp.ID, "req-1")
	}
}

func TestMarshalResponseError(t *testing.T) {
	_, err := MarshalResponse("req-1", unmarshallable{})
	if err == nil {
		t.Fatal("expected error marshaling unmarshallable payload")
	}
}

func TestMarshalErrorResponse(t *testing.T) {
	data, err := MarshalErrorResponse("req-2", ErrorPayload{
		Code:    "AUTH_FAILED",
		Message: "invalid token",
	})
	if err != nil {
		t.Fatal(err)
	}
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.OK {
		t.Error("ok = true, want false")
	}
	if resp.Error == nil {
		t.Fatal("error is nil")
	}
	if resp.Error.Code != "AUTH_FAILED" {
		t.Errorf("error.code = %q, want %q", resp.Error.Code, "AUTH_FAILED")
	}
}

func TestMarshalEvent(t *testing.T) {
	data, err := MarshalEvent("connect.challenge", ConnectChallenge{
		Nonce: "abc123",
		Ts:    1737264000000,
	})
	if err != nil {
		t.Fatal(err)
	}
	var ev Event
	if err := json.Unmarshal(data, &ev); err != nil {
		t.Fatal(err)
	}
	if ev.Type != FrameTypeEvent {
		t.Errorf("type = %q, want %q", ev.Type, FrameTypeEvent)
	}
	if ev.EventName != "connect.challenge" {
		t.Errorf("event = %q, want %q", ev.EventName, "connect.challenge")
	}
}

func TestMarshalEventError(t *testing.T) {
	_, err := MarshalEvent("test", unmarshallable{})
	if err == nil {
		t.Fatal("expected error marshaling unmarshallable payload")
	}
}

func TestParseFrame(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType FrameType
	}{
		{"request", `{"type":"req","id":"1","method":"connect","params":{}}`, FrameTypeRequest},
		{"response", `{"type":"res","id":"1","ok":true}`, FrameTypeResponse},
		{"event", `{"type":"event","event":"tick"}`, FrameTypeEvent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ParseFrame([]byte(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if f.Type != tt.wantType {
				t.Errorf("type = %q, want %q", f.Type, tt.wantType)
			}
		})
	}
}

func TestParseFrameError(t *testing.T) {
	_, err := ParseFrame([]byte("not json"))
	if err == nil {
		t.Fatal("expected error parsing invalid JSON")
	}
}

func TestUnmarshalRequest(t *testing.T) {
	input := `{"type":"req","id":"r1","method":"ping","params":{}}`
	req, err := UnmarshalRequest([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if req.Type != FrameTypeRequest {
		t.Errorf("type = %q, want %q", req.Type, FrameTypeRequest)
	}
	if req.ID != "r1" {
		t.Errorf("id = %q, want %q", req.ID, "r1")
	}
	if req.Method != "ping" {
		t.Errorf("method = %q, want %q", req.Method, "ping")
	}
}

func TestUnmarshalRequestError(t *testing.T) {
	_, err := UnmarshalRequest([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUnmarshalResponse(t *testing.T) {
	input := `{"type":"res","id":"r1","ok":true,"payload":{"key":"val"}}`
	resp, err := UnmarshalResponse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Type != FrameTypeResponse {
		t.Errorf("type = %q, want %q", resp.Type, FrameTypeResponse)
	}
	if !resp.OK {
		t.Error("ok = false, want true")
	}
	if resp.ID != "r1" {
		t.Errorf("id = %q, want %q", resp.ID, "r1")
	}
}

func TestUnmarshalResponseError(t *testing.T) {
	_, err := UnmarshalResponse([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUnmarshalEvent(t *testing.T) {
	input := `{"type":"event","event":"tick","payload":{"ts":123}}`
	ev, err := UnmarshalEvent([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if ev.Type != FrameTypeEvent {
		t.Errorf("type = %q, want %q", ev.Type, FrameTypeEvent)
	}
	if ev.EventName != "tick" {
		t.Errorf("event = %q, want %q", ev.EventName, "tick")
	}
}

func TestUnmarshalEventError(t *testing.T) {
	_, err := UnmarshalEvent([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConnectParamsRoundTrip(t *testing.T) {
	params := ConnectParams{
		MinProtocol: 3,
		MaxProtocol: 3,
		Client: ClientInfo{
			ID:              ClientIDCLI,
			DisplayName:     "My CLI",
			Version:         "1.2.3",
			Platform:        "macos",
			DeviceFamily:    "mac",
			ModelIdentifier: "Mac14,2",
			Mode:            ClientModeCLI,
			InstanceID:      "inst-1",
		},
		Role:    RoleOperator,
		Scopes:  []Scope{ScopeOperatorRead, ScopeOperatorWrite},
		Caps:    nil,
		PathEnv: "/usr/local/bin",
		Auth:    AuthParams{Token: "secret"},
		Device: &DeviceIdentity{
			ID:        "dev-1",
			PublicKey: "pk",
			Signature: "sig",
			SignedAt:  1737264000000,
			Nonce:     "nonce",
		},
		Locale:    "en-US",
		UserAgent: "test/1.0",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	var got ConnectParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.MinProtocol != 3 || got.MaxProtocol != 3 {
		t.Errorf("protocol = %d/%d, want 3/3", got.MinProtocol, got.MaxProtocol)
	}
	if got.Role != RoleOperator {
		t.Errorf("role = %q, want %q", got.Role, RoleOperator)
	}
	if got.Device == nil || got.Device.ID != "dev-1" {
		t.Error("device identity not round-tripped")
	}
	if got.PathEnv != "/usr/local/bin" {
		t.Errorf("pathEnv = %q, want /usr/local/bin", got.PathEnv)
	}
	if got.Client.DisplayName != "My CLI" {
		t.Errorf("displayName = %q, want 'My CLI'", got.Client.DisplayName)
	}
	if got.Client.DeviceFamily != "mac" {
		t.Errorf("deviceFamily = %q, want 'mac'", got.Client.DeviceFamily)
	}
	if got.Client.ModelIdentifier != "Mac14,2" {
		t.Errorf("modelIdentifier = %q, want 'Mac14,2'", got.Client.ModelIdentifier)
	}
	if got.Client.InstanceID != "inst-1" {
		t.Errorf("instanceId = %q, want 'inst-1'", got.Client.InstanceID)
	}
}

func TestHelloOKRoundTrip(t *testing.T) {
	issuedAt := int64(1700000000000)
	hello := HelloOK{
		Type:     "hello-ok",
		Protocol: 3,
		Server: HelloServer{
			Version: "1.0.0",
			Commit:  "abc123",
			Host:    "gateway-1",
			ConnID:  "conn-42",
		},
		Features: HelloFeatures{
			Methods: []string{"chat.send", "system-presence"},
			Events:  []string{"chat", "tick"},
		},
		Snapshot: Snapshot{
			Presence:     []PresenceEntry{{Ts: 1234, DeviceID: "d1"}},
			Health:       json.RawMessage(`{}`),
			StateVersion: StateVersion{Presence: 1, Health: 2},
			UptimeMs:     5000,
			ConfigPath:   "/etc/openclaw/config.yaml",
			StateDir:     "/var/lib/openclaw",
			SessionDefaults: &SessionDefaults{
				DefaultAgentID: "default",
				MainKey:        "main",
				MainSessionKey: "main",
				Scope:          "per-sender",
			},
			AuthMode: "token",
		},
		CanvasHostURL: "https://canvas.example.com",
		Auth: &HelloAuth{
			DeviceToken: "tok-123",
			Role:        "operator",
			Scopes:      []string{"operator.read"},
			IssuedAtMs:  &issuedAt,
		},
		Policy: HelloPolicy{
			MaxPayload:       MaxPayloadBytes,
			MaxBufferedBytes: MaxBufferedBytes,
			TickIntervalMs:   DefaultTickIntervalMs,
		},
	}
	data, err := json.Marshal(hello)
	if err != nil {
		t.Fatal(err)
	}
	var got HelloOK
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Protocol != 3 {
		t.Errorf("protocol = %d, want 3", got.Protocol)
	}
	if got.Server.Version != "1.0.0" {
		t.Errorf("server.version = %q, want '1.0.0'", got.Server.Version)
	}
	if got.Server.ConnID != "conn-42" {
		t.Errorf("server.connId = %q, want 'conn-42'", got.Server.ConnID)
	}
	if len(got.Features.Methods) != 2 {
		t.Errorf("features.methods len = %d, want 2", len(got.Features.Methods))
	}
	if got.Snapshot.UptimeMs != 5000 {
		t.Errorf("snapshot.uptimeMs = %d, want 5000", got.Snapshot.UptimeMs)
	}
	if got.Snapshot.SessionDefaults == nil || got.Snapshot.SessionDefaults.DefaultAgentID != "default" {
		t.Error("sessionDefaults not round-tripped")
	}
	if got.Snapshot.AuthMode != "token" {
		t.Errorf("authMode = %q, want 'token'", got.Snapshot.AuthMode)
	}
	if got.CanvasHostURL != "https://canvas.example.com" {
		t.Errorf("canvasHostUrl = %q", got.CanvasHostURL)
	}
	if got.Auth == nil || got.Auth.DeviceToken != "tok-123" {
		t.Error("auth not round-tripped")
	}
	if got.Auth.IssuedAtMs == nil || *got.Auth.IssuedAtMs != 1700000000000 {
		t.Error("auth.issuedAtMs not round-tripped")
	}
	if got.Policy.MaxPayload != MaxPayloadBytes {
		t.Errorf("policy.maxPayload = %d", got.Policy.MaxPayload)
	}
	if got.Policy.MaxBufferedBytes != MaxBufferedBytes {
		t.Errorf("policy.maxBufferedBytes = %d", got.Policy.MaxBufferedBytes)
	}
}

func TestNodeConnectParams(t *testing.T) {
	params := ConnectParams{
		MinProtocol: 3,
		MaxProtocol: 3,
		Client: ClientInfo{
			ID:       ClientIDIOS,
			Version:  "1.2.3",
			Platform: "ios",
			Mode:     ClientModeNode,
		},
		Role:        RoleNode,
		Scopes:      nil,
		Caps:        []string{"camera", "canvas", "screen"},
		Commands:    []string{"camera.snap", "canvas.navigate"},
		Permissions: map[string]bool{"camera.capture": true, "screen.record": false},
		Auth:        AuthParams{Token: "node-token"},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}
	var got ConnectParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Role != RoleNode {
		t.Errorf("role = %q, want %q", got.Role, RoleNode)
	}
	if len(got.Caps) != 3 {
		t.Errorf("caps len = %d, want 3", len(got.Caps))
	}
	if len(got.Commands) != 2 {
		t.Errorf("commands len = %d, want 2", len(got.Commands))
	}
	if !got.Permissions["camera.capture"] {
		t.Error("permissions[camera.capture] = false, want true")
	}
}

func TestExecApprovalTypes(t *testing.T) {
	req := ExecApprovalRequested{
		ID:      "approval-1",
		Command: "ls -la",
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var got ExecApprovalRequested
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != "approval-1" {
		t.Errorf("id = %q, want 'approval-1'", got.ID)
	}

	resolve := ExecApprovalResolveParams{
		ID:       "approval-1",
		Decision: "approved",
	}
	data, err = json.Marshal(resolve)
	if err != nil {
		t.Fatal(err)
	}
	var gotR ExecApprovalResolveParams
	if err := json.Unmarshal(data, &gotR); err != nil {
		t.Fatal(err)
	}
	if gotR.Decision != "approved" {
		t.Errorf("decision = %q, want 'approved'", gotR.Decision)
	}
}

func TestInvokeRoundTrip(t *testing.T) {
	inv := Invoke{
		Type:    "invoke",
		ID:      "inv-1",
		Command: "camera.snap",
		Params:  json.RawMessage(`{"format":"jpeg"}`),
	}
	data, err := json.Marshal(inv)
	if err != nil {
		t.Fatal(err)
	}
	var got Invoke
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Command != "camera.snap" {
		t.Errorf("command = %q, want %q", got.Command, "camera.snap")
	}

	res := InvokeResponse{
		Type:    "invoke-res",
		ID:      "inv-1",
		OK:      true,
		Payload: json.RawMessage(`{"url":"file:///photo.jpg"}`),
	}
	data, err = json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	var gotR InvokeResponse
	if err := json.Unmarshal(data, &gotR); err != nil {
		t.Fatal(err)
	}
	if !gotR.OK {
		t.Error("ok = false, want true")
	}
}

func TestExecFinishedRoundTrip(t *testing.T) {
	exitCode := 0
	success := true
	timedOut := false
	ef := ExecFinished{
		SessionKey: "main",
		RunID:      "run-1",
		Command:    "ls",
		ExitCode:   &exitCode,
		TimedOut:   &timedOut,
		Success:    &success,
		Output:     "file.txt",
	}
	data, err := json.Marshal(ef)
	if err != nil {
		t.Fatal(err)
	}
	var got ExecFinished
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.SessionKey != "main" {
		t.Errorf("sessionKey = %q, want main", got.SessionKey)
	}
	if *got.ExitCode != 0 {
		t.Errorf("exitCode = %d, want 0", *got.ExitCode)
	}
	if !*got.Success {
		t.Error("success = false, want true")
	}
}

func TestExecDeniedRoundTrip(t *testing.T) {
	ed := ExecDenied{
		SessionKey: "main",
		RunID:      "run-2",
		Command:    "rm -rf /",
		Reason:     "dangerous command",
	}
	data, err := json.Marshal(ed)
	if err != nil {
		t.Fatal(err)
	}
	var got ExecDenied
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Reason != "dangerous command" {
		t.Errorf("reason = %q, want 'dangerous command'", got.Reason)
	}
}

func TestPresenceEntryRoundTrip(t *testing.T) {
	lastInput := 42
	pe := PresenceEntry{
		Host:             "host-1",
		IP:               "192.168.1.1",
		Version:          "1.0.0",
		Platform:         "macos",
		DeviceFamily:     "mac",
		ModelIdentifier:  "Mac14,2",
		Mode:             "operator",
		LastInputSeconds: &lastInput,
		Reason:           "connected",
		Tags:             []string{"primary"},
		Text:             "hello",
		Ts:               1700000000000,
		DeviceID:         "dev-1",
		Roles:            []string{"operator"},
		Scopes:           []string{"operator.read"},
		InstanceID:       "inst-1",
	}
	data, err := json.Marshal(pe)
	if err != nil {
		t.Fatal(err)
	}
	var got PresenceEntry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.DeviceID != "dev-1" {
		t.Errorf("deviceId = %q, want dev-1", got.DeviceID)
	}
	if got.Host != "host-1" {
		t.Errorf("host = %q, want host-1", got.Host)
	}
	if got.IP != "192.168.1.1" {
		t.Errorf("ip = %q", got.IP)
	}
	if got.LastInputSeconds == nil || *got.LastInputSeconds != 42 {
		t.Error("lastInputSeconds not round-tripped")
	}
	if len(got.Roles) != 1 || got.Roles[0] != "operator" {
		t.Errorf("roles = %v", got.Roles)
	}
	if len(got.Scopes) != 1 {
		t.Errorf("scopes = %v", got.Scopes)
	}
	if got.InstanceID != "inst-1" {
		t.Errorf("instanceId = %q", got.InstanceID)
	}
}

func TestEventWithSeqAndStateVersion(t *testing.T) {
	seq := int64(42)
	sv := StateVersion{Presence: 7, Health: 3}
	ev := Event{
		Type:         FrameTypeEvent,
		EventName:    "test",
		Payload:      json.RawMessage(`{}`),
		Seq:          &seq,
		StateVersion: &sv,
	}
	data, err := json.Marshal(ev)
	if err != nil {
		t.Fatal(err)
	}
	var got Event
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Seq == nil || *got.Seq != 42 {
		t.Errorf("seq = %v, want 42", got.Seq)
	}
	if got.StateVersion == nil || got.StateVersion.Presence != 7 {
		t.Errorf("stateVersion.presence = %v, want 7", got.StateVersion)
	}
	if got.StateVersion.Health != 3 {
		t.Errorf("stateVersion.health = %d, want 3", got.StateVersion.Health)
	}
}

func TestErrorPayloadFull(t *testing.T) {
	retryable := true
	retryAfter := 5000
	ep := ErrorPayload{
		Code:         "AGENT_TIMEOUT",
		Message:      "agent timed out",
		Details:      map[string]string{"key": "val"},
		Retryable:    &retryable,
		RetryAfterMs: &retryAfter,
	}
	data, err := json.Marshal(ep)
	if err != nil {
		t.Fatal(err)
	}
	var got ErrorPayload
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Code != "AGENT_TIMEOUT" {
		t.Errorf("code = %q, want AGENT_TIMEOUT", got.Code)
	}
	if got.Retryable == nil || !*got.Retryable {
		t.Error("retryable not round-tripped")
	}
	if got.RetryAfterMs == nil || *got.RetryAfterMs != 5000 {
		t.Error("retryAfterMs not round-tripped")
	}
	if got.Details == nil {
		t.Error("details is nil")
	}
}

func TestProtocolVersionConst(t *testing.T) {
	if ProtocolVersion != 4 {
		t.Errorf("ProtocolVersion = %d, want 4", ProtocolVersion)
	}
	if MinProtocolVersion != 3 {
		t.Errorf("MinProtocolVersion = %d, want 3", MinProtocolVersion)
	}
	if MinProtocolVersion > ProtocolVersion {
		t.Errorf("MinProtocolVersion %d must not exceed ProtocolVersion %d", MinProtocolVersion, ProtocolVersion)
	}
}

func TestServerConstants(t *testing.T) {
	if MaxPayloadBytes != 25*1024*1024 {
		t.Errorf("MaxPayloadBytes = %d", MaxPayloadBytes)
	}
	if MaxBufferedBytes != 50*1024*1024 {
		t.Errorf("MaxBufferedBytes = %d", MaxBufferedBytes)
	}
	if DefaultTickIntervalMs != 30000 {
		t.Errorf("DefaultTickIntervalMs = %d", DefaultTickIntervalMs)
	}
	if DefaultHandshakeTimeoutMs != 10000 {
		t.Errorf("DefaultHandshakeTimeoutMs = %d", DefaultHandshakeTimeoutMs)
	}
	if DedupeTTLMs != 300000 {
		t.Errorf("DedupeTTLMs = %d", DedupeTTLMs)
	}
	if DedupeMax != 1000 {
		t.Errorf("DedupeMax = %d", DedupeMax)
	}
	if DefaultMaxChatHistoryMessagesBytes != 6*1024*1024 {
		t.Errorf("DefaultMaxChatHistoryMessagesBytes = %d", DefaultMaxChatHistoryMessagesBytes)
	}
	if HealthRefreshIntervalMs != 60000 {
		t.Errorf("HealthRefreshIntervalMs = %d", HealthRefreshIntervalMs)
	}
	if SessionLabelMaxLength != 64 {
		t.Errorf("SessionLabelMaxLength = %d", SessionLabelMaxLength)
	}
}

func TestErrorCodeConstants(t *testing.T) {
	codes := map[string]string{
		"NOT_LINKED":      ErrorCodeNotLinked,
		"NOT_PAIRED":      ErrorCodeNotPaired,
		"AGENT_TIMEOUT":   ErrorCodeAgentTimeout,
		"INVALID_REQUEST": ErrorCodeInvalidRequest,
		"UNAVAILABLE":     ErrorCodeUnavailable,
	}
	for want, got := range codes {
		if got != want {
			t.Errorf("ErrorCode = %q, want %q", got, want)
		}
	}
}

func TestClientIDConstants(t *testing.T) {
	ids := []struct{ got, want string }{
		{ClientIDWebchatUI, "webchat-ui"},
		{ClientIDControlUI, "openclaw-control-ui"},
		{ClientIDWebchat, "webchat"},
		{ClientIDCLI, "cli"},
		{ClientIDGateway, "gateway-client"},
		{ClientIDMacOS, "openclaw-macos"},
		{ClientIDIOS, "openclaw-ios"},
		{ClientIDAndroid, "openclaw-android"},
		{ClientIDNodeHost, "node-host"},
		{ClientIDTest, "test"},
		{ClientIDFingerprint, "fingerprint"},
		{ClientIDProbe, "openclaw-probe"},
	}
	for _, tt := range ids {
		if tt.got != tt.want {
			t.Errorf("ClientID = %q, want %q", tt.got, tt.want)
		}
	}
}

func TestClientModeConstants(t *testing.T) {
	modes := []struct{ got, want string }{
		{ClientModeWebchat, "webchat"},
		{ClientModeCLI, "cli"},
		{ClientModeUI, "ui"},
		{ClientModeBackend, "backend"},
		{ClientModeNode, "node"},
		{ClientModeProbe, "probe"},
		{ClientModeTest, "test"},
	}
	for _, tt := range modes {
		if tt.got != tt.want {
			t.Errorf("ClientMode = %q, want %q", tt.got, tt.want)
		}
	}
}

func TestClientCapConstants(t *testing.T) {
	if ClientCapToolEvents != "tool-events" {
		t.Errorf("ClientCapToolEvents = %q", ClientCapToolEvents)
	}
}

func TestScopeConstants(t *testing.T) {
	scopes := []Scope{
		ScopeOperatorRead,
		ScopeOperatorWrite,
		ScopeOperatorAdmin,
		ScopeOperatorApprovals,
		ScopeOperatorPairing,
	}
	expected := []string{
		"operator.read",
		"operator.write",
		"operator.admin",
		"operator.approvals",
		"operator.pairing",
	}
	for i, s := range scopes {
		if string(s) != expected[i] {
			t.Errorf("scope %d = %q, want %q", i, s, expected[i])
		}
	}
}

func TestRoleConstants(t *testing.T) {
	if string(RoleOperator) != "operator" {
		t.Errorf("RoleOperator = %q", RoleOperator)
	}
	if string(RoleNode) != "node" {
		t.Errorf("RoleNode = %q", RoleNode)
	}
}

func TestFrameTypeConstants(t *testing.T) {
	if string(FrameTypeRequest) != "req" {
		t.Errorf("FrameTypeRequest = %q", FrameTypeRequest)
	}
	if string(FrameTypeResponse) != "res" {
		t.Errorf("FrameTypeResponse = %q", FrameTypeResponse)
	}
	if string(FrameTypeEvent) != "event" {
		t.Errorf("FrameTypeEvent = %q", FrameTypeEvent)
	}
	if string(FrameTypeInvoke) != "invoke" {
		t.Errorf("FrameTypeInvoke = %q", FrameTypeInvoke)
	}
	if string(FrameTypeInvokeResponse) != "invoke-res" {
		t.Errorf("FrameTypeInvokeResponse = %q", FrameTypeInvokeResponse)
	}
}

func TestAuthParamsWithPassword(t *testing.T) {
	ap := AuthParams{
		Token:    "",
		Password: "secret123",
	}
	data, err := json.Marshal(ap)
	if err != nil {
		t.Fatal(err)
	}
	var got AuthParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Password != "secret123" {
		t.Errorf("password = %q, want 'secret123'", got.Password)
	}
}

func TestStateVersionRoundTrip(t *testing.T) {
	sv := StateVersion{Presence: 5, Health: 10}
	data, err := json.Marshal(sv)
	if err != nil {
		t.Fatal(err)
	}
	var got StateVersion
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Presence != 5 {
		t.Errorf("presence = %d, want 5", got.Presence)
	}
	if got.Health != 10 {
		t.Errorf("health = %d, want 10", got.Health)
	}
}

func TestSnapshotRoundTrip(t *testing.T) {
	snap := Snapshot{
		Presence: []PresenceEntry{
			{Ts: 1234, DeviceID: "d1"},
		},
		Health:       json.RawMessage(`{"cpu":0.5}`),
		StateVersion: StateVersion{Presence: 1, Health: 2},
		UptimeMs:     60000,
		ConfigPath:   "/etc/config.yaml",
		StateDir:     "/var/state",
		SessionDefaults: &SessionDefaults{
			DefaultAgentID: "agent-1",
			MainKey:        "main",
			MainSessionKey: "main",
			Scope:          "global",
		},
		AuthMode: "password",
	}
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	var got Snapshot
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Presence) != 1 {
		t.Errorf("presence len = %d", len(got.Presence))
	}
	if got.UptimeMs != 60000 {
		t.Errorf("uptimeMs = %d", got.UptimeMs)
	}
	if got.SessionDefaults == nil || got.SessionDefaults.Scope != "global" {
		t.Error("sessionDefaults not round-tripped")
	}
}

// --- Chat types ---

func TestChatSendParamsRoundTrip(t *testing.T) {
	deliver := true
	timeout := 30000
	p := ChatSendParams{
		SessionKey:     "main",
		Message:        "hello",
		Thinking:       "low",
		Deliver:        &deliver,
		Attachments:    json.RawMessage(`[{"type":"file"}]`),
		TimeoutMs:      &timeout,
		IdempotencyKey: "key-1",
	}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var got ChatSendParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.SessionKey != "main" {
		t.Errorf("sessionKey = %q", got.SessionKey)
	}
	if got.Deliver == nil || !*got.Deliver {
		t.Error("deliver not round-tripped")
	}
	if got.IdempotencyKey != "key-1" {
		t.Errorf("idempotencyKey = %q", got.IdempotencyKey)
	}
}

func TestChatHistoryParamsRoundTrip(t *testing.T) {
	limit := 50
	p := ChatHistoryParams{SessionKey: "main", Limit: &limit}
	data, _ := json.Marshal(p)
	var got ChatHistoryParams
	json.Unmarshal(data, &got)
	if got.Limit == nil || *got.Limit != 50 {
		t.Error("limit not round-tripped")
	}
}

func TestChatAbortParamsRoundTrip(t *testing.T) {
	p := ChatAbortParams{SessionKey: "main", RunID: "run-1"}
	data, _ := json.Marshal(p)
	var got ChatAbortParams
	json.Unmarshal(data, &got)
	if got.RunID != "run-1" {
		t.Errorf("runId = %q", got.RunID)
	}
}

func TestChatInjectParamsRoundTrip(t *testing.T) {
	p := ChatInjectParams{SessionKey: "main", Message: "injected", Label: "test"}
	data, _ := json.Marshal(p)
	var got ChatInjectParams
	json.Unmarshal(data, &got)
	if got.Label != "test" {
		t.Errorf("label = %q", got.Label)
	}
}

func TestChatSendResultRoundTrip(t *testing.T) {
	r := ChatSendResult{RunID: "run-1", Status: "started"}
	data, _ := json.Marshal(r)
	var got ChatSendResult
	json.Unmarshal(data, &got)
	if got.RunID != "run-1" {
		t.Errorf("runId = %q", got.RunID)
	}
	if got.Status != "started" {
		t.Errorf("status = %q", got.Status)
	}
}

func TestChatEventRoundTrip(t *testing.T) {
	e := ChatEvent{
		RunID:      "run-1",
		SessionKey: "main",
		Seq:        5,
		State:      "delta",
		StopReason: "end_turn",
	}
	data, _ := json.Marshal(e)
	var got ChatEvent
	json.Unmarshal(data, &got)
	if got.State != "delta" {
		t.Errorf("state = %q", got.State)
	}
	if got.StopReason != "end_turn" {
		t.Errorf("stopReason = %q", got.StopReason)
	}
}

// --- Agent types ---

func TestAgentParamsRoundTrip(t *testing.T) {
	p := AgentParams{
		Message:        "hello agent",
		AgentID:        "agent-1",
		IdempotencyKey: "key-1",
		InputProvenance: &InputProvenance{
			Kind:             "external_user",
			SourceSessionKey: "session-1",
		},
		Label:     "test-label",
		SpawnedBy: "parent",
	}
	data, _ := json.Marshal(p)
	var got AgentParams
	json.Unmarshal(data, &got)
	if got.AgentID != "agent-1" {
		t.Errorf("agentId = %q", got.AgentID)
	}
	if got.InputProvenance == nil || got.InputProvenance.Kind != "external_user" {
		t.Error("inputProvenance not round-tripped")
	}
}

func TestAgentEventRoundTrip(t *testing.T) {
	e := AgentEvent{
		RunID:  "run-1",
		Seq:    0,
		Stream: "main",
		Ts:     1700000000000,
		Data:   map[string]any{"key": "val"},
	}
	data, _ := json.Marshal(e)
	var got AgentEvent
	json.Unmarshal(data, &got)
	if got.Stream != "main" {
		t.Errorf("stream = %q", got.Stream)
	}
}

func TestAgentIdentityRoundTrip(t *testing.T) {
	p := AgentIdentityParams{AgentID: "a1", SessionKey: "s1"}
	data, _ := json.Marshal(p)
	var gotP AgentIdentityParams
	json.Unmarshal(data, &gotP)
	if gotP.AgentID != "a1" {
		t.Errorf("agentId = %q", gotP.AgentID)
	}

	r := AgentIdentityResult{AgentID: "a1", Name: "Test", Avatar: "av", Emoji: "🤖"}
	data, _ = json.Marshal(r)
	var gotR AgentIdentityResult
	json.Unmarshal(data, &gotR)
	if gotR.Name != "Test" {
		t.Errorf("name = %q", gotR.Name)
	}
}

func TestAgentWaitParamsRoundTrip(t *testing.T) {
	timeout := 5000
	p := AgentWaitParams{RunID: "run-1", TimeoutMs: &timeout}
	data, _ := json.Marshal(p)
	var got AgentWaitParams
	json.Unmarshal(data, &got)
	if got.TimeoutMs == nil || *got.TimeoutMs != 5000 {
		t.Error("timeoutMs not round-tripped")
	}
}

// --- Session types ---

func TestSessionsListParamsRoundTrip(t *testing.T) {
	limit := 10
	p := SessionsListParams{Limit: &limit, AgentID: "a1", Label: "test"}
	data, _ := json.Marshal(p)
	var got SessionsListParams
	json.Unmarshal(data, &got)
	if got.Limit == nil || *got.Limit != 10 {
		t.Error("limit not round-tripped")
	}
	if got.AgentID != "a1" {
		t.Errorf("agentId = %q", got.AgentID)
	}
}

func TestSessionsPreviewParamsRoundTrip(t *testing.T) {
	p := SessionsPreviewParams{Keys: []string{"k1", "k2"}}
	data, _ := json.Marshal(p)
	var got SessionsPreviewParams
	json.Unmarshal(data, &got)
	if len(got.Keys) != 2 {
		t.Errorf("keys len = %d", len(got.Keys))
	}
}

func TestSessionsPatchParamsRoundTrip(t *testing.T) {
	label := "new-label"
	p := SessionsPatchParams{Key: "main", Label: &label}
	data, _ := json.Marshal(p)
	var got SessionsPatchParams
	json.Unmarshal(data, &got)
	if got.Label == nil || *got.Label != "new-label" {
		t.Error("label not round-tripped")
	}
}

func TestSessionsResetParamsRoundTrip(t *testing.T) {
	p := SessionsResetParams{Key: "main", Reason: "new"}
	data, _ := json.Marshal(p)
	var got SessionsResetParams
	json.Unmarshal(data, &got)
	if got.Reason != "new" {
		t.Errorf("reason = %q", got.Reason)
	}
}

func TestSessionsDeleteParamsRoundTrip(t *testing.T) {
	del := true
	p := SessionsDeleteParams{Key: "main", DeleteTranscript: &del}
	data, _ := json.Marshal(p)
	var got SessionsDeleteParams
	json.Unmarshal(data, &got)
	if got.DeleteTranscript == nil || !*got.DeleteTranscript {
		t.Error("deleteTranscript not round-tripped")
	}
}

func TestSessionsCompactParamsRoundTrip(t *testing.T) {
	max := 100
	p := SessionsCompactParams{Key: "main", MaxLines: &max}
	data, _ := json.Marshal(p)
	var got SessionsCompactParams
	json.Unmarshal(data, &got)
	if got.MaxLines == nil || *got.MaxLines != 100 {
		t.Error("maxLines not round-tripped")
	}
}

func TestSessionsUsageParamsRoundTrip(t *testing.T) {
	p := SessionsUsageParams{Key: "main", StartDate: "2024-01-01", EndDate: "2024-12-31"}
	data, _ := json.Marshal(p)
	var got SessionsUsageParams
	json.Unmarshal(data, &got)
	if got.StartDate != "2024-01-01" {
		t.Errorf("startDate = %q", got.StartDate)
	}
}

func TestSessionsResolveParamsRoundTrip(t *testing.T) {
	p := SessionsResolveParams{Key: "main", AgentID: "a1"}
	data, _ := json.Marshal(p)
	var got SessionsResolveParams
	json.Unmarshal(data, &got)
	if got.AgentID != "a1" {
		t.Errorf("agentId = %q", got.AgentID)
	}
}

// --- Node types ---

func TestNodePairRequestParamsRoundTrip(t *testing.T) {
	p := NodePairRequestParams{
		NodeID:      "node-1",
		DisplayName: "My Node",
		Platform:    "ios",
		Version:     "1.0.0",
		Caps:        []string{"camera"},
	}
	data, _ := json.Marshal(p)
	var got NodePairRequestParams
	json.Unmarshal(data, &got)
	if got.NodeID != "node-1" {
		t.Errorf("nodeId = %q", got.NodeID)
	}
}

func TestNodeInvokeParamsRoundTrip(t *testing.T) {
	p := NodeInvokeParams{
		NodeID:         "node-1",
		Command:        "camera.snap",
		Params:         json.RawMessage(`{}`),
		IdempotencyKey: "key-1",
	}
	data, _ := json.Marshal(p)
	var got NodeInvokeParams
	json.Unmarshal(data, &got)
	if got.Command != "camera.snap" {
		t.Errorf("command = %q", got.Command)
	}
}

func TestNodeInvokeResultParamsRoundTrip(t *testing.T) {
	p := NodeInvokeResultParams{
		ID:     "inv-1",
		NodeID: "node-1",
		OK:     false,
		Error:  &NodeInvokeResultError{Code: "TIMEOUT", Message: "timed out"},
	}
	data, _ := json.Marshal(p)
	var got NodeInvokeResultParams
	json.Unmarshal(data, &got)
	if got.Error == nil || got.Error.Code != "TIMEOUT" {
		t.Error("error not round-tripped")
	}
}

func TestNodeEventParamsRoundTrip(t *testing.T) {
	p := NodeEventParams{Event: "status", Payload: json.RawMessage(`{"ok":true}`)}
	data, _ := json.Marshal(p)
	var got NodeEventParams
	json.Unmarshal(data, &got)
	if got.Event != "status" {
		t.Errorf("event = %q", got.Event)
	}
}

func TestNodeInvokeRequestEventRoundTrip(t *testing.T) {
	e := NodeInvokeRequestEvent{
		ID:      "inv-1",
		NodeID:  "node-1",
		Command: "snap",
	}
	data, _ := json.Marshal(e)
	var got NodeInvokeRequestEvent
	json.Unmarshal(data, &got)
	if got.Command != "snap" {
		t.Errorf("command = %q", got.Command)
	}
}

func TestNodePairApproveRejectVerifyRoundTrip(t *testing.T) {
	ap := NodePairApproveParams{RequestID: "r1"}
	rp := NodePairRejectParams{RequestID: "r1"}
	vp := NodePairVerifyParams{NodeID: "n1", Token: "tok"}
	rn := NodeRenameParams{NodeID: "n1", DisplayName: "New Name"}
	dp := NodeDescribeParams{NodeID: "n1"}

	for _, v := range []any{ap, rp, vp, rn, dp} {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Error("empty marshal")
		}
	}
}

// --- Device pairing types ---

func TestDevicePairTypesRoundTrip(t *testing.T) {
	types := []any{
		DevicePairApproveParams{RequestID: "r1"},
		DevicePairRejectParams{RequestID: "r1"},
		DevicePairRemoveParams{DeviceID: "d1"},
		DeviceTokenRotateParams{DeviceID: "d1", Role: "operator", Scopes: []string{"operator.read"}},
		DeviceTokenRevokeParams{DeviceID: "d1", Role: "operator"},
	}
	for _, v := range types {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Error("empty marshal")
		}
	}
}

func TestDevicePairRequestedEventRoundTrip(t *testing.T) {
	e := DevicePairRequestedEvent{
		RequestID: "r1",
		DeviceID:  "d1",
		PublicKey: "pk",
		Ts:        1700000000000,
		Roles:     []string{"operator"},
		Scopes:    []string{"operator.read"},
	}
	data, _ := json.Marshal(e)
	var got DevicePairRequestedEvent
	json.Unmarshal(data, &got)
	if got.PublicKey != "pk" {
		t.Errorf("publicKey = %q", got.PublicKey)
	}
}

func TestDevicePairResolvedEventRoundTrip(t *testing.T) {
	e := DevicePairResolvedEvent{RequestID: "r1", DeviceID: "d1", Decision: "approved", Ts: 123}
	data, _ := json.Marshal(e)
	var got DevicePairResolvedEvent
	json.Unmarshal(data, &got)
	if got.Decision != "approved" {
		t.Errorf("decision = %q", got.Decision)
	}
}

// --- Config types ---

func TestConfigTypesRoundTrip(t *testing.T) {
	types := []any{
		ConfigGetParams{},
		ConfigSetParams{Raw: "yaml", BaseHash: "abc"},
		ConfigApplyParams{Raw: "yaml", Note: "test"},
		ConfigPatchParams{Raw: "yaml", SessionKey: "main"},
	}
	for _, v := range types {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Error("empty marshal")
		}
	}
}

func TestConfigSchemaResponseRoundTrip(t *testing.T) {
	r := ConfigSchemaResponse{
		Schema:      json.RawMessage(`{}`),
		UIHints:     map[string]ConfigUIHint{"field": {Label: "Field", Help: "help"}},
		Version:     "1.0.0",
		GeneratedAt: "2024-01-01T00:00:00Z",
	}
	data, _ := json.Marshal(r)
	var got ConfigSchemaResponse
	json.Unmarshal(data, &got)
	if got.Version != "1.0.0" {
		t.Errorf("version = %q", got.Version)
	}
	if got.UIHints["field"].Label != "Field" {
		t.Error("uiHints not round-tripped")
	}
}

// --- Agents CRUD types ---

func TestAgentsCRUDTypesRoundTrip(t *testing.T) {
	list := AgentsListResult{
		DefaultID: "default",
		MainKey:   "main",
		Scope:     "per-sender",
		Agents: []AgentSummary{
			{ID: "a1", Name: "Agent 1", Identity: &AgentIdentity{Name: "A1", Emoji: "🤖"}},
		},
	}
	data, _ := json.Marshal(list)
	var got AgentsListResult
	json.Unmarshal(data, &got)
	if got.Scope != "per-sender" {
		t.Errorf("scope = %q", got.Scope)
	}
	if len(got.Agents) != 1 || got.Agents[0].Identity == nil {
		t.Error("agents not round-tripped")
	}

	create := AgentsCreateResult{OK: true, AgentID: "a1", Name: "Test", Workspace: "/ws"}
	data, _ = json.Marshal(create)
	var gotC AgentsCreateResult
	json.Unmarshal(data, &gotC)
	if !gotC.OK || gotC.AgentID != "a1" {
		t.Error("create result not round-tripped")
	}

	del := AgentsDeleteResult{OK: true, AgentID: "a1", RemovedBindings: 3}
	data, _ = json.Marshal(del)
	var gotD AgentsDeleteResult
	json.Unmarshal(data, &gotD)
	if gotD.RemovedBindings != 3 {
		t.Errorf("removedBindings = %d", gotD.RemovedBindings)
	}
}

func TestAgentsFileEntryRoundTrip(t *testing.T) {
	size := 1024
	updAt := int64(1700000000000)
	e := AgentsFileEntry{
		Name: "agent.yaml", Path: "/ws/agent.yaml", Missing: false,
		Size: &size, UpdatedAtMs: &updAt, Content: "content",
	}
	data, _ := json.Marshal(e)
	var got AgentsFileEntry
	json.Unmarshal(data, &got)
	if got.Size == nil || *got.Size != 1024 {
		t.Error("size not round-tripped")
	}
}

func TestAgentsFilesListResultRoundTrip(t *testing.T) {
	r := AgentsFilesListResult{
		AgentID:   "a1",
		Workspace: "/ws",
		Files:     []AgentsFileEntry{{Name: "f1", Path: "/ws/f1", Missing: false}},
	}
	data, _ := json.Marshal(r)
	var got AgentsFilesListResult
	json.Unmarshal(data, &got)
	if len(got.Files) != 1 {
		t.Errorf("files len = %d", len(got.Files))
	}
}

func TestAgentsFilesGetSetRoundTrip(t *testing.T) {
	getP := AgentsFilesGetParams{AgentID: "a1", Name: "f1"}
	setP := AgentsFilesSetParams{AgentID: "a1", Name: "f1", Content: "data"}
	for _, v := range []any{getP, setP} {
		data, _ := json.Marshal(v)
		if len(data) == 0 {
			t.Error("empty")
		}
	}
}

// --- Models types ---

func TestModelsListResultRoundTrip(t *testing.T) {
	ctxWindow := 128000
	reasoning := true
	r := ModelsListResult{
		Models: []ModelChoice{
			{ID: "m1", Name: "Model 1", Provider: "openai", ContextWindow: &ctxWindow, Reasoning: &reasoning},
		},
	}
	data, _ := json.Marshal(r)
	var got ModelsListResult
	json.Unmarshal(data, &got)
	if len(got.Models) != 1 || got.Models[0].Provider != "openai" {
		t.Error("models not round-tripped")
	}
}

// --- Logs types ---

func TestLogsTailRoundTrip(t *testing.T) {
	p := LogsTailParams{Limit: intPtr(100)}
	data, _ := json.Marshal(p)
	var gotP LogsTailParams
	json.Unmarshal(data, &gotP)
	if gotP.Limit == nil || *gotP.Limit != 100 {
		t.Error("limit not round-tripped")
	}

	trunc := true
	r := LogsTailResult{
		File: "/var/log/test.log", Cursor: 42, Size: 1000,
		Lines: []string{"line1", "line2"}, Truncated: &trunc,
	}
	data, _ = json.Marshal(r)
	var gotR LogsTailResult
	json.Unmarshal(data, &gotR)
	if gotR.Cursor != 42 {
		t.Errorf("cursor = %d", gotR.Cursor)
	}
}

// --- Cron types ---

func TestCronJobRoundTrip(t *testing.T) {
	job := CronJob{
		ID: "j1", Name: "Test Job", Enabled: true,
		CreatedAtMs: 1700000000000, UpdatedAtMs: 1700000000000,
		Schedule:      CronSchedule{Kind: "cron", Expr: "0 * * * *", Tz: "UTC"},
		SessionTarget: "main", WakeMode: "now",
		Payload:  CronPayload{Kind: "agentTurn", Message: "hello"},
		Delivery: &CronDelivery{Mode: "announce", Channel: "last"},
		State:    CronJobState{LastStatus: "ok"},
	}
	data, _ := json.Marshal(job)
	var got CronJob
	json.Unmarshal(data, &got)
	if got.Schedule.Expr != "0 * * * *" {
		t.Errorf("schedule.expr = %q", got.Schedule.Expr)
	}
	if got.Delivery == nil || got.Delivery.Mode != "announce" {
		t.Error("delivery not round-tripped")
	}
}

func TestCronParamsRoundTrip(t *testing.T) {
	types := []any{
		CronListParams{},
		CronAddParams{Name: "j1", Schedule: CronSchedule{Kind: "at", At: "2024-01-01T00:00:00Z"}, SessionTarget: "main", WakeMode: "now", Payload: CronPayload{Kind: "systemEvent", Text: "test"}},
		CronUpdateParams{ID: "j1", Patch: CronJobPatch{Name: "updated"}},
		CronRemoveParams{ID: "j1"},
		CronRunParams{JobID: "j1", Mode: "force"},
		CronRunsParams{ID: "j1"},
	}
	for _, v := range types {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Error("empty")
		}
	}
}

func TestCronRunLogEntryRoundTrip(t *testing.T) {
	in := 42
	out := 10
	total := 52
	e := CronRunLogEntry{
		Ts: 1700000000000, JobID: "j1", Action: "finished",
		Status: "ok", SessionKey: "main", JobName: "daily",
		Model: "gpt-4", Provider: "anthropic",
		Usage: &CronRunLogUsage{InputTokens: &in, OutputTokens: &out, TotalTokens: &total},
	}
	data, _ := json.Marshal(e)
	var got CronRunLogEntry
	json.Unmarshal(data, &got)
	if got.Action != "finished" {
		t.Errorf("action = %q", got.Action)
	}
	if got.JobName != "daily" {
		t.Errorf("jobName = %q", got.JobName)
	}
	if got.Usage == nil || got.Usage.TotalTokens == nil || *got.Usage.TotalTokens != 52 {
		t.Error("usage not round-tripped")
	}
}

func TestCronListResultRoundTrip(t *testing.T) {
	nextOff := 2
	r := CronListResult{
		Jobs: []CronJob{
			{ID: "j1", Name: "daily", Enabled: true, Schedule: CronSchedule{Kind: "cron", Expr: "0 9 * * *"}, SessionTarget: "main", WakeMode: "now", Payload: CronPayload{Kind: "systemEvent", Text: "hi"}},
		},
		Total: 5, Offset: 0, Limit: 1, HasMore: true, NextOffset: &nextOff,
	}
	data, _ := json.Marshal(r)
	var got CronListResult
	json.Unmarshal(data, &got)
	if len(got.Jobs) != 1 || got.Jobs[0].ID != "j1" {
		t.Errorf("jobs not round-tripped: %+v", got.Jobs)
	}
	if got.Total != 5 || !got.HasMore {
		t.Errorf("pagination not round-tripped: total=%d hasMore=%v", got.Total, got.HasMore)
	}
	if got.NextOffset == nil || *got.NextOffset != 2 {
		t.Error("nextOffset not round-tripped")
	}
}

func TestCronRunsResultRoundTrip(t *testing.T) {
	r := CronRunsResult{
		Entries: []CronRunLogEntry{
			{Ts: 1000, JobID: "j1", Action: "finished", Status: "ok"},
		},
		Total: 10, Offset: 0,
	}
	data, _ := json.Marshal(r)
	var got CronRunsResult
	json.Unmarshal(data, &got)
	if len(got.Entries) != 1 || got.Entries[0].JobID != "j1" {
		t.Errorf("entries not round-tripped: %+v", got.Entries)
	}
	if got.Total != 10 {
		t.Errorf("total not round-tripped: %d", got.Total)
	}
}

// --- Channels / Talk types ---

func TestTalkTypesRoundTrip(t *testing.T) {
	mp := TalkModeParams{Enabled: true, Phase: "listening"}
	data, _ := json.Marshal(mp)
	var gotM TalkModeParams
	json.Unmarshal(data, &gotM)
	if !gotM.Enabled {
		t.Error("enabled not round-tripped")
	}

	r := TalkConfigResult{Config: TalkConfigData{
		Talk:    &TalkConfigTalk{VoiceID: "alloy"},
		Session: &TalkConfigSession{MainKey: "main"},
		UI:      &TalkConfigUI{SeamColor: "#000"},
	}}
	data, _ = json.Marshal(r)
	var gotR TalkConfigResult
	json.Unmarshal(data, &gotR)
	if gotR.Config.Talk == nil || gotR.Config.Talk.VoiceID != "alloy" {
		t.Error("talk config not round-tripped")
	}
}

func TestChannelsStatusRoundTrip(t *testing.T) {
	r := ChannelsStatusResult{
		Ts:                      1700000000000,
		ChannelOrder:            []string{"slack", "discord"},
		ChannelLabels:           map[string]string{"slack": "Slack"},
		Channels:                map[string]json.RawMessage{"slack": json.RawMessage(`{}`)},
		ChannelAccounts:         map[string][]ChannelAccountSnapshot{"slack": {{AccountID: "a1"}}},
		ChannelDefaultAccountID: map[string]string{"slack": "a1"},
		ChannelMeta:             []ChannelUIMeta{{ID: "slack", Label: "Slack", DetailLabel: "Slack Workspace"}},
	}
	data, _ := json.Marshal(r)
	var got ChannelsStatusResult
	json.Unmarshal(data, &got)
	if len(got.ChannelOrder) != 2 {
		t.Errorf("channelOrder len = %d", len(got.ChannelOrder))
	}
}

// --- Skills types ---

func TestSkillsTypesRoundTrip(t *testing.T) {
	types := []any{
		SkillsStatusParams{AgentID: "a1"},
		SkillsBinsResult{Bins: []string{"bin1"}},
		SkillsInstallParams{Name: "skill1", InstallID: "inst-1"},
		SkillsUpdateParams{SkillKey: "sk1", Env: map[string]string{"KEY": "val"}},
	}
	for _, v := range types {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Error("empty")
		}
	}
}

// --- Wizard types ---

func TestWizardTypesRoundTrip(t *testing.T) {
	startR := WizardStartResult{
		SessionID: "s1", Done: false,
		Step:   &WizardStep{ID: "step1", Type: "select", Title: "Choose", Options: []WizardStepOption{{Value: json.RawMessage(`"a"`), Label: "Option A"}}},
		Status: "running",
	}
	data, _ := json.Marshal(startR)
	var got WizardStartResult
	json.Unmarshal(data, &got)
	if got.Step == nil || got.Step.ID != "step1" {
		t.Error("step not round-tripped")
	}
	if len(got.Step.Options) != 1 {
		t.Errorf("options len = %d", len(got.Step.Options))
	}

	types := []any{
		WizardStartParams{Mode: "local"},
		WizardNextParams{SessionID: "s1", Answer: &WizardAnswer{StepID: "step1", Value: json.RawMessage(`"a"`)}},
		WizardCancelParams{SessionID: "s1"},
		WizardStatusParams{SessionID: "s1"},
		WizardStatusResult{Status: "done"},
		WizardNextResult{Done: true, Status: "done"},
	}
	for _, v := range types {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Error("empty")
		}
	}
}

// --- Push types ---

func TestPushTypesRoundTrip(t *testing.T) {
	p := PushTestParams{NodeID: "n1", Title: "Test", Environment: "sandbox"}
	data, _ := json.Marshal(p)
	var gotP PushTestParams
	json.Unmarshal(data, &gotP)
	if gotP.Environment != "sandbox" {
		t.Errorf("environment = %q", gotP.Environment)
	}

	r := PushTestResult{OK: true, Status: 200, TokenSuffix: "abc", Topic: "com.test", Environment: "production"}
	data, _ = json.Marshal(r)
	var gotR PushTestResult
	json.Unmarshal(data, &gotR)
	if !gotR.OK {
		t.Error("ok not round-tripped")
	}
}

// --- Send / Poll / Wake types ---

func TestSendParamsRoundTrip(t *testing.T) {
	p := SendParams{To: "user-1", Message: "hello", IdempotencyKey: "key-1", MediaURLs: []string{"http://example.com/img.jpg"}}
	data, _ := json.Marshal(p)
	var got SendParams
	json.Unmarshal(data, &got)
	if got.To != "user-1" {
		t.Errorf("to = %q", got.To)
	}
	if len(got.MediaURLs) != 1 {
		t.Errorf("mediaUrls len = %d", len(got.MediaURLs))
	}
}

func TestPollParamsRoundTrip(t *testing.T) {
	p := PollParams{To: "user-1", Question: "?", Options: []string{"a", "b"}, IdempotencyKey: "key-1"}
	data, _ := json.Marshal(p)
	var got PollParams
	json.Unmarshal(data, &got)
	if len(got.Options) != 2 {
		t.Errorf("options len = %d", len(got.Options))
	}
}

func TestWakeParamsRoundTrip(t *testing.T) {
	p := WakeParams{Mode: "now", Text: "wake up"}
	data, _ := json.Marshal(p)
	var got WakeParams
	json.Unmarshal(data, &got)
	if got.Mode != "now" {
		t.Errorf("mode = %q", got.Mode)
	}
}

// --- Update / Misc event types ---

func TestUpdateRunParamsRoundTrip(t *testing.T) {
	timeout := 5000
	p := UpdateRunParams{SessionKey: "main", Note: "test", TimeoutMs: &timeout}
	data, _ := json.Marshal(p)
	var got UpdateRunParams
	json.Unmarshal(data, &got)
	if got.TimeoutMs == nil || *got.TimeoutMs != 5000 {
		t.Error("timeoutMs not round-tripped")
	}
}

func TestTickEventRoundTrip(t *testing.T) {
	e := TickEvent{Ts: 1700000000000}
	data, _ := json.Marshal(e)
	var got TickEvent
	json.Unmarshal(data, &got)
	if got.Ts != 1700000000000 {
		t.Errorf("ts = %d", got.Ts)
	}
}

func TestShutdownEventRoundTrip(t *testing.T) {
	restart := int64(5000)
	e := ShutdownEvent{Reason: "update", RestartExpectedMs: &restart}
	data, _ := json.Marshal(e)
	var got ShutdownEvent
	json.Unmarshal(data, &got)
	if got.Reason != "update" {
		t.Errorf("reason = %q", got.Reason)
	}
	if got.RestartExpectedMs == nil || *got.RestartExpectedMs != 5000 {
		t.Error("restartExpectedMs not round-tripped")
	}
}

// --- Exec approvals admin types ---

func TestExecApprovalsTypesRoundTrip(t *testing.T) {
	file := ExecApprovalsFile{
		Version: 1,
		Socket:  &ExecApprovalsSocket{Path: "/tmp/sock", Token: "tok"},
		Defaults: &ExecApprovalsDefaults{
			Security: "strict", Ask: "always",
		},
		Agents: map[string]ExecApprovalsAgent{
			"agent-1": {
				Security: "relaxed",
				Allowlist: []ExecApprovalsAllowlistEntry{
					{Pattern: "ls *", ID: "e1"},
				},
			},
		},
	}
	snap := ExecApprovalsSnapshot{Path: "/etc/approvals", Exists: true, Hash: "abc", File: file}
	data, _ := json.Marshal(snap)
	var got ExecApprovalsSnapshot
	json.Unmarshal(data, &got)
	if !got.Exists {
		t.Error("exists not round-tripped")
	}
	if got.File.Socket == nil || got.File.Socket.Path != "/tmp/sock" {
		t.Error("socket not round-tripped")
	}
	if len(got.File.Agents) != 1 {
		t.Error("agents not round-tripped")
	}

	setP := ExecApprovalsSetParams{File: file, BaseHash: "abc"}
	nodeGetP := ExecApprovalsNodeGetParams{NodeID: "n1"}
	nodeSetP := ExecApprovalsNodeSetParams{NodeID: "n1", File: file}
	reqP := ExecApprovalRequestParams{Command: "ls", ID: "r1"}
	for _, v := range []any{ExecApprovalsGetParams{}, setP, nodeGetP, nodeSetP, reqP} {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Error("empty")
		}
	}
}

// --- Node pair events ---

func TestNodePairEventsRoundTrip(t *testing.T) {
	req := NodePairRequestedEvent{
		RequestID: "r1", NodeID: "n1", DisplayName: "Node", Ts: 123,
	}
	data, _ := json.Marshal(req)
	var gotReq NodePairRequestedEvent
	json.Unmarshal(data, &gotReq)
	if gotReq.NodeID != "n1" {
		t.Errorf("nodeId = %q", gotReq.NodeID)
	}

	res := NodePairResolvedEvent{
		RequestID: "r1", NodeID: "n1", Decision: "approved", Ts: 456,
	}
	data, _ = json.Marshal(res)
	var gotRes NodePairResolvedEvent
	json.Unmarshal(data, &gotRes)
	if gotRes.Decision != "approved" {
		t.Errorf("decision = %q", gotRes.Decision)
	}
}

// --- Web login types ---

func TestWebLoginTypesRoundTrip(t *testing.T) {
	force := true
	sp := WebLoginStartParams{Force: &force, AccountID: "a1"}
	data, _ := json.Marshal(sp)
	var gotS WebLoginStartParams
	json.Unmarshal(data, &gotS)
	if gotS.Force == nil || !*gotS.Force {
		t.Error("force not round-tripped")
	}

	wp := WebLoginWaitParams{AccountID: "a1"}
	data, _ = json.Marshal(wp)
	var gotW WebLoginWaitParams
	json.Unmarshal(data, &gotW)
	if gotW.AccountID != "a1" {
		t.Errorf("accountId = %q", gotW.AccountID)
	}
}

// --- Event payload types ---

func TestPresenceEventRoundTrip(t *testing.T) {
	lastInput := 5.5
	ev := PresenceEvent{
		Presence: []SystemPresence{
			{
				Text: "idle", Ts: 123456, Host: "mac-pro", IP: "192.168.1.10",
				Version: "1.0", Platform: "darwin", DeviceFamily: "mac",
				ModelIdentifier: "Mac14,5", LastInputSeconds: &lastInput,
				Mode: "idle", Reason: "listening", DeviceID: "d1",
				Roles: []string{"operator"}, Scopes: []string{"admin"},
				InstanceID: "i1",
			},
		},
	}
	data, _ := json.Marshal(ev)
	var got PresenceEvent
	json.Unmarshal(data, &got)
	if len(got.Presence) != 1 || got.Presence[0].Host != "mac-pro" {
		t.Errorf("presence mismatch")
	}
	if got.Presence[0].LastInputSeconds == nil || *got.Presence[0].LastInputSeconds != 5.5 {
		t.Error("lastInputSeconds mismatch")
	}
}

func TestHealthEventRoundTrip(t *testing.T) {
	ev := HealthEvent{
		OK: true, Ts: 1000, DurationMs: 50,
		Channels:         map[string]ChannelHealthSummary{"sms": {AccountID: "a1", Configured: boolPtr(true)}},
		ChannelOrder:     []string{"sms"},
		ChannelLabels:    map[string]string{"sms": "SMS"},
		HeartbeatSeconds: 30,
		DefaultAgentID:   "main",
		Agents: []AgentHealthSummary{
			{AgentID: "main", Name: "Main Agent", IsDefault: true, Heartbeat: json.RawMessage(`{}`), Sessions: HealthSessionsSummary{Path: "/s", Count: 1, Recent: []HealthRecentSession{{Key: "k1"}}}},
		},
		Sessions: HealthSessionsSummary{Path: "/sessions", Count: 5, Recent: []HealthRecentSession{{Key: "k2", UpdatedAt: int64Ptr(999), Age: int64Ptr(100)}}},
	}
	data, _ := json.Marshal(ev)
	var got HealthEvent
	json.Unmarshal(data, &got)
	if !got.OK || got.DefaultAgentID != "main" {
		t.Error("health event mismatch")
	}
	if len(got.Agents) != 1 || got.Agents[0].AgentID != "main" {
		t.Error("agents mismatch")
	}
	if got.Sessions.Count != 5 {
		t.Errorf("sessions.count = %d", got.Sessions.Count)
	}
}

func TestHeartbeatEventRoundTrip(t *testing.T) {
	dur := int64(100)
	ev := HeartbeatEvent{
		Ts: 1000, Status: "sent", To: "user1", AccountID: "a1",
		Preview: "Hello", DurationMs: &dur, HasMedia: boolPtr(false),
		Reason: "scheduled", Channel: "sms", Silent: boolPtr(true),
		IndicatorType: "ok",
	}
	data, _ := json.Marshal(ev)
	var got HeartbeatEvent
	json.Unmarshal(data, &got)
	if got.Status != "sent" || got.To != "user1" {
		t.Error("heartbeat event mismatch")
	}
}

func TestVoicewakeChangedEventRoundTrip(t *testing.T) {
	ev := VoicewakeChangedEvent{Triggers: []string{"hey openclaw", "ok computer"}}
	data, _ := json.Marshal(ev)
	var got VoicewakeChangedEvent
	json.Unmarshal(data, &got)
	if len(got.Triggers) != 2 || got.Triggers[0] != "hey openclaw" {
		t.Error("voicewake changed event mismatch")
	}
}

func TestCronEventRoundTrip(t *testing.T) {
	runAt := int64(5000)
	dur := int64(100)
	nextRun := int64(10000)
	ev := CronEvent{
		JobID: "j1", Action: "finished", RunAtMs: &runAt, DurationMs: &dur,
		Status: "ok", Summary: "done", SessionID: "s1", SessionKey: "sk1",
		NextRunAtMs: &nextRun, Model: "gpt-4", Provider: "openai",
		Usage: &CronUsageSummary{InputTokens: intPtr(100), OutputTokens: intPtr(200), TotalTokens: intPtr(300)},
	}
	data, _ := json.Marshal(ev)
	var got CronEvent
	json.Unmarshal(data, &got)
	if got.JobID != "j1" || got.Action != "finished" {
		t.Error("cron event mismatch")
	}
	if got.Usage == nil || *got.Usage.InputTokens != 100 {
		t.Error("usage mismatch")
	}

	// Test with error field.
	evErr := CronEvent{JobID: "j2", Action: "finished", Status: "error", Error: "timeout"}
	data, _ = json.Marshal(evErr)
	var gotErr CronEvent
	json.Unmarshal(data, &gotErr)
	if gotErr.Error != "timeout" {
		t.Errorf("error = %q", gotErr.Error)
	}

	// CronUsageSummary with all optional fields.
	usage := CronUsageSummary{
		InputTokens: intPtr(10), OutputTokens: intPtr(20), TotalTokens: intPtr(30),
		CacheReadTokens: intPtr(5), CacheWriteTokens: intPtr(2),
	}
	data, _ = json.Marshal(usage)
	var gotUsage CronUsageSummary
	json.Unmarshal(data, &gotUsage)
	if gotUsage.CacheReadTokens == nil || *gotUsage.CacheReadTokens != 5 {
		t.Error("cache read tokens mismatch")
	}
}

func TestExecApprovalResolvedEventRoundTrip(t *testing.T) {
	ev := ExecApprovalResolvedEvent{
		ID: "a1", Decision: "allow-once", ResolvedBy: "operator", Ts: 99999,
	}
	data, _ := json.Marshal(ev)
	var got ExecApprovalResolvedEvent
	json.Unmarshal(data, &got)
	if got.Decision != "allow-once" {
		t.Errorf("decision = %q", got.Decision)
	}
}

// --- TTS types ---

func TestTTSTypesRoundTrip(t *testing.T) {
	// TTSStatusResult
	fallback := "elevenlabs"
	status := TTSStatusResult{
		Enabled:           true,
		Auto:              "on",
		Provider:          "openai",
		FallbackProvider:  &fallback,
		FallbackProviders: []string{"elevenlabs", "edge"},
		PrefsPath:         "/tmp/tts.json",
		HasOpenAIKey:      true,
		HasElevenLabsKey:  false,
		EdgeEnabled:       true,
	}
	data, _ := json.Marshal(status)
	var gotStatus TTSStatusResult
	json.Unmarshal(data, &gotStatus)
	if gotStatus.Provider != "openai" {
		t.Errorf("provider = %q", gotStatus.Provider)
	}
	if gotStatus.FallbackProvider == nil || *gotStatus.FallbackProvider != "elevenlabs" {
		t.Error("fallbackProvider mismatch")
	}

	// TTSProvidersResult
	providers := TTSProvidersResult{
		Providers: []TTSProviderInfo{
			{ID: "openai", Name: "OpenAI", Configured: true, Models: []string{"tts-1"}, Voices: []string{"alloy"}},
			{ID: "edge", Name: "Edge TTS", Configured: true, Models: []string{}},
		},
		Active: "openai",
	}
	data, _ = json.Marshal(providers)
	var gotProviders TTSProvidersResult
	json.Unmarshal(data, &gotProviders)
	if len(gotProviders.Providers) != 2 {
		t.Errorf("providers = %d", len(gotProviders.Providers))
	}
	if gotProviders.Active != "openai" {
		t.Errorf("active = %q", gotProviders.Active)
	}

	// TTSEnableResult / TTSDisableResult
	enable := TTSEnableResult{Enabled: true}
	data, _ = json.Marshal(enable)
	var gotEnable TTSEnableResult
	json.Unmarshal(data, &gotEnable)
	if !gotEnable.Enabled {
		t.Error("enabled not true")
	}

	disable := TTSDisableResult{Enabled: false}
	data, _ = json.Marshal(disable)
	var gotDisable TTSDisableResult
	json.Unmarshal(data, &gotDisable)
	if gotDisable.Enabled {
		t.Error("enabled not false")
	}

	// TTSConvertParams / TTSConvertResult
	convert := TTSConvertParams{Text: "hello", Channel: "default"}
	data, _ = json.Marshal(convert)
	var gotConvert TTSConvertParams
	json.Unmarshal(data, &gotConvert)
	if gotConvert.Text != "hello" {
		t.Errorf("text = %q", gotConvert.Text)
	}

	compat := true
	convertResult := TTSConvertResult{
		AudioPath:       "/tmp/out.mp3",
		Provider:        "openai",
		OutputFormat:    "mp3",
		VoiceCompatible: &compat,
	}
	data, _ = json.Marshal(convertResult)
	var gotResult TTSConvertResult
	json.Unmarshal(data, &gotResult)
	if gotResult.AudioPath != "/tmp/out.mp3" {
		t.Errorf("audioPath = %q", gotResult.AudioPath)
	}

	// TTSSetProviderParams / TTSSetProviderResult
	setP := TTSSetProviderParams{Provider: "edge"}
	data, _ = json.Marshal(setP)
	var gotSetP TTSSetProviderParams
	json.Unmarshal(data, &gotSetP)
	if gotSetP.Provider != "edge" {
		t.Errorf("provider = %q", gotSetP.Provider)
	}

	setR := TTSSetProviderResult{Provider: "edge"}
	data, _ = json.Marshal(setR)
	var gotSetR TTSSetProviderResult
	json.Unmarshal(data, &gotSetR)
	if gotSetR.Provider != "edge" {
		t.Errorf("provider = %q", gotSetR.Provider)
	}
}

// --- exec.approval types ---

func TestExecApprovalWaitDecisionTypesRoundTrip(t *testing.T) {
	// Params
	params := ExecApprovalWaitDecisionParams{ID: "a1"}
	data, _ := json.Marshal(params)
	var gotParams ExecApprovalWaitDecisionParams
	json.Unmarshal(data, &gotParams)
	if gotParams.ID != "a1" {
		t.Errorf("id = %q", gotParams.ID)
	}

	// Result with decision
	decision := "allow-once"
	created := int64(1000)
	expires := int64(2000)
	result := ExecApprovalWaitDecisionResult{
		ID: "a1", Decision: &decision, CreatedAtMs: &created, ExpiresAtMs: &expires,
	}
	data, _ = json.Marshal(result)
	var gotResult ExecApprovalWaitDecisionResult
	json.Unmarshal(data, &gotResult)
	if gotResult.ID != "a1" {
		t.Errorf("id = %q", gotResult.ID)
	}
	if gotResult.Decision == nil || *gotResult.Decision != "allow-once" {
		t.Error("decision mismatch")
	}

	// Result with null decision (timeout)
	resultNull := ExecApprovalWaitDecisionResult{ID: "a2", Decision: nil}
	data, _ = json.Marshal(resultNull)
	var gotNull ExecApprovalWaitDecisionResult
	json.Unmarshal(data, &gotNull)
	if gotNull.Decision != nil {
		t.Error("expected nil decision")
	}

	// ExecApprovalRequestResult
	reqResult := ExecApprovalRequestResult{
		ID: "r1", Status: "accepted", CreatedAtMs: 1000, ExpiresAtMs: 2000,
	}
	data, _ = json.Marshal(reqResult)
	var gotReqResult ExecApprovalRequestResult
	json.Unmarshal(data, &gotReqResult)
	if gotReqResult.Status != "accepted" {
		t.Errorf("status = %q", gotReqResult.Status)
	}

	// ExecApprovalResolveResult
	resolveResult := ExecApprovalResolveResult{OK: true}
	data, _ = json.Marshal(resolveResult)
	var gotResolve ExecApprovalResolveResult
	json.Unmarshal(data, &gotResolve)
	if !gotResolve.OK {
		t.Error("ok not true")
	}

	// ExecApprovalResolvedEvent
	resolved := ExecApprovalResolvedEvent{
		ID: "a1", Decision: "deny", ResolvedBy: "user1", Ts: 12345,
	}
	data, _ = json.Marshal(resolved)
	var gotResolved ExecApprovalResolvedEvent
	json.Unmarshal(data, &gotResolved)
	if gotResolved.Decision != "deny" {
		t.Errorf("decision = %q", gotResolved.Decision)
	}
}

func TestSessionGatewayExtendedTypesRoundTrip(t *testing.T) {
	create := SessionsCreateParams{Key: "main", AgentID: "agent1", Message: "hi"}
	data, _ := json.Marshal(create)
	var gotCreate SessionsCreateParams
	json.Unmarshal(data, &gotCreate)
	if gotCreate.Key != "main" || gotCreate.AgentID != "agent1" {
		t.Errorf("sessions.create mismatch: %#v", gotCreate)
	}

	send := SessionsSendParams{Key: "main", Message: "hello"}
	data, _ = json.Marshal(send)
	var gotSend SessionsSendParams
	json.Unmarshal(data, &gotSend)
	if gotSend.Message != "hello" {
		t.Errorf("sessions.send message = %q", gotSend.Message)
	}

	abs := SessionsAbortParams{Key: "main", RunID: "r1"}
	data, _ = json.Marshal(abs)
	var gotAbort SessionsAbortParams
	json.Unmarshal(data, &gotAbort)
	if gotAbort.RunID != "r1" {
		t.Errorf("sessions.abort runId = %q", gotAbort.RunID)
	}
}

func TestNodePendingExtendedTypesRoundTrip(t *testing.T) {
	ack := NodePendingAckParams{IDs: []string{"a1", "a2"}}
	data, _ := json.Marshal(ack)
	var gotAck NodePendingAckParams
	json.Unmarshal(data, &gotAck)
	if len(gotAck.IDs) != 2 {
		t.Errorf("ids len = %d", len(gotAck.IDs))
	}

	pull := NodePendingPullResult{NodeID: "n1", Actions: []NodePendingAction{{ID: "a1", Command: "status"}}}
	data, _ = json.Marshal(pull)
	var gotPull NodePendingPullResult
	json.Unmarshal(data, &gotPull)
	if gotPull.NodeID != "n1" || len(gotPull.Actions) != 1 {
		t.Errorf("node.pending.pull mismatch: %#v", gotPull)
	}

	canvas := NodeCanvasCapabilityRefreshResult{CanvasCapability: "cap", CanvasCapabilityExpiresAtMs: 123, CanvasHostURL: "https://x"}
	data, _ = json.Marshal(canvas)
	var gotCanvas NodeCanvasCapabilityRefreshResult
	json.Unmarshal(data, &gotCanvas)
	if gotCanvas.CanvasCapability != "cap" {
		t.Errorf("canvas capability = %q", gotCanvas.CanvasCapability)
	}
}

func TestConfigSecretsAndSystemTypesRoundTrip(t *testing.T) {
	lookup := ConfigSchemaLookupResult{Path: "session.mainKey", Schema: json.RawMessage(`{"type":"string"}`)}
	data, _ := json.Marshal(lookup)
	var gotLookup ConfigSchemaLookupResult
	json.Unmarshal(data, &gotLookup)
	if gotLookup.Path != "session.mainKey" {
		t.Errorf("lookup path = %q", gotLookup.Path)
	}

	resolve := SecretsResolveResult{Assignments: []SecretsResolveAssignment{{PathSegments: []string{"x"}, Value: json.RawMessage(`"y"`)}}}
	data, _ = json.Marshal(resolve)
	var gotResolve SecretsResolveResult
	json.Unmarshal(data, &gotResolve)
	if len(gotResolve.Assignments) != 1 {
		t.Errorf("assignments len = %d", len(gotResolve.Assignments))
	}

	ident := GatewayIdentityResult{DeviceID: "d1", PublicKey: "pk"}
	data, _ = json.Marshal(ident)
	var gotIdent GatewayIdentityResult
	json.Unmarshal(data, &gotIdent)
	if gotIdent.DeviceID != "d1" {
		t.Errorf("deviceId = %q", gotIdent.DeviceID)
	}
}

// --- Helpers ---

func intPtr(v int) *int       { return &v }
func int64Ptr(v int64) *int64 { return &v }
func boolPtr(v bool) *bool    { return &v }
