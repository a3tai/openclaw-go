# Upstream Release Automation Spec

## Goal

For every upstream OpenClaw release, run an automated downstream sync pipeline:

1. detect upstream release
2. create/update `docs/specs/<version>.md`
3. scaffold AI coder execution
4. enforce required PR review gates with final HITL
5. publish matching downstream tag/release after merge

## Workflow map

- `check-upstream-releases.yml`
  - scans recent upstream stable releases and creates missing tracking issues
  - includes previous stable version + compare URL in each tracking issue for explicit release-diff context
- `upstream-release-spec.yml`
  - on `upstream-release` issue events, creates/updates `docs/specs/<version>.md` in a PR
  - generates and embeds upstream release diff evidence (`scripts/upstream_release_diff.py`) into the spec
  - auto-dispatches `release-sync-a3t.yml` after spec PR creation/update
- `release-sync-a3t.yml`
  - trigger to run OpenCode release-sync automation from tracking issues
  - restores OpenCode auth from repository secrets (`OPENCODE_AUTH_JSON`, optional `OPENCODE_MCP_AUTH_JSON`)
  - validates release coverage using `scripts/validate_release_sync.py`
- `ci.yml`
  - runs `Release Smoke` (core package smoke tests)
  - runs `Release E2E` (mock server + examples/client + examples/chat)
  - runs `Release Diff Coverage` for release PRs to validate release-delta method/test coverage
- `pr-guard.yml` (`Release Gates` job)
  - enforces release checklist completion for release-sync PRs
- `release-publish.yml`
  - on merged release-sync PR: tag, publish release, comment+close tracking issue
- `upstream-api-drift-report.yml`
  - compares upstream server methods to openclaw-go gateway surface and files an automated drift report issue
  - optionally sends a Resend email notification when `RESEND_API_KEY`, `RESEND_FROM_EMAIL`, and `RESEND_TO_EMAIL` are configured
- `forks-report.yml`
  - publishes a scheduled forks activity report issue assigned to maintainers for review
  - optionally sends a Resend email notification when `RESEND_API_KEY`, `RESEND_FROM_EMAIL`, and `RESEND_TO_EMAIL` are configured

## Upstream-release issue format

Expected issue title format: `Upstream release: openclaw vYYYY.M.D[-N]`

## Spec file requirements

Each generated spec must include:

- upstream release URL and tracking issue URL
- detected protocol/RPC deltas relevant to `openclaw-go`
- required `protocol/` and `gateway/` updates
- review gates:
  - Architecture review
  - Go standards code review
  - API coverage review
  - Security review (Kingpin)

## Required review gates

Release implementation PRs must not merge until all are complete:

1. Architecture review complete.
2. Go standards code review complete.
3. API coverage review complete.
4. Security review complete and signed off by Kingpin (`@slantview` currently).
5. Final HITL review complete.

## Enforcement

- PR template includes mandatory checkbox block.
- PR guard fails release-sync PRs when required checklist items are not checked.
- Release publish workflow validates checklist completion before tagging.
- Release publish workflow also validates green check runs on merge commit:
  - `Test`
  - `Release Gates`
  - `Release Smoke`
  - `Release E2E`
- Branch protection should require CODEOWNERS review and required checks.
- Drift detection runs every 6 hours and raises an issue when upstream methods diverge.
- Drift report classifies upstream compat/internal methods separately from true extras to reduce false positives.

## Deferred

- Discord notifications are intentionally deferred and will be added in a later phase.
- Native a3t platform integration is deferred; current automation uses OpenCode CLI runner mode.

## Non-goals

- Fully automatic code implementation without reviewer approval.
- Direct pushes or tags from unreviewed branches.
