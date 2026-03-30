## Summary

-

## Why

-

## Validation

- [ ] `go test ./... -race`

## Release Review Gates (required for release-sync PRs)

> ⚠️ Release-sync PRs must have the `release-sync` label applied **before merge** to trigger
> automatic publishing. The `release-sync-claude` workflow sets this automatically; if filing
> manually, add it yourself or the publish workflow will fall back to title-pattern detection.

- [ ] Architecture review
- [ ] Go standards code review
- [ ] API coverage review
- [ ] Security review (Kingpin / @slantview)
- [ ] Final HITL review

## Notes

-
