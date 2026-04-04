# ADR-003: Subprocess Harness System

**Status:** Accepted  
**Date:** 2026-04-04  
**Author:** KooshaPari  
**Reviewers:** Architecture Team  
**Supersedes:** None  

---

## Context

AgentAPI++ must control CLI-based AI agents (Claude Code, Codex, Aider, etc.) programmatically. Each agent:

- Runs as a separate OS process
- Expects interactive terminal input (PTY)
- Produces ANSI-colored output
- Streams responses in real-time
- Reports token usage and costs in output
- Has unique command-line interfaces
- Requires specific environment variables
- May need working directory configuration

The challenge: **How do we reliably spawn, communicate with, and capture output from diverse CLI agents in a unified way?**

### Forces

| Force | Weight | Description |
|-------|--------|-------------|
| **PTY requirement** | Critical | Agents need terminal emulation for interactive features |
| **Cross-platform** | High | Must work on macOS, Linux, Windows (WSL) |
| **Output parsing** | High | Must extract structured data from agent output |
| **Streaming** | High | Real-time response delivery required |
| **Resource cleanup** | Critical | Must prevent zombie processes |
| **Timeout handling** | High | Agents may hang indefinitely |
| **Signal propagation** | Medium | Ctrl+C should reach agent correctly |
| **Security** | High | Prevent command injection |

---

## Decision

We will implement a **Subprocess Harness System** with the following architecture:

1. **Runner Interface** - Abstract contract for all agent harnesses
2. **Base Runner** - Shared PTY, process management, ANSI stripping
3. **Agent-Specific Harnesses** - Claude, Codex, Aider implementations
4. **Resource Manager** - Process lifecycle, timeout enforcement, cleanup

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Subprocess Harness System                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────┐       │
│  │                    Runner Interface                  │       │
│  │  Run(ctx, prompt, opts) (*Result, error)            │       │
│  │  Stream(ctx, prompt, opts) (<-chan Event, error)   │       │
│  └─────────────────────────────────────────────────────┘       │
│                              │                                   │
│          ┌───────────────────┼───────────────────┐               │
│          │                   │                   │               │
│          ▼                   ▼                   ▼               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
│  │ClaudeHarness │    │ CodexHarness │    │ GenericHarness│     │
│  │              │    │              │    │               │     │
│  │--print       │    │--full-auto   │    │ configurable  │     │
│  │--output-fmt  │    │--model       │    │ flags         │     │
│  └──────────────┘    └──────────────┘    └──────────────┘     │
│          │                   │                   │               │
│          └───────────────────┼───────────────────┘               │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────┐       │
│  │                    baseRunner                         │       │
│  │  • PTY allocation (creack/pty)                       │       │
│  │  • Process spawning (os/exec)                        │       │
│  │  • ANSI stripping (stripansi)                        │       │
│  │  • Token/cost parsing                                │       │
│  │  • Timeout enforcement (context.Context)           │       │
│  │  • Resource cleanup (defer + Kill)                   │       │
│  └─────────────────────────────────────────────────────┘       │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────┐       │
│  │                   OS Process Layer                    │       │
│  │              (claude, codex, aider CLI)              │       │
│  └─────────────────────────────────────────────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Core Types

```go
// Runner is the interface for all agent harnesses
type Runner interface {
    // Run executes a single prompt and returns complete result
    Run(ctx context.Context, prompt string, opts RunOptions) (*Result, error)
    
    // Stream executes a prompt and returns streaming events
    Stream(ctx context.Context, prompt string, opts RunOptions) (<-chan Event, error)
    
    // Name returns the agent identifier
    Name() string
    
    // Capabilities returns supported features
    Capabilities() Capabilities
}

// RunOptions configures execution
type RunOptions struct {
    WorkingDir    string            // Working directory
    Env           map[string]string // Additional environment variables
    Timeout       time.Duration     // Execution timeout
    Model         string            // Model override
    SystemMessage string            // System prompt
    Temperature   float64           // Sampling temperature
    MaxTokens     int               // Token limit
}

// Result contains execution output
type Result struct {
    Content      string        // Agent response
    Usage        TokenUsage    // Token counts
    Cost         CostEstimate  // Cost estimate
    Duration     time.Duration // Execution time
    ExitCode     int           // Process exit code
    FinishReason string        // Why the agent stopped
}

// TokenUsage tracks API usage
type TokenUsage struct {
    InputTokens  int
    OutputTokens int
    TotalTokens  int
}

// CostEstimate provides pricing info
type CostEstimate struct {
    InputCost  float64 // Cost for input tokens
    OutputCost float64 // Cost for output tokens
    TotalCost  float64 // Total cost in USD
}

// Event represents a streaming update
type Event struct {
    Type    EventType // content, status, tool_call, tool_result, error, done
    Content string    // Event payload
    Usage   *TokenUsage
}
```

### Base Runner Implementation

```go
type baseRunner struct {
    binary       string        // Path to agent CLI
    args         []string      // Base arguments
    env          []string      // Required env vars
    stripANSI    bool          // Whether to strip ANSI codes
    parseTokens  bool          // Whether to parse usage
    timeout      time.Duration // Default timeout
}

func (r *baseRunner) spawn(ctx context.Context, prompt string, opts RunOptions) (*process, error) {
    // 1. Start PTY
    pty, tty, err := pty.Open()
    if err != nil {
        return nil, fmt.Errorf("pty open: %w", err)
    }
    defer tty.Close()
    
    // 2. Configure command
    cmd := exec.CommandContext(ctx, r.binary, r.buildArgs(prompt, opts)...)
    cmd.Stdin = tty
    cmd.Stdout = tty
    cmd.Stderr = tty
    cmd.Dir = opts.WorkingDir
    cmd.Env = r.buildEnv(opts)
    
    // 3. Start process
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("start: %w", err)
    }
    
    // 4. Return process handle
    return &process{
        cmd:    cmd,
        pty:    pty,
        reader: bufio.NewReader(pty),
    }, nil
}

func (r *baseRunner) readOutput(proc *process) (*Result, error) {
    var output strings.Builder
    var usage TokenUsage
    
    scanner := bufio.NewScanner(proc.pty)
    for scanner.Scan() {
        line := scanner.Text()
        
        // Strip ANSI codes
        if r.stripANSI {
            line = stripansi.Strip(line)
        }
        
        // Parse token usage if present
        if r.parseTokens {
            if u := parseTokenLine(line); u != nil {
                usage = *u
            }
        }
        
        // Accumulate content
        output.WriteString(line)
        output.WriteByte('\n')
    }
    
    return &Result{
        Content: output.String(),
        Usage:   usage,
    }, scanner.Err()
}
```

### Agent-Specific Harnesses

#### Claude Harness

```go
type ClaudeHarness struct {
    baseRunner
}

func NewClaudeHarness() *ClaudeHarness {
    return &ClaudeHarness{
        baseRunner: baseRunner{
            binary:      "claude",
            args:        []string{"--print", "--output-format", "stream-json"},
            stripANSI:   true,
            parseTokens: true,
            timeout:     5 * time.Minute,
        },
    }
}

func (h *ClaudeHarness) buildArgs(prompt string, opts RunOptions) []string {
    args := []string{
        "--print",
        "--output-format", "stream-json",
    }
    
    if opts.Model != "" {
        args = append(args, "--model", opts.Model)
    }
    
    // Add prompt via stdin
    return append(args, "-")
}

func (h *ClaudeHarness) parseOutput(reader io.Reader) (*Result, error) {
    // Claude outputs JSON lines:
    // {"type": "message", "content": "..."}
    // {"type": "usage", "input_tokens": 100, "output_tokens": 50}
    
    decoder := json.NewDecoder(reader)
    var content strings.Builder
    var usage TokenUsage
    
    for {
        var msg claudeMessage
        if err := decoder.Decode(&msg); err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }
        
        switch msg.Type {
        case "message":
            content.WriteString(msg.Content)
        case "usage":
            usage.InputTokens = msg.InputTokens
            usage.OutputTokens = msg.OutputTokens
        }
    }
    
    return &Result{
        Content: content.String(),
        Usage:   usage,
        Cost:    calculateClaudeCost(usage),
    }, nil
}
```

#### Codex Harness

```go
type CodexHarness struct {
    baseRunner
}

func NewCodexHarness() *CodexHarness {
    return &CodexHarness{
        baseRunner: baseRunner{
            binary:      "codex",
            args:        []string{"--full-auto"},
            stripANSI:   true,
            parseTokens: true,
            timeout:     5 * time.Minute,
        },
    }
}

func (h *CodexHarness) buildArgs(prompt string, opts RunOptions) []string {
    args := []string{"--full-auto"}
    
    if opts.Model != "" {
        args = append(args, "--model", opts.Model)
    }
    
    // Codex accepts prompt as argument
    return append(args, prompt)
}
```

#### Generic Harness

```go
// GenericHarness for agents not yet fully supported
type GenericHarness struct {
    baseRunner
    config GenericConfig
}

type GenericConfig struct {
    Binary        string
    ArgsTemplate  string // Go template for args
    OutputFormat  string // json, text, stream-json
    TokenRegex    string // Regex to extract token counts
    CostPer1K     float64
}

func NewGenericHarness(config GenericConfig) *GenericHarness {
    return &GenericHarness{
        baseRunner: baseRunner{
            binary:      config.Binary,
            stripANSI:   true,
            parseTokens: config.TokenRegex != "",
            timeout:     5 * time.Minute,
        },
        config: config,
    }
}
```

### Resource Management

```go
// ProcessManager tracks and cleans up processes
type ProcessManager struct {
    processes map[int]*exec.Cmd  // pid -> cmd
    mu        sync.RWMutex
    timeout   time.Duration
}

func (pm *ProcessManager) Track(cmd *exec.Cmd) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    pm.processes[cmd.Process.Pid] = cmd
}

func (pm *ProcessManager) Release(pid int) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    delete(pm.processes, pid)
}

func (pm *ProcessManager) Shutdown() {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    for pid, cmd := range pm.processes {
        log.Printf("Killing orphaned process %d", pid)
        cmd.Process.Kill()
    }
}

// Global cleanup on exit
var defaultProcessManager = &ProcessManager{
    processes: make(map[int]*exec.Cmd),
    timeout:   5 * time.Minute,
}

func init() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        defaultProcessManager.Shutdown()
        os.Exit(0)
    }()
}
```

---

## Consequences

### Positive

1. **Unified Interface** - Same code path for all agents via Runner interface
2. **Clean Output** - ANSI stripping ensures parseable content
3. **Resource Safety** - Guaranteed cleanup via defer + context cancellation
4. **Extensibility** - New agents added by implementing Runner
5. **Testability** - Interface enables mock implementations
6. **Cross-Platform** - PTY abstraction works on macOS/Linux

### Negative

1. **PTY Complexity** - Platform-specific PTY handling required
2. **Output Parsing Fragility** - Depends on agent output format stability
3. **Process Overhead** - Spawning processes slower than in-process calls
4. **Security Surface** - Must sanitize all inputs to prevent injection
5. **Debugging Difficulty** - Subprocess issues harder to diagnose

### Neutral

1. **Memory Overhead** - Each subprocess has its own memory (isolation benefit)
2. **Startup Latency** - ~500ms to spawn agent process (acceptable for long operations)

---

## Alternatives Considered

### Alternative 1: Native API Integration

**Description:** Use official HTTP APIs instead of CLI agents.

**Pros:**
- No subprocess complexity
- Faster (no process spawn)
- Official support
- Better error handling

**Cons:**
- Doesn't support all agents (some CLI-only)
- Different API per provider
- Rate limits more restrictive
- No local file system access

**Rejected:** Many desired agents (Claude Code, Cursor, etc.) are CLI-only.

### Alternative 2: Docker Containers

**Description:** Run each agent in isolated Docker container.

**Pros:**
- Perfect isolation
- Reproducible environments
- Security sandboxing

**Cons:**
- Docker dependency for deployment
- Container startup latency (~2s)
- Resource overhead per container
- Complex volume mounting for file access

**Rejected:** Deployment complexity outweighs isolation benefits for current use case.

### Alternative 3: WebSocket to Agent

**Description:** Modify agents to speak WebSocket protocol directly.

**Pros:**
- Native streaming
- Bidirectional communication
- No parsing required

**Cons:**
- Requires agent modification (not possible for proprietary)
- Complex protocol design
- Agent ecosystem fragmentation

**Rejected:** Cannot modify proprietary agents (Claude Code, Cursor, etc.).

### Alternative 4: exec.Command Only (No PTY)

**Description:** Use simple exec.Command without terminal emulation.

**Pros:**
- Simpler implementation
- No PTY dependencies
- Better Windows support

**Cons:**
- Many agents require PTY for interactive features
- No ANSI code handling
- Some agents refuse non-TTY input

**Rejected:** PTY required by target agents for full functionality.

---

## Security Considerations

### Threat: Command Injection

**Mitigation:**
- Use `exec.Command` with argument array (not shell)
- Never concatenate user input into shell commands
- Validate prompt content for injection patterns

```go
// SAFE: Argument array
cmd := exec.Command("claude", "--print", "--output-format", "stream-json", "-")

// UNSAFE: Shell concatenation (NEVER DO THIS)
cmd := exec.Command("sh", "-c", "claude " + userInput)  // ❌ INJECTION RISK
```

### Threat: Resource Exhaustion

**Mitigation:**
- Context with timeout on all operations
- Maximum concurrent process limits
- Process memory limits (where supported)

### Threat: Path Traversal

**Mitigation:**
- Validate working directory paths
- Use `filepath.Clean()` on all paths
- Restrict to allowed directories

---

## Implementation Checklist

- [x] Runner interface definition
- [x] baseRunner implementation
- [x] PTY allocation (creack/pty)
- [x] ANSI stripping (stripansi)
- [x] Token parsing framework
- [x] Claude harness
- [x] Codex harness
- [x] Aider harness
- [x] Generic harness
- [x] Process manager
- [x] Timeout enforcement
- [x] Resource cleanup
- [x] Signal handling
- [ ] Windows PTY support (optional)
- [ ] Container harness (future)

---

## References

1. [creack/pty - Go PTY library](https://github.com/creack/pty)
2. [stripansi - ANSI escape stripping](https://github.com/acarl005/stripansi)
3. [Go os/exec Documentation](https://pkg.go.dev/os/exec)
4. [PTY - Wikipedia](https://en.wikipedia.org/wiki/Pseudoterminal)
5. [ANSI Escape Codes](https://en.wikipedia.org/wiki/ANSI_escape_code)
6. [Command Injection Prevention - OWASP](https://owasp.org/www-community/attacks/Command_Injection)

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-04-04 | Initial accepted version | KooshaPari |

---

*This ADR follows the nanovms-style decision record format.*
