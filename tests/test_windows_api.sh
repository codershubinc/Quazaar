#!/bin/bash

# Test Windows-specific API endpoints
echo "🧪 Testing Windows API endpoints..."
echo ""

# Test Windows player info endpoint
echo "Testing GET /api/v0.1/player/info/windows"
curl -s -X GET "http://localhost:8765/api/v0.1/player/info/windows" | jq . 2>/dev/null || curl -s -X GET "http://localhost:8765/api/v0.1/player/info/windows"
echo ""
echo ""

# Test Windows active players list endpoint
echo "Testing GET /api/v0.1/player/windows/list"
curl -s -X GET "http://localhost:8765/api/v0.1/player/windows/list" | jq . 2>/dev/null || curl -s -X GET "http://localhost:8765/api/v0.1/player/windows/list"
echo ""
echo ""

# Test player-specific endpoint
echo "Testing GET /api/v0.1/player/info/player?player=spotify"
curl -s -X GET "http://localhost:8765/api/v0.1/player/info/player?player=spotify" | jq . 2>/dev/null || curl -s -X GET "http://localhost:8765/api/v0.1/player/info/player?player=spotify"
echo ""
echo ""

echo "✅ Windows API endpoint tests completed!"