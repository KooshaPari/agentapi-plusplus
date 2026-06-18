#!/usr/bin/env bash
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

command="${1:-verify}"

case "$command" in
  verify)
    GOFLAGS=-mod=mod go generate ./...
    ./check_unstaged.sh
    ;;
  *)
    echo "usage: $0 verify" >&2
    exit 1
    ;;
esac
