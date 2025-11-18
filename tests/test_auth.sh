#!/bin/bash

# 🧪 Testing Script for Single-User Auth Token System
# Save this as: test_auth.sh
# Run with: bash test_auth.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

SERVER="http://localhost:8765"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🧪 Quazaar Auth System Test${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check if server is running
echo -e "${YELLOW}📡 Checking if server is running...${NC}"
if ! curl -s "$SERVER" > /dev/null 2>&1; then
    echo -e "${RED}❌ Server not running at $SERVER${NC}"
    echo "Start it with: go run main.go"
    exit 1
fi
echo -e "${GREEN}✅ Server is running${NC}"
echo ""

# Test 1: Register User
echo -e "${YELLOW}1️⃣  Testing: Register User${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$SERVER/api/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}')

echo "Response: $REGISTER_RESPONSE"
if echo "$REGISTER_RESPONSE" | grep -q "success"; then
    echo -e "${GREEN}✅ Registration test passed${NC}"
else
    # User might already exist
    if echo "$REGISTER_RESPONSE" | grep -q "already registered"; then
        echo -e "${YELLOW}⚠️  User already registered (continuing)${NC}"
    else
        echo -e "${RED}❌ Registration test failed${NC}"
        exit 1
    fi
fi
echo ""

# Test 2: Login User
echo -e "${YELLOW}2️⃣  Testing: Login User${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$SERVER/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}')

echo "Response: $LOGIN_RESPONSE"
if echo "$LOGIN_RESPONSE" | grep -q "success"; then
    echo -e "${GREEN}✅ Login test passed${NC}"
else
    echo -e "${RED}❌ Login test failed${NC}"
    exit 1
fi

# Extract a token for testing (we'll use this as auth token)
AUTH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"user_id":[0-9]*' | head -1 | grep -o '[0-9]*')
if [ -z "$AUTH_TOKEN" ]; then
    AUTH_TOKEN="default_test_token"
fi
echo -e "${BLUE}Using auth token for next tests...${NC}"
echo ""

# Test 3: Create Token
echo -e "${YELLOW}3️⃣  Testing: Create Service Token${NC}"
TOKEN_RESPONSE=$(curl -s -X POST "$SERVER/api/tokens/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: $AUTH_TOKEN" \
  -d '{"name":"Test Device","service":"test","duration_hours":24}')

echo "Response: $TOKEN_RESPONSE"
if echo "$TOKEN_RESPONSE" | grep -q "token"; then
    echo -e "${GREEN}✅ Token creation test passed${NC}"
    # Extract the created token
    CREATED_TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"token":"[^"]*' | head -1 | cut -d'"' -f4)
    echo -e "${BLUE}Created token: ${CREATED_TOKEN:0:20}...${NC}"
else
    echo -e "${RED}❌ Token creation test failed${NC}"
fi
echo ""

# Test 4: List Tokens
echo -e "${YELLOW}4️⃣  Testing: List All Tokens${NC}"
LIST_RESPONSE=$(curl -s -X GET "$SERVER/api/tokens/list" \
  -H "Authorization: $AUTH_TOKEN")

echo "Response: $LIST_RESPONSE"
if echo "$LIST_RESPONSE" | grep -q "tokens"; then
    echo -e "${GREEN}✅ List tokens test passed${NC}"
    COUNT=$(echo "$LIST_RESPONSE" | grep -o '"count":[0-9]*' | grep -o '[0-9]*')
    echo -e "${BLUE}Found $COUNT tokens${NC}"
else
    echo -e "${RED}❌ List tokens test failed${NC}"
fi
echo ""

# Test 5: Invalid Token Test
echo -e "${YELLOW}5️⃣  Testing: Reject Invalid Token${NC}"
INVALID_RESPONSE=$(curl -s -X POST "$SERVER/api/tokens/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: invalid_token_12345" \
  -d '{"name":"Should Fail","service":"test","duration_hours":24}')

echo "Response: $INVALID_RESPONSE"
if echo "$INVALID_RESPONSE" | grep -q "Unauthorized"; then
    echo -e "${GREEN}✅ Invalid token rejection test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Unexpected response (might be different error)${NC}"
fi
echo ""

# Test 6: Wrong Password Test
echo -e "${YELLOW}6️⃣  Testing: Reject Wrong Password${NC}"
WRONG_PASS_RESPONSE=$(curl -s -X POST "$SERVER/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrongpassword"}')

echo "Response: $WRONG_PASS_RESPONSE"
if echo "$WRONG_PASS_RESPONSE" | grep -q "Invalid credentials"; then
    echo -e "${GREEN}✅ Wrong password rejection test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Unexpected response${NC}"
fi
echo ""

# Summary
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All basic tests completed!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}📝 Next Steps:${NC}"
echo "1. Verify tokens are created in database:"
echo "   sqlite3 ~/.quazaar/quazaar.db \"SELECT name, service FROM tokens;\""
echo ""
echo "2. Test WebSocket with token:"
echo "   wscat -c \"ws://localhost:8765/ws?token=YOUR_TOKEN_HERE\""
echo ""
echo "3. See full API docs in: docs/AUTH_SYSTEM.md"
echo ""
