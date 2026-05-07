# Specifications

## Acceptance Criteria

- Active-branch README contains the Sladge badge.
- Session docs record the current-head refresh and superseded stale branch.
- Validation covers diff hygiene, badge presence, and repo-local task targets
  or records exact blockers.
- Canonical checkout remains unchanged until gates and branch state are safe
  for integration.

## Assumptions, Risks, Uncertainties

- Assumption: the active `ci/add-golangci-lint` branch is the correct truth
  surface for current governance evidence.
- Risk: the previous prepared branch could be mistaken as current evidence.
- Mitigation: record the new worktree, branch, and commit in both downstream
  and projects-landing ledgers.
