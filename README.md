# agentapi

Agent API Layer - Sits between thegent and cliproxy+bifrost.

## Architecture

```
thegent → heliosHarness → agentapi → cliproxy+bifrost
```

## Purpose

- **Intermediary layer** between thegent and proxy
- **Custom routing rules** per agent
- **Session-aware** load balancing
- **Agent-specific governance**
- **Dynamic benchmark data** from tokenledger

## Usage

```bash
# Start agentapi
go run ./cmd/agentapi/main.go --port 8318 --cliproxy http://127.0.0.1:8317
```

## Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Health check |
| `POST /v1/chat completions` | Chat with agent routing |
| `GET /admin/rules` | List routing rules |
| `POST /admin/rules` | Set routing rule |
| `GET /admin/sessions` | List active sessions |

## Routing Rules

```json
{
  "agent": "claude",
  "preferred_model": "claude-3-5-sonnet-20241022",
  "fallback_models": ["gpt-4o", "gemini-1.5-pro"],
  "max_retries": 3,
  "timeout_seconds": 30
}
```

## Supported Agents

| Agent | Type | Status |
|-------|------|--------|
| Claude Code | claude | ✅ |
| Amazon Q | amazon-q | ✅ |
| Opencode | opencode | ✅ |
| Goose | goose | ✅ |
| Aider | aider | ✅ |
| Gemini CLI | gemini | ✅ |
| GitHub Copilot | github-copilot | ✅ |
| Sourcegraph Amp | amp | ✅ |
| Codex | codex | ✅ |
| Auggie | auggie | ✅ |
| Cursor | cursor | ✅ |

## Build

```bash
go build -o agentapi ./cmd/agentapi
```

## Fork Differences

This fork includes:
- Custom message formatters for Kush agents
- MCP server integration
- cliproxy routing integration
- Enhanced session management
- Bifrost extension for agent-specific routing

## License

MIT License - see LICENSE file
