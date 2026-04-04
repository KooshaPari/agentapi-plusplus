# ADR-001: HTTP API Gateway Architecture

**Status:** Accepted  
**Date:** 2026-04-04  
**Author:** KooshaPari  
**Reviewers:** Architecture Team  

---

## Context

The AI coding agent ecosystem is fragmented. Each agent (Claude Code, Cursor, Aider, Codex, Gemini CLI, Copilot, Amazon Q, Augment Code, Goose, Sourcegraph Amp) provides its own CLI interface with unique:

- Command-line flags and arguments
- Output formats (JSON, stream-json, plain text)
- Authentication mechanisms
- Session management approaches
- Tool use protocols
- Streaming capabilities

Organizations building orchestration systems, CI/CD pipelines, and IDE integrations must implement custom adapters for each agent. This creates:

1. **Integration duplication** - Same adapter logic repeated across projects
2. **Maintenance burden** - Each agent CLI update requires adapter updates
3. **Inconsistent interfaces** - Different patterns across integrations
4. **Barrier to adoption** - High effort to add new agents

The problem is: **How do we provide unified programmatic control of heterogeneous CLI agents through a standard interface?**

### Forces

| Force | Weight | Description |
|-------|--------|-------------|
| **Agent heterogeneity** | High | 10+ agents with different interfaces |
| **Integration ease** | High | Must be simpler than direct CLI control |
| **Language independence** | High | Clients in any language must be able to use |
| **Real-time requirements** | High | Must support streaming responses |
| **Operational simplicity** | Medium | Single binary deployment preferred |
| **Performance** | Medium | Sub-100ms latency for control operations |
| **Security** | High | Must validate and sanitize all inputs |

---

## Decision

We will implement **AgentAPI++ as an HTTP API Gateway** that:

1. **Presents a unified REST API** across all supported agents
2. **Handles subprocess execution** via PTY/terminal emulation
3. **Normalizes agent outputs** to a common format
4. **Provides streaming responses** via Server-Sent Events (SSE)
5. **Implements session management** for stateful conversations
6. **Tracks benchmark telemetry** for cost optimization

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Applications                       │
│  (CI/CD, IDEs, Orchestrators, Chat UIs)                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/REST + SSE
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     AgentAPI++ HTTP Gateway                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  REST API   │  │   Admin     │  │   Health    │             │
│  │  Handlers   │  │   API       │  │   Checks    │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AgentBifrost Router                         │
│  - Routing rules per agent                                       │
│  - Model fallback chains                                         │
│  - Benchmark-informed decisions                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Subprocess Harnesses                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │
│  │  Claude  │ │  Codex   │ │  Aider   │ │ Generic  │           │
│  │ Harness  │ │ Harness  │ │ Harness  │ │ Harness  │           │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        CLI Agents                                │
│  (claude, codex, aider, goose, cursor, etc.)                     │
└─────────────────────────────────────────────────────────────────┘
```

### Key Components

#### 1. HTTP Layer (go-chi/chi)

- **Router:** `go-chi/chi/v5` for HTTP routing
- **Middleware:** Recovery, logging, CORS, rate limiting
- **Content Types:** JSON for API, text/event-stream for SSE
- **Error Format:** RFC 7807 Problem Details

```go
// Route structure
r.Route("/api/v0", func(r chi.Router) {
    r.Get("/agents", listAgents)
    r.Route("/sessions", func(r chi.Router) {
        r.Post("/", createSession)
        r.Get("/{id}", getSession)
        r.Post("/{id}/messages", sendMessage)
    })
})

r.Route("/admin", func(r chi.Router) {
    r.Get("/rules", listRules)
    r.Post("/rules", createRule)
})

r.Get("/events", sseHandler)
```

#### 2. AgentBifrost Router

- **Routing Rules:** Per-agent model preferences and fallback chains
- **Session Registry:** In-memory session state
- **Benchmark Integration:** Query telemetry for routing decisions
- **Decision Logic:** ML-informed model selection (future)

#### 3. Subprocess Harnesses

- **PTY Allocation:** Virtual terminal for interactive agents
- **ANSI Stripping:** Clean output for parsing
- **Token Parsing:** Extract usage from agent output
- **Timeout Enforcement:** Context-based cancellation

---

## Consequences

### Positive

1. **Unified Interface** - Single API for 10+ agents reduces integration effort by 80%
2. **Language Agnostic** - Any HTTP client can interact with agents
3. **Streaming Support** - SSE provides real-time updates without WebSocket complexity
4. **Observable** - HTTP layer enables standard monitoring (Prometheus, tracing)
5. **Deployable** - Single binary with no external dependencies
6. **Extensible** - New agents added by implementing harness interface

### Negative

1. **Latency Overhead** - HTTP adds 5-10ms vs direct subprocess calls
2. **Session Volatility** - In-memory sessions lost on restart
3. **Resource Overhead** - HTTP server adds ~50MB baseline memory
4. **Operational Complexity** - New service to monitor and maintain
5. **Scaling Limits** - Single instance limits to ~1000 concurrent sessions

### Neutral

1. **Binary Size** - ~18MB compiled binary (acceptable for deployment)
2. **Language Choice** - Go provides performance but smaller ecosystem than Python

---

## Alternatives Considered

### Alternative 1: Library/SDK Approach

**Description:** Provide language-specific SDKs (Python, TypeScript, Go) that wrap agents.

**Pros:**
- Native language integration
- Type safety per language
- No operational overhead of HTTP service

**Cons:**
- Must maintain SDK per language
- Each SDK must implement subprocess control
- Version synchronization complexity
- Cannot support languages without SDK

**Rejected:** Would require maintaining 5+ SDKs vs single HTTP service.

### Alternative 2: gRPC-First

**Description:** Use gRPC as primary transport instead of REST.

**Pros:**
- Binary protocol efficiency
- Strong typing via protobuf
- Bidirectional streaming
- Code generation for clients

**Cons:**
- Browser/client JavaScript support requires grpc-web
- Harder to debug (binary payloads)
- Additional toolchain complexity
- Less universal than HTTP

**Status:** Accepted as future enhancement (ADR-011), not initial implementation.

### Alternative 3: Message Queue (NATS/Kafka)

**Description:** Use message queue for agent communication.

**Pros:**
- Natural fit for async agent processing
- Built-in persistence options
- Excellent for distributed systems
- Decouples client from agent lifecycle

**Cons:**
- Requires message broker infrastructure
- Adds operational complexity
- Higher latency for simple requests
- Overkill for single-instance deployments

**Rejected:** Adds unnecessary infrastructure for current use cases.

### Alternative 4: Direct Shell Integration

**Description:** Skip API layer, provide shell scripts that wrap agents.

**Pros:**
- Zero additional infrastructure
- Minimal overhead
- Direct agent control

**Cons:**
- Language-specific (shell)
- No streaming support
- Hard to integrate with applications
- Security concerns with shell execution

**Rejected:** Does not meet requirement for language-agnostic integration.

---

## Implementation

### Phase 1: Core HTTP Layer

- [x] chi router setup with middleware
- [x] Session CRUD endpoints
- [x] Message send/receive
- [x] SSE streaming
- [x] RFC 7807 error responses

### Phase 2: Agent Integration

- [x] Claude harness implementation
- [x] Codex harness implementation
- [x] Aider harness implementation
- [x] Generic harness for other agents

### Phase 3: Routing & Telemetry

- [x] AgentBifrost router implementation
- [x] Routing rules API
- [x] Benchmark store
- [x] Model fallback chains

### Phase 4: Production Hardening

- [ ] Rate limiting per host
- [ ] Persistent sessions (Redis/SQLite)
- [ ] Distributed tracing
- [ ] Kubernetes operator

---

## References

1. [REST API Design Best Practices](https://restfulapi.net/)
2. [RFC 7807 Problem Details](https://datatracker.ietf.org/doc/html/rfc7807)
3. [Server-Sent Events Specification](https://html.spec.whatwg.org/multipage/server-sent-events.html)
4. [go-chi/chi Documentation](https://github.com/go-chi/chi)
5. [OpenAPI 3.0 Specification](https://swagger.io/specification/)

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-04-04 | Initial accepted version | KooshaPari |

---

*This ADR follows the nanovms-style decision record format.*
