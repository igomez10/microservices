#!/bin/bash

# Simple test script for new CLI commands
# Tests role, scope, and URL management commands

set -e

CLI="./socialapp-cli"
HOST="http://localhost:8086"

echo "=== Testing New CLI Commands ==="
echo ""

# Test 1: List roles
echo "Test 1: List roles"
$CLI list roles --host $HOST --limit 10
echo "✓ List roles passed"
echo ""

# Test 2: List scopes
echo "Test 2: List scopes"
$CLI list scopes --host $HOST --limit 10
echo "✓ List scopes passed"
echo ""

# Test 3: Create a role
echo "Test 3: Create a test role"
ROLE_OUTPUT=$($CLI create role --name "test_role_$(date +%s)" --description "Test role" --host $HOST)
echo "$ROLE_OUTPUT"
ROLE_ID=$(echo "$ROLE_OUTPUT" | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | head -1)
echo "✓ Created role with ID: $ROLE_ID"
echo ""

# Test 4: Get role
if [ ! -z "$ROLE_ID" ]; then
    echo "Test 4: Get role by ID"
    $CLI get role $ROLE_ID --host $HOST
    echo "✓ Get role passed"
    echo ""
fi

# Test 5: Create a scope
echo "Test 5: Create a test scope"
SCOPE_OUTPUT=$($CLI create scope --name "test.scope.$(date +%s)" --description "Test scope" --host $HOST)
echo "$SCOPE_OUTPUT"
SCOPE_ID=$(echo "$SCOPE_OUTPUT" | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | head -1)
echo "✓ Created scope with ID: $SCOPE_ID"
echo ""

# Test 6: Get scope
if [ ! -z "$SCOPE_ID" ]; then
    echo "Test 6: Get scope by ID"
    $CLI get scope $SCOPE_ID --host $HOST
    echo "✓ Get scope passed"
    echo ""
fi

# Test 7: Get role scopes (should be empty for new role)
if [ ! -z "$ROLE_ID" ]; then
    echo "Test 7: Get scopes for role"
    $CLI get role-scopes $ROLE_ID --host $HOST
    echo "✓ Get role-scopes passed"
    echo ""
fi

# Test 8: Add scope to role
if [ ! -z "$ROLE_ID" ] && [ ! -z "$SCOPE_ID" ]; then
    echo "Test 8: Add scope to role"
    $CLI add scope-to-role $ROLE_ID "test.scope.$(date +%s)" --host $HOST
    echo "✓ Add scope to role passed"
    echo ""
fi

# Test 9: Get user roles (for admin user)
echo "Test 9: Get user roles"
$CLI get user-roles admin --host $HOST
echo "✓ Get user-roles passed"
echo ""

# Test 10: Create shortened URL
echo "Test 10: Create shortened URL"
URL_ALIAS="test_$(date +%s)"
$CLI create url --url "https://example.com" --alias "$URL_ALIAS" --host $HOST
echo "✓ Create URL passed"
echo ""

# Test 11: Get URL data
echo "Test 11: Get URL data"
$CLI get url-data "$URL_ALIAS" --host $HOST
echo "✓ Get URL data passed"
echo ""

# Cleanup
echo "=== Cleanup ==="

# Delete URL
if [ ! -z "$URL_ALIAS" ]; then
    echo "Deleting URL: $URL_ALIAS"
    $CLI delete url "$URL_ALIAS" --host $HOST
    echo "✓ Deleted URL"
fi

# Delete scope
if [ ! -z "$SCOPE_ID" ]; then
    echo "Deleting scope: $SCOPE_ID"
    $CLI delete scope $SCOPE_ID --host $HOST
    echo "✓ Deleted scope"
fi

# Delete role
if [ ! -z "$ROLE_ID" ]; then
    echo "Deleting role: $ROLE_ID"
    $CLI delete role $ROLE_ID --host $HOST
    echo "✓ Deleted role"
fi

echo ""
echo "=== All Tests Passed! ==="
