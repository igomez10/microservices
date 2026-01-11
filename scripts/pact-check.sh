#!/usr/bin/env bash
set -euo pipefail

PACT_GO_BIN="${PACT_GO_BIN:-pact-go}"
lib_dir="${PACT_GO_LIB_DIR:-/tmp}"

if ! command -v "$PACT_GO_BIN" >/dev/null 2>&1; then
  echo "pact-go CLI is missing; install via 'go install github.com/pact-foundation/pact-go/v2@latest'"
  exit 1
fi

"$PACT_GO_BIN" check --libDir "$lib_dir"
