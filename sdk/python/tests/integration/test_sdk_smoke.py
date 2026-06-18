"""
Integration tests for the agentapi-plusplus Python SDK.

These tests exercise the AgentAPI client against a real or stubbed server.
They are marked as integration tests so they can be selected/deselected
via pytest's `-m integration` / `-m "not integration"` flags.
"""
# pytestmark = pytest.mark.integration

from __future__ import annotations

from agentapi import (
    AgentAPI,
    AgentStatus,
    AgentType,
    MessageType,
    Status,
)


def test_client_initializes_with_default_base_url() -> None:
    """Smoke test: SDK constructs a client with the documented defaults."""
    client = AgentAPI()
    try:
        assert client.base_url.endswith(":8318")
        assert client.timeout == 30
    finally:
        client.close()


def test_message_type_enum_values() -> None:
    """Verify MessageType enum exposes both USER and RAW message kinds."""
    assert MessageType.USER.value == "user"
    assert MessageType.RAW.value == "raw"


def test_agent_type_enum_contains_supported_agents() -> None:
    """Verify AgentType enum includes all supported agent backends."""
    supported = {t.value for t in AgentType}
    assert {"claude", "goose", "aider", "gemini", "amp", "codex"}.issubset(supported)


def test_status_is_idle_property() -> None:
    """Status.is_idle should be true only when state is STABLE."""
    stable = Status(status=AgentStatus.STABLE, agent_type=AgentType.CLAUDE)
    running = Status(status=AgentStatus.RUNNING, agent_type=AgentType.CLAUDE)
    assert stable.is_idle is True
    assert running.is_idle is False
