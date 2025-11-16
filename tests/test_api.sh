#!/bin/bash

# Quazaar Auth API Test Script
# This script tests all the authentication endpoints

API_URL="http://192.168.1.109:8765"
USERNAME="testuser"
PASSWORD="password123"

echo "🚀 Quazaar Auth API Test Script"
echo "================================"
echo "API URL: $API_URL"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 1. Register User
echo -e "${BLUE}1. Registering user...${NC}"
REG_RESPONSE=$(curl -s -X POST "$API_URL/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

echo "Response: $REG_RESPONSE"
echo ""

# Check if registration was successful or if user already exists
if echo "$REG_RESPONSE" | grep -q "success.*true\|already registered"; then
    echo -e "${GREEN}✅ Registration complete (user may already exist)${NC}"
else
    echo -e "${RED}❌ Registration failed${NC}"
    exit 1
fi
echo ""

# 2. Login
echo -e "${BLUE}2. Logging in...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

echo "Response: $LOGIN_RESPONSE"
echo ""

# Extract auth token from response
AUTH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$AUTH_TOKEN" ]; then
    AUTH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"auth_token":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$AUTH_TOKEN" ]; then
    echo -e "${RED}❌ Could not extract auth token from login response${NC}"
    echo "Full response: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Login successful${NC}"
echo "Auth Token: ${AUTH_TOKEN:0:20}..."
echo ""

# 3. Create Service Token
echo -e "${BLUE}3. Creating service token...${NC}"
TOKEN_RESPONSE=$(curl -s -X POST "$API_URL/api/tokens/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: $AUTH_TOKEN" \
  -d "{\"name\":\"Mobile App\",\"service\":\"mobile\",\"duration_hours\":720}")

echo "Response: $TOKEN_RESPONSE"
echo ""

SERVICE_TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$SERVICE_TOKEN" ]; then
    echo -e "${GREEN}✅ Service token created${NC}"
    echo "Service Token: ${SERVICE_TOKEN:0:20}..."
    echo ""
else
    echo -e "${RED}❌ Failed to create service token${NC}"
    echo "Response: $TOKEN_RESPONSE"
fi
echo ""

# 4. List Tokens
echo -e "${BLUE}4. Listing all tokens...${NC}"
LIST_RESPONSE=$(curl -s -X GET "$API_URL/api/tokens/list" \
  -H "Authorization: $AUTH_TOKEN")

echo "Response: $LIST_RESPONSE"
echo ""

if echo "$LIST_RESPONSE" | grep -q "tokens"; then
    TOKEN_COUNT=$(echo "$LIST_RESPONSE" | grep -o '"count":[0-9]*' | cut -d':' -f2)
    echo -e "${GREEN}✅ Listed tokens (count: $TOKEN_COUNT)${NC}"
else
    echo -e "${RED}❌ Failed to list tokens${NC}"
fi
echo ""

# 5. WebSocket Connection Test (if SERVICE_TOKEN was created)
if [ -n "$SERVICE_TOKEN" ]; then
    echo -e "${BLUE}5. Testing WebSocket connection...${NC}"
    echo "WebSocket URL: ws://192.168.1.109:8765/ws?token=$SERVICE_TOKEN"
    echo ""
    echo -e "${YELLOW}Use this token in the WebSocket connection:${NC}"
    echo "$SERVICE_TOKEN"
    echo ""
fi

# Summary
echo -e "${GREEN}================================${NC}"
echo "✅ Test Complete!"
echo ""
echo "Summary:"
echo "  Username: $USERNAME"
echo "  Auth Token: ${AUTH_TOKEN:0:30}..."
if [ -n "$SERVICE_TOKEN" ]; then
    echo "  Service Token: ${SERVICE_TOKEN:0:30}..."
fi
echo ""
echo "Next steps:"
echo "  1. Open http://192.168.1.109:8765/auth.html in browser"
echo "  2. Or use the tokens printed above to test API"
echo ""
