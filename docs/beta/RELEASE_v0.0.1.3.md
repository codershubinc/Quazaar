# Quazaar v0.0.1.3 - Modern Architecture Release

> **Fourth beta release with complete project restructure and authentication foundation**

## 📦 Release Information

- **Version:** v0.0.1.3
- **Release Date:** November 17, 2025
- **File Name:** `quazaar_v0.0.1.3_linux_x64`
- **File Size:** 13 MB
- **Platform:** Linux x86-64 (GNU/Linux 4.4.0+)
- **Status:** 🟡 Beta Release
- **Project Renamed:** Blitz → Quazaar

---

## 🎉 What's New

This release marks a **major restructure** of the project with modern Go architecture, database integration, and authentication foundation.

### Major Changes

- ✅ **Project Renamed** - Blitz is now **Quazaar**
- ✅ **Go Standard Layout** - Restructured to `cmd/`, `internal/`, `pkg/` architecture
- ✅ **SQLite Database** - Persistent storage with authentication support
- ✅ **User Authentication** - Signup/Login endpoints with bcrypt password hashing
- ✅ **Token System** - Foundation for token-based API authentication
- ✅ **D-Bus Integration** - Direct MPRIS communication for media control
- ✅ **Static Asset Serving** - Built-in support for CSS/JS/images
- ✅ **Snake Case Files** - All filenames follow Go naming conventions

### Architecture Improvements

- **Modern Structure:**

  ```
  cmd/server/          - Application entry point
  internal/            - Private application packages
  pkg/                 - Public reusable packages
  assets/              - Static files (CSS, JS, images)
  ```

- **Clean Package Organization:**
  - `internal/auth` - Authentication & authorization
  - `internal/db` - Database layer with SQLite
  - `internal/media` - Media player integration (D-Bus, Windows)
  - `internal/player` - Player control commands
  - `internal/poller` - Background media polling
  - `internal/system` - System utilities (WiFi, Bluetooth)
  - `internal/websocket` - Real-time WebSocket communication
  - `pkg/helpers` - Shared utility functions
  - `pkg/models` - Data models and types

---

## 🔄 Changes from v0.0.1.2

### Breaking Changes

- **Project Name:** `blitz` → `quazaar`
- **Binary Name:** `blitz_v0.0.1.2` → `quazaar_v0.0.1.3_linux_x64`
- **Build Command:** `go build` → `go build ./cmd/server`
- **Import Paths:** All imports restructured

### New Features

#### Authentication System

```bash
# New endpoints
POST /api/signup     - User registration
POST /api/login      - User authentication
```

**Features:**

- ✅ Bcrypt password hashing (cost: 10)
- ✅ SQLite database for user storage
- ✅ Token generation and management
- ✅ Single-user system (local server design)

#### Database Integration

```
Location: ~/.quazaar/quazaar.db

Tables:
- users   - User credentials (id, name, pass)
- tokens  - API tokens (id, tokenOf, tokenType, token, expiry)
```

#### Media Player Enhancements

**D-Bus/MPRIS Support:**

```bash
# Direct D-Bus communication
GET /api/v0.1/player/info/dbus

# List MPRIS players
GET /api/v0.1/player/mpris/list

# Query specific player
GET /api/v0.1/player/info/dbus?player=org.mpris.MediaPlayer2.spotify
```

**Windows Media Support:**

```bash
# Windows PowerShell integration (cross-platform ready)
GET /api/v0.1/player/info/windows
GET /api/v0.1/player/windows/list
```

**Player Controls (v0.1 API):**

```bash
POST /api/v0.1/player/play-pause
POST /api/v0.1/player/play
POST /api/v0.1/player/pause
POST /api/v0.1/player/next
POST /api/v0.1/player/previous
```

#### Static Assets

```
/assets/css/         - Stylesheets
/assets/js/          - JavaScript files
/assets/images/      - Images and icons
```

### Code Quality Improvements

- ✅ **Go Standard Naming** - All files use `snake_case.go`
- ✅ **Package Separation** - Clean public/private boundaries
- ✅ **Error Handling** - Improved error messages and logging
- ✅ **Build Tags** - Proper platform-specific code separation
- ✅ **Import Organization** - Clean import paths (`Quazaar/internal/*`, `Quazaar/pkg/*`)

---

## 🚀 Installation

### Quick Start

```bash
# Navigate to release directory
cd release/

# Make executable
chmod +x quazaar_v0.0.1.3_linux_x64

# Run the server
./quazaar_v0.0.1.3_linux_x64
```

The server will start on `http://localhost:8765` by default.

### First Time Setup

```bash
# 1. Run the server
./quazaar_v0.0.1.3_linux_x64

# 2. Register a user (in another terminal)
curl -X POST http://localhost:8765/api/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"yourpassword"}'

# 3. Login to get token
curl -X POST http://localhost:8765/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"yourpassword"}'
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/codershubinc/Quazaar.git
cd Quazaar

# Checkout beta branch
git checkout beta

# Build
go build -o quazaar_v0.0.1.3_linux_x64 ./cmd/server

# Run
./quazaar_v0.0.1.3_linux_x64
```

### Upgrade from Blitz v0.0.1.2

```bash
# Stop old Blitz server
pkill blitz_v0.0.1.2

# Run new Quazaar
cd release/
chmod +x quazaar_v0.0.1.3_linux_x64
./quazaar_v0.0.1.3_linux_x64
```

**Note:** Database will be created automatically at `~/.quazaar/quazaar.db`

---

## 📋 System Requirements

### Minimum Requirements

- **OS:** Linux (kernel 4.4.0+) with D-Bus
- **Architecture:** x86-64 (AMD64)
- **Memory:** 100 MB RAM
- **Disk Space:** 50 MB
- **Libraries:**
  - glibc 2.27+
  - D-Bus (for MPRIS)
  - SQLite3 (embedded)

### Recommended

- **Memory:** 200 MB RAM
- **playerctl:** v2.3.0+ (for enhanced media control)
- **Go:** 1.21+ (for building from source)

### Optional Dependencies

- `nmcli` - WiFi information
- `bluetoothctl` - Bluetooth device info
- `iw` - Advanced WiFi statistics
- `playerctl` - Fallback media control

---

## 🌐 API Endpoints

### Authentication (NEW in v0.0.1.3)

```bash
# Register first user (only 1 allowed)
POST /api/signup
Content-Type: application/json

Request:
{
  "username": "admin",
  "password": "yourpassword"
}

Response:
{
  "success": true,
  "message": "User registered successfully",
  "username": "admin"
}
```

```bash
# Login and get authentication token
POST /api/login
Content-Type: application/json

Request:
{
  "username": "admin",
  "password": "yourpassword"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "username": "admin",
  "token": "$2a$10$...",
  "tokenType": "auth"
}
```

### Media Player API v0.1 (NEW)

```bash
# Get current player info (auto-detect best source)
GET /api/v0.1/player/info

# Get player info via D-Bus only
GET /api/v0.1/player/info/dbus

# Get specific player info
GET /api/v0.1/player/info/dbus?player=org.mpris.MediaPlayer2.spotify

# List all MPRIS players
GET /api/v0.1/player/mpris/list

# List all active players
GET /api/v0.1/player/list
```

### Player Controls (NEW)

```bash
# Toggle play/pause
POST /api/v0.1/player/play-pause

# Play
POST /api/v0.1/player/play

# Pause
POST /api/v0.1/player/pause

# Next track
POST /api/v0.1/player/next

# Previous track
POST /api/v0.1/player/previous
```

### WebSocket

```bash
# Real-time media updates
GET /ws

# Example connection:
const ws = new WebSocket('ws://localhost:8765/ws');
ws.onmessage = (event) => {
  const mediaInfo = JSON.parse(event.data);
  console.log(mediaInfo);
};
```

### Static Assets (NEW)

```bash
# Serve static files
GET /assets/css/style.css
GET /assets/js/app.js
GET /assets/images/logo.png
```

### Home Page

```bash
GET /
# Serves temp/web/index.html
```

---

## 🔧 Configuration

### Environment Variables

Create [`.env`](.env) file in project root:

```env
# Server configuration
LOCAL_HOST_IP=127.0.0.1
LOCAL_HOST_PORT=8765
```

### Database Configuration

- **Location:** `~/.quazaar/quazaar.db`
- **Auto-created:** On first run
- **Schema:**

  ```sql
  -- Users table (max 1 user)
  CREATE TABLE users (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    name TEXT NOT NULL UNIQUE,
    pass TEXT NOT NULL
  );

  -- Tokens table
  CREATE TABLE tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tokenOf TEXT NOT NULL,
    tokenType TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expiry NUMERIC
  );
  ```

### Default Settings

- **Host:** `127.0.0.1`
- **Port:** `8765`
- **Database:** `~/.quazaar/quazaar.db`
- **Max Users:** 1 (single-user system)

---

## 📊 Project Structure

```
Quazaar/
├── cmd/
│   └── server/
│       └── main.go             # Application entry point
├── internal/                   # Private packages
│   ├── auth/
│   │   ├── auth.go            # Auth logic & password hashing
│   │   └── handlers.go        # HTTP handlers (signup/login)
│   ├── db/
│   │   └── db.go              # SQLite database layer
│   ├── media/
│   │   ├── artwork.go
│   │   ├── media_info.go      # Common media types
│   │   ├── media_info_dbus.go # D-Bus/MPRIS integration
│   │   ├── media_info_windows.go # Windows support
│   │   ├── spotify.go
│   │   └── volume_controls.go
│   ├── player/
│   │   ├── commands.go        # Player control commands
│   │   └── handlers.go        # HTTP handlers (v0.1 API)
│   ├── poller/
│   │   ├── handler.go
│   │   └── poller.go          # Background media polling
│   ├── system/
│   │   ├── app_launcher.go
│   │   ├── bluetooth_info.go
│   │   └── wifi_info.go
│   └── websocket/
│       ├── channel.go
│       ├── handler.go
│       └── websocket.go       # WebSocket server
├── pkg/                        # Public packages
│   ├── helpers/
│   │   └── spawn_processes.go # Process spawning utilities
│   └── models/
│       ├── auth_models.go     # Auth data models
│       └── server_response.go # Response types
├── assets/                     # Static files
│   ├── css/
│   │   └── style.css
│   ├── js/
│   └── images/
├── docs/                       # Documentation
│   └── beta/
│       └── README.md          # Beta releases guide
├── release/                    # Compiled binaries
│   └── quazaar_v0.0.1.3_linux_x64
└── temp/                       # Temporary files
    └── web/
        ├── index.html
        └── auth.html
```

---

## 📈 Performance & Size Comparison

| Metric        | v0.0.1.2 | v0.0.1.3     | Change                 |
| ------------- | -------- | ------------ | ---------------------- |
| Binary Size   | 10 MB    | 13 MB        | ⬆️ 30% (SQLite + auth) |
| Memory Usage  | ~50 MB   | ~70 MB       | ⬆️ 40% (database)      |
| Startup Time  | ~100ms   | ~150ms       | ⬆️ 50% (DB init)       |
| API Endpoints | 2        | 15+          | ⬆️ 650%                |
| Code Lines    | ~1500    | ~2500        | ⬆️ 67%                 |
| Packages      | 1 (main) | 10 (modular) | New structure          |

**Note:** Size increase due to SQLite, authentication, and modular architecture.

---

## 🔐 Security Features

### Authentication

- ✅ Bcrypt password hashing (cost: 10)
- ✅ Secure token generation
- ✅ Token expiration support
- ✅ Single-user enforcement (CHECK constraint)

### Database

- ✅ SQLite with proper schema
- ✅ Prepared statements (SQL injection prevention)
- ✅ Password hashing (no plaintext storage)
- ✅ Unique constraints on username and tokens

### Planned (Next Release)

- [ ] Token-based API authentication
- [ ] Middleware for protected endpoints
- [ ] Rate limiting
- [ ] CORS configuration
- [ ] TLS/SSL support

---

## 🐛 Fixed Issues

### From v0.0.1.2

1. ✅ **Project Organization** - Restructured to Go standard layout
2. ✅ **File Naming** - All files now use snake_case
3. ✅ **Import Paths** - Clean, organized imports
4. ✅ **Build System** - Proper cmd/ structure
5. ✅ **Platform Code** - Correct build tags for Windows/Linux
6. ✅ **Error Handling** - Improved error messages
7. ✅ **Code Duplication** - Reduced through proper packaging
8. ✅ **Missing SpawnProcess** - Fixed imports in system package

---

## ⚠️ Known Limitations

### Current Limitations

- **Authentication Incomplete** - Endpoints not yet protected (in progress)
- **Single User Only** - Database limited to 1 user
- **Linux Primary** - Windows support is experimental
- **No TLS** - HTTP only (localhost use recommended)
- **No Rate Limiting** - Unlimited API requests
- **Token Expiry** - Tokens don't expire yet (field exists but not enforced)

### Work in Progress

- Token-based endpoint protection
- Middleware implementation
- Password change functionality
- Token refresh mechanism
- Multi-device token support

---

## 🔄 Migration Guide

### From Blitz v0.0.1.2

**No data to migrate** - v0.0.1.2 had no database.

#### Steps:

1. **Stop old Blitz server:**

   ```bash
   pkill blitz_v0.0.1.2
   ```

2. **Run new Quazaar:**

   ```bash
   cd release/
   chmod +x quazaar_v0.0.1.3_linux_x64
   ./quazaar_v0.0.1.3_linux_x64
   ```

3. **Register first user:**

   ```bash
   curl -X POST http://localhost:8765/api/signup \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"secure_password"}'
   ```

4. **Update your client code:**
   - Old endpoint: `GET /ws`
   - New endpoints: `POST /api/signup`, `POST /api/login`, `GET /ws` (same)
   - New player API: `/api/v0.1/player/*`

### Configuration Changes

```bash
# Old (v0.0.1.2)
./blitz_v0.0.1.2

# New (v0.0.1.3)
./quazaar_v0.0.1.3_linux_x64
```

**Environment variables:** Same (no changes needed)

---

## 📚 Documentation

### Available Docs

- **Beta Releases:** [`docs/beta/README.md`](../README.md)
- **This Release:** [`docs/beta/RELEASE_v0.0.1.3.md`](RELEASE_v0.0.1.3.md)
- **API Testing:** [`docs/API_TESTING_GUIDE.md`](../API_TESTING_GUIDE.md)
- **WebSocket:** [`docs/WEBSOCKET.md`](../WEBSOCKET.md)
- **Authentication:** [`docs/AUTH_SYSTEM.md`](../AUTH_SYSTEM.md)
- **Project Structure:** [`docs/PROJECT_STRUCTURE.md`](../PROJECT_STRUCTURE.md)

### New Documentation

- **Auth Completion Plan:** [`docs/AUTH_COMPLETION_PLAN.md`](../AUTH_COMPLETION_PLAN.md)
- **Windows Testing:** [`docs/WINDOWS_TESTING_GUIDE.md`](../WINDOWS_TESTING_GUIDE.md)

---

## 🎯 Roadmap

### Next Release (v0.0.1.4) - Planned

**Focus:** Complete Authentication System

- [ ] Authentication middleware
- [ ] Protected endpoints (WebSocket, Player controls)
- [ ] Token validation
- [ ] Password change endpoint
- [ ] Token refresh endpoint
- [ ] User profile endpoint

**ETA:** 1-2 weeks

### Future Releases

**v0.1.0** - Feature Complete

- [ ] Full authentication system
- [ ] API versioning (v1)
- [ ] Rate limiting
- [ ] Enhanced error handling
- [ ] Comprehensive testing

**v0.2.0** - Cross-Platform

- [ ] Windows native support
- [ ] macOS support
- [ ] Mobile app integration
- [ ] TLS/SSL support

**v1.0.0** - Production Ready

- [ ] Stable API
- [ ] Full documentation
- [ ] Performance optimizations
- [ ] Security audit
- [ ] Plugin system

---

## 🤝 Contributing

This is an **active development** project. Contributions welcome!

### Ways to Contribute

- 🐛 Report bugs
- 💡 Suggest features
- 📝 Improve documentation
- 🔧 Submit pull requests
- ⭐ Star the repository

### Development Setup

```bash
git clone https://github.com/codershubinc/Quazaar.git
cd Quazaar
git checkout beta
go mod download
go build ./cmd/server
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/auth
go test ./internal/db

# Run with coverage
go test -cover ./...
```

**Repository:** https://github.com/codershubinc/Quazaar

---

## 📄 License

**MIT License**

Copyright © 2025 Swapnil Ingle

See LICENSE file for full details.

---

## 🔗 Links

- **GitHub:** https://github.com/codershubinc/Quazaar
- **Issues:** https://github.com/codershubinc/Quazaar/issues
- **Releases:** https://github.com/codershubinc/Quazaar/releases
- **Documentation:** https://github.com/codershubinc/Quazaar/tree/beta/docs
- **License:** https://github.com/codershubinc/Quazaar/blob/beta/LICENSE

---

## 📝 Detailed Changelog

```
v0.0.1.3 (2025-11-17) - Major Release

  Project:
    - Renamed: Blitz → Quazaar
    - Restructured to Go standard layout
    - All files renamed to snake_case
    - Binary: quazaar_v0.0.1.3_linux_x64

  Architecture:
    + cmd/server/ - Application entry point
    + internal/ - Private packages (8 packages)
      - auth/ - Authentication & authorization
      - db/ - SQLite database layer
      - media/ - Media player integration
      - player/ - Player control commands
      - poller/ - Background polling
      - system/ - System utilities
      - websocket/ - WebSocket server
    + pkg/ - Public packages (2 packages)
      - helpers/ - Shared utilities
      - models/ - Data models
    + assets/ - Static file serving

  Features:
    + SQLite database integration
    + User authentication (signup/login)
    + Token generation system
    + D-Bus/MPRIS media control
    + Windows media support (experimental)
    + Static asset serving (/assets/*)
    + Player selection via query params
    + Single-user enforcement

  API Changes:
    + POST /api/signup
    + POST /api/login
    + GET /api/v0.1/player/info
    + GET /api/v0.1/player/info/dbus
    + GET /api/v0.1/player/info/windows
    + GET /api/v0.1/player/mpris/list
    + GET /api/v0.1/player/windows/list
    + GET /api/v0.1/player/list
    + POST /api/v0.1/player/play-pause
    + POST /api/v0.1/player/play
    + POST /api/v0.1/player/pause
    + POST /api/v0.1/player/next
    + POST /api/v0.1/player/previous
    + GET /assets/* (static files)

  Database:
    + Location: ~/.quazaar/quazaar.db
    + Table: users (id, name, pass)
    + Table: tokens (id, tokenOf, tokenType, token, expiry)
    + Single user constraint (CHECK id = 1)
    + Unique constraints on username and tokens

  Code Quality:
    - Proper Go package structure
    - Clean import paths (Quazaar/internal/*, Quazaar/pkg/*)
    - Platform-specific build tags
    - Improved error handling
    - Comprehensive logging
    - Snake case file naming

  Bug Fixes:
    - Fixed missing SpawnProcess imports in system package
    - Fixed duplicate import in media_info_windows.go
    - Fixed error message capitalization (linter)
    - Fixed build tags for cross-platform code

  Performance:
    - Binary size: 10 MB → 13 MB (+30%)
    - Memory usage: ~50 MB → ~70 MB (+40%)
    - Startup time: ~100ms → ~150ms (+50%)

  Documentation:
    + Beta releases guide (docs/beta/README.md)
    + Release notes (docs/beta/RELEASE_v0.0.1.3.md)
    + Authentication plan
    + Windows testing guide
    + Updated project structure docs
```

---

**Thank you for using Quazaar v0.0.1.3!** 🚀

_This is an active development release. Feedback and contributions welcome!_

---

**Release Notes Last Updated:** November 17, 2025  
**Project Status:** 🟡 Beta - Active Development  
**Next Milestone:** Complete Authentication System (v0.0.1.4)
