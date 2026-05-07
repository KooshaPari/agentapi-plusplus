# Testing Strategy

## Planned Checks

- `git diff --check` passed.
- README badge search with `rg` passed.
- `env BUN_TMPDIR=/tmp/agentapi-bun-tmp BUN_INSTALL=/tmp/agentapi-bun-install task lint`
  passed with one non-fatal existing oxlint warning.
- `task test` is blocked by sandbox-denied loopback listener creation in
  `internal/routing.TestForwardToCliproxy_Success`.
- `env BUN_TMPDIR=/tmp/agentapi-bun-tmp BUN_INSTALL=/tmp/agentapi-bun-install task build`
  is blocked by the chat Next.js build fetching the Google-hosted `Geist` font.

## Scope

This is a README/session-doc governance refresh. Failures from unrelated
pre-existing source, generated artifacts, missing tools, or sandbox limits are
recorded as blockers rather than broadened into this change.
