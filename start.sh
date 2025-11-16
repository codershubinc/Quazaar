#!/bin/bash

# Quick Start Script - Run this to set everything up

cd ~/Github/Quazaar

echo "🚀 Quazaar Quick Start"
echo "===================="
echo ""

# 1. Build
echo "📦 Building server..."
go build -o quazaar 2>&1 | tail -5

if [ ! -f quazaar ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"
echo ""

# 2. Show URLs
echo "📋 Access URLs:"
echo "  🌐 Web UI: http://192.168.1.109:8765/auth.html"
echo "  🔌 WebSocket: http://192.168.1.109:8765/"
echo "  ❤️  Health: http://192.168.1.109:8765/api/health"
echo ""

# 3. Run server
echo "▶️  Starting server..."
echo "   (Press Ctrl+C to stop)"
echo ""
./quazaar
