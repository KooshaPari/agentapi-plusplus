# Known Issues

## Superseded Branch

Older prepared evidence at `032364f` on
`docs/agentapi-plusplus-sladge-current` diverged from the active branch and is
superseded by this current-head refresh.

## Validation Blockers

`task test` reaches Go package tests but fails when
`internal/routing.TestForwardToCliproxy_Success` tries to create an
`httptest` listener and the sandbox denies binding to `[::1]:0`.

`task build` reaches the chat production build but fails because `next/font`
cannot fetch the Google-hosted `Geist` font while network access is unavailable.

`task lint` passes with `BUN_TMPDIR` and `BUN_INSTALL` pinned under `/tmp`, but
it reports a non-fatal existing oxlint warning in
`chat/src/components/message-list.tsx` for an invalid
`removeEventListener("scrollend", () => ...)` listener reference.
