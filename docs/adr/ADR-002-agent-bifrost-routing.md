# ADR-002: AgentBifrost Routing Layer

**Status:** Accepted  
**Date:** 2026-04-04  
**Author:** KooshaPari  
**Reviewers:** Architecture Team  
**Supersedes:** None  

---

## Context

Multiple AI coding agents (Claude, Codex, Gemini, Copilot, etc.) need to be addressable through AgentAPI++. Each agent has:

- **Preferred LLM models** (e.g., Claude Code prefers Sonnet)
- **Fallback models** when preferred is unavailable (rate limits, outages)
- **Cost characteristics** (input/output token pricing varies by model)
- **Latency profiles** (some models faster than others)
- **Context window limits** (varying by model)
- **Capability differences** (some models better at certain tasks)

Organizations want to:
1. Use cost-effective models by default
2. Fall back to premium models when needed
3. Route based on task requirements
4. Optimize for latency or quality based on context
5. Avoid vendor lock-in by supporting multiple providers

The naive approach (static model selection) leads to:
- **Cost overruns** - Always using expensive models
- **Availability issues** - Single point of failure per agent
- **Suboptimal latency** - No adaptation to network conditions
- **Wasted capacity** - No intelligence in model selection

The problem is: **How do we intelligently route requests to the optimal model while handling failures gracefully?**

### Forces

| Force | Weight | Description |
|-------|--------|-------------|
| **Cost optimization** | High | Minimize API spend without sacrificing quality |
| **Availability** | Critical | Must handle model outages gracefully |
| **Latency** | High | Fast response times for interactive use |
| **Quality** | High | Appropriate model capability for task |
| **Extensibility** | Medium | Easy to add new models and providers |
| **Operational simplicity** | Medium | No ML infrastructure required initially |
| **Data-driven** | Medium | Use historical performance for decisions |

---

## Decision

We will implement **AgentBifrost** as a dedicated routing layer with the following characteristics:

1. **Per-agent routing rules** - Each agent type has configurable preferences
2. **Fallback chains** - Ordered list of models to try on failure
3. **Benchmark integration** - Query telemetry for cost/latency data
4. **Session affinity** - Consistent model per session (optional)
5. **Failure learning** - Track which models fail for which agents
6. **Pluggable selection** - Interface for custom routing algorithms

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        AgentBifrost                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────┐    ┌──────────────────┐                  │
│  │   Routing Rules  │    │   Session State  │                  │
│  │   (sync.RWMutex) │    │   (sync.RWMutex) │                  │
│  │                  │    │                  │                  │
│  │  claude: {       │    │  session-abc: {  │                  │
│  │    preferred:    │    │    agent: claude │                  │
│  │      sonnet,     │    │    model: sonnet │                  │
│  │    fallback:     │    │    history: [...]│                  │
│ │      [opus,haiku]│    │  }               │                  │
│  │  }               │    │                  │                  │
│  └──────────────────┘    └──────────────────┘                  │
│           │                       │                            │
│           └───────────┬───────────┘                            │
│                       │                                        │
│                       ▼                                        │
│  ┌──────────────────────────────────┐                        │
│  │       RouteRequest()             │                        │
│  │  1. Lookup session               │                        │
│  │  2. Get routing rule             │                        │
│  │  3. Query benchmarks             │                        │
│  │  4. Select model                 │                        │
│  │  5. Execute with fallback        │                        │
│  └──────────────────────────────────┘                        │
│                       │                                        │
│                       ▼                                        │
│  ┌──────────────────┐    ┌──────────────────┐                  │
│  │  Benchmark Store │    │   Harness Layer  │                  │
│  │  (ring buffer)   │───►│   (subprocess)   │                  │
│  │                  │    │                  │                  │
│  │  latency: 150ms  │    │  → claude CLI    │                  │
│  │  cost: $0.002    │    │  → parse output  │                  │
│  │  success: true   │    │  → tokens: 342   │                  │
│  └──────────────────┘    └──────────────────┘                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Core Types

```go
// AgentBifrost is the router instance
type AgentBifrost struct {
    cliproxyClient *http.Client
    rules          map[string]*RoutingRule    // agent -> rule
    sessions       map[string]*AgentSession   // sessionID -> session
    benchmarks     *benchmarks.Store
    mu             sync.RWMutex
}

// RoutingRule defines preferences for an agent type
type RoutingRule struct {
    AgentName       string
    PreferredModel  ModelID
    FallbackChain   []ModelID
    MaxRetries      int
    Timeout         time.Duration
    RateLimit       RateLimit
    SessionAffinity bool  // Keep same model per session
}

// AgentSession tracks stateful conversation
type AgentSession struct {
    ID          string
    AgentType   string
    CurrentModel ModelID
    Messages    []Message
    CreatedAt   time.Time
    LastActive  time.Time
    Metadata    map[string]any
}
```

### Routing Algorithm

```
function RouteRequest(ctx, agent, prompt, sessionID):
    // 1. Session lookup
    session = getSession(sessionID)
    if session == nil:
        session = createSession(agent)
    
    // 2. Get routing rule
    rule = getRule(agent)
    if rule == nil:
        rule = defaultRule(agent)
    
    // 3. Determine candidate models
    candidates = [rule.PreferredModel] + rule.FallbackChain
    
    // 4. Query benchmarks if available
    if hasBenchmarkData(agent):
        candidates = rankByPerformance(candidates, agent)
    
    // 5. Try models in order
    for model in candidates:
        if attempts >= rule.MaxRetries:
            break
        
        result, err = executeModel(ctx, agent, model, prompt)
        
        if err == nil:
            // Success - record benchmark
            recordBenchmark(agent, model, result.latency, result.cost, true)
            return result
        
        // Failure - record and continue
        recordBenchmark(agent, model, 0, 0, false)
        logFailure(agent, model, err)
    
    // All models failed
    return error("all models unavailable")
```

### Benchmark Integration

The benchmark store provides historical data for routing decisions:

```go
type Benchmark struct {
    AgentID      string
    ModelID      string
    Timestamp    time.Time
    Latency      time.Duration
    InputTokens  int
    OutputTokens int
    CostUSD      float64
    Success      bool
    ErrorType    string  // If failed
}

// Query for routing decision
func (s *Store) GetModelPerformance(agent string, model string, window time.Duration) Performance {
    benchmarks := s.query(agent, model, window)
    return Performance{
        AvgLatency:   weightedAverage(benchmarks, Latency),
        AvgCost:      weightedAverage(benchmarks, Cost),
        SuccessRate:  successRate(benchmarks),
        SampleCount:  len(benchmarks),
    }
}
```

---

## Consequences

### Positive

1. **Cost Optimization** - Historical data enables intelligent model selection
2. **High Availability** - Automatic fallback prevents outages
3. **Extensibility** - New models added via configuration
4. **Testability** - Router can be unit tested independently
5. **Observability** - Routing decisions logged for analysis
6. **Session Consistency** - Optional affinity prevents model switching mid-conversation

### Negative

1. **Memory Overhead** - In-memory maps require RAM (~2MB per 1000 sessions)
2. **Cold Start Problem** - No benchmark data initially requires fallback to static rules
3. **Complexity** - Routing adds ~500 lines of code vs direct delegation
4. **Failure Amplification** - Retry storms if multiple models failing
5. **State Loss** - Process restart loses session state and recent benchmarks

### Neutral

1. **Decision Latency** - ~2ms overhead for benchmark queries (acceptable)
2. **Learning Curve** - Users must understand routing rule concepts

---

## Alternatives Considered

### Alternative 1: Static Configuration

**Description:** Hardcoded model selection per agent.

```go
// No router - direct mapping
func route(agent string) string {
    switch agent {
    case "claude": return "claude-3-5-sonnet"
    case "codex": return "gpt-4-turbo"
    // ...
    }
}
```

**Pros:**
- Simple implementation
- Predictable behavior
- No state management

**Cons:**
- No cost optimization
- Manual updates for outages
- No adaptation to conditions
- Requires code changes for new models

**Rejected:** Does not meet cost optimization or availability requirements.

### Alternative 2: External Load Balancer

**Description:** Use nginx/HAProxy for model selection.

**Pros:**
- Battle-tested technology
- Health check integration
- No code complexity

**Cons:**
- Cannot integrate benchmark data
- No cost-aware decisions
- Additional infrastructure
- HTTP-only (not subprocess)

**Rejected:** Does not support intelligent cost/latency optimization.

### Alternative 3: ML-Based Model Selection

**Description:** Train ML model to predict optimal model per request.

**Pros:**
- Optimal decisions theoretically possible
- Can consider prompt content
- Sophisticated feature engineering

**Cons:**
- Requires ML infrastructure
- Training data collection complexity
- Explainability challenges
- Overkill for initial implementation

**Status:** Accepted as Phase 2 enhancement, not initial implementation.

### Alternative 4: Client-Driven Selection

**Description:** Let client specify model in request.

**Pros:**
- Maximum client control
- Simple server implementation
- No routing state

**Cons:**
- Clients must know model landscape
- No global optimization
- Harder to manage costs centrally

**Rejected:** Pushes complexity to clients, doesn't enable organization-wide optimization.

---

## Future Enhancements

### Phase 2: ML-Enhanced Routing

```go
// Future: Content-aware routing
type MLRouter struct {
    model    *tf.SavedModel  // or ONNX
    features FeatureExtractor
}

func (r *MLRouter) Predict(prompt string, context Context) ModelID {
    features := r.features.Extract(prompt, context)
    scores := r.model.Predict(features)
    return selectBest(scores)
}
```

**Features to consider:**
- Prompt complexity (token count, structure)
- Historical task patterns
- Time of day (model availability patterns)
- Cost/latency tradeoff preference

### Phase 3: Global Optimization

- **Cross-session learning** - Aggregate patterns across all sessions
- **Predictive pre-warming** - Anticipate model needs
- **Cost budgeting** - Per-organization spend limits with enforcement

---

## Implementation Checklist

- [x] AgentBifrost struct definition
- [x] RoutingRule and AgentSession types
- [x] RouteRequest() core algorithm
- [x] Fallback chain execution
- [x] Benchmark store integration
- [x] Session registry
- [x] Admin API for rule management
- [x] Unit tests for routing logic
- [ ] ML-based model selection (Phase 2)
- [ ] Global optimization (Phase 3)

---

## References

1. [Load Balancing Algorithms](https://www.nginx.com/resources/glossary/load-balancing/)
2. [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
3. [Cost Optimization for LLMs](https://platform.openai.com/docs/guides/production-best-practices)
4. [Go sync.RWMutex Documentation](https://pkg.go.dev/sync#RWMutex)
5. [Token Bucket Rate Limiting](https://en.wikipedia.org/wiki/Token_bucket)

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-04-04 | Initial accepted version | KooshaPari |

---

*This ADR follows the nanovms-style decision record format.*
