#!/bin/bash

# Test script for socialapp-cli
# This script tests all implemented CLI commands
# Prerequisites:
# - socialapp server running on localhost:8086
# - socialapp-cli binary built and in PATH or in current directory

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
HOST="http://localhost:8086"
CLI_BINARY="./socialapp-cli"
TEST_USERNAME="testuser_$(date +%s)"
TEST_PASSWORD="TestPass123!"
TEST_EMAIL="test_$(date +%s)@example.com"
TEST_FIRSTNAME="Test"
TEST_LASTNAME="User"
ADMIN_USERNAME="${ADMIN_USERNAME:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin}"

# Check if CLI binary exists
if [ ! -f "$CLI_BINARY" ]; then
    echo -e "${RED}Error: CLI binary not found at $CLI_BINARY${NC}"
    echo "Please build the CLI first: make build"
    exit 1
fi

# Function to print test header
print_test() {
    echo -e "\n${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}TEST: $1${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Function to wait for user input
wait_for_key() {
    if [ "${AUTO_CONTINUE:-}" != "true" ]; then
        echo -e "\n${YELLOW}Press any key to continue...${NC}"
        read -n 1 -s
    fi
}

echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║         Socialapp CLI Test Suite                          ║${NC}"
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "Host: $HOST"
echo -e "CLI: $CLI_BINARY"
echo -e "Test User: $TEST_USERNAME"
echo ""

# ============================================================================
# TEST 1: Create User
# ============================================================================
print_test "1. Create User (POST /v1/users)"
echo "Command: $CLI_BINARY create user -u $TEST_USERNAME -p $TEST_PASSWORD -e $TEST_EMAIL -f $TEST_FIRSTNAME -l $TEST_LASTNAME --host $HOST"
if $CLI_BINARY create user \
    -u "$TEST_USERNAME" \
    -p "$TEST_PASSWORD" \
    -e "$TEST_EMAIL" \
    -f "$TEST_FIRSTNAME" \
    -l "$TEST_LASTNAME" \
    --host "$HOST"; then
    print_success "User created successfully"
else
    print_error "Failed to create user"
    exit 1
fi
wait_for_key

# ============================================================================
# TEST 2: List Users
# ============================================================================
print_test "2. List Users (GET /v1/users)"
echo "Command: $CLI_BINARY list users -u $ADMIN_USERNAME -p $ADMIN_PASSWORD --host $HOST --pagesize 10"
if $CLI_BINARY list users \
    -u "$ADMIN_USERNAME" \
    -p "$ADMIN_PASSWORD" \
    --host "$HOST" \
    --pagesize 10; then
    print_success "Users listed successfully"
else
    print_error "Failed to list users"
fi
wait_for_key

# ============================================================================
# TEST 3: Get User by Username
# ============================================================================
print_test "3. Get User by Username (GET /v1/users/{username})"
echo "Command: $CLI_BINARY get user $TEST_USERNAME -u $ADMIN_USERNAME -p $ADMIN_PASSWORD --host $HOST"
if $CLI_BINARY get user "$TEST_USERNAME" \
    -u "$ADMIN_USERNAME" \
    -p "$ADMIN_PASSWORD" \
    --host "$HOST"; then
    print_success "User retrieved successfully"
else
    print_error "Failed to get user"
fi
wait_for_key

# ============================================================================
# TEST 4: Update User
# ============================================================================
print_test "4. Update User (PUT /v1/users/{username})"
NEW_EMAIL="updated_${TEST_EMAIL}"
echo "Command: $CLI_BINARY update user $TEST_USERNAME -u $ADMIN_USERNAME -p $ADMIN_PASSWORD -e $NEW_EMAIL --host $HOST"
if $CLI_BINARY update user "$TEST_USERNAME" \
    -u "$ADMIN_USERNAME" \
    -p "$ADMIN_PASSWORD" \
    -e "$NEW_EMAIL" \
    --host "$HOST"; then
    print_success "User updated successfully"
else
    print_error "Failed to update user"
fi
wait_for_key

# ============================================================================
# TEST 5: Create Comment
# ============================================================================
print_test "5. Create Comment (POST /v1/comments)"
COMMENT_CONTENT="This is a test comment from CLI at $(date)"
echo "Command: $CLI_BINARY create comment -u $TEST_USERNAME -p $TEST_PASSWORD --content \"$COMMENT_CONTENT\" --author $TEST_USERNAME --host $HOST"
if COMMENT_OUTPUT=$($CLI_BINARY create comment \
    -u "$TEST_USERNAME" \
    -p "$TEST_PASSWORD" \
    --content "$COMMENT_CONTENT" \
    --author "$TEST_USERNAME" \
    --host "$HOST"); then
    print_success "Comment created successfully"
    # Extract comment ID for later use
    COMMENT_ID=$(echo "$COMMENT_OUTPUT" | grep -o '"id": [0-9]*' | head -1 | awk '{print $2}')
    echo "Comment ID: $COMMENT_ID"
else
    print_error "Failed to create comment"
    COMMENT_ID=""
fi
wait_for_key

# ============================================================================
# TEST 6: Get Comment by ID
# ============================================================================
if [ -n "$COMMENT_ID" ]; then
    print_test "6. Get Comment by ID (GET /v1/comments/{id})"
    echo "Command: $CLI_BINARY get comment $COMMENT_ID -u $TEST_USERNAME -p $TEST_PASSWORD --host $HOST"
    if $CLI_BINARY get comment "$COMMENT_ID" \
        -u "$TEST_USERNAME" \
        -p "$TEST_PASSWORD" \
        --host "$HOST"; then
        print_success "Comment retrieved successfully"
    else
        print_error "Failed to get comment"
    fi
    wait_for_key
else
    print_error "Skipping Test 6: No comment ID available"
fi

# ============================================================================
# TEST 7: Get User Comments
# ============================================================================
print_test "7. Get User Comments (GET /v1/users/{username}/comments)"
echo "Command: $CLI_BINARY get user-comments $TEST_USERNAME -u $TEST_USERNAME -p $TEST_PASSWORD --host $HOST"
if $CLI_BINARY get user-comments "$TEST_USERNAME" \
    -u "$TEST_USERNAME" \
    -p "$TEST_PASSWORD" \
    --host "$HOST"; then
    print_success "User comments retrieved successfully"
else
    print_error "Failed to get user comments"
fi
wait_for_key

# ============================================================================
# TEST 8: Get User Feed
# ============================================================================
print_test "8. Get User Feed (GET /v1/feed)"
echo "Command: $CLI_BINARY get feed -u $TEST_USERNAME -p $TEST_PASSWORD --host $HOST"
if $CLI_BINARY get feed \
    -u "$TEST_USERNAME" \
    -p "$TEST_PASSWORD" \
    --host "$HOST"; then
    print_success "User feed retrieved successfully"
else
    print_error "Failed to get user feed"
fi
wait_for_key

# ============================================================================
# TEST 9: Get URL (if URL shortener is available)
# ============================================================================
print_test "9. Get URL Data (GET /v1/urls/{alias}/data) - Optional"
echo "Command: $CLI_BINARY get url example -u $TEST_USERNAME -p $TEST_PASSWORD --host $HOST"
echo "(This test may fail if no URLs exist - that's okay)"
if $CLI_BINARY get url example \
    -u "$TEST_USERNAME" \
    -p "$TEST_PASSWORD" \
    --host "$HOST" 2>/dev/null; then
    print_success "URL retrieved successfully"
else
    echo -e "${YELLOW}⚠ URL test skipped or failed (may not exist)${NC}"
fi
wait_for_key

# ============================================================================
# TEST 10: Delete User (Cleanup)
# ============================================================================
print_test "10. Delete User (DELETE /v1/users/{username})"
echo "Command: $CLI_BINARY delete user $TEST_USERNAME -u $ADMIN_USERNAME -p $ADMIN_PASSWORD --host $HOST"
if $CLI_BINARY delete user "$TEST_USERNAME" \
    -u "$ADMIN_USERNAME" \
    -p "$ADMIN_PASSWORD" \
    --host "$HOST"; then
    print_success "User deleted successfully"
else
    print_error "Failed to delete user"
fi
wait_for_key

# ============================================================================
# Summary
# ============================================================================
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                   Test Suite Complete                      ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo -e "\n${YELLOW}All tests completed!${NC}"
echo -e "\nNote: Some commands may have been skipped if prerequisites weren't met."
echo -e "To run in auto mode (no pauses): ${GREEN}AUTO_CONTINUE=true $0${NC}"
echo -e "To set custom admin credentials:"
echo -e "  ${GREEN}ADMIN_USERNAME=youradmin ADMIN_PASSWORD=yourpass $0${NC}"
