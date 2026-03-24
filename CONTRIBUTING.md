# Contributing to openclaw-go

Thanks for your interest in contributing. This project is independently maintained
by [a3t.ai](https://a3t.ai). We welcome thoughtful, well-tested contributions.

## Before You Start

- Open an issue first to discuss non-trivial changes
- One focused change per PR -- don't bundle unrelated work
- If you're adding a new package or API surface, discuss the design in an issue first

## Requirements

All PRs must pass the following automated checks:

### DCO Sign-off (Required)

Every commit must include a Developer Certificate of Origin sign-off line:

```
Signed-off-by: Your Name <your@email.com>
```

Use `git commit -s` to add this automatically. This certifies that you wrote the
code or have the right to submit it. See [developercertificate.org](https://developercertificate.org/).

### Tests (Required)

- All existing tests must pass: `go test ./... -race`
- New code must include tests targeting 100% statement coverage
- Tests must pass with the race detector enabled

### Code Quality (Required)

- `go vet ./...` must pass with no warnings
- `golangci-lint` must pass
- `go mod tidy` must result in no changes to `go.mod` or `go.sum`
- No new dependencies unless discussed and approved in an issue first

### PR Standards

- **Title**: Descriptive, not generic (e.g., "Add TTSConvert retry logic" not "Update")
- **Description**: Explain what the PR does and why, at minimum 30 characters
- **Size**: PRs adding more than 1,000 lines of Go get flagged for extra review.
  PRs over 3,000 lines are blocked -- split them up

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use the functional options pattern for configuration (see `gateway/options.go`)
- Exported types need doc comments
- Error messages should be lowercase, no punctuation (standard Go style)
- Use `context.Context` as the first parameter for functions that do I/O
- Prefer returning errors over panicking

## AI-Generated Code Policy

We use automated checks to detect low-quality generated code. Your PR will be
flagged or rejected if it contains:

- `panic("not implemented")` or similar stubs
- Empty function bodies
- Placeholder TODO comments without actual implementation
- Highly repetitive boilerplate comments
- Large undifferentiated code dumps

AI tools can be useful for writing code, but the submitter is responsible for
reviewing, testing, and understanding every line. **If you can't explain what
the code does, don't submit it.**

## Running Checks Locally

```bash
# Tests with race detector
go test ./... -race

# Vet
go vet ./...

# Lint (install: https://golangci-lint.run/usage/install/)
golangci-lint run

# Tidy
go mod tidy

# Build everything including examples
go build ./...
go build ./examples/...
```

## License

By contributing, you agree that your contributions will be licensed under the
MIT License. See [LICENSE](LICENSE).

## Releasing

This section is for maintainers.

**Policy:** releases must come from reviewed code merged via PR. Do not push
release commits directly to `main`.

### How to Cut a Release

1. **Ensure everything is green**

   ```bash
   go test ./... -race
   go vet ./...
   go build ./...
   go build ./examples/...
   go mod tidy
   ```

   All tests must pass with the race detector enabled. There should be no
   uncommitted changes after `go mod tidy`.

2. **Update CHANGELOG.md**

   Add a new version section following the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
   format. Move items from `[Unreleased]` if present, or document the changes
   under the appropriate headings (`Added`, `Changed`, `Deprecated`, `Removed`,
   `Fixed`, `Security`).

3. **Create a release PR**

   ```bash
   git checkout -b release/vX.Y.Z
   git add CHANGELOG.md
   git commit -s -m "release: vX.Y.Z"
   git push -u origin release/vX.Y.Z
   gh pr create --base main --title "release: vX.Y.Z" --body "Release prep for vX.Y.Z"
   ```

   Wait for review and required checks, then merge the PR.

   Required release reviews:

   - Architecture review
   - Go standards code review
   - API coverage review
   - Security review by Kingpin
   - Final HITL review

4. **Tag from the reviewed merge commit**

   ```bash
   git checkout main
   git pull --ff-only
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

5. **Create the GitHub release**

   ```bash
   gh release create vX.Y.Z \
     --title "vX.Y.Z" \
     --notes "$(awk '/^## \[X.Y.Z\]/,/^## \[/' CHANGELOG.md | head -n -1)" \
     --repo a3tai/openclaw-go
   ```

   Or use the GitHub UI: go to **Releases → Draft a new release**, select the
   tag, paste the CHANGELOG section as release notes, and publish.

### Semantic Versioning Commitment (post-1.0.0)

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
strictly from v1.0.0 onward:

- **Patch** (`vX.Y.Z+1`): backwards-compatible bug fixes only. No new exported
  symbols, no signature changes.
- **Minor** (`vX.Y+1.0`): new backwards-compatible functionality. New exported
  functions, types, or methods. Existing callers are unaffected.
- **Major** (`vX+1.0.0`): breaking changes to the public API. This requires a
  module path change (e.g., `github.com/a3tai/openclaw-go/v2`) per Go module
  conventions. Breaking changes should be rare and discussed in issues first.

When in doubt, bump minor rather than major — adding is non-breaking, changing
or removing is breaking.
