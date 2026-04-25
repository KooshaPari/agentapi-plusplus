# agentapi-plusplus Charter

## 1. Mission Statement

**agentapi-plusplus** is the enhanced fork of `coder/agentapi`—a battle-tested PTY/terminal-emulation layer for wrapping CLI agents over HTTP, SSE, and structured message models—extended with Phenotype-specific routing, benchmarking, and harness capabilities. The mission is to provide a robust, scalable, and extensible agent API infrastructure that enables seamless integration of multiple AI coding agents (Claude, Codex, Gemini, Copilot, etc.) through a unified HTTP endpoint with intelligent routing, session management, and performance benchmarking.

The project exists to bridge the gap between raw PTY-based agent communication and production-grade API infrastructure—transforming terminal-based agent interactions into structured, observable, and manageable HTTP services while preserving the full fidelity of terminal emulation.

---

## 2. Tenets (Unless You Know Better Ones)

### Tenet 1: Upstream Compatibility

Preserve compatibility with the upstream `coder/agentapi` project. Rebase periodically to incorporate upstream improvements. Maintain the upstream module path (`github.com/coder/agentapi`) to avoid breaking SDK-generated clients. Extensions are isolated in internal packages without modifying upstream-derived code paths where avoidable.

### Tenet 2. Intelligent Agent Routing

Multiple AI agents must be addressable through a single HTTP endpoint. Each agent may have preferred LLM models, fallback models, and per-agent retry policies. Routing logic is isolated in `AgentBifrost`—a dedicated routing struct that owns HTTP clients, routing rules, agent sessions, and benchmarking stores. HTTP servers delegate all routing decisions to the bifrost layer.

### Tenet 3. Session State Management

Agent sessions are tracked in-memory with proper concurrency control. Session state includes agent identity, current model, conversation history, and performance metrics. While session state is in-memory only (process restart loses sessions), this is an accepted trade-off for simplicity. Future versions may introduce persistent session backends.

### Tenet 4. Deterministic Benchmarking

All agent interactions are benchmarked with deterministic metrics: latency, token throughput, error rates, and routing decisions. Benchmark data enables performance regression detection, capacity planning, and cost optimization. The `benchmarks.Store` provides thread-safe metric collection and export capabilities.

### Tenet 5. Fallback Model Chaining

When primary models fail or rate-limit, automatic fallback to secondary models ensures service continuity. Fallback chains are configurable per agent with customizable retry policies, timeout values, and error thresholds. Clients are unaware of fallback retries—only the final response matters.

### Tenet 6. Harness Portability

Agent subprocess harnesses are portable across languages. The Go implementation ports harness capabilities from Python (thegent) while maintaining behavioral parity. Harnesses manage agent lifecycle: startup, health checks, graceful shutdown, and crash recovery.

### Tenet 7. Observable Operations

All routing decisions, session mutations, and benchmark events are observable through structured logging and metrics export. Admin endpoints enable runtime inspection of routing tables, session states, and performance statistics without requiring process restarts.

---

## 3. Scope & Boundaries

### In Scope

**Core Agent API:**
- PTY/terminal-emulation over HTTP and SSE
- Structured message model for agent communication
- Session lifecycle management
- Health check and readiness probes

**Routing Infrastructure:**
- AgentBifrost routing middleware
- Per-agent routing rules with RWMutex guards
- Dynamic routing table updates via admin endpoints
- Fallback model chain configuration
- Retry policies with exponential backoff

**Benchmarking System:**
- Request/response latency tracking
- Token throughput measurement
- Error rate monitoring
- Routing decision logging
- Benchmark data export (JSON, Prometheus)

**Harness Management:**
- Agent subprocess lifecycle management
- Health monitoring and crash recovery
- Resource limit enforcement
- Cross-platform harness compatibility

**Extensions:**
- Phenotype-specific internal packages
- Custom routing logic
- Integration with Phenotype ecosystem components

### Out of Scope

- LLM model training or fine-tuning
- Custom terminal emulation implementations (use upstream)
- Persistent session storage (in-memory only)
- WebSocket transport (SSE only)
- Multi-region deployment orchestration

### Boundaries

- Upstream code is preserved intact in designated packages
- Extensions live in `internal/*` packages only
- Session state remains ephemeral by design
- Benchmarking adds overhead; acceptable for observability value
- Routing decisions are synchronous; no async queue

---

## 4. Target Users & Personas

### Primary Persona: Platform Engineer Patricia

**Role:** Infrastructure engineer managing AI agent infrastructure
**Goals:** Deploy reliable agent APIs, monitor performance, manage capacity
**Pain Points:** Agent downtime, unclear performance characteristics, manual routing configuration
**Needs:** Health checks, metrics export, dynamic routing updates, benchmarking data
**Tech Comfort:** Very high, expert in distributed systems

### Secondary Persona: Agent Developer Alex

**Role:** Developer building AI agents that expose CLI interfaces
**Goals:** Expose agent CLI as HTTP API, integrate with agent router
**Pain Points:** PTY handling complexity, HTTP protocol details, session management
**Needs:** Simple harness integration, clear API contracts, debugging tools
**Tech Comfort:** High, comfortable with CLI and HTTP

### Tertiary Persona: Integration Engineer Ingrid

**Role:** Engineer integrating agent API with existing systems
**Goals:** Route requests to appropriate agents, handle failures gracefully
**Pain Points:** Model downtime, rate limiting, agent selection logic
**Needs:** Fallback chains, retry policies, routing visibility
**Tech Comfort:** High, integration specialist

### Quaternary Persona: Benchmarker Ben

**Role:** Performance engineer evaluating agent performance
**Goals:** Measure agent latency, compare models, detect regressions
**Pain Points:** Missing metrics, inconsistent measurement, export difficulties
**Needs:** Comprehensive benchmarks, export formats, trend analysis
**Tech Comfort:** High, data-driven

---

## 5. Success Criteria (Measurable)

### Reliability Metrics

- **Uptime Target:** 99.9% availability for routing layer
- **Failover Success Rate:** 95%+ successful fallback to secondary models
- **Session Stability:** <0.1% session corruption rate
- **Harness Recovery:** 99%+ automatic recovery from agent crashes

### Performance Metrics

- **Routing Latency:** P99 <10ms for routing decisions (excluding LLM time)
- **Session Lookup:** P99 <5ms for session state retrieval
- **Benchmark Overhead:** <5% overhead from metrics collection
- **Concurrent Sessions:** Support for 1000+ concurrent agent sessions

### Compatibility Metrics

- **Upstream Sync:** Rebase with upstream at least quarterly
- **Breaking Changes:** Zero breaking changes to upstream-compatible endpoints
- **Extension Isolation:** 100% extension code in `internal/*` packages
- **Module Path:** Maintain `github.com/coder/agentapi` module path

### Adoption Metrics

- **Agent Integration:** 5+ distinct AI agents routed through platform
- **Routing Rules:** Dynamic routing updates without restarts
- **Benchmark Coverage:** 100% of agent requests benchmarked
- **Documentation:** Complete API docs for all public endpoints

---

## 6. Governance Model

### Component Organization

```
agentapi-plusplus/
├── agentapi/            # Upstream-derived code (preserved)
├── chat/                # Upstream chat interfaces
├── lib/                 # Upstream library code
├── sdk/                 # Upstream SDK
├── internal/
│   ├── routing/         # AgentBifrost and routing logic
│   ├── harness/         # Agent subprocess harnesses
│   ├── phenotype/       # Phenotype-specific extensions
│   └── benchmarks/      # Benchmarking and metrics
├── cmd/                 # CLI entry points
└── docs/                # Extended documentation
```

### Upstream Synchronization Process

1. Track upstream `coder/agentapi` repository
2. Quarterly rebase assessment
3. Merge upstream improvements to preserved packages
4. Extension compatibility validation
5. Integration testing before merge

### Extension Development Process

**New Routing Features:**
- RFC for routing logic changes
- Compatibility impact assessment
- Performance benchmark requirements
- Admin API documentation updates

**New Harness Types:**
- Portability requirements from Python reference
- Cross-platform testing matrix
- Resource limit enforcement validation
- Graceful shutdown verification

---

## 7. Charter Compliance Checklist

### For Upstream Rebase

- [ ] All preserved packages compile without changes
- [ ] Extension packages still compile after rebase
- [ ] Integration tests pass
- [ ] Benchmarks show no regression
- [ ] Admin endpoints functional

### For New Routing Features

- [ ] RWMutex protection for all shared state
- [ ] Fallback chain tested with simulated failures
- [ ] Metrics collection enabled
- [ ] Admin endpoint exposed for configuration
- [ ] Documentation updated

### For New Harness Implementations

- [ ] Cross-platform compatibility verified
- [ ] Health check implemented
- [ ] Graceful shutdown handling
- [ ] Crash recovery tested
- [ ] Resource limits enforced

---

## 8. Decision Authority Levels

### Level 1: Extension Maintainer Authority

**Scope:** Bug fixes in extensions, benchmark improvements, documentation updates
**Process:** Maintainer approval, no RFC required
**Examples:** Session timeout tuning, metric naming improvements

### Level 2: Routing Logic Authority

**Scope:** New routing rules, fallback chain modifications, admin API changes
**Process:** Technical review, RFC for significant changes
**Examples:** New routing algorithms, admin endpoint additions

### Level 3: Upstream Sync Authority

**Scope:** Rebase decisions, upstream merge strategy, compatibility commitments
**Process:** Technical steering review, quarterly assessment
**Examples:** Major upstream version adoption, breaking change mitigation

### Level 4: Architecture Authority

**Scope:** Harness architecture changes, session storage redesign, transport changes
**Process:** Written ADR, steering committee approval
**Examples:** Persistent session backends, WebSocket addition

### Level 5: Strategic Authority

**Scope:** Fork maintenance decisions, upstream relationship strategy, project sunset
**Process:** Executive decision with technical input
**Examples:** Fork discontinuation, upstream contribution strategy

---

## 9. Security & Compliance Considerations

### Agent Isolation

- Agent subprocesses run in isolated environments
- Resource limits prevent resource exhaustion attacks
- No agent can access other agent sessions
- Agent credentials managed through pheno-credentials integration

### Session Security

- Session IDs are cryptographically random
- Session state accessible only to authorized agents
- Session timeouts prevent indefinite resource holding
- Admin endpoints require authentication

### Transport Security

- TLS required for all HTTP endpoints
- SSE connections authenticated
- Admin endpoints on separate port with stricter auth
- Benchmark data excludes sensitive content

---

## 10. Operational Guidelines

### Deployment Patterns

- Stateless routing layer enables horizontal scaling
- Session affinity required (sticky sessions)
- Health checks for load balancer integration
- Metrics export to Prometheus-compatible endpoints

### Monitoring Requirements

- Alert on routing latency P99 >50ms
- Alert on failover rate >10%
- Alert on session count approaching limit
- Dashboard for real-time routing visualization

### Capacity Planning

- 1000 sessions per instance baseline
- Benchmark data retention: 7 days default
- Session timeout: 1 hour default
- Memory budget: 512MB per 1000 sessions

---

## 11. Integration Points

### Phenotype Ecosystem

- **pheno-credentials:** Secure agent credential storage
- **AgilePlus:** Feature tracking and work management
- **cliproxy:** CLI proxy integration for local agents
- **pheno-evaluation:** Benchmark result aggregation

### External Systems

- **Prometheus:** Metrics export format
- **Grafana:** Dashboard visualization
- **Kubernetes:** Health probe endpoints
- **OpenTelemetry:** Distributed tracing (future)

---

## 12. Development Workflows

### Local Development

1. Fork and clone the repository
2. Install Go 1.21 or later
3. Run `go mod download` to fetch dependencies
4. Run `go test ./...` to verify setup
5. Make changes in extension packages
6. Ensure tests pass and add new tests for new functionality
7. Run benchmarks to check for regressions
8. Update documentation for any API changes

### Testing Strategy

- Unit tests for routing logic and benchmarks
- Integration tests for HTTP handlers
- Harness tests for agent lifecycle
- End-to-end tests for complete flows
- Benchmark tests for performance validation
- Race detection tests with `-race` flag

### Release Process

1. Version bump in version file
2. Update CHANGELOG.md with all changes
3. Run full test suite including integration tests
4. Create release tag with semantic versioning
5. Build release binaries for all target platforms
6. Publish release notes with migration guide if needed
7. Notify downstream consumers of breaking changes

---

## 13. Quality Standards

### Code Quality

- All code must pass `go vet` and `golint`
- Test coverage minimum 80% for new code
- No race conditions (test with `-race`)
- Benchmarks must not regress
- Error handling must be explicit

### Documentation Quality

- All public APIs documented with examples
- README kept current with quick start
- Architecture decisions recorded in ADRs
- Examples provided for all major features
- Troubleshooting guide maintained

### Performance Standards

- Routing latency P99 <10ms
- No memory leaks in long-running processes
- CPU usage <10% at idle
- Graceful degradation under load

---

*This charter governs agentapi-plusplus, the enhanced agent routing infrastructure. Intelligent routing enables seamless multi-agent orchestration.*

*Last Updated: April 2026*
*Next Review: July 2026*
