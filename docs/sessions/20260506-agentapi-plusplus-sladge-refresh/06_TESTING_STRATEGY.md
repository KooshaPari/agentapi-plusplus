# Testing Strategy

## Targeted Checks

- `git diff --check`
- `rg "sladge|AI Slop" README.md`
- `go test -run "^$" ./...`

## Result Notes

The Go compile-only check is attempted with an isolated `/tmp` cache. In the
current environment it is blocked by local disk exhaustion while writing build
artifacts, so the generated cache is removed and the blocker is carried into
the projects-landing ledger.
