# agentapi-plusplus — Go Agent API

This repository provides the agent API server and CLI for Coder workspaces.

## Stack

- **Go** 1.24+ (root module: `github.com/coder/agentapi`)
- **Bun** frontend in `chat/` (Next.js, React, Radix UI)
- **VitePress** docs in `docs/`
- **Taskfile.yml** for build/lint/test (task runner; not Go's `go test`)

## Repository Structure

```
./
  main.go                    # Entrypoint
  cmd/                       # CLI commands (cobra)
    agentapi/                # Agent subcommand
    attach/                  # Attach subcommand
    server/                  # Server subcommand + platform process mgmt
    root.go                  # Root command
  internal/                  # Private packages
  agentapi-plusplus/         # Sync'd upstream module (standalone go.mod)
  chat/                      # Next.js frontend
  docs/                      # VitePress docs site
  config/                    # Runtime config
  e2e/                       # End-to-end test helpers
```

## Development Commands

```bash
task build          # Build all targets (Go + chat + docs)
task test           # Run tests
task lint           # Run linters
task clean          # Remove build artifacts
```

## Go Conventions

- **Formatter**: `gofumpt` (via `go vet` in task lint:go)
- **Linter**: `golangci-lint` v2 (config in `.golangci.yml`)
- **Module path**: `github.com/coder/agentapi`
- **Test**: `go test -count=1 -v ./...` (CGO_ENABLED=0)
- **No `vendor/`** directory; use module proxy

## Branch Discipline

- `main` is protected. All changes via PR.
- Branch naming: `feat/`, `fix/`, `chore/`, `ci/`, `docs/` prefixes.
- Feature work in worktrees: `repos/agentapi-plusplus-wtrees/<topic>/`
