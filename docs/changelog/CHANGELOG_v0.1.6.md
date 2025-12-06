# Changelog v0.1.6 - Patch Release

**Date:** December 6, 2025

## 🚀 New Features

### System Controls

- **Volume Control**: Added full system volume control via WebSocket.
  - API: `{"type": "system", "msg_of": "volume", "action": "..."}` (inc, dec, set, mute).
  - Integration: Uses `pactl` for Linux audio control.
- **Brightness Control**: Added screen brightness control.
  - API: `{"type": "system", "msg_of": "brightness", "action": "..."}` (inc, dec, set).
  - Integration: Uses `brightnessctl` for hardware backlight control.

### Web Client Overhaul

- **Modular Architecture**: Refactored the monolithic `script.js` into ES6 modules (`main.js`, `websocket.js`, `ui.js`, `state.js`).
- **New UI Components**: Added interactive sliders for Volume and Brightness in the web dashboard.
- **Authentication UI**: Implemented a dedicated Login Page (`login.html`) with JWT token handling.
- **Dynamic Connection**: WebSocket connection now uses dynamic authentication tokens instead of hardcoded values.

## 🐛 Bug Fixes

- **Network Accessibility**: Fixed the server bind address. It now defaults to `0.0.0.0:8765` (was `127.0.0.1`), allowing access from other devices on the LAN.
- **Go Module Version**: Downgraded `go.mod` version from `1.25.3` (invalid) to `1.23.0` (stable) to ensure build compatibility.
- **Static File Serving**: Fixed 404 errors for web assets by correctly routing `/web/` to `statics/web/`.
- **Command Protocol**: Reverted JSON command payloads to match the legacy protocol expected by the backend (e.g., `{ command: 'play' }`).

## 🛠️ Improvements

- **Project Structure**: Moved web frontend assets from `temp/web` to a structured `statics/web` directory.
- **Audit & Planning**: Established a comprehensive architectural audit plan (`todo.md`) for future stability and performance improvements.

## 🔒 Security

- **Login Flow**: Enforced authentication on the web client. Users are now redirected to the login page if no valid token is found.
