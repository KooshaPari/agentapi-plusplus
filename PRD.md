# Product Requirements Document — AgentAPI++

**Module:** `github.com/coder/agentapi` (KooshaPari fork — `agentapi-plusplus`)
**Baseline commit:** `ddaedc2`
**Status:** Active Development

---

## 1. Overview

AgentAPI++ is an HTTP-API layer that wraps CLI-based AI coding agents (Claude Code, Codex, Aider, Gemini CLI, GitHub Copilot, Cursor CLI, Goose, Sourcegraph Amp, Amazon Q, Augment Code) in a unified, programmable interface. It removes the need to manually drive a terminal emulator and exposes agent interaction over REST + SSE so any language or automation platform can programmatically control agents.

The `++` additions over the upstream `coder/agentapi` codebase are:

- `AgentBifrost` routing layer (`internal/routing/`) — per-agent model selection, fallback chaining, and session-aware load balancing sitting between callers and the `cliproxy+bifrost` proxy.
- `harness` package (`internal/harness/`) — Go port of the agent subprocess abstractions from `thegent` (Python), enabling direct subprocess control of Claude Code, Codex, and generic CLIs with stdin injection, ANSI stripping, timeout enforcement, and token/cost telemetry parsing.
- Phenotype SDK init hook (`internal/phenotype/`) — lightweight `.phenotype/` directory bootstrap for Phenotype-org workspace integration.
- Benchmark telemetry store (`internal/benchmarks/`) — token counts and estimated cost from every harness run, surfaced to the `AgentBifrost` for routing decisions.

---

## 2. Goals

| ID | Goal |
|----|------|
| G-1 | Provide a stable HTTP API to send prompts to any supported CLI agent and receive structured responses. |
| G-2 | Maintain per-agent routing rules (preferred model, fallback chain, retry policy) without requiring callers to change. |
| G-3 | Track agent sessions with metadata so multi-turn conversations and load-balancing decisions are session-aware. |
| G-4 | Enable subprocess-level agent control via the `harness` package for callers that need direct process access rather than the HTTP proxy path. |
| G-5 | Capture benchmark telemetry (token counts, cost, duration) from every run for downstream routing and cost accounting. |
| G-6 | Integrate into the Phenotype ecosystem workspace via the `phenotype` SDK init hook. |

---

## 3. Non-Goals

- AgentAPI++ does not implement a model provider itself; it delegates all LLM calls to agents via the `cliproxy+bifrost` proxy or direct CLI subprocess.
- It does not persist conversation history to a database; session state is in-memory.
- It does not provide a user-facing product UI beyond the read-only `/chat` debug interface served on port 3284.

---

## 4. Epics and User Stories

### E1 — HTTP Agent Control

**Description:** External callers (CI pipelines, orchestrators, thegent) send prompts via HTTP and receive structured agent output.

| Story ID | Story | Acceptance Criteria |
|----------|-------|---------------------|
| E1.S1 | As an orchestrator, I can POST a prompt to `/v1/chat/completions` and receive the agent's response. | HTTP 200 with agent output JSON on success; HTTP 4xx/5xx on bad input or agent error. |
| E1.S2 | As a caller, I can GET `/messages` to see all messages in the current session. | Returns ordered list of messages with role and content fields. |
| E1.S3 | As a caller, I can GET `/status` to determine whether the agent is idle or processing. | Returns `{"status": "stable" | "running"}`. |
| E1.S4 | As a caller, I can GET `/events` to receive a real-time SSE stream of message and status events. | SSE stream delivers `message` and `status` event types; connection held open until agent is stable. |

### E2 — Multi-Agent Routing (AgentBifrost)

**Description:** The `AgentBifrost` component selects the correct model and proxy path for each named agent.

| Story ID | Story | Acceptance Criteria |
|----------|-------|---------------------|
| E2.S1 | As an admin, I can GET `/admin/rules` to inspect current per-agent routing rules. | Returns all configured `RoutingRule` entries as JSON. |
| E2.S2 | As an admin, I can POST `/admin/rules` to set or update a routing rule for a named agent. | Rule persisted in memory; subsequent requests for that agent use the new rule immediately. |
| E2.S3 | As a caller, when a preferred model fails, the system retries with fallback models in configured order. | After exhausting `max_retries` across the fallback chain, the final error is returned to the caller. |
| E2.S4 | Requests for unknown agents use a safe default rule: Claude Sonnet as preferred, GPT-4o and Gemini as fallbacks. | Default `RoutingRule` applied when no rule is registered for the given agent name. |

### E3 — Session Management

**Description:** Session state tracks which agent is active, which models have been used, and arbitrary metadata per session.

| Story ID | Story | Acceptance Criteria |
|----------|-------|---------------------|
| E3.S1 | Sessions are created automatically on first request for a given agent name. | `AgentSession` with a unique ID exists in the sessions map after the first `RouteRequest` call. |
| E3.S2 | As an admin, I can GET `/admin/sessions` to see all active sessions. | Returns list of `AgentSession` structs with ID, agent name, start time, and models used. |
| E3.S3 | Sessions are safe for concurrent goroutine access. | `sync.RWMutex` guards the sessions map; no data races under `go test -race`. |

### E4 — Subprocess Agent Harness

**Description:** The `harness` package runs agent CLIs as managed subprocesses.

| Story ID | Story | Acceptance Criteria |
|----------|-------|---------------------|
| E4.S1 | I can invoke `RunHarness("claude", opts)` and receive a `RunResult` with stdout, stderr, exit code, token counts, and cost. | `RunResult.ExitCode == 0` on success; `PromptTokens`, `CompletionTokens`, `CostUSD` populated by parsing agent output. |
| E4.S2 | A timeout-exceeded run sets `RunResult.TimedOut = true` and terminates the subprocess cleanly. | Process killed within 2 seconds of `Timeout` expiry; `TimedOut` field is `true` in the returned result. |
| E4.S3 | Claude harness invokes `claude --print --output-format stream-json --verbose [--model <m>] < <prompt>` with correct flags per `Mode`. | `ModeReadOnly` omits `--dangerously-skip-permissions`; `ModeWrite`/`ModeFull` add it. |
| E4.S4 | Codex harness invokes `codex --full-auto [--model <m>] <prompt>` passing the prompt as an argument rather than stdin. | `RunResult` populated with parsed Codex output; `usesStdin == false` for Codex. |
| E4.S5 | Generic harness handles cursor-agent, copilot, gemini, and opencode via a configurable command template. | Template substitution (model, prompt, mode flags) is applied correctly; `RunResult` populated with parsed output. |

### E5 — Telemetry and Benchmarks

**Description:** Token counts and estimated cost from every harness run are stored and made available for routing decisions.

| Story ID | Story | Acceptance Criteria |
|----------|-------|---------------------|
| E5.S1 | `benchmarks.Store` accumulates `RunResult` telemetry entries without blocking callers. | `Store.Record(result)` returns immediately; data is queryable after the call. |
| E5.S2 | `AgentBifrost` reads benchmark data when making model selection decisions. | Routing logic can factor in average latency and cost per agent/model pair from the store. |

### E6 — Phenotype Workspace Integration

**Description:** The `phenotype` package initialises the `.phenotype/` directory at server startup.

| Story ID | Story | Acceptance Criteria |
|----------|-------|---------------------|
| E6.S1 | `phenotype.Init(repoRoot)` creates `.phenotype/` if absent, without error. | Directory exists after the call; function is idempotent on repeated invocations. |
| E6.S2 | If `repoRoot` is empty, `Init` falls back to the process working directory. | `os.Getwd()` used as fallback; no panic on empty string input. |

---

## 5. Constraints and Assumptions

- Agents must be installed and accessible on `PATH` (or supplied via explicit `cliPath` parameter to harness constructors).
- The `cliproxy+bifrost` URL is required for the `AgentBifrost` HTTP routing path; the harness path bypasses it.
- Session state is ephemeral (in-memory); a process restart loses all session history.
- The default listen port is `3284`; configurable via CLI flag.
- Allowed-hosts enforcement defaults to `localhost` only; overridable via `--allowed-hosts` flag or `AGENTAPI_ALLOWED_HOSTS` environment variable.
- GitHub Actions CI is blocked by billing limits on the KooshaPari account; quality must be verified locally.

## Project Description

High-performance API framework for agent orchestration and management.

See [README.md](./README.md) for detailed project overview.

## Key Features & Epics

### E1: Core Functionality
Primary features and capabilities.

### E2: Integration & Extensibility
System integration points and extension mechanisms.

### E3: Operations & Quality
Operational tooling and quality assurance.

## Success Criteria

- [ ] Core features implemented and tested
- [ ] Documentation complete
- [ ] Performance requirements met
- [ ] Security validated
- [ ] User acceptance passed

## Future Roadmap

- **Phase 2**: Advanced capabilities
- **Phase 3**: Performance optimization
- **Phase 4**: Enterprise features

---

**Status**: ACTIVE
**Owner**: Engineering Team
**Last Updated**: 2026-03-25
