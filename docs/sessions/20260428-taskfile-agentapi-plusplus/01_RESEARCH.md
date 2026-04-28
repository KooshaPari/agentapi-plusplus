# Research

- Detected languages:
  - Go via `/tmp/wt-taskfile-agentapi-plusplus/go.mod`
  - Frontend docs via `/tmp/wt-taskfile-agentapi-plusplus/chat/package.json`
  - Documentation site via `/tmp/wt-taskfile-agentapi-plusplus/docs/package.json`
- Repo scripts inspected:
  - `chat/package.json` exposes `build` and `lint`
  - `docs/package.json` exposes `build`
- Validation:
  - `task --dir /tmp/wt-taskfile-agentapi-plusplus --list` succeeds and exposes `build`, `test`, `lint`, and `clean`

