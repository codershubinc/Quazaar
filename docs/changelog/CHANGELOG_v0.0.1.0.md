# Changelog - v0.0.1.0

**Release Date:** Early November 2025  
**Project Name:** Blitz  
**Status:** 🟢 Initial Beta Release

---

## Overview

First beta release of Blitz - a media control and system manager for Linux. Initial implementation with basic WebSocket functionality and media player integration.

---

## ✨ New Features

### Core Features
- **WebSocket Server** - Real-time communication for media updates
- **Media Player Integration** - Basic playerctl integration for Linux media control
- **System Information** - WiFi and Bluetooth information retrieval
- **HTTP Server** - Basic HTTP endpoint serving

### Media Control
- Get current player information
- Track metadata retrieval
- Basic playback information

### System Features
- WiFi information (SSID, signal strength, etc.)
- Bluetooth device enumeration
- System information APIs

---

## 🏗️ Architecture

### Project Structure
```
blitz/
├── main.go                 # Application entry point
├── mediaInfo.go           # Media player integration
├── spotifyArtwork.go      # Spotify artwork retrieval
├── wifiInfo.go            # WiFi information
├── bluetoothInfo.go       # Bluetooth information
└── websocket.go           # WebSocket implementation
```

### Technology Stack
- **Language:** Go 1.21+
- **Dependencies:**
  - `gorilla/websocket` - WebSocket support
  - `playerctl` - Media player control
  - `nmcli` - Network management
  - `bluetoothctl` - Bluetooth management

---

## 📋 API Endpoints

### WebSocket
```bash
GET /ws
# Real-time media information updates
```

### Home
```bash
GET /
# Basic home page
```

---

## 🔧 Configuration

### Default Settings
- **Host:** 127.0.0.1
- **Port:** 8765
- **Polling Interval:** 1 second

---

## 📦 Release Assets

- `blitz_v0.0.1.0` - Linux x86-64 binary
- Size: ~8 MB
- Platform: Linux (kernel 4.4.0+)

---

## 🐛 Known Issues

- No authentication system
- No database integration
- Limited error handling
- No API versioning
- Single file architecture
- Platform-specific code not separated

---

## 📝 Notes

This is the initial beta release focused on core functionality. The project is designed for local use on Linux systems with D-Bus/MPRIS support.

---

## 🔗 Links

- **Repository:** github.com/codershubinc/Blitz
- **License:** MIT

---

**Contributors:** Swapnil Ingle  
**Project Status:** 🟢 Beta - Initial Release
