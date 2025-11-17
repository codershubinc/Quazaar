# Changelog - v0.0.1.2

**Release Date:** November 16, 2025  
**Project Name:** Blitz  
**Status:** 🟡 Beta Enhancement

---

## Overview

Third beta release with enhanced features and expanded API endpoints. Major improvements to media control and system information APIs.

---

## ✨ What's New

### New Features

#### Player Controls
- **Play/Pause Toggle** - `POST /player/play-pause`
- **Explicit Play** - `POST /player/play`
- **Explicit Pause** - `POST /player/pause`
- **Next Track** - `POST /player/next`
- **Previous Track** - `POST /player/previous`

#### Media Information
- **Player Info** - `GET /player/info`
- **Active Players List** - `GET /player/list`
- Enhanced metadata retrieval
- Better artwork handling

#### System Information
- **WiFi Information** - `GET /system/wifi`
- **Bluetooth Devices** - `GET /system/bluetooth`
- Improved system stats

### Code Organization
- Split utilities into separate files
- Created `utils/` directory
- Better code separation
- Improved maintainability

---

## 🔄 Changes from v0.0.1.1

### New API Endpoints
```bash
# Player Controls
POST /player/play-pause
POST /player/play
POST /player/pause
POST /player/next
POST /player/previous

# Player Information
GET /player/info
GET /player/list

# System Information
GET /system/wifi
GET /system/bluetooth
```

### Code Structure
```
blitz/
├── main.go
├── mediaInfo.go
├── spotifyArtwork.go
├── websocket.go
├── playerCommands.go        # NEW
└── utils/                   # NEW
    ├── wifiInfo.go
    ├── bluetoothInfo.go
    └── spawnProcesses.go
```

### Improvements
- Added player control commands
- Enhanced error handling
- Better logging
- Improved WebSocket reliability
- More robust player detection

---

## 📋 API Reference

### Player Controls

**Play/Pause Toggle**
```bash
curl -X POST http://localhost:8765/player/play-pause
```

**Play**
```bash
curl -X POST http://localhost:8765/player/play
```

**Pause**
```bash
curl -X POST http://localhost:8765/player/pause
```

**Next Track**
```bash
curl -X POST http://localhost:8765/player/next
```

**Previous Track**
```bash
curl -X POST http://localhost:8765/player/previous
```

### Information Endpoints

**Player Info**
```bash
curl http://localhost:8765/player/info
```

**Active Players**
```bash
curl http://localhost:8765/player/list
```

**WiFi Information**
```bash
curl http://localhost:8765/system/wifi
```

**Bluetooth Devices**
```bash
curl http://localhost:8765/system/bluetooth
```

---

## 🐛 Fixed Issues

1. WebSocket connection stability
2. Player state synchronization
3. Missing player detection
4. Artwork retrieval errors
5. System info parsing issues

---

## 📦 Release Assets

- `blitz_v0.0.1.2` - Linux x86-64 binary
- Size: ~10 MB
- Platform: Linux (kernel 4.4.0+)

---

## 🐛 Known Issues

- No authentication system
- No database integration
- No API versioning
- File naming not following Go conventions (camelCase vs snake_case)
- Utils in separate directory but no proper package structure
- No cross-platform support

---

## ⚠️ Deprecation Notice

The `utils/` directory structure will be replaced with proper Go package structure in the next major release (v0.0.1.3).

---

## 📝 Notes

This release significantly expands the API surface with player controls and system information endpoints. The code organization was improved with the introduction of the `utils/` directory, though proper Go package structure will come in v0.0.1.3.

---

## 🔗 Links

- **Repository:** github.com/codershubinc/Blitz
- **License:** MIT

---

**Contributors:** Swapnil Ingle  
**Project Status:** 🟡 Beta - Feature Enhancement
