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

func intPtr(v int) *int { return &v }

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
	// Test with typed response matching real server shape.
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "chat.send" {
			result := protocol.ChatSendResult{RunID: "run-1", Status: "started"}
			respData, _ := protocol.MarshalResponse(req.ID, result)
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

	result, err := client.ChatSend(ctx, protocol.ChatSendParams{
		SessionKey: "main", Message: "hello", IdempotencyKey: "k1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.RunID != "run-1" {
		t.Errorf("runId = %q", result.RunID)
	}
	if result.Status != "started" {
		t.Errorf("status = %q, want started", result.Status)
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

func TestSessionsCreate(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.create", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsCreate(ctx, protocol.SessionsCreateParams{Key: "main"})
		return err
	}}
	tm.run()
}

func TestSessionsSend(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.send", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsSend(ctx, protocol.SessionsSendParams{Key: "main", Message: "hello"})
		return err
	}}
	tm.run()
}

func TestSessionsSteer(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.steer", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsSteer(ctx, protocol.SessionsSendParams{Key: "main", Message: "hello"})
		return err
	}}
	tm.run()
}

func TestSessionsAbort(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.abort", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsAbort(ctx, protocol.SessionsAbortParams{Key: "main"})
		return err
	}}
	tm.run()
}

func TestSessionsSubscribe(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.subscribe", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsSubscribe(ctx)
		return err
	}}
	tm.run()
}

func TestSessionsUnsubscribe(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.unsubscribe", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsUnsubscribe(ctx)
		return err
	}}
	tm.run()
}

func TestSessionsMessagesSubscribe(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.messages.subscribe", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsMessagesSubscribe(ctx, protocol.SessionsMessagesSubscribeParams{Key: "main"})
		return err
	}}
	tm.run()
}

func TestSessionsMessagesUnsubscribe(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.messages.unsubscribe", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsMessagesUnsubscribe(ctx, protocol.SessionsMessagesUnsubscribeParams{Key: "main"})
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

func TestSessionsCompactionList(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.compaction.list", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsCompactionList(ctx, protocol.SessionsCompactionListParams{Key: "main"})
		return err
	}}
	tm.run()
}

func TestSessionsCompactionGet(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.compaction.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsCompactionGet(ctx, protocol.SessionsCompactionGetParams{Key: "main", CheckpointID: "cp1"})
		return err
	}}
	tm.run()
}

func TestSessionsCompactionBranch(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.compaction.branch", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsCompactionBranch(ctx, protocol.SessionsCompactionBranchParams{Key: "main", CheckpointID: "cp1"})
		return err
	}}
	tm.run()
}

func TestSessionsCompactionRestore(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.compaction.restore", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsCompactionRestore(ctx, protocol.SessionsCompactionRestoreParams{Key: "main", CheckpointID: "cp1"})
		return err
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

func TestSessionsUsageTimeseries(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.usage.timeseries", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsUsageTimeseries(ctx, protocol.SessionsUsageTimeseriesParams{Key: "main"})
		return err
	}}
	tm.run()
}

func TestSessionsUsageLogs(t *testing.T) {
	tm := &testMethod{t: t, method: "sessions.usage.logs", success: func(c *Client, ctx context.Context) error {
		_, err := c.SessionsUsageLogs(ctx, protocol.SessionsUsageLogsParams{Key: "main"})
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

func TestConfigSchemaLookup(t *testing.T) {
	tm := &testMethod{t: t, method: "config.schema.lookup", success: func(c *Client, ctx context.Context) error {
		_, err := c.ConfigSchemaLookup(ctx, protocol.ConfigSchemaLookupParams{Path: "session.mainKey"})
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

func TestNodePendingEnqueue(t *testing.T) {
	m := &testMethod{t: t, method: "node.pending.enqueue", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodePendingEnqueue(ctx, protocol.NodePendingEnqueueParams{NodeID: "n1", Type: "status.request"})
		return err
	}}
	m.run()
}

func TestNodePendingDrain(t *testing.T) {
	m := &testMethod{t: t, method: "node.pending.drain", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodePendingDrain(ctx, protocol.NodePendingDrainParams{})
		return err
	}}
	m.run()
}

func TestNodePendingPull(t *testing.T) {
	m := &testMethod{t: t, method: "node.pending.pull", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodePendingPull(ctx)
		return err
	}}
	m.run()
}

func TestNodePendingAck(t *testing.T) {
	m := &testMethod{t: t, method: "node.pending.ack", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodePendingAck(ctx, protocol.NodePendingAckParams{IDs: []string{"a1"}})
		return err
	}}
	m.run()
}

func TestNodeCanvasCapabilityRefresh(t *testing.T) {
	m := &testMethod{t: t, method: "node.canvas.capability.refresh", success: func(c *Client, ctx context.Context) error {
		_, err := c.NodeCanvasCapabilityRefresh(ctx)
		return err
	}}
	m.run()
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
	tm := &testMethod{t: t, method: "cron.list", successPayload: json.RawMessage(`{"jobs":[],"total":0,"offset":0,"limit":50,"hasMore":false,"nextOffset":null}`), success: func(c *Client, ctx context.Context) error {
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

func TestCronListPaginated(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "cron.list" {
			page := protocol.CronListResult{
				Jobs: []protocol.CronJob{
					{ID: "j1", Name: "daily", Enabled: true, Schedule: protocol.CronSchedule{Kind: "cron", Expr: "0 9 * * *"}, SessionTarget: "main", WakeMode: "now", Payload: protocol.CronPayload{Kind: "systemEvent", Text: "hi"}},
					{ID: "j2", Name: "weekly", Enabled: false, Schedule: protocol.CronSchedule{Kind: "cron", Expr: "0 0 * * 0"}, SessionTarget: "main", WakeMode: "now", Payload: protocol.CronPayload{Kind: "systemEvent", Text: "hi"}},
				},
				Total:      5,
				Offset:     0,
				Limit:      2,
				HasMore:    true,
				NextOffset: intPtr(2),
			}
			respData, _ := protocol.MarshalResponse(req.ID, page)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx := context.Background()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	r, err := client.CronList(ctx, protocol.CronListParams{})
	if err != nil {
		t.Fatalf("CronList: %v", err)
	}
	if len(r.Jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(r.Jobs))
	}
	if r.Jobs[0].ID != "j1" || r.Jobs[1].ID != "j2" {
		t.Errorf("unexpected job IDs: %s, %s", r.Jobs[0].ID, r.Jobs[1].ID)
	}
	if r.Total != 5 {
		t.Errorf("expected total=5, got %d", r.Total)
	}
	if !r.HasMore {
		t.Error("expected hasMore=true")
	}
	if r.NextOffset == nil || *r.NextOffset != 2 {
		t.Error("expected nextOffset=2")
	}
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
	tm := &testMethod{t: t, method: "cron.runs", successPayload: json.RawMessage(`{"entries":[],"total":0,"offset":0}`), success: func(c *Client, ctx context.Context) error {
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

func TestCronRunsPaginated(t *testing.T) {
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "cron.runs" {
			page := protocol.CronRunsResult{
				Entries: []protocol.CronRunLogEntry{
					{Ts: 1000, JobID: "j1", Action: "finished", Status: "ok", Summary: "done", JobName: "daily"},
					{Ts: 2000, JobID: "j1", Action: "finished", Status: "error", Error: "timeout"},
				},
				Total:  10,
				Offset: 0,
			}
			respData, _ := protocol.MarshalResponse(req.ID, page)
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}

	client := NewClient(WithToken("tok"), WithConnectTimeout(5*time.Second))
	defer client.Close()
	ctx := context.Background()
	if err := client.Connect(ctx, wsURL); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	r, err := client.CronRuns(ctx, protocol.CronRunsParams{ID: "j1"})
	if err != nil {
		t.Fatalf("CronRuns: %v", err)
	}
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
	if r.Entries[0].Status != "ok" || r.Entries[1].Status != "error" {
		t.Errorf("unexpected statuses: %s, %s", r.Entries[0].Status, r.Entries[1].Status)
	}
	if r.Entries[0].JobName != "daily" {
		t.Errorf("expected jobName=daily, got %s", r.Entries[0].JobName)
	}
	if r.Total != 10 {
		t.Errorf("expected total=10, got %d", r.Total)
	}
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

func TestExecApprovalResolveMethod(t *testing.T) {
	tm := &testMethod{t: t, method: "exec.approval.resolve", success: func(c *Client, ctx context.Context) error {
		_, err := c.ExecApprovalResolve(ctx, protocol.ExecApprovalResolveParams{ID: "a1", Decision: "approve"})
		return err
	}}
	tm.run()
}

func TestExecApprovalList(t *testing.T) {
	tm := &testMethod{t: t, method: "exec.approval.list", successPayload: json.RawMessage(`[]`), success: func(c *Client, ctx context.Context) error {
		_, err := c.ExecApprovalList(ctx)
		return err
	}}
	tm.run()
}

// --- Plugin Approvals ---

func TestPluginApprovalRequestMethod(t *testing.T) {
	tm := &testMethod{t: t, method: "plugin.approval.request", success: func(c *Client, ctx context.Context) error {
		_, err := c.PluginApprovalRequest(ctx, protocol.PluginApprovalRequestParams{Title: "Allow tool", Description: "Plugin wants to run a tool"})
		return err
	}}
	tm.run()
}

func TestPluginApprovalWaitDecisionMethod(t *testing.T) {
	tm := &testMethod{t: t, method: "plugin.approval.waitDecision", success: func(c *Client, ctx context.Context) error {
		_, err := c.PluginApprovalWaitDecision(ctx, protocol.PluginApprovalWaitDecisionParams{ID: "p1"})
		return err
	}}
	tm.run()
}

func TestPluginApprovalResolveMethod(t *testing.T) {
	tm := &testMethod{t: t, method: "plugin.approval.resolve", success: func(c *Client, ctx context.Context) error {
		_, err := c.PluginApprovalResolve(ctx, protocol.PluginApprovalResolveParams{ID: "p1", Decision: "allow-once"})
		return err
	}}
	tm.run()
}

func TestPluginApprovalList(t *testing.T) {
	tm := &testMethod{t: t, method: "plugin.approval.list", successPayload: json.RawMessage(`[]`), success: func(c *Client, ctx context.Context) error {
		_, err := c.PluginApprovalList(ctx)
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

func TestTalkSpeak(t *testing.T) {
	m := &testMethod{t: t, method: "talk.speak", success: func(c *Client, ctx context.Context) error {
		_, err := c.TalkSpeak(ctx, protocol.TalkSpeakParams{Text: "hello"})
		return err
	}}
	m.run()
}

func TestWebLoginStart(t *testing.T) {
	tm := &testMethod{t: t, method: "web.login.start", success: func(c *Client, ctx context.Context) error {
		_, err := c.WebLoginStart(ctx, protocol.WebLoginStartParams{})
		return err
	}}
	tm.run()
}

func TestWebLoginWait(t *testing.T) {
	tm := &testMethod{t: t, method: "web.login.wait", success: func(c *Client, ctx context.Context) error {
		_, err := c.WebLoginWait(ctx, protocol.WebLoginWaitParams{})
		return err
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

func TestGatewayIdentityGet(t *testing.T) {
	tm := &testMethod{t: t, method: "gateway.identity.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.GatewayIdentityGet(ctx)
		return err
	}}
	tm.run()
}

func TestSystemPresence(t *testing.T) {
	tm := &testMethod{t: t, method: "system-presence", successPayload: json.RawMessage(`[]`), success: func(c *Client, ctx context.Context) error {
		_, err := c.Presence(ctx)
		return err
	}}
	tm.run()
}

func TestDoctorMemoryStatus(t *testing.T) {
	tm := &testMethod{t: t, method: "doctor.memory.status", success: func(c *Client, ctx context.Context) error {
		_, err := c.DoctorMemoryStatus(ctx)
		return err
	}}
	tm.run()
}

func TestDoctorMemoryDreamDiary(t *testing.T) {
	tm := &testMethod{t: t, method: "doctor.memory.dreamDiary", success: func(c *Client, ctx context.Context) error {
		_, err := c.DoctorMemoryDreamDiary(ctx)
		return err
	}}
	tm.run()
}

// --- Commands ---

func TestCommandsList(t *testing.T) {
	tm := &testMethod{t: t, method: "commands.list", successPayload: json.RawMessage(`{"commands":[]}`), success: func(c *Client, ctx context.Context) error {
		_, err := c.CommandsList(ctx, protocol.CommandsListParams{})
		return err
	}}
	tm.run()
}

func TestMessageAction(t *testing.T) {
	tm := &testMethod{t: t, method: "message.action", successPayload: json.RawMessage(`{"ok":true}`), success: func(c *Client, ctx context.Context) error {
		_, err := c.MessageAction(ctx, protocol.MessageActionParams{
			Channel:        "slack",
			Action:         "open-thread",
			Params:         map[string]any{"messageId": "123"},
			IdempotencyKey: "idem-message-action",
		})
		return err
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

func TestSecretsReload(t *testing.T) {
	tm := &testMethod{t: t, method: "secrets.reload", success: func(c *Client, ctx context.Context) error {
		return c.SecretsReload(ctx)
	}}
	tm.run()
}

func TestSecretsResolve(t *testing.T) {
	tm := &testMethod{t: t, method: "secrets.resolve", success: func(c *Client, ctx context.Context) error {
		_, err := c.SecretsResolve(ctx, protocol.SecretsResolveParams{CommandName: "send", TargetIDs: []string{"slack"}})
		return err
	}}
	tm.run()
}

func TestToolsCatalog(t *testing.T) {
	tm := &testMethod{t: t, method: "tools.catalog", success: func(c *Client, ctx context.Context) error {
		_, err := c.ToolsCatalog(ctx, protocol.ToolsCatalogParams{})
		return err
	}}
	tm.run()
}

func TestToolsEffective(t *testing.T) {
	tm := &testMethod{t: t, method: "tools.effective", success: func(c *Client, ctx context.Context) error {
		_, err := c.ToolsEffective(ctx, protocol.ToolsEffectiveParams{SessionKey: "main"})
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

// --- skills.search ---

func TestSkillsSearch(t *testing.T) {
	// Typed round-trip: verify the response fields decode correctly.
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "skills.search" {
			result := protocol.SkillsSearchResult{
				Results: []protocol.SkillsSearchItem{
					{Score: 0.95, Slug: "my-skill", DisplayName: "My Skill", Summary: "does stuff", Version: "1.0.0"},
				},
			}
			respData, _ := protocol.MarshalResponse(req.ID, result)
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

	result, err := client.SkillsSearch(ctx, protocol.SkillsSearchParams{Query: "my-skill"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
	if result.Results[0].Slug != "my-skill" {
		t.Errorf("slug = %q, want my-skill", result.Results[0].Slug)
	}
	if result.Results[0].Score != 0.95 {
		t.Errorf("score = %v, want 0.95", result.Results[0].Score)
	}

	// Error paths.
	tm := &testMethod{t: t, method: "skills.search", success: func(c *Client, ctx context.Context) error {
		_, err := c.SkillsSearch(ctx, protocol.SkillsSearchParams{Query: "foo"})
		return err
	}}
	tm.run()
}

// --- skills.detail ---

func TestSkillsDetail(t *testing.T) {
	// Typed round-trip: verify nullable fields decode correctly.
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	handle := "author-handle"
	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "skills.detail" {
			result := protocol.SkillsDetailResult{
				Skill: &protocol.SkillsDetailSkill{
					Slug:        "my-skill",
					DisplayName: "My Skill",
					CreatedAt:   1700000000,
					UpdatedAt:   1700000001,
				},
				LatestVersion: &protocol.SkillsDetailVersion{
					Version:   "1.2.3",
					CreatedAt: 1700000001,
				},
				Owner: &protocol.SkillsDetailOwner{Handle: &handle},
			}
			respData, _ := protocol.MarshalResponse(req.ID, result)
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

	result, err := client.SkillsDetail(ctx, protocol.SkillsDetailParams{Slug: "my-skill"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Skill == nil || result.Skill.Slug != "my-skill" {
		t.Errorf("skill.slug = %v, want my-skill", result.Skill)
	}
	if result.LatestVersion == nil || result.LatestVersion.Version != "1.2.3" {
		t.Errorf("latestVersion.version = %v, want 1.2.3", result.LatestVersion)
	}
	if result.Owner == nil || result.Owner.Handle == nil || *result.Owner.Handle != "author-handle" {
		t.Errorf("owner.handle = %v, want author-handle", result.Owner)
	}

	// Null skill field — server returns null for unknown slug.
	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "skills.detail" {
			respData, _ := protocol.MarshalResponse(req.ID, protocol.SkillsDetailResult{Skill: nil})
			conn.WriteMessage(websocket.TextMessage, respData)
		}
	}
	nullResult, err := client.SkillsDetail(ctx, protocol.SkillsDetailParams{Slug: "unknown"})
	if err != nil {
		t.Fatal(err)
	}
	if nullResult.Skill != nil {
		t.Errorf("expected nil skill for unknown slug, got %v", nullResult.Skill)
	}

	// Error paths.
	tm := &testMethod{t: t, method: "skills.detail", success: func(c *Client, ctx context.Context) error {
		_, err := c.SkillsDetail(ctx, protocol.SkillsDetailParams{Slug: "foo"})
		return err
	}}
	tm.run()
}

// --- exec.approval.get ---

func TestExecApprovalGet(t *testing.T) {
	// Typed round-trip: verify response fields decode correctly.
	mg, wsURL, cleanup := startMockGateway(t)
	defer cleanup()

	agentID := "agent-123"
	mg.onRequest = func(conn *websocket.Conn, req protocol.Request) {
		if req.Method == "exec.approval.get" {
			result := protocol.ExecApprovalGetResult{
				ID:               "approval-abc",
				CommandText:      "rm -rf /tmp/test",
				AllowedDecisions: []string{"allow", "deny"},
				AgentID:          &agentID,
				ExpiresAtMs:      9999999999,
			}
			respData, _ := protocol.MarshalResponse(req.ID, result)
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

	result, err := client.ExecApprovalGet(ctx, protocol.ExecApprovalGetParams{ID: "approval-abc"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "approval-abc" {
		t.Errorf("id = %q, want approval-abc", result.ID)
	}
	if result.CommandText != "rm -rf /tmp/test" {
		t.Errorf("commandText = %q", result.CommandText)
	}
	if len(result.AllowedDecisions) != 2 {
		t.Errorf("allowedDecisions = %v, want [allow deny]", result.AllowedDecisions)
	}
	if result.AgentID == nil || *result.AgentID != "agent-123" {
		t.Errorf("agentId = %v, want agent-123", result.AgentID)
	}

	// Error paths.
	tm := &testMethod{t: t, method: "exec.approval.get", success: func(c *Client, ctx context.Context) error {
		_, err := c.ExecApprovalGet(ctx, protocol.ExecApprovalGetParams{ID: "approval-abc"})
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
