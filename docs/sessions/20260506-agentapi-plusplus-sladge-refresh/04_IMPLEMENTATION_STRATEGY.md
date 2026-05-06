# Implementation Strategy

## Badge Placement

Place the badge directly under the README H1 so the disclosure is visible before
the project summary and agent-control feature list.

## Isolation

Use an isolated worktree even though the canonical checkout is clean, because
the active branch already has local commits and the stale badge evidence lives
on a different branch.
