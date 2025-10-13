#!/bin/bash
# Quick test to verify CLI and server connectivity

set -e

CLI="./socialapp-cli"
HOST="http://localhost:8086"

echo "=== Quick CLI Test ==="
echo "1. Testing CLI binary..."
$CLI --help > /dev/null && echo "✓ CLI binary works"

echo ""
echo "2. Testing server connectivity..."
if curl -s "$HOST/" > /dev/null 2>&1; then
    echo "✓ Server is responding at $HOST"
else
    echo "✗ Server is NOT responding at $HOST"
    echo "Please start the server first: make start-dev-server"
    exit 1
fi

echo ""
echo "3. Testing list users command..."
if $CLI list users -u admin -p admin --host "$HOST" --pagesize 5 2>&1 | head -10; then
    echo ""
    echo "✓ List users command works"
else
    echo "✗ List users command failed"
    echo "Check admin credentials (default: admin/admin)"
fi

echo ""
echo "=== Quick test complete ==="
