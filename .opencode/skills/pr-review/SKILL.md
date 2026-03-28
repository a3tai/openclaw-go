---
name: pr-review
description: Review pull requests against repo standards — Go conventions, test coverage, race detector, CHANGELOG, and docs
license: MIT
metadata:
  audience: reviewers
  workflow: github
---

## What I do

Perform a structured code review of a pull request against the openclaw-go repository standards.

## When to use me

Use this skill when reviewing a PR, or when asked to check if a branch is ready to merge.

## Review checklist

### 1. Go standards
- [ ] Exported types and functions have GoDoc comments
- [ ] Error handling follows `if err != nil { return ... }` pattern — no swallowed errors
- [ ] Naming is idiomatic Go (MixedCaps, not underscores; receiver names are short)
- [ ] No `panic` in library code — return errors instead
- [ ] `context.Context` is the first parameter where applicable
- [ ] JSON struct tags use `omitempty` where the zero value is valid to omit

### 2. Test coverage
- [ ] New public functions have corresponding test cases
- [ ] Tests use table-driven patterns consistent with existing tests
- [ ] Both success and error paths are covered
- [ ] `go test ./... -race` passes cleanly
- [ ] No test files import packages they shouldn't (e.g., no circular deps)

### 3. Protocol correctness
- [ ] New types match the upstream OpenClaw wire format exactly
- [ ] Method strings match the upstream RPC names (e.g., `"sessions.get"`)
- [ ] No breaking changes to existing exported types without a major version bump discussion

### 4. CHANGELOG
- [ ] New features, fixes, or breaking changes have a `[Unreleased]` entry
- [ ] Entries follow Keep a Changelog format
- [ ] Links to upstream version or issue are included where relevant

### 5. Documentation
- [ ] `docs-site/packages/` pages updated if public API surface changed
- [ ] `docs/` markdown files updated if applicable
- [ ] README badges and links still work

### 6. Security
- [ ] No hardcoded tokens, passwords, or secrets
- [ ] TLS defaults are safe (no `InsecureSkipVerify: true` in production code)
- [ ] No user input passed unsanitized into shell commands or SQL

### 7. Git hygiene
- [ ] Commits have clear messages describing the "why"
- [ ] No unrelated changes bundled into the PR
- [ ] Branch is based on current `main`
- [ ] No force-pushes to shared branches

## Output format

Provide the review as a structured summary:
1. **Verdict** — Approve / Request Changes / Comment
2. **Summary** — 1-2 sentence overview
3. **Findings** — list each issue with file:line reference and severity (blocking / suggestion / nit)
4. **Tests** — confirm `go test ./... -race` results
