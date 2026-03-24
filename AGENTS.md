# AGENTS.md

This file defines repository-specific instructions for coding agents working in `openclaw-go`.

## Non-negotiable workflow rules

- Never push directly to `main`.
- Never open releases from unreviewed commits.
- Always use feature branches and pull requests for all code and docs changes.
- Treat branch protection as policy, even if push bypass is technically possible.

## Required PR flow

1. Create a branch from `main`.
2. Make changes on that branch only.
3. Run validation locally (at minimum: `go test ./... -race`).
4. Commit with clear message(s).
5. Push the branch.
6. Open a PR to `main` with summary, rationale, and test results.
7. Wait for review and required checks to pass.
8. Merge via normal PR merge flow.

## Release skill: upstream-release issues

When handling `upstream-release` issues, use this sequence:

1. Read issue changelog excerpt and identify API/protocol deltas.
2. Create or update `docs/specs/<version>.md` with upstream deltas and required work.
3. Complete review gates for the release PR:
   - architecture review
   - Go standards code review
   - API coverage review
   - security review by Kingpin
   - final HITL review
4. Update `protocol/` and `gateway/` typed surfaces as needed.
5. Update tests (usually `gateway/methods_test.go`).
6. Update `CHANGELOG.md`.
7. Run `go test ./... -race`.
8. Open PR with the release-sync changes and review evidence.
9. Only after PR merge: create tag and GitHub release with `gh release create`.
10. Close issue with release link.

## Guardrails for agents

- If asked to "just push," still use PR workflow unless explicitly told to bypass policy by a maintainer.
- If branch protections are bypassable in the environment, do not use that path.
- If direct pushes already occurred by mistake, notify the user and propose a cleanup plan (revert or follow-up PR), but do not rewrite history without explicit instruction.
