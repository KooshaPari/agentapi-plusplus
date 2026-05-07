# Implementation Strategy

## Approach

Use a new current-head isolated worktree instead of reusing the stale prepared
branch. Keep the change limited to README badge evidence and session
documentation.

## Boundaries

- Do not modify canonical agentapi-plusplus during the refresh.
- Do not reuse stale `docs/agentapi-plusplus-sladge-current` evidence.
- Do not alter generated dependency or build artifacts unless validation
  requires a deliberate tracked update.
