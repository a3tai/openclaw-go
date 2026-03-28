---
name: release-tag
description: Create a GitHub release with tag after a PR merge — follows the repo's tag and gh release create workflow
license: MIT
metadata:
  audience: maintainers
  workflow: github
---

## What I do

Create a Git tag and GitHub release after a release PR has been merged to `main`.

## When to use me

Use this skill only after a release PR has been merged to `main` and all CI checks have passed. Never use this before the PR is merged.

## Pre-conditions

Before proceeding, verify ALL of these:

1. The release PR is merged to `main`
2. CI passed on the merge commit
3. `go test ./... -race` passes locally on `main`
4. `CHANGELOG.md` has entries under the version being released
5. No other release is in-flight

## Workflow

1. **Checkout and pull main**:
   ```sh
   git checkout main
   git pull origin main
   ```

2. **Determine version** — read from `CHANGELOG.md`:
   - Find the first `## [x.y.z]` heading that isn't `[Unreleased]`
   - The tag will be `v<x.y.z>`

3. **Verify the tag doesn't exist**:
   ```sh
   git tag -l "v<x.y.z>"
   ```

4. **Build release notes** from the CHANGELOG section for this version. Include:
   - Summary of what changed
   - Link to the upstream OpenClaw version if this is a sync release
   - Link to the PR that was merged

5. **Create the tag and release**:
   ```sh
   git tag v<x.y.z>
   git push origin v<x.y.z>
   gh release create v<x.y.z> --title "v<x.y.z>" --notes "<release notes>"
   ```

6. **Verify the release**:
   ```sh
   gh release view v<x.y.z>
   ```

7. **Close related issues** — if there's an upstream-release issue, close it with a link to the release.

## Version format

This project uses semantic versioning:
- **Patch** (`x.y.Z`) — bug fixes, test improvements, doc updates
- **Minor** (`x.Y.0`) — new features (new RPC methods, new packages), backward compatible
- **Major** (`X.0.0`) — breaking changes to exported API (should be rare)

## Rules

- Never create a tag from an unmerged branch
- Never create a release without passing tests
- Never skip the CHANGELOG — if there's no entry, something was missed
- One release per merged PR — do not batch
- The docs-deploy workflow triggers automatically on `v*` tags
