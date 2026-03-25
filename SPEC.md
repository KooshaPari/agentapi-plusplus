# AgentAPI++ Specification

## Repository Overview

AgentAPI++ provides unified API interfaces for agent operations using Go with FastMCP and provider pattern.

## Architecture

```
agentapi-plusplus/
├── agent-core/               # Core agent logic (no external deps)
│   ├── src/
│   │   ├── lib.rs
│   │   ├── agent.rs         # Agent entity
│   │   ├── tool.rs          # Tool definitions
│   │   ├── executor.rs      # Execution engine
│   │   └── state.rs         # Agent state
│   └── Cargo.toml
├── adapters/                 # API adapters
│   ├── http/               # HTTP handlers
│   ├── mcp/               # MCP protocol
│   └── grpc/               # gRPC services
├── ports/                   # Port interfaces (traits)
│   ├── agent_port.rs       # Agent trait
│   ├── tool_port.rs        # Tool trait
│   └── executor_port.rs    # Executor trait
├── api/                    # API definitions
│   ├── proto/              # Protobuf definitions
│   └── openapi/            # OpenAPI specs
├── cmd/                    # Binary entry points
├── internal/               # Internal packages
│   ├── server/            # HTTP server
│   ├── config/            # Configuration
│   ├── middleware/        # Middleware
│   └── harness/           # Test harness
├── sdk/                    # Client SDKs
├── e2e/                   # End-to-end tests
├── test/                   # Integration tests
├── scripts/               # Build scripts
├── docs/                  # Documentation
└── README.md
```

## Domain Model

### Bounded Contexts

1. **Agent Execution** - Agent lifecycle, tool execution
2. **Provider Management** - Provider registration, routing
3. **Harness Testing** - Test harness, benchmarks

### Core Entities

- `Agent` - Agent configuration and state
- `Tool` - Tool definition and metadata
- `Execution` - Tool execution record
- `Provider` - AI provider configuration

## xDD Methodologies Checklist

### TDD (Test-Driven Development)

- [ ] Red-Green-Refactor cycles
- [ ] Unit tests first
- [ ] Test coverage > 80%
- [ ] Table-driven tests
- [ ] Subcommand-based fuzzing

### BDD (Behavior-Driven Development)

- [ ] Feature files `*.feature`
- [ ] Godog/Cucumber scenarios
- [ ] Step definitions
- [ ] Scenario outlines

### DDD (Domain-Driven Design)

- [ ] Bounded contexts identified
- [ ] Aggregates defined
- [ ] Value objects created
- [ ] Domain events modeled
- [ ] Repository patterns

### ATDD (Acceptance TDD)

- [ ] Acceptance criteria first
- [ ] Executable specs
- [ ] Customer-readable documentation

### Clean Architecture

- [ ] Domain layer isolated
- [ ] Application layer use cases
- [ ] Infrastructure adapters
- [ ] Ports define interfaces

### Hexagonal/Ports & Adapters

- [ ] Primary (driving) adapters: HTTP, gRPC
- [ ] Secondary (driven) adapters: Database, External APIs
- [ ] Ports: Trait definitions

### Architecture Tests

```go
// internal/server/server_test.go
func TestAgentCoreNoOuterDependencies(t *testing.T) {
    // Agent core should not depend on server/middleware
    // Enforce at build level
}
```

## Design Principles

### SOLID (Go-style)

- [ ] Interface Segregation: Small interfaces > large interfaces
- [ ] Dependency Inversion: Depend on interfaces, not implementations
- [ ] Single Responsibility: Package per concern

### GRASP

- [ ] Controller: HTTP handlers coordinate
- [ ] Creator: Factory functions
- [ ] Expert: Services own their data
- [ ] High Cohesion: Related functions in packages
- [ ] Low Coupling: Minimize package dependencies

### Other Principles

- [ ] KISS: Keep it simple
- [ ] DRY: Share common code
- [ ] YAGNI: Don't over-engineer
- [ ] Law of Demeter: Minimize method chains

## Quality Gates

### CI/CD Pipeline

```yaml
# .github/workflows/quality.yml
- name: Go Quality Gates
  run: |
    - go fmt ./...
    - go vet ./...
    - golangci-lint run
    - go test -race -cover ./...
    - gosec ./...
```

### Code Quality Tools

- [ ] gofmt: Formatting
- [ ] go vet: Static analysis
- [ ] golangci-lint: Multi-linter
- [ ] gosec: Security scanning
- [ ] staticcheck: Static analysis

## Module Rules

1. **agent-core**: No dependencies on adapters, server, or external services
2. **ports**: Interface definitions, zero implementation
3. **adapters**: Implement ports, depend on agent-core
4. **internal**: Application logic, depends on ports

## Testing Strategy

```go
// internal/harness/harness_test.go
func TestAgentExecution(t *testing.T) {
    agent := NewAgent(cfg)
    result, err := agent.Execute(context.Background(), req)
    require.NoError(t, err)
    require.Equal(t, expected, result)
}

// Table-driven tests
func TestPricingRules(t *testing.T) {
    cases := []struct {
        name     string
        input    int
        expected float64
    }{
        {"flat rate", 100, 1.00},
        {"volume discount", 1000, 8.50},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := Calculate(c.input)
            assert.Equal(t, c.expected, got)
        })
    }
}
```

## File Organization Rules

```
agent-core/
├── src/
│   ├── lib.rs           # Entry point
│   ├── domain/          # Pure domain (no deps)
│   │   ├── agent.rs     # Agent entity
│   │   └── events.rs    # Domain events
│   └── ports/          # Interface traits
│       └── mod.rs
adapters/
├── http/               # HTTP handlers
│   └── handler.go
└── grpc/              # gRPC servers
    └── server.go
```

## Next Steps

1. Add comprehensive unit tests
2. Implement architecture tests
3. Add BDD scenarios
4. Set up contract testing
5. Create integration test suite
