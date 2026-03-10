package gateway

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/a3tai/openclaw-go/protocol"
	"github.com/gorilla/websocket"
)

// testMethod is a helper that tests a convenience method against a mock gateway.
// It sets up a mock, connects, and exercises success/error/no-detail paths.
type testMethod struct {
	t              *testing.T
	method         string
	success        func(client *Client, ctx context.Context) error
	successPayload json.RawMessage // optional: custom success response payload (default: `{}`)
}

func (tm *testMethod) run() {
	tm.t.Helper()
	tm.t.Run(tm.method+"/success", func(t *testing.T) {
		mg, wsURL, cleanup := startMockGateway(t)
		defer cleanup()

		payload := tm.successPayload
		if payload == nil {
			payload = json.RawMessage(`{}`)
		}
		mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
			if req.Method == tm.method {
				respData, _ := protocol.MarshalResponse(req.ID, payload)
				conn.WriteMessage(websocket.TextMessage, respData)
			}
		}

		client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Connect(ctx, wsURL); err != nil {
			t.Fatalf("Connect: %v", err)
		}

		if err := tm.success(client, ctx); err != nil {
			t.Fatalf("%s: %v", tm.method, err)
		}
	})

	tm.t.Run(tm.method+"/error_with_payload", func(t *testing.T) {
		mg, wsURL, cleanup := startMockGateway(t)
		defer cleanup()

		mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
			if req.Method == tm.method {
				respData, _ := protocol.MarshalErrorResponse(req.ID, protocol.ErrorPayload{
					Code: "FORBIDDEN", Message: "not allowed",
				})
				conn.WriteMessage(websocket.TextMessage, respData)
			}
		}

		client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Connect(ctx, wsURL); err != nil {
			t.Fatalf("Connect: %v", err)
		}

		err := tm.success(client, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "FORBIDDEN") {
			t.Errorf("error = %q, want to contain 'FORBIDDEN'", err.Error())
		}
	})

	tm.t.Run(tm.method+"/error_no_detail", func(t *testing.T) {
		mg, wsURL, cleanup := startMockGateway(t)
		defer cleanup()

		mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
			if req.Method == tm.method {
				resp := protocol.Response{Type: protocol.FrameTypeResponse, ID: req.ID, OK: false}
				data, _ := json.Marshal(resp)
				conn.WriteMessage(websocket.TextMessage, data)
			}
		}

		client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Connect(ctx, wsURL); err != nil {
			t.Fatalf("Connect: %v", err)
		}

		err := tm.success(client, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "request failed") {
			t.Errorf("error = %q, want to contain 'request failed'", err.Error())
		}
	})
}

// --- Chat methods ---

func TestChatSend(t *testing.T) {
	// Test with typed response
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "chat.send" {
			ev := protocol.ChatEvent{RunID: "run-1", SessionKey: "main", Seq: 0, State: "final"}
			respData, _ := protocol.MarshalResponse(req.ID, ev)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	ev, err := client.ChatSend(ctx, protocol.ChatSendParams{
		SessionKey: "main", Message: "hello", IdempotencyKey: "k1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ev.RunID != "run-1" {
		t.Errorf("runId = %q", ev.RunID)
	}

	// Test error paths
	tm := &testMethod{t: t, method: "chat.send", success: func(c *Client, ctx context.Context) error {
		_, err := c.ChatSend(ctx, protocol.ChatSendParams{SessionKey: "main", Message: "hi", IdempotencyKey: "k1"})
		return err
	}}
	tm.run()
}

func TestChatHistory(t *testing.T) {
	tm := &testMethod{t: t, method: "chat.history", success: func(c *Client, ctx context.Context) error {
		_, err := c.ChatHistory(ctx, protocol.ChatHistoryParams{SessionKey: "main"})
		return err
	}}
	tm.run()
}

func TestChatAbort(t *testing.T) {
	tm := &testMethod{t: t, method: "chat.abort", success: func(c *Client, ctx context.Context) error {
		return c.ChatAbort(ctx, protocol.ChatAbortParams{SessionKey: "main"})
	}}
	tm.run()
}

func TestChatInject(t *testing.T) {
	tm := &testMethod{t: t, method: "chat.inject", success: func(c *Client, ctx context.Context) error {
		return c.ChatInject(ctx, protocol.ChatInjectParams{SessionKey: "main", Message: "injected"})
	}}
	tm.run()
}

// --- Agent methods ---

func TestAgent(t *testing.T) {
	tm := &testMethod{t: t, method: "agent", success: func(c *Client, ctx context.Context) error {
		_, err := c.Agent(ctx, protocol.AgentParams{Message: "hello", IdempotencyKey: "k1"})
		return err
	}}
	tm.run()
}

func TestAgentIdentity(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "agent.identity.get" {
			r := protocol.AgentIdentityResult{AgentID: "a1", Name: "Test"}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.AgentIdentity(ctx, protocol.AgentIdentityParams{AgentID: "a1"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Name != "Test" {
		t.Errorf("name = %q", r.Name)
	}

	tm := &testMethod{t: t, method: "agent.identity.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentIdentity(ctx, protocol.AgentIdentityParams{})
		return err
	}}
	tm.run()
}

func TestAgentWait(t *testing.T) {
	tm := &testMethod{t: t, method: "agent.wait", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentWait(ctx, protocol.AgentWaitParams{RunID: "run-1"})
		return err
	}}
	tm.run()
}

// --- Session methods ---

func TestSessionsGet(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsGet(ctx, protocol.SessionsGetParams{Key: "main"})
		return err
	}}
	tm.run()
}

func TestSessionsList(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsList(ctx, protocol.SessionsListParams{})
		return err
	}}
	tm.run()
}

func TestSessionsPreview(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.preview", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsPreview(ctx, protocol.SessionsPreviewParams{Keys: []string{"main"}})
		return err
	}}
	tm.run()
}

func TestSessionsResolve(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.resolve", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsResolve(ctx, protocol.SessionsResolveParams{Key: "main"})
		return err
	}}
	tm.run()
}

func TestSessionsPatch(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.patch", success: func(c *Client, ctx context.Context) error {
		return c.SessionsPatch(ctx, protocol.SessionsPatchParams{Key: "main"})
	}}
	tm.run()
}

func TestSessionsReset(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.reset", success: func(c *Client, ctx context.Context) error {
		return c.SessionsReset(ctx, protocol.SessionsResetParams{Key: "main"})
	}}
	tm.run()
}

func TestSessionsDelete(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.delete", success: func(c *Client, ctx context.Context) error {
		return c.SessionsDelete(ctx, protocol.SessionsDeleteParams{Key: "main"})
	}}
	tm.run()
}

func TestSessionsCompact(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.compact", success: func(c *Client, ctx context.Context) error {
		return c.SessionsCompact(ctx, protocol.SessionsCompactParams{Key: "main"})
	}}
	tm.run()
}

func TestSessionsUsage(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.usage", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsUsage(ctx, protocol.SessionsUsageParams{Key: "main"})
		return err
	}}
	tm.run()
}

// --- Config methods ---

func TestConfigGet(t *testing.T) {
	tm := &testMethod{t: t, method: "config.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.ConfigGet(ctx)
		return err
	}}
	tm.run()
}

func TestConfigSet(t *testing.T) {
	tm := &testMethod{t: t, method: "config.set", success: func(c *Client, ctx context.Context) error {
		return c.ConfigSet(ctx, protocol.ConfigSetParams{Raw: "yaml"})
	}}
	tm.run()
}

func TestConfigApply(t *testing.T) {
	tm := &testMethod{t: t, method: "config.apply", success: func(c *Client, ctx context.Context) error {
		return c.ConfigApply(ctx, protocol.ConfigApplyParams{Raw: "yaml"})
	}}
	tm.run()
}

func TestConfigPatch(t *testing.T) {
	tm := &testMethod{t: t, method: "config.patch", success: func(c *Client, ctx context.Context) error {
		return c.ConfigPatch(ctx, protocol.ConfigPatchParams{Raw: "yaml"})
	}}
	tm.run()
}

func TestConfigSchema(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "config.schema" {
			r := protocol.ConfigSchemaResponse{
				Schema: json.RawMessage(`{}`), UIHints: map[string]protocol.ConfigUIHint{},
				Version: "1.0.0", GeneratedAt: "now",
			}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.ConfigSchema(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if r.Version != "1.0.0" {
		t.Errorf("version = %q", r.Version)
	}

	tm := &testMethod{t: t, method: "config.schema", success: func(c *Client, ctx context.Context) error {
		_, err := c.ConfigSchema(ctx)
		return err
	}}
	tm.run()
}

// --- Agents CRUD methods ---

func TestAgentsList(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "agents.list" {
			r := protocol.AgentsListResult{DefaultID: "d", MainKey: "main", Scope: "per-sender", Agents: []protocol.AgentSummary{}}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.AgentsList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if r.Scope != "per-sender" {
		t.Errorf("scope = %q", r.Scope)
	}

	tm := &testMethod{t: t, method: "agents.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentsList(ctx)
		return err
	}}
	tm.run()
}

func TestAgentsCreate(t *testing.T) {
	tm := &testMethod{t: t, method: "agents.create", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentsCreate(ctx, protocol.AgentsCreateParams{Name: "a", Workspace: "/ws"})
		return err
	}}
	tm.run()
}

func TestAgentsUpdate(t *testing.T) {
	tm := &testMethod{t: t, method: "agents.update", success: func(c *Client, ctx context.Context) error {
		return c.AgentsUpdate(ctx, protocol.AgentsUpdateParams{AgentID: "a1"})
	}}
	tm.run()
}

func TestAgentsDelete(t *testing.T) {
	tm := &testMethod{t: t, method: "agents.delete", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentsDelete(ctx, protocol.AgentsDeleteParams{AgentID: "a1"})
		return err
	}}
	tm.run()
}

func TestAgentsFilesList(t *testing.T) {
	tm := &testMethod{t: t, method: "agents.files.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentsFilesList(ctx, protocol.AgentsFilesListParams{AgentID: "a1"})
		return err
	}}
	tm.run()
}

func TestAgentsFilesGet(t *testing.T) {
	tm := &testMethod{t: t, method: "agents.files.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentsFilesGet(ctx, protocol.AgentsFilesGetParams{AgentID: "a1", Name: "f1"})
		return err
	}}
	tm.run()
}

func TestAgentsFilesSet(t *testing.T) {
	tm := &testMethod{t: t, method: "agents.files.set", success: func(c *Client, ctx context.Context) error {
		_, err := c.AgentsFilesSet(ctx, protocol.AgentsFilesSetParams{AgentID: "a1", Name: "f1", Content: "data"})
		return err
	}}
	tm.run()
}

// --- Models ---

func TestModelsList(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "models.list" {
			r := protocol.ModelsListResult{Models: []protocol.ModelChoice{{ID: "m1", Name: "Model", Provider: "openai"}}}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.ModelsList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Models) != 1 {
		t.Errorf("models len = %d", len(r.Models))
	}

	tm := &testMethod{t: t, method: "models.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.ModelsList(ctx)
		return err
	}}
	tm.run()
}

// --- Health ---

func TestHealth(t *testing.T) {
	tm := &testMethod{t: t, method: "health", success: func(c *Client, ctx context.Context) error {
		_, err := c.Health(ctx)
		return err
	}}
	tm.run()
}

func TestStatus(t *testing.T) {
	tm := &testMethod{t: t, method: "status", success: func(c *Client, ctx context.Context) error {
		_, err := c.Status(ctx)
		return err
	}}
	tm.run()
}

// --- Logs ---

func TestLogsTail(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "logs.tail" {
			r := protocol.LogsTailResult{File: "test.log", Cursor: 0, Size: 0, Lines: []string{}}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.LogsTail(ctx, protocol.LogsTailParams{})
	if err != nil {
		t.Fatal(err)
	}
	if r.File != "test.log" {
		t.Errorf("file = %q", r.File)
	}

	tm := &testMethod{t: t, method: "logs.tail", success: func(c *Client, ctx context.Context) error {
		_, err := c.LogsTail(ctx, protocol.LogsTailParams{})
		return err
	}}
	tm.run()
}

// --- Node methods ---

func TestNodeList(t *testing.T) {
	tm := &testMethod{t: t, method: "node.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodeList(ctx)
		return err
	}}
	tm.run()
}

func TestNodeDescribe(t *testing.T) {
	tm := &testMethod{t: t, method: "node.describe", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodeDescribe(ctx, protocol.NodeDescribeParams{NodeID: "n1"})
		return err
	}}
	tm.run()
}

func TestNodeInvoke(t *testing.T) {
	tm := &testMethod{t: t, method: "node.invoke", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodeInvoke(ctx, protocol.NodeInvokeParams{NodeID: "n1", Command: "snap", IdempotencyKey: "k1"})
		return err
	}}
	tm.run()
}

func TestNodeInvokeResult(t *testing.T) {
	tm := &testMethod{t: t, method: "node.invoke.result", success: func(c *Client, ctx context.Context) error {
		return c.NodeInvokeResult(ctx, protocol.NodeInvokeResultParams{ID: "i1", NodeID: "n1", OK: true})
	}}
	tm.run()
}

func TestNodeEventMethod(t *testing.T) {
	tm := &testMethod{t: t, method: "node.event", success: func(c *Client, ctx context.Context) error {
		return c.NodeEvent(ctx, protocol.NodeEventParams{Event: "status"})
	}}
	tm.run()
}

func TestNodeRename(t *testing.T) {
	tm := &testMethod{t: t, method: "node.rename", success: func(c *Client, ctx context.Context) error {
		return c.NodeRename(ctx, protocol.NodeRenameParams{NodeID: "n1", DisplayName: "New"})
	}}
	tm.run()
}

// --- Node pairing ---

func TestNodePairRequest(t *testing.T) {
	tm := &testMethod{t: t, method: "node.pair.request", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodePairRequest(ctx, protocol.NodePairRequestParams{NodeID: "n1"})
		return err
	}}
	tm.run()
}

func TestNodePairList(t *testing.T) {
	tm := &testMethod{t: t, method: "node.pair.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodePairList(ctx)
		return err
	}}
	tm.run()
}

func TestNodePairApprove(t *testing.T) {
	tm := &testMethod{t: t, method: "node.pair.approve", success: func(c *Client, ctx context.Context) error {
		return c.NodePairApprove(ctx, protocol.NodePairApproveParams{RequestID: "r1"})
	}}
	tm.run()
}

func TestNodePairReject(t *testing.T) {
	tm := &testMethod{t: t, method: "node.pair.reject", success: func(c *Client, ctx context.Context) error {
		return c.NodePairReject(ctx, protocol.NodePairRejectParams{RequestID: "r1"})
	}}
	tm.run()
}

func TestNodePairVerify(t *testing.T) {
	tm := &testMethod{t: t, method: "node.pair.verify", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodePairVerify(ctx, protocol.NodePairVerifyParams{NodeID: "n1", Token: "tok"})
		return err
	}}
	tm.run()
}

// --- Device pairing ---

func TestDevicePairList(t *testing.T) {
	tm := &testMethod{t: t, method: "device.pair.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.DevicePairList(ctx)
		return err
	}}
	tm.run()
}

func TestDevicePairApprove(t *testing.T) {
	tm := &testMethod{t: t, method: "device.pair.approve", success: func(c *Client, ctx context.Context) error {
		return c.DevicePairApprove(ctx, protocol.DevicePairApproveParams{RequestID: "r1"})
	}}
	tm.run()
}

func TestDevicePairReject(t *testing.T) {
	tm := &testMethod{t: t, method: "device.pair.reject", success: func(c *Client, ctx context.Context) error {
		return c.DevicePairReject(ctx, protocol.DevicePairRejectParams{RequestID: "r1"})
	}}
	tm.run()
}

func TestDevicePairRemove(t *testing.T) {
	tm := &testMethod{t: t, method: "device.pair.remove", success: func(c *Client, ctx context.Context) error {
		return c.DevicePairRemove(ctx, protocol.DevicePairRemoveParams{DeviceID: "d1"})
	}}
	tm.run()
}

func TestDeviceTokenRotate(t *testing.T) {
	tm := &testMethod{t: t, method: "device.token.rotate", success: func(c *Client, ctx context.Context) error {
		_, err := c.DeviceTokenRotate(ctx, protocol.DeviceTokenRotateParams{DeviceID: "d1", Role: "operator"})
		return err
	}}
	tm.run()
}

func TestDeviceTokenRevoke(t *testing.T) {
	tm := &testMethod{t: t, method: "device.token.revoke", success: func(c *Client, ctx context.Context) error {
		return c.DeviceTokenRevoke(ctx, protocol.DeviceTokenRevokeParams{DeviceID: "d1", Role: "operator"})
	}}
	tm.run()
}

// --- Cron methods ---

func TestCronList(t *testing.T) {
	tm := &testMethod{t: t, method: "cron.list", successPayload: json.RawMessage(`[]`), success: func(c *Client, ctx context.Context) error {
		r, err := c.CronList(ctx, protocol.CronListParams{})
		if err != nil {
			return err
		}
		if r == nil {
			t.Error("expected non-nil result")
		}
		return nil
	}}
	tm.run()
}

func TestCronStatus(t *testing.T) {
	tm := &testMethod{t: t, method: "cron.status", success: func(c *Client, ctx context.Context) error {
		_, err := c.CronStatus(ctx)
		return err
	}}
	tm.run()
}

func TestCronAdd(t *testing.T) {
	tm := &testMethod{t: t, method: "cron.add", success: func(c *Client, ctx context.Context) error {
		_, err := c.CronAdd(ctx, protocol.CronAddParams{
			Name: "j1", Schedule: protocol.CronSchedule{Kind: "at", At: "2024-01-01"},
			SessionTarget: "main", WakeMode: "now",
			Payload: protocol.CronPayload{Kind: "systemEvent", Text: "test"},
		})
		return err
	}}
	tm.run()
}

func TestCronUpdate(t *testing.T) {
	tm := &testMethod{t: t, method: "cron.update", success: func(c *Client, ctx context.Context) error {
		return c.CronUpdate(ctx, protocol.CronUpdateParams{ID: "j1", Patch: protocol.CronJobPatch{Name: "updated"}})
	}}
	tm.run()
}

func TestCronRemove(t *testing.T) {
	tm := &testMethod{t: t, method: "cron.remove", success: func(c *Client, ctx context.Context) error {
		return c.CronRemove(ctx, protocol.CronRemoveParams{ID: "j1"})
	}}
	tm.run()
}

func TestCronRun(t *testing.T) {
	tm := &testMethod{t: t, method: "cron.run", success: func(c *Client, ctx context.Context) error {
		return c.CronRun(ctx, protocol.CronRunParams{ID: "j1"})
	}}
	tm.run()
}

func TestCronRuns(t *testing.T) {
	tm := &testMethod{t: t, method: "cron.runs", successPayload: json.RawMessage(`[]`), success: func(c *Client, ctx context.Context) error {
		r, err := c.CronRuns(ctx, protocol.CronRunsParams{ID: "j1"})
		if err != nil {
			return err
		}
		if r == nil {
			t.Error("expected non-nil result")
		}
		return nil
	}}
	tm.run()
}

// --- Exec approvals admin ---

func TestExecApprovalsGet(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "exec.approvals.get" {
			r := protocol.ExecApprovalsSnapshot{Path: "/etc/a", Exists: true, Hash: "abc", File: protocol.ExecApprovalsFile{Version: 1}}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.ExecApprovalsGet(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Exists {
		t.Error("exists = false")
	}

	tm := &testMethod{t: t, method: "exec.approvals.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.ExecApprovalsGet(ctx)
		return err
	}}
	tm.run()
}

func TestExecApprovalsSet(t *testing.T) {
	tm := &testMethod{t: t, method: "exec.approvals.set", success: func(c *Client, ctx context.Context) error {
		return c.ExecApprovalsSet(ctx, protocol.ExecApprovalsSetParams{File: protocol.ExecApprovalsFile{Version: 1}})
	}}
	tm.run()
}

func TestExecApprovalsNodeGet(t *testing.T) {
	tm := &testMethod{t: t, method: "exec.approvals.node.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.ExecApprovalsNodeGet(ctx, protocol.ExecApprovalsNodeGetParams{NodeID: "n1"})
		return err
	}}
	tm.run()
}

func TestExecApprovalsNodeSet(t *testing.T) {
	tm := &testMethod{t: t, method: "exec.approvals.node.set", success: func(c *Client, ctx context.Context) error {
		return c.ExecApprovalsNodeSet(ctx, protocol.ExecApprovalsNodeSetParams{NodeID: "n1", File: protocol.ExecApprovalsFile{Version: 1}})
	}}
	tm.run()
}

func TestExecApprovalRequestMethod(t *testing.T) {
	tm := &testMethod{t: t, method: "exec.approval.request", success: func(c *Client, ctx context.Context) error {
		_, err := c.ExecApprovalRequest(ctx, protocol.ExecApprovalRequestParams{Command: "ls"})
		return err
	}}
	tm.run()
}

// --- Skills ---

func TestSkillsStatus(t *testing.T) {
	tm := &testMethod{t: t, method: "skills.status", success: func(c *Client, ctx context.Context) error {
		_, err := c.SkillsStatus(ctx, protocol.SkillsStatusParams{})
		return err
	}}
	tm.run()
}

func TestSkillsBins(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "skills.bins" {
			r := protocol.SkillsBinsResult{Bins: []string{"bin1"}}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.SkillsBins(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Bins) != 1 {
		t.Errorf("bins len = %d", len(r.Bins))
	}

	tm := &testMethod{t: t, method: "skills.bins", success: func(c *Client, ctx context.Context) error {
		_, err := c.SkillsBins(ctx)
		return err
	}}
	tm.run()
}

func TestSkillsInstall(t *testing.T) {
	tm := &testMethod{t: t, method: "skills.install", success: func(c *Client, ctx context.Context) error {
		_, err := c.SkillsInstall(ctx, protocol.SkillsInstallParams{Name: "sk1", InstallID: "i1"})
		return err
	}}
	tm.run()
}

func TestSkillsUpdate(t *testing.T) {
	tm := &testMethod{t: t, method: "skills.update", success: func(c *Client, ctx context.Context) error {
		return c.SkillsUpdate(ctx, protocol.SkillsUpdateParams{SkillKey: "sk1"})
	}}
	tm.run()
}

// --- Wizard ---

func TestWizardStart(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "wizard.start" {
			r := protocol.WizardStartResult{SessionID: "s1", Done: false, Status: "running"}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.WizardStart(ctx, protocol.WizardStartParams{})
	if err != nil {
		t.Fatal(err)
	}
	if r.SessionID != "s1" {
		t.Errorf("sessionId = %q", r.SessionID)
	}

	tm := &testMethod{t: t, method: "wizard.start", success: func(c *Client, ctx context.Context) error {
		_, err := c.WizardStart(ctx, protocol.WizardStartParams{})
		return err
	}}
	tm.run()
}

func TestWizardNext(t *testing.T) {
	tm := &testMethod{t: t, method: "wizard.next", success: func(c *Client, ctx context.Context) error {
		_, err := c.WizardNext(ctx, protocol.WizardNextParams{SessionID: "s1"})
		return err
	}}
	tm.run()
}

func TestWizardCancel(t *testing.T) {
	tm := &testMethod{t: t, method: "wizard.cancel", success: func(c *Client, ctx context.Context) error {
		return c.WizardCancel(ctx, protocol.WizardCancelParams{SessionID: "s1"})
	}}
	tm.run()
}

func TestWizardStatus(t *testing.T) {
	tm := &testMethod{t: t, method: "wizard.status", success: func(c *Client, ctx context.Context) error {
		_, err := c.WizardStatus(ctx, protocol.WizardStatusParams{SessionID: "s1"})
		return err
	}}
	tm.run()
}

// --- Channels / Talk ---

func TestChannelsStatus(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "channels.status" {
			r := protocol.ChannelsStatusResult{
				Ts: 1, ChannelOrder: []string{}, ChannelLabels: map[string]string{},
				Channels: map[string]json.RawMessage{}, ChannelAccounts: map[string][]protocol.ChannelAccountSnapshot{},
				ChannelDefaultAccountID: map[string]string{},
			}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.ChannelsStatus(ctx, protocol.ChannelsStatusParams{})
	if err != nil {
		t.Fatal(err)
	}
	if r.Ts != 1 {
		t.Errorf("ts = %d", r.Ts)
	}

	tm := &testMethod{t: t, method: "channels.status", success: func(c *Client, ctx context.Context) error {
		_, err := c.ChannelsStatus(ctx, protocol.ChannelsStatusParams{})
		return err
	}}
	tm.run()
}

func TestChannelsLogout(t *testing.T) {
	tm := &testMethod{t: t, method: "channels.logout", success: func(c *Client, ctx context.Context) error {
		return c.ChannelsLogout(ctx, protocol.ChannelsLogoutParams{Channel: "slack"})
	}}
	tm.run()
}

func TestTalkConfig(t *testing.T) {
	tm := &testMethod{t: t, method: "talk.config", success: func(c *Client, ctx context.Context) error {
		_, err := c.TalkConfig(ctx, protocol.TalkConfigParams{})
		return err
	}}
	tm.run()
}

func TestTalkMode(t *testing.T) {
	tm := &testMethod{t: t, method: "talk.mode", success: func(c *Client, ctx context.Context) error {
		return c.TalkMode(ctx, protocol.TalkModeParams{Enabled: true})
	}}
	tm.run()
}

// --- Send / Wake / System ---

func TestSendMessage(t *testing.T) {
	tm := &testMethod{t: t, method: "send", success: func(c *Client, ctx context.Context) error {
		_, err := c.SendMessage(ctx, protocol.SendParams{To: "user", IdempotencyKey: "k1"})
		return err
	}}
	tm.run()
}

func TestWake(t *testing.T) {
	tm := &testMethod{t: t, method: "wake", success: func(c *Client, ctx context.Context) error {
		return c.Wake(ctx, protocol.WakeParams{Mode: "now", Text: "hello"})
	}}
	tm.run()
}

func TestLastHeartbeat(t *testing.T) {
	tm := &testMethod{t: t, method: "last-heartbeat", success: func(c *Client, ctx context.Context) error {
		_, err := c.LastHeartbeat(ctx)
		return err
	}}
	tm.run()
}

func TestSetHeartbeats(t *testing.T) {
	tm := &testMethod{t: t, method: "set-heartbeats", success: func(c *Client, ctx context.Context) error {
		return c.SetHeartbeats(ctx, true)
	}}
	tm.run()
}

func TestSystemEvent(t *testing.T) {
	tm := &testMethod{t: t, method: "system-event", success: func(c *Client, ctx context.Context) error {
		return c.SystemEvent(ctx, map[string]string{"type": "test"})
	}}
	tm.run()
}

// --- Misc ---

func TestUpdateRun(t *testing.T) {
	tm := &testMethod{t: t, method: "update.run", success: func(c *Client, ctx context.Context) error {
		_, err := c.UpdateRun(ctx, protocol.UpdateRunParams{})
		return err
	}}
	tm.run()
}

func TestPushTest(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "push.test" {
			r := protocol.PushTestResult{OK: true, Status: 200, TokenSuffix: "abc", Topic: "com.test", Environment: "sandbox"}
			respData, _ := protocol.MarshalResponse(req.ID, r)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	r, err := client.PushTest(ctx, protocol.PushTestParams{NodeID: "n1"})
	if err != nil {
		t.Fatal(err)
	}
	if !r.OK {
		t.Error("ok = false")
	}

	tm := &testMethod{t: t, method: "push.test", success: func(c *Client, ctx context.Context) error {
		_, err := c.PushTest(ctx, protocol.PushTestParams{NodeID: "n1"})
		return err
	}}
	tm.run()
}

func TestBrowserRequest(t *testing.T) {
	tm := &testMethod{t: t, method: "browser.request", success: func(c *Client, ctx context.Context) error {
		_, err := c.BrowserRequest(ctx, map[string]string{"url": "https://example.com"})
		return err
	}}
	tm.run()
}

func TestVoiceWakeGet(t *testing.T) {
	tm := &testMethod{t: t, method: "voicewake.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.VoiceWakeGet(ctx)
		return err
	}}
	tm.run()
}

func TestVoiceWakeSet(t *testing.T) {
	tm := &testMethod{t: t, method: "voicewake.set", success: func(c *Client, ctx context.Context) error {
		return c.VoiceWakeSet(ctx, map[string]bool{"enabled": true})
	}}
	tm.run()
}

func TestUsageStatus(t *testing.T) {
	tm := &testMethod{t: t, method: "usage.status", success: func(c *Client, ctx context.Context) error {
		_, err := c.UsageStatus(ctx)
		return err
	}}
	tm.run()
}

func TestUsageCost(t *testing.T) {
	tm := &testMethod{t: t, method: "usage.cost", success: func(c *Client, ctx context.Context) error {
		_, err := c.UsageCost(ctx, map[string]string{"period": "2024-01"})
		return err
	}}
	tm.run()
}

func TestTTSStatus(t *testing.T) {
	tm := &testMethod{t: t, method: "tts.status", success: func(c *Client, ctx context.Context) error {
		_, err := c.TTSStatus(ctx)
		return err
	}}
	tm.run()
}

func TestPoll(t *testing.T) {
	tm := &testMethod{t: t, method: "poll", success: func(c *Client, ctx context.Context) error {
		_, err := c.Poll(ctx, protocol.PollParams{To: "user", Question: "?", Options: []string{"a", "b"}, IdempotencyKey: "k1"})
		return err
	}}
	tm.run()
}

// --- TTS Methods ---

func TestTTSProviders(t *testing.T) {
	tm := &testMethod{t: t, method: "tts.providers", success: func(c *Client, ctx context.Context) error {
		_, err := c.TTSProviders(ctx)
		return err
	}}
	tm.run()
}

func TestTTSEnable(t *testing.T) {
	tm := &testMethod{t: t, method: "tts.enable", success: func(c *Client, ctx context.Context) error {
		_, err := c.TTSEnable(ctx)
		return err
	}}
	tm.run()
}

func TestTTSDisable(t *testing.T) {
	tm := &testMethod{t: t, method: "tts.disable", success: func(c *Client, ctx context.Context) error {
		_, err := c.TTSDisable(ctx)
		return err
	}}
	tm.run()
}

func TestTTSConvert(t *testing.T) {
	tm := &testMethod{t: t, method: "tts.convert", success: func(c *Client, ctx context.Context) error {
		_, err := c.TTSConvert(ctx, protocol.TTSConvertParams{Text: "hello"})
		return err
	}}
	tm.run()
}

func TestTTSSetProvider(t *testing.T) {
	tm := &testMethod{t: t, method: "tts.setProvider", success: func(c *Client, ctx context.Context) error {
		_, err := c.TTSSetProvider(ctx, protocol.TTSSetProviderParams{Provider: "openai"})
		return err
	}}
	tm.run()
}

// --- exec.approval.waitDecision ---

func TestExecApprovalWaitDecision(t *testing.T) {
	tm := &testMethod{t: t, method: "exec.approval.waitDecision", success: func(c *Client, ctx context.Context) error {
		_, err := c.ExecApprovalWaitDecision(ctx, protocol.ExecApprovalWaitDecisionParams{ID: "approval-1"})
		return err
	}}
	tm.run()
}

// --- sendRPCTyped unmarshal error ---

func TestSendRPCTypedUnmarshalError(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		// Return a payload that can't be unmarshalled into the expected type
		raw := `{"type":"res","id":"` + req.ID + `","ok":true,"payload":"not an object"}`
		conn.WriteMessage(websocket.TextMessage, []byte(raw))
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatal(err)
	}

	_, err := client.ModelsList(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("error = %q, want to contain 'unmarshal'", err.Error())
	}
}
