# Kilo Gastown Methodology Specification

## Overview

Kilo Gastown is an agent orchestration system designed for AI-native software development. It provides a structured framework for coordinating multiple AI agents across work items called "beads," organized into collective units called "convoys." Gastown enables scalable, traceable, and self-organizing agent workflows with explicit lifecycle management and inter-agent communication primitives.

## Core Concepts

### 1. Beads (Work Items)

Beads are the fundamental unit of work in Gastown. Each bead represents a discrete task or issue that an agent can own and execute.

| Field | Description |
|-------|-------------|
| `bead_id` | Unique identifier (UUID) |
| `type` | `issue`, `convoy`, `task`, `triage` |
| `status` | Lifecycle state (see Bead Lifecycle) |
| `title` | Short description |
| `body` | Detailed requirements or description |
| `assignee_agent_bead_id` | Agent currently working the bead |
| `parent_bead_id` | Hierarchical parent (if any) |
| `priority` | `critical`, `high`, `medium`, `low` |
| `labels` | Categorization tags |
| `metadata` | Arbitrary key-value data |

### 2. Convoys

Convoys are thematic groupings of related beads that share a feature branch and common objective. They enable parallel execution of related work items while maintaining structural coherence.

| Field | Description |
|-------|-------------|
| `convoy_id` | Unique identifier |
| `feature_branch` | Git branch name for all beads in convoy |
| `status` | `open`, `in_progress`, `merged`, `closed` |
| `ready_to_land` | Boolean indicating merge readiness |

**Convoy Benefits:**
- Groups related beads for coordinated review and merge
- Enables `gt_list_convoys` progress tracking across all constituent beads
- Shared branch lifecycle simplifies CI/CD and release management

### 3. Towns

Towns are top-level organizational units (workspaces) that contain rigs. Each town has a unique identity and provides isolation for related projects.

| Field | Description |
|-------|-------------|
| `town_id` | Unique identifier (UUID) |
| `name` | Display name |
| `rigs` | Collection of rigs within the town |

## Bead Lifecycle

```
   ┌─────────┐
   │ triage  │
   └────┬────┘
        │
        ▼
   ┌─────────┐     ┌─────────────┐
   │  open   │────►│ in_progress  │
   └─────────┘     └──────┬───────┘
                          │
                          ▼
                    ┌───────────┐
                    │ in_review │
                    └─────┬─────┘
                          │
            ┌─────────────┼─────────────┐
            ▼             ▼             ▼
       ┌────────┐   ┌──────────┐  ┌─────────┐
       │ closed │   │  merged  │  │  rework │
       └────────┘   └──────────┘  └─────────┘
```

| Status | Description |
|--------|-------------|
| `triage` | New, unclassified work requiring initial assessment |
| `open` | Queued, not yet assigned or started |
| `in_progress` | Actively being worked by an agent |
| `in_review` | Submitted for review, awaiting merge decision |
| `merged` | Successfully integrated (for convoy beads) |
| `closed` | Completed without merge (rejected, cancelled) |
| `rework` | Returned for revision after review |

## Agent Roles

### Polecat

The primary implementation agent. Polecats execute code, write tests, and produce deliverables.

| Responsibility | Description |
|---------------|-------------|
| Hook bead | Claim and begin work on assigned beads |
| Execute | Write code, tests, and documentation |
| Push | Commit and push changes frequently |
| Quality gates | Run lint, test, and format checks |
| Signal done | Call `gt_done` when bead complete |

### Architect/PM

Planner roles responsible for decomposition, planning, and dependency management.

| Responsibility | Description |
|---------------|-------------|
| Create beads | Break down work into discrete items |
| Define dependencies | Establish DAG relationships between beads |
| Assign priority | Order work based on critical path |
| Monitor progress | Track via `gt_list_convoys` |

### Refinery

Merge and review agent. Refinery agents evaluate completed work and make merge decisions.

| Responsibility | Description |
|---------------|-------------|
| Review code | Evaluate against quality standards |
| Request changes | Use `gt_request_changes` for feedback |
| Merge | Approve and integrate when criteria met |
| Escalate | Flag issues requiring human intervention |

## Inter-Agent Communication

### gt_mail_send / gt_mail_check

Formal persistent messaging between agents.

```go
gt_mail_send(to_agent_id, subject, body)
// to_agent_id: UUID of recipient agent
// subject: Brief topic summary
// body: Detailed message content
```

### gt_nudge

Real-time notification for time-sensitive coordination.

| Mode | Behavior |
|------|----------|
| `immediate` | Deliver at next agent idle moment |
| `wait-idle` | Wait until agent is idle |
| `queue` | Queue for later delivery |

## Delegation Primitives

### gt_sling

Single-task delegation to a subagent.

```go
gt_sling(
  task="implement feature X",
  subagent_type="code",
  context={...}
)
```

### gt_sling_batch

Multi-task delegation with parallel or sequential execution.

```go
gt_sling_batch(
  tasks=[...],
  mode="parallel",  // or "sequential"
  max_concurrent=10
)
```

## Key Orchestration Commands

| Command | Purpose |
|---------|---------|
| `gt_prime` | Get full context: identity, hooked bead, mail, open beads |
| `gt_done` | Complete bead, push branch, transition to in_review |
| `gt_bead_status` | Inspect current state of any bead |
| `gt_bead_close` | Mark bead as completed |
| `gt_list_convoys` | Track progress across all convoys |
| `gt_checkpoint` | Write crash-recovery state |
| `gt_escalate` | Create escalation bead for blocked issues |
| `gt_status` | Emit dashboard-visible status update |

## Workflow Integration with AgilePlus

Kilo Gastown provides the orchestration layer while AgilePlus defines process and quality standards:

| Gastown Responsibility | AgilePlus Responsibility |
|------------------------|--------------------------|
| Agent orchestration | Iteration cadence |
| Bead lifecycle | Quality gates |
| Inter-agent messaging | Review pipeline |
| Convoy management | Documentation standards |
| Branch strategy | Versioning |

## Phased WBS with DAG Dependencies

All bead work follows a phased structure with explicit dependencies:

| Phase | Focus | Output | Dependencies |
|-------|-------|--------|--------------|
| 1 - Discovery | Scope, requirements | PRD, user stories | None |
| 2 - Design | Architecture, approach | ADRs, specs | Phase 1 |
| 3 - Build | Implementation | Code, unit tests | Phase 2 |
| 4 - Test | Validation | Integration tests | Phase 3 |
| 5 - Deploy | Release | Deployed artifacts | Phase 4 |

**Dependency Rules:**
- No cyclic dependencies (DAG)
- Beads can only start when predecessors complete
- Use `parent_bead_id` for hierarchical relationships
- Use `metadata.depends_on` for explicit dependencies

## Quality Enforcement

### Pre-Submission Gates

Before calling `gt_done`, all work must pass:

1. **Lint**: Zero golangci-lint errors
2. **Tests**: 80%+ code coverage
3. **Format**: gofmt compliance
4. **Vet**: go vet clean

### Violation Handling

| Violation | Action |
|-----------|--------|
| `revive` (unused params) | Rename to `_paramName` |
| `goconst` (repeated strings) | Extract to named constants |
| `mnd` (magic numbers) | Extract to named constants with units |
| `gocognit` (complexity) | Extract helper functions |
| `funlen` (long functions) | Split into focused helpers |
| `gosec` (security) | Fix immediately |

## Multi-Actor Coordination

### Command Debouncing

High-impact commands (`make lint`, `make test`, `make quality`) MUST use `smart-command.sh` to prevent conflicts between concurrent agents.

### Lock Files

Active command locks stored in `.process-compose/locks/`. Always check before running heavy tasks.

### Shared Services

Use `process-compose` for orchestration. Use `make dev-status` and `make dev-restart` for service management.

## Anti-Patterns (Forbidden)

- Human checkpoint gates in plans
- "Schedule external audit" task steps
- Silent degradation for required dependencies
- Creating v2 files instead of refactoring
- Duplicate implementations

---

**Document Version**: 1.0  
**Last Updated**: 2026-03-31  
**Owners**: Kilo Architecture Team
