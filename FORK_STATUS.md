# Fork Status — KooshaPari/agentapi-plusplus

**Last updated:** 2026-09-01 (G1 forensic pass)
**Authoritative:** see [README.md](README.md) for canonical provenance

## TL;DR

This fork is **source preservation**, not active development. Work continues
in [KooshaPari/substrate](https://github.com/KooshaPari/substrate) under the
`engine-agentapi` adapter. Do not open issues or PRs assuming active maintenance
of this repository — review the README "Status" section first.

## Source preservation contents

| Surface | State |
|---|---|
| Go HTTP gateway | Preserved (from `coder/agentapi` upstream) |
| Multi-agent integration | **Moved to substrate's `engine-agentapi`** |
| HTTP endpoints (`/messages`, `/message`, `/status`, `/events`, `/health`, `/version`, `/info`) | Preserved, not modified |
| OpenAPI schema (`openapi.json`) | Preserved |
| CI workflows (`.github/`) | Preserved but **noise-prone** — gates are global, not per-PR |

## Branches

- 110 branches exist; **103 have unique commits not in main** (airlock-recovery
  replay snapshots from 2026-05–07 sessions and 2026-06-08 audit day).
- 2 branches merged into main and deleted on 2026-09-01 (G1 sweep).
- 3 orphan branches (no common ancestor) deleted on 2026-09-01.
- `gh-pages` branch preserved (intentional docs hosting).
- `sync/upstream-v0.12.2` branch preserved (deliberately 78 commits behind main,
  not merged by operator decision).

## Upstream

- **Parent:** [`coder/agentapi`](https://github.com/coder/agentapi) — alive, 1,493 stars, 137 forks, last push 2026-05-27.
- **This fork is 243 commits ahead of upstream.**
- **Sync branch:** `sync/upstream-v0.12.2` exists but is 78 commits behind this fork's main — deliberately not merged.

## Phenotype governance

This fork is part of the [Phenotype](https://github.com/KooshaPari) ecosystem.
Governance overlay: `AGENTS.md`, `CLAUDE.md`, `CODEOWNERS`, `CONTRIBUTING.md`,
`FUNDING.yml`, `Justfile`, `.mergify.yml`, `.pre-commit-config.yaml`,
`.golangci.yml`, `.gitleaks.toml`, `.coderabbit.yaml`. These are **fork-local**
and do not exist upstream.

## License

Fork inherits `MIT` from upstream. See [LICENSE](LICENSE) and upstream
[coder/agentapi LICENSE](https://github.com/coder/agentapi/blob/main/LICENSE).

## Recent activity log

- **2026-09-01** — G1 forensic pass; 110 branches inventoried; 5 branches deleted (2 merged + 3 orphan); 4 stale PRs closed (2 closed, 2 blocked by branch protection — parked pending operator decision on admin merge).
- **2026-08-12** — Last open PR opened (#555).
- **2026-08-11** — 4 PRs opened (#552–555) for CI/governance fixes.
- **2026-06-08** — Audit day; ~25 chore/* branches created.
- **2026-05–07** — airlock-recovery session; ~60 branches created for replay.

## Contact

For questions, open an issue or contact via the [Phenotype org](https://github.com/KooshaPari).
