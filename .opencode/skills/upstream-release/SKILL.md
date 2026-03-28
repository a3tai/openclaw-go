---
name: upstream-release
description: Handle upstream-release issues — read changelog deltas, update spec docs, sync protocol and gateway types, run tests, open PR, tag release
license: MIT
metadata:
  audience: maintainers
  workflow: github
---

## What I do

Automate the full upstream release sync workflow for openclaw-go when a new OpenClaw gateway version is published.

## When to use me

Use this skill when handling GitHub issues labeled `upstream-release` or when a new OpenClaw gateway version needs to be synced into this library.

## Workflow

1. **Read the issue** — extract the changelog excerpt and identify API/protocol deltas (new RPCs, changed params, removed fields, new event types).

2. **Create or update the spec doc** — write `docs/specs/<version>.md` summarizing upstream deltas and required work items.

3. **Update protocol types** — add or modify types in `protocol/` to match the upstream wire format:
   - New `*Params` and `*Result` structs
   - New constants (methods, event names, scopes)
   - Updated `json` tags and field types
   - Run `go vet ./protocol/...` after changes

4. **Update gateway methods** — add typed convenience methods in `gateway/methods.go`:
   - Follow the existing pattern: marshal params → `client.Send()` → unmarshal result
   - Method name matches the RPC name in PascalCase
   - Include GoDoc comment referencing the RPC method string

5. **Update tests** — add test cases in `gateway/methods_test.go`:
   - Follow existing table-driven test pattern
   - Mock the expected request/response exchange
   - Cover both success and error paths

6. **Update CHANGELOG.md** — add entries under `[Unreleased] > Added` following Keep a Changelog format.

7. **Update docs-site package pages** — add new methods to `docs-site/packages/gateway.md` and any new types to `docs-site/packages/protocol.md`.

8. **Run validation**:
   ```sh
   go test ./... -race
   go vet ./...
   ```

9. **Complete review gates** before marking the PR ready:
   - Architecture review — confirm type design matches upstream semantics
   - Go standards review — idiomatic naming, error handling, documentation
   - API coverage review — every new upstream RPC has a typed method and test
   - Security review — no credential leaks, safe defaults

10. **Open PR** to `main` with:
    - Summary of upstream changes
    - Link to the upstream-release issue
    - Test results pasted in the PR body

11. **After PR merge only** — create tag and GitHub release:
    ```sh
    git tag v<version>
    gh release create v<version> --title "v<version>" --notes "..."
    ```

12. **Close the issue** with a link to the release.

## Key files

- `protocol/types.go` — wire types
- `protocol/constants.go` — method strings, scopes, roles
- `gateway/methods.go` — typed RPC methods
- `gateway/methods_test.go` — method tests
- `CHANGELOG.md` — release notes
- `docs/specs/` — version spec documents
- `docs-site/packages/` — published package docs

## Rules

- Never push directly to `main` — always use feature branches and PRs.
- Branch naming: `issue/<number>-upstream-release-v<version>`
- One upstream version per PR — do not batch multiple releases.
- Every new RPC must have a typed method, a test, a CHANGELOG entry, and a docs update.
