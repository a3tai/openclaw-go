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
- `upstream-release-spec.yml`
  - on `upstream-release` issue events, creates/updates `docs/specs/<version>.md` in a PR
- `release-sync-a3t.yml`
  - scaffold trigger for manual OpenCode/a3t release-sync execution
- `pr-guard.yml` (`Release Gates` job)
  - enforces release checklist completion for release-sync PRs
- `release-publish.yml`
  - on merged release-sync PR: tag, publish release, comment+close tracking issue

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
- Branch protection should require CODEOWNERS review and required checks.

## Deferred

- Discord notifications are intentionally deferred and will be added in a later phase.

## Non-goals

- Fully automatic code implementation without reviewer approval.
- Direct pushes or tags from unreviewed branches.
