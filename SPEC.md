# AgentAPI++ Specification

## Repository Overview

AgentAPI++ is a Go-based agent API platform using FastMCP and provider patterns.

## Architecture

```
agentapi-plusplus/
├── cmd/                    # CLI entry points
├── internal/                # Core business logic
│   ├── adapters/           # Primary/Driving adapters
│   ├── config/             # Configuration
│   ├── middleware/         # HTTP middleware
│   ├── phenotype/          # Phenotype integration
│   ├── routing/            # Request routing
│   ├── server/             # HTTP server
│   └── harness/            # Test harness
├── lib/                    # Shared libraries
├── sdk/                    # SDK implementations
├── chat/                   # Chat interfaces
├── e2e/                   # E2E tests
└── agentapi-plusplus/      # Main binary
```

## Bounded Contexts

### Core Agent Context
- Agent management
- Provider abstraction
- Tool execution

### Harness Context
- Claude harness
- Codex harness
- Test execution

### Routing Context
- Agent Bifrost
- Request routing

## xDD Methodologies Checklist

### TDD (Test-Driven Development)
- [ ] Red-Green-Refactor cycles
- [ ] Unit tests first in `*_test.go`
- [ ] Test coverage > 80%
- [ ] Property-based tests with `go-check.v2`
- [ ] Table-driven tests

### BDD (Behavior-Driven Development)
- [ ] Gherkin scenarios
- [ ] Godog/Cucumber integration
- [ ] Feature files in `e2e/`
- [ ] Step definitions

### DDD (Domain-Driven Design)
- [ ] Bounded contexts defined
- [ ] Entities and value objects
- [ ] Domain services
- [ ] Repository interfaces

### ATDD (Acceptance TDD)
- [ ] Acceptance criteria first
- [ ] Executable specs
- [ ] Living documentation

### CQRS (Command Query Responsibility Segregation)
- [ ] Separate command/query handlers
- [ ] Read models optimized
- [ ] Write models for mutations

### EDA (Event-Driven Architecture)
- [ ] Domain events
- [ ] Event handlers
- [ ] Async messaging

### Clean Architecture
- [ ] Inner layers have no dependencies
- [ ] Dependencies point inward
- [ ] Adapters depend on ports

### Hexagonal/Ports & Adapters
- [ ] Primary ports (driving)
- [ ] Secondary ports (driven)
- [ ] Adapters implement ports

## Design Principles

### SOLID
- [ ] Single Responsibility: Each package one reason to change
- [ ] Open/Closed: Open for extension, closed for modification
- [ ] Liskov Substitution: Subtypes substitutable for base types
- [ ] Interface Segregation: Many specific interfaces > one general
- [ ] Dependency Inversion: Depend on abstractions, not concretions

### GRASP
- [ ] Controller: Handle UI/system event
- [ ] Creator: Class creates another
- [ ] Information Expert: Assign responsibility to class with info
- [ ] Indirection: Introduce mediator
- [ ] Low Coupling: Minimize dependencies
- [ ] High Cohesion: Related responsibilities together
- [ ] Polymorphism: Different behaviors based on type

### Other Principles
- [ ] DRY: No duplicate code
- [ ] KISS: Keep simple
- [ ] YAGNI: You aren't gonna need it
- [ ] Law of Demeter: Talk to friends only
- [ ] SoC: Separation of concerns
- [ ] PoLA: Principle of Least Astonishment

## CI/CD Quality Gates

```yaml
lint:
  - golangci-lint run
  - go vet ./...
  - staticcheck ./...

test:
  - go test -race -cover ./...
  - go test -bench=. ./...
  - mutation coverage

security:
  - gosec ./...
  - trivy scan
  - trivy fs .

format:
  - gofmt -s
  - goimports
```

## Testing Strategies

### Unit Tests
```go
func TestAgentCreation(t *testing.T) {
    // Arrange
    config := Config{Provider: "claude"}

    // Act
    agent, err := NewAgent(config)

    // Assert
    require.NoError(t, err)
    require.NotNil(t, agent)
}
```

### Integration Tests
```go
func TestHarnessExecution(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    // Full harness test
}
```

### E2E Tests
```gherkin
Feature: Agent Execution
  Scenario: Execute simple command
    Given an agent with claude provider
    When I send "echo hello"
    Then I receive "hello"
```

## Architecture Tests

```go
// internal/core/no_external_deps.go
func TestCoreHasNoExternalDeps(t *testing.T) {
    deps, _ := analyzeImports("internal/core")
    for _, dep := range deps {
        require.False(t, isExternal(dep),
            "core package should not depend on %s", dep)
    }
}
```

## Project-Specific Patterns

### Provider Pattern
```go
type Provider interface {
    Execute(ctx context.Context, req Request) (*Response, error)
}

type ClaudeProvider struct{}
type CodexProvider struct{}
```

### Harness Pattern
```go
type Harness interface {
    Run(ctx context.Context, scenario Scenario) Result
}
```

### Config Pattern
```go
type Config struct {
    Provider    string
    Timeout    time.Duration
    MaxRetries int
}
```

## Documentation

- [x] ADR-001: FastMCP + Provider pattern
- [ ] ADR-002: Error handling strategy
- [ ] ADR-003: Authentication approach
- [ ] Domain README for each context
- [ ] API documentation

## File Organization Rules

1. `internal/` - No external dependencies except stdlib
2. `adapters/` - Implementation details
3. `ports/` - Interface definitions
4. `*_test.go` - Test files co-located
5. `e2e/` - End-to-end tests
6. `cmd/` - Executable entry points
