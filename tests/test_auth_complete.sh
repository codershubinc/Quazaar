#!/bin/bash

# Quazaar Authentication Flow Test Script
# Tests all authentication endpoints including utility functions

BASE_URL="http://127.0.0.1:8765"
API_BASE="$BASE_URL/api/v0.1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Test data
TEST_USERNAME="testuser_$(date +%s)"
TEST_PASSWORD="testpass123"
NEW_PASSWORD="newpass456"
TOKEN=""

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  Quazaar Authentication Flow Test Suite${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# Helper function to run a test
run_test() {
    local test_name=$1
    local expected_status=$2
    shift 2
    local curl_command=("$@")
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${YELLOW}Test $TOTAL_TESTS:${NC} $test_name"
    
    # Execute curl command and capture response
    response=$(eval "${curl_command[@]}" 2>&1)
    status_code=$(echo "$response" | grep -oP 'HTTP/\d+\.?\d?\s+\K\d+' | tail -1)
    body=$(echo "$response" | sed -n '/^{/,/^}/p')
    
    # Check if status code matches expected
    if [ "$status_code" == "$expected_status" ]; then
        echo -e "${GREEN}✓ PASSED${NC} - Status: $status_code"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}✗ FAILED${NC} - Expected: $expected_status, Got: $status_code"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo "$body"
    fi
    echo ""
}

# Test 1: Signup
echo -e "${BLUE}=== 1. User Registration ===${NC}"
run_test "Signup new user" "201" \
    curl -s -i -X POST "$API_BASE/signup" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"password\":\"$TEST_PASSWORD\"}"

# Test 2: Login
echo -e "${BLUE}=== 2. User Login ===${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"password\":\"$TEST_PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

run_test "Login with correct credentials" "200" \
    curl -s -i -X POST "$API_BASE/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"password\":\"$TEST_PASSWORD\"}"

echo -e "${GREEN}Token received: ${TOKEN:0:50}...${NC}"
echo ""

# Test 3: Get User Info
echo -e "${BLUE}=== 3. Get User Info ===${NC}"
run_test "Get user info with valid token (Bearer)" "200" \
    curl -s -i -X GET "$API_BASE/auth/user" \
    -H "Authorization: Bearer $TOKEN"

run_test "Get user info with valid token (query param)" "200" \
    curl -s -i -X GET "$API_BASE/auth/user?token=$TOKEN"

run_test "Get user info without token" "401" \
    curl -s -i -X GET "$API_BASE/auth/user"

# Test 4: Get Active Tokens
echo -e "${BLUE}=== 4. Get Active Tokens ===${NC}"
run_test "List user tokens" "200" \
    curl -s -i -X GET "$API_BASE/auth/tokens" \
    -H "Authorization: Bearer $TOKEN"

# Test 5: Protected Endpoints
echo -e "${BLUE}=== 5. Protected Endpoints ===${NC}"
run_test "Access WebSocket with valid token" "200" \
    curl -s -i -X GET "$BASE_URL/ws?token=$TOKEN"

run_test "Access WebSocket without token" "401" \
    curl -s -i -X GET "$BASE_URL/ws"

run_test "Access Player Info with valid token" "200" \
    curl -s -i -X GET "$API_BASE/player/info" \
    -H "Authorization: Bearer $TOKEN"

run_test "Access Player Info without token" "401" \
    curl -s -i -X GET "$API_BASE/player/info"

run_test "Access WiFi Info with valid token" "200" \
    curl -s -i -X GET "$API_BASE/system/wifi" \
    -H "Authorization: Bearer $TOKEN"

run_test "Access Bluetooth Info with valid token" "200" \
    curl -s -i -X GET "$API_BASE/system/bluetooth" \
    -H "Authorization: Bearer $TOKEN"

# Test 6: Change Password
echo -e "${BLUE}=== 6. Change Password ===${NC}"
run_test "Change password with valid credentials" "200" \
    curl -s -i -X POST "$API_BASE/auth/change-password" \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"$TOKEN\",\"old_password\":\"$TEST_PASSWORD\",\"new_password\":\"$NEW_PASSWORD\"}"

# Test 7: Login with new password
echo -e "${BLUE}=== 7. Login with New Password ===${NC}"
run_test "Login with old password (should fail)" "401" \
    curl -s -i -X POST "$API_BASE/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"password\":\"$TEST_PASSWORD\"}"

NEW_LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"password\":\"$NEW_PASSWORD\"}")

NEW_TOKEN=$(echo "$NEW_LOGIN_RESPONSE" | jq -r '.token')

run_test "Login with new password (should succeed)" "200" \
    curl -s -i -X POST "$API_BASE/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USERNAME\",\"password\":\"$NEW_PASSWORD\"}"

echo -e "${GREEN}New token received: ${NEW_TOKEN:0:50}...${NC}"
echo ""

# Test 8: Token Refresh
echo -e "${BLUE}=== 8. Token Refresh ===${NC}"
REFRESH_RESPONSE=$(curl -s -X POST "$API_BASE/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"$NEW_TOKEN\"}")

REFRESHED_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.token')

run_test "Refresh token" "200" \
    curl -s -i -X POST "$API_BASE/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"$NEW_TOKEN\"}"

echo -e "${GREEN}Refreshed token: ${REFRESHED_TOKEN:0:50}...${NC}"
echo ""

run_test "Old token should be invalid after refresh" "401" \
    curl -s -i -X GET "$API_BASE/auth/user" \
    -H "Authorization: Bearer $NEW_TOKEN"

run_test "New refreshed token should work" "200" \
    curl -s -i -X GET "$API_BASE/auth/user" \
    -H "Authorization: Bearer $REFRESHED_TOKEN"

# Test 9: Logout
echo -e "${BLUE}=== 9. Logout ===${NC}"
run_test "Logout with valid token" "200" \
    curl -s -i -X POST "$API_BASE/auth/logout" \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"$REFRESHED_TOKEN\"}"

run_test "Token should be invalid after logout" "401" \
    curl -s -i -X GET "$API_BASE/auth/user" \
    -H "Authorization: Bearer $REFRESHED_TOKEN"

# Test 10: Edge Cases
echo -e "${BLUE}=== 10. Edge Cases ===${NC}"
run_test "Signup with short username" "400" \
    curl -s -i -X POST "$API_BASE/signup" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"ab\",\"password\":\"$TEST_PASSWORD\"}"

run_test "Signup with short password" "400" \
    curl -s -i -X POST "$API_BASE/signup" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"validuser\",\"password\":\"short\"}"

run_test "Login with invalid credentials" "401" \
    curl -s -i -X POST "$API_BASE/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"nonexistent\",\"password\":\"wrongpass\"}"

run_test "Refresh with invalid token" "401" \
    curl -s -i -X POST "$API_BASE/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"invalid_token_xyz\"}"

run_test "Change password with wrong old password" "401" \
    curl -s -i -X POST "$API_BASE/auth/change-password" \
    -H "Content-Type: application/json" \
    -d "{\"token\":\"$TOKEN\",\"old_password\":\"wrongpass\",\"new_password\":\"newpass789\"}"

# Summary
echo ""
echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  Test Summary${NC}"
echo -e "${BLUE}================================================${NC}"
echo -e "Total Tests:  $TOTAL_TESTS"
echo -e "${GREEN}Passed:       $PASSED_TESTS${NC}"
echo -e "${RED}Failed:       $FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed!${NC}"
    exit 1
fi
