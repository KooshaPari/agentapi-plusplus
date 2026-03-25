#!/usr/bin/env bash
# ci-local.sh -- Local CI for agentapi-plusplus (Go, multi-directory).
set -euo pipefail

passed=0
failed=0
results=()

run_step() {
  local name="$1"; shift
  printf "\n==> %s\n" "$name"
  if "$@"; then
    results+=("PASS  $name")
    ((passed++))
  else
    results+=("FAIL  $name")
    ((failed++))
    return 1
  fi
}

# Root module
run_step "go vet ./..." go vet ./... || exit 1
run_step "go build ./..." go build ./... || exit 1
run_step "go test ./..." go test ./... || exit 1

# agentapi-plusplus subdirectory (if it has its own go.mod)
if [ -f "agentapi-plusplus/go.mod" ]; then
  pushd agentapi-plusplus >/dev/null
  run_step "agentapi-plusplus: go vet ./..." go vet ./... || { popd >/dev/null; exit 1; }
  run_step "agentapi-plusplus: go build ./..." go build ./... || { popd >/dev/null; exit 1; }
  run_step "agentapi-plusplus: go test ./..." go test ./... || { popd >/dev/null; exit 1; }
  popd >/dev/null
fi

# gofmt check (repo-wide)
run_step "gofmt -l ." bash -c '
  bad=$(gofmt -l .)
  if [ -n "$bad" ]; then
    echo "Files need formatting:"
    echo "$bad"
    exit 1
  fi
' || exit 1

printf "\n========== CI Summary ==========\n"
for r in "${results[@]}"; do echo "  $r"; done
printf "Passed: %d  Failed: %d\n" "$passed" "$failed"
[ "$failed" -eq 0 ] && echo "ALL CHECKS PASSED" || { echo "SOME CHECKS FAILED"; exit 1; }
