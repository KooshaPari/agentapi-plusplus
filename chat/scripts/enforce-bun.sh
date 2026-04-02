#!/usr/bin/env bash
#!/usr/bin/env sh
# Skips heavy pre-push checks for docs-only changes
_changed_files=$(git diff --name-only origin/main...HEAD 2>/dev/null || git diff --name-only HEAD^..HEAD 2>/dev/null || git diff --name-only)
if [ -z "$_changed_files" ]; then _changed_files=$(git diff --name-only); fi
skip=1
for f in $_changed_files; do
  case "$f" in
    docs/*|.github/*|package.json|docs/.vitepress/*) ;;
    *) skip=0; break ;;
  esac
done
if [ "$skip" -eq 1 ]; then
  echo "Docs-only changes detected — skipping heavy pre-push checks."
  exit 0
fi


set -euo pipefail
if [ -f package-lock.json ] || [ -f yarn.lock ] || [ -f pnpm-lock.yaml ]; then
  echo "ERROR: Non-bun lockfile found. Use bun exclusively."
  exit 1
fi
if [ -n "${npm_execpath:-}" ] && echo "$npm_execpath" | grep -qv "bun"; then
  echo "ERROR: Use bun, not npm/yarn/pnpm."
  exit 1
fi
