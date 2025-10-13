#!/bin/bash

# Comprehensive Test Suite for socialapp-cli
# Tests all OpenAPI endpoints through CLI commands
# Prerequisites:
# - socialapp server running on localhost:8086
# - socialapp-cli binary built
# - Database seeded with initial admin user

set -e

# Configuration
HOST="${HOST:-http://localhost:8086}"
CLI="${CLI:-./socialapp-cli}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASS="${ADMIN_PASS:-admin}"

# Test data
TIMESTAMP=$(date +%s)
TEST_USER="testuser_${TIMESTAMP}"
TEST_PASS="TestPass123!"
TEST_EMAIL="test_${TIMESTAMP}@example.com"
TEST_USER2="testuser2_${TIMESTAMP}"
TEST_PASS2="TestPass456!"
TEST_EMAIL2="test2_${TIMESTAMP}@example.com"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Test result tracking
declare -a FAILED_TEST_NAMES=()

run_test() {
    local test_name="$1"
    local command="$2"
    local should_succeed="${3:-true}"
    
    ((TOTAL_TESTS++))
    echo -e "\n${BLUE}[TEST $TOTAL_TESTS] $test_name${NC}"
    echo -e "${YELLOW}Command: $command${NC}"
    
    if [ "$should_succeed" = "true" ]; then
        if eval "$command" > /tmp/cli_test_output_$$.txt 2>&1; then
            echo -e "${GREEN}✓ PASSED${NC}"
            ((PASSED_TESTS++))
            cat /tmp/cli_test_output_$$.txt
            return 0
        else
            echo -e "${RED}✗ FAILED${NC}"
            cat /tmp/cli_test_output_$$.txt
            ((FAILED_TESTS++))
            FAILED_TEST_NAMES+=("$test_name")
            return 1
        fi
    else
        if eval "$command" > /tmp/cli_test_output_$$.txt 2>&1; then
            echo -e "${RED}✗ FAILED (should have failed but succeeded)${NC}"
            cat /tmp/cli_test_output_$$.txt
            ((FAILED_TESTS++))
            FAILED_TEST_NAMES+=("$test_name")
            return 1
        else
            echo -e "${GREEN}✓ PASSED (failed as expected)${NC}"
            ((PASSED_TESTS++))
            return 0
        fi
    fi
}

skip_test() {
    local test_name="$1"
    local reason="$2"
    
    ((TOTAL_TESTS++))
    ((SKIPPED_TESTS++))
    echo -e "\n${YELLOW}[TEST $TOTAL_TESTS - SKIPPED] $test_name${NC}"
    echo -e "${YELLOW}Reason: $reason${NC}"
}

extract_json_field() {
    local json="$1"
    local field="$2"
    echo "$json" | grep -o "\"$field\": *[^,}]*" | sed 's/.*: *//' | tr -d '"'
}

echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║       Socialapp CLI Comprehensive Test Suite                ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo -e "Host: $HOST"
echo -e "CLI: $CLI"
echo -e "Start time: $(date)"
echo ""

if [ ! -f "$CLI" ]; then
    echo -e "${RED}Error: CLI binary not found at $CLI${NC}"
    exit 1
fi

# ============================================================================
# USER MANAGEMENT TESTS
# ============================================================================
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}USER MANAGEMENT TESTS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"

# Test 1: Create User
run_test "Create User" \
    "$CLI create user -u $TEST_USER -p $TEST_PASS -e $TEST_EMAIL -f Test -l User --host $HOST"

# Test 2: Create Second User (for following tests)
run_test "Create Second User" \
    "$CLI create user -u $TEST_USER2 -p $TEST_PASS2 -e $TEST_EMAIL2 -f Test2 -l User2 --host $HOST"

# Test 3: List Users
run_test "List Users" \
    "$CLI list users -u $ADMIN_USER -p $ADMIN_PASS --host $HOST --pagesize 5"

# Test 4: Get User by Username
run_test "Get User by Username" \
    "$CLI get user $TEST_USER -u $ADMIN_USER -p $ADMIN_PASS --host $HOST"

# Test 5: Update User
run_test "Update User" \
    "$CLI update user $TEST_USER -u $ADMIN_USER -p $ADMIN_PASS -e updated_$TEST_EMAIL --host $HOST"

# ============================================================================
# COMMENT TESTS
# ============================================================================
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}COMMENT TESTS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"

# Test 6: Create Comment (using admin credentials for authentication)
if run_test "Create Comment" \
    "$CLI create comment -u \"$ADMIN_USER\" -p \"$ADMIN_PASS\" --content \"Test comment from CLI\" --author \"$TEST_USER\" --host \"$HOST\""; then
    # Extract comment ID from the last successful output
    COMMENT_ID=$(cat /tmp/cli_test_output_$$.txt | grep -o '"id": *[0-9]*' | head -1 | awk '{print $2}')
    if [ -n "$COMMENT_ID" ]; then
        echo -e "${GREEN}Comment ID extracted: $COMMENT_ID${NC}"
    fi
else
    COMMENT_ID=""
fi

# Test 7: Get Comment by ID
if [ -n "$COMMENT_ID" ]; then
    run_test "Get Comment by ID" \
        "$CLI get comment $COMMENT_ID -u $ADMIN_USER -p $ADMIN_PASS --host $HOST"
else
    skip_test "Get Comment by ID" "No comment ID available"
fi

# Test 8: Get User Comments
run_test "Get User Comments" \
    "$CLI get user-comments $TEST_USER -u $ADMIN_USER -p $ADMIN_PASS --host $HOST"

# Test 9: Get User Feed
run_test "Get User Feed" \
    "$CLI get feed -u $ADMIN_USER -p $ADMIN_PASS --host $HOST"

# ============================================================================
# FOLLOWING TESTS (if implemented)
# ============================================================================
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}FOLLOWING TESTS (if implemented)${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"

# These tests will be enabled once the commands are implemented
skip_test "Follow User" "Command not yet implemented"
skip_test "Get User Followers" "Command not yet implemented"
skip_test "Get User Following" "Command not yet implemented"
skip_test "Unfollow User" "Command not yet implemented"

# ============================================================================
# ROLE & SCOPE TESTS
# ============================================================================
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}ROLE & SCOPE TESTS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"

# Test: Create Role
if run_test "Create Role" \
    "$CLI create role --name \"test_role_${TIMESTAMP}\" --description \"Test role from CLI\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"; then
    TEST_ROLE=$(cat /tmp/cli_test_output_$$.txt | grep -o '"name": *"[^"]*"' | head -1 | sed 's/"name": *"\([^"]*\)"/\1/')
    echo -e "${GREEN}Role name extracted: $TEST_ROLE${NC}"
else
    TEST_ROLE=""
fi

# Test: List Roles
run_test "List Roles" \
    "$CLI list roles -u $ADMIN_USER -p $ADMIN_PASS --host $HOST"

# Test: Get Role
if [ -n "$TEST_ROLE" ]; then
    run_test "Get Role" \
        "$CLI get role 2 --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Get Role" "No role name available"
fi

# Test: Update Role
if [ -n "$TEST_ROLE" ]; then
    run_test "Update Role" \
        "$CLI update role \"$TEST_ROLE\" --description \"Updated test role\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Update Role" "No role name available"
fi

# Test: Create Scope
if run_test "Create Scope" \
    "$CLI create scope --name \"test.scope.${TIMESTAMP}\" --description \"Test scope from CLI\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"; then
    TEST_SCOPE=$(cat /tmp/cli_test_output_$$.txt | grep -o '"name": *"[^"]*"' | head -1 | sed 's/"name": *"\([^"]*\)"/\1/')
    echo -e "${GREEN}Scope name extracted: $TEST_SCOPE${NC}"
else
    TEST_SCOPE=""
fi

# Test: List Scopes
run_test "List Scopes" \
    "$CLI list scopes --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"

# Test: Get Scope
# TODO: Fix hardcoded ID once create scope returns ID
if [ -n "$TEST_SCOPE" ]; then
    run_test "Get Scope" \
        "$CLI get scope 2 --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Get Scope" "No scope name available"
fi

# Test: Update Scope
if [ -n "$TEST_SCOPE" ]; then
    run_test "Update Scope" \
        "$CLI update scope 2 --description \"Updated test scope\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Update Scope" "No scope name available"
fi

# Test: Add Scope to Role
if [ -n "$TEST_ROLE" ] && [ -n "$TEST_SCOPE" ]; then
    run_test "Add Scope to Role" \
        "$CLI add scope-to-role \"$TEST_ROLE\" \"$TEST_SCOPE\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Add Scope to Role" "No role or scope name available"
fi

# Test: Get Role Scopes
if [ -n "$TEST_ROLE" ]; then
    run_test "Get Role Scopes" \
        "$CLI get role-scopes \"$TEST_ROLE\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Get Role Scopes" "No role name available"
fi

# Test: Get User Roles
run_test "Get User Roles" \
    "$CLI get user-roles $TEST_USER --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"

# Test: Update User Roles
if [ -n "$TEST_ROLE" ]; then
    run_test "Update User Roles" \
        "$CLI update user-roles $TEST_USER --roles \"$TEST_ROLE\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Update User Roles" "No role name available"
fi

# ============================================================================
# URL SHORTENER TESTS
# ============================================================================
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}URL SHORTENER TESTS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"

# Test: Create Short URL
if run_test "Create Short URL" \
    "$CLI create url --url \"https://example.com/test/${TIMESTAMP}\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"; then
    SHORT_URL=$(cat /tmp/cli_test_output_$$.txt | grep -o '"short_url": *"[^"]*"' | head -1 | sed 's/"short_url": *"\([^"]*\)"/\1/')
    echo -e "${GREEN}Short URL extracted: $SHORT_URL${NC}"
else
    SHORT_URL=""
fi

# Test: Get URL Data
if [ -n "$SHORT_URL" ]; then
    run_test "Get URL Data" \
        "$CLI get url-data \"$SHORT_URL\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Get URL Data" "No short URL available"
fi

# Test: Delete URL (in cleanup section)
# Will be added to cleanup section below

# ============================================================================
# CLEANUP
# ============================================================================
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}CLEANUP${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"

# Delete Short URL
if [ -n "$SHORT_URL" ]; then
    run_test "Delete Short URL" \
        "$CLI delete url \"$SHORT_URL\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Delete Short URL" "No short URL to delete"
fi

# Remove Scope from Role
if [ -n "$TEST_ROLE" ] && [ -n "$TEST_SCOPE" ]; then
    run_test "Remove Scope from Role" \
        "$CLI remove scope-from-role \"$TEST_ROLE\" \"$TEST_SCOPE\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Remove Scope from Role" "No role or scope to clean up"
fi

# Delete Scope
# if [ -n "$TEST_SCOPE" ]; then
#     run_test "Delete Scope" \
#         "$CLI delete scope 2 --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
# else
#     skip_test "Delete Scope" "No scope to delete"
# fi

# Delete Role
if [ -n "$TEST_ROLE" ]; then
    run_test "Delete Role" \
        "$CLI delete role \"$TEST_ROLE\" --username $ADMIN_USER --password $ADMIN_PASS --host $HOST"
else
    skip_test "Delete Role" "No role to delete"
fi

# Delete Test Users
run_test "Delete Test User 1" \
    "$CLI delete user $TEST_USER -u $ADMIN_USER -p $ADMIN_PASS --host $HOST"

run_test "Delete Test User 2" \
    "$CLI delete user $TEST_USER2 -u $ADMIN_USER -p $ADMIN_PASS --host $HOST"

# ============================================================================
# TEST SUMMARY
# ============================================================================
echo -e "\n${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                      TEST SUMMARY                            ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo -e "Total Tests:   $TOTAL_TESTS"
echo -e "${GREEN}Passed:        $PASSED_TESTS${NC}"
echo -e "${RED}Failed:        $FAILED_TESTS${NC}"
echo -e "${YELLOW}Skipped:       $SKIPPED_TESTS${NC}"
echo -e "End time:      $(date)"

if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "\n${RED}Failed Tests:${NC}"
    for test_name in "${FAILED_TEST_NAMES[@]}"; do
        echo -e "  ${RED}✗${NC} $test_name"
    done
    exit 1
else
    echo -e "\n${GREEN}All tests passed! 🎉${NC}"
    exit 0
fi

# Cleanup temp file
rm -f /tmp/cli_test_output_$$.txt
