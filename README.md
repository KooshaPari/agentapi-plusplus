# agentapi-plusplus

Multi-model AI routing gateway for Phenotype, providing unified access to Claude, Gemini, and other LLM providers with automatic model selection, load balancing, and fallback strategies.

## Overview

agentapi-plusplus is a production-grade Go service that abstracts away the complexity of managing multiple AI model providers. It intelligently routes requests based on model capabilities, cost, latency, and availability. Agents use agentapi++ for seamless model orchestration without knowing which underlying provider handles each request.

## Technology Stack

- **Language**: Go 1.21+ (async-native)
- **Frameworks**: Gin (HTTP), gRPC for service-to-service
- **Key Dependencies**: OpenAI SDK, Google Cloud Vertex AI, anthropic-sdk
- **Architecture**: Distributed service mesh; integrates with Phenotype infrastructure
- **Deployment**: Docker, Kubernetes, Cloud Run

## Key Features

- Unified API for multiple LLM providers (Claude, Gemini, GPT-4, Llama)
- Intelligent model selection based on task requirements
- Automatic failover and fallback routing
- Request/response caching and deduplication
- Token counting and cost tracking
- Rate limiting per provider and per-user
- Streaming support with connection pooling
- Circuit breaker for provider resilience
- Comprehensive audit logging

## Quick Start

```bash
# Clone and setup
git clone https://github.com/KooshaPari/agentapi-plusplus.git
cd agentapi-plusplus

# Review project governance
cat CLAUDE.md
cat docs/ARCHITECTURE.md

# Install dependencies
go mod download

# Configure API keys
cp config/secrets.yaml.example config/secrets.yaml
# Edit with your provider API keys

# Run tests
go test ./...

# Start the gateway
go run main.go

# Test endpoint
curl -X POST http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"auto","messages":[{"role":"user","content":"Hello"}]}'
```

## Project Structure

```
agentapi-plusplus/
├── cmd/
│   └── gateway/
│       └── main.go             # Service entrypoint
├── internal/
│   ├── api/
│   │   ├── handlers.go         # HTTP request handlers
│   │   ├── middleware.go       # Auth, rate limiting
│   │   └── openai/             # OpenAI-compatible routes
│   ├── models/
│   │   ├── registry.go         # Model capability registry
│   │   ├── selector.go         # Model selection logic
│   │   └── types.go            # Data types
│   ├── providers/
│   │   ├── claude.go           # Claude provider
│   │   ├── gemini.go           # Gemini provider
│   │   ├── openai.go           # OpenAI provider
│   │   └── common.go           # Shared provider logic
│   ├── cache/
│   │   ├── memory.go           # In-memory cache
│   │   ├── redis.go            # Redis cache (optional)
│   │   └── dedup.go            # Request deduplication
│   ├── routing/
│   │   ├── load_balancer.go    # Load balancing
│   │   ├── circuit_breaker.go  # Failure handling
│   │   └── failover.go         # Provider failover
│   └── telemetry/
│       ├── logging.go          # Structured logging
│       ├── metrics.go          # Prometheus metrics
│       └── tracing.go          # Distributed tracing
├── tests/
│   ├── integration/             # Integration tests
│   └── e2e/                     # End-to-end scenarios
├── config/
│   ├── config.yaml             # Configuration
│   └── secrets.yaml.example    # Secrets template
├── docs/
│   ├── ARCHITECTURE.md         # System design
│   ├── MODEL_SELECTION.md      # Selection algorithm
│   └── INTEGRATION_GUIDE.md    # Consumer guide
└── go.mod                      # Go module manifest
```

## Related Phenotype Projects

- **[AgilePlus](../AgilePlus)** — Work tracking and feature specs
- **[Tracera](../Tracera)** — Observability platform (telemetry integration)
- **[cheap-llm-mcp](../cheap-llm-mcp)** — Cost-optimized LLM routing for MCP

## Governance & Documentation

- **CLAUDE.md** — Agent operating instructions and patterns
- **docs/ARCHITECTURE.md** — Detailed system architecture
- **docs/MODEL_SELECTION.md** — Model selection algorithm specification
- **License**: MIT

---

**Status**: Active development  
**Maintained by**: Phenotype Org  
**Last Updated**: 2026-04-24
