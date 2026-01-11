#!/usr/bin/env bash
set -euo pipefail

PACT_BROKER_BASE_URL="${PACT_BROKER_BASE_URL:-}"
PACT_BROKER_VERSION="${PACT_BROKER_VERSION:-}"

if [[ -z "$PACT_BROKER_BASE_URL" ]]; then
  echo "PACT_BROKER_BASE_URL must be set"
  exit 1
fi

if [[ -z "$PACT_BROKER_VERSION" ]]; then
  echo "PACT_BROKER_VERSION must be set"
  exit 1
fi

PACT_TOKEN="${PACT_BROKER_TOKEN:-}"
PACT_USER="${PACT_BROKER_USERNAME:-}"
PACT_PASS="${PACT_BROKER_PASSWORD:-}"

auth_header=""
if [[ -n "$PACT_TOKEN" ]]; then
  auth_header="Authorization: Bearer $PACT_TOKEN"
elif [[ -n "$PACT_USER" && -n "$PACT_PASS" ]]; then
  credentials=$(printf "%s:%s" "$PACT_USER" "$PACT_PASS" | base64)
  auth_header="Authorization: Basic $credentials"
else
  echo "Either PACT_BROKER_TOKEN or both PACT_BROKER_USERNAME and PACT_BROKER_PASSWORD must be set"
  exit 1
fi

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 pact-file [another-pact-file ...]"
  exit 1
fi

publish_pact() {
  local pact_file="$1"
  local pact_name
  pact_name=$(basename "$pact_file")

  local consumer provider
  case "$pact_name" in
  socialapp-client-socialapp.json)
    consumer="socialapp-client"
    provider="socialapp"
    ;;
  socialapp-urlshortener.json)
    consumer="socialapp"
    provider="urlshortener"
    ;;
  *)
    echo "Unknown pact artifact: $pact_name"
    exit 1
    ;;
  esac

  local url="${PACT_BROKER_BASE_URL%/}/pacts/provider/${provider}/consumer/${consumer}/version/${PACT_BROKER_VERSION}"

  curl -fsS --retry 3 -X PUT \
    -H "Content-Type: application/json" \
    -H "$auth_header" \
    "$url" --data-binary @"$pact_file"
}

for pact in "$@"; do
  publish_pact "$pact"
done
