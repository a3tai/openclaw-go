# Upstream Release Specs

This directory tracks upstream `openclaw` release analysis for `openclaw-go`.

Each file is named `vYYYY.M.D[-N].md` and captures:

- upstream release source links
- protocol/RPC deltas relevant to `openclaw-go`
- required `openclaw-go` changes
- validation and release status

Automation target:

- when an `upstream-release` issue is opened, generate or update the matching
  spec file in this directory automatically via CI PR
- human review confirms required typed-surface changes before merge and tag
