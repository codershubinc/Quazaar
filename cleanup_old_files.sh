#!/bin/bash
# Script to remove old duplicate files after refactoring

echo "Removing old duplicate files..."

# Main spotify package
rm -f internal/spotify/conf.go
rm -f internal/spotify/handler.go

# Auth package
rm -f internal/spotify/auth/handler.go

# Devices package
rm -f internal/spotify/devices/handler.go
rm -f internal/spotify/devices/devices.go

# Tokens package
rm -f internal/spotify/tokens/tokens.go
rm -f internal/spotify/tokens/helper.go

echo "Cleanup complete!"
echo "Running build test..."
go build -o quazaar ./cmd/server

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi
