#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

"$ROOT_DIR/scripts/pact-check.sh"

(
  cd "$ROOT_DIR/socialapp"
  make contract-test
)

(
  cd "$ROOT_DIR/urlshortener"
  make contract-test
)
