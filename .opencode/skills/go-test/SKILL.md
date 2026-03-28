---
name: go-test
description: Run go test with race detector, analyze failures, fix issues, and re-run until green
license: MIT
metadata:
  audience: developers
  workflow: testing
---

## What I do

Run the Go test suite with the race detector, analyze any failures, apply fixes, and iterate until all tests pass.

## When to use me

Use this skill before committing changes, before opening a PR, or when asked to fix failing tests.

## Workflow

1. **Run the full test suite**:
   ```sh
   go test ./... -race -v
   ```

2. **Analyze failures** — for each failing test:
   - Identify the test function and file (e.g., `gateway/methods_test.go:TestChatSend`)
   - Read the test code and the code under test
   - Determine root cause: type mismatch, missing mock, logic error, race condition, etc.

3. **Fix the issue** — apply the minimal correct fix:
   - If the test expectation is wrong (API changed), update the test
   - If the implementation is wrong, fix the implementation
   - If a new type/method is missing, add it
   - Never delete or skip tests to make the suite pass

4. **Re-run tests**:
   ```sh
   go test ./... -race
   ```

5. **Repeat** steps 2-4 until all tests pass.

6. **Run vet** as a final check:
   ```sh
   go vet ./...
   ```

## Test patterns in this repo

- **Table-driven tests** — most tests use `[]struct{ name string; ... }` with `t.Run()`
- **Mock gateway** — `gateway/methods_test.go` uses a mock WebSocket server that validates request frames and sends response frames
- **Race detector** — always use `-race`; the gateway client is concurrent and race bugs are real
- **Coverage target** — aim for 100% statement coverage on library packages

## Key test files

| File | What it covers |
|------|---------------|
| `gateway/methods_test.go` | All 96+ typed RPC methods |
| `gateway/client_test.go` | Connection lifecycle, handshake, reconnect |
| `protocol/types_test.go` | Serialization round-trips |
| `protocol/parse_test.go` | Frame parsing edge cases |
| `chatcompletions/client_test.go` | HTTP client, SSE streaming |
| `openresponses/client_test.go` | Responses API, SSE events |
| `toolsinvoke/client_test.go` | Tools invoke HTTP client |
| `discovery/browser_test.go` | mDNS discovery |
| `acp/server_test.go` | JSON-RPC dispatch, handler interface |

## Rules

- Never use `-count=0` or `t.Skip()` to hide failures
- Always include `-race` flag
- Fix the root cause, not the symptom
- If a test is flaky, fix the flakiness rather than retrying
