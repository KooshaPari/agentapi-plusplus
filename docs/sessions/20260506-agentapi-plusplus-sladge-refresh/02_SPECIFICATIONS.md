# Specifications

## Requirement

The active AgentAPI++ branch must disclose heavy agent/LLM runtime ownership
with the Sladge badge near the README title.

## Acceptance Criteria

- `README.md` contains `https://sladge.net/badge.svg`.
- The change is prepared in `agentapi-plusplus-wtrees/sladge-current`.
- Validation includes whitespace/diff hygiene, README badge proof, and any
  repo-native compile blocker evidence.

## ARUs

- Assumption: Badge-only documentation work does not require full integration
  test coverage.
- Risk: Go compile/test checks can be blocked by local disk exhaustion or
  sandbox-limited listener/network behavior.
- Mitigation: Run the lightweight compile-only package check and document the
  observed blocker instead of reading it as a badge-regression signal.
