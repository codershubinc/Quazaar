# Changelog - v0.0.1.3

**Release Date:** November 17, 2025  
**Project Name:** Quazaar (renamed from Blitz)  
**Status:** 🟡 Beta - Major Restructure

---

## Overview

Fourth beta release with complete project restructure following Go standard layout. Major architectural changes including database integration, authentication system foundation, and modern package organization.

---

## 🎉 Major Changes

### Project Renamed
- **Old Name:** Blitz
- **New Name:** Quazaar
- **Binary:** `blitz_v0.0.1.2` → `quazaar_v0.0.1.3_linux_x64`

### Architecture Overhaul
- **Go Standard Layout** - Restructured to `cmd/`, `internal/`, `pkg/`
- **Package Organization** - 10 well-organized packages
- **File Naming** - All files renamed to `snake_case.go`
- **Import Paths** - Clean, organized imports

### New Features
- ✅ **SQLite Database** - Persistent storage at `~/.quazaar/quazaar.db`
- ✅ **User Authentication** - Signup/Login with bcrypt hashing
- ✅ **Token System** - Foundation for API authentication
- ✅ **D-Bus Integration** - Direct MPRIS communication
- ✅ **Static Assets** - Built-in CSS/JS/images serving
- ✅ **API Versioning** - v0.1 API endpoints
- ✅ **Startup Banner** - 7 customizable banner variants

---

## 📦 Project Structure

### New Architecture
```
Quazaar/
├── cmd/
│   └── server/
│       └── main.go             # Entry point
├── internal/                   # Private packages
│   ├── api/
│   │   └── router.go           # Centralized routing
│   ├── auth/
│   │   ├── auth.go            # Auth logic
│   │   └── handlers.go        # Signup/Login
│   ├── banner/
│   │   └── banner.go          # Startup banners (7 variants)
│   ├── db/
│   │   └── db.go              # SQLite layer
│   ├── media/
│   │   ├── artwork.go
│   │   ├── media_info.go
│   │   ├── media_info_dbus.go
│   │   ├── media_info_windows.go
│   │   ├── spotify.go
│   │   └── volume_controls.go
│   ├── player/
│   │   ├── commands.go
│   │   └── handlers.go
│   ├── poller/
│   │   ├── handler.go
│   │   └── poller.go
│   ├── system/
│   │   ├── app_launcher.go
│   │   ├── bluetooth_info.go
│   │   ├── handlers.go        # WiFi/Bluetooth endpoints
│   │   └── wifi_info.go
│   └── websocket/
│       ├── channel.go
│       ├── handler.go
│       └── websocket.go
├── pkg/                        # Public packages
│   ├── helpers/
│   │   └── spawn_processes.go
│   └── models/
│       ├── auth_models.go
│       └── server_response.go
├── assets/                     # Static files
│   ├── css/
│   │   └── style.css
│   ├── js/
│   └── images/
└── docs/
    ├── changelog/              # NEW
    │   ├── CHANGELOG_v0.0.1.0.md
    │   ├── CHANGELOG_v0.0.1.1.md
    │   ├── CHANGELOG_v0.0.1.2.md
    │   └── CHANGELOG_v0.0.1.3.md
    └── beta/
        └── RELEASE_v0.0.1.3.md
```

---

## ✨ New Features

### Database Integration

**Location:** `~/.quazaar/quazaar.db`

**Schema:**
```sql
-- Users table (single user system)
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

### Authentication System

**Endpoints:**
```bash
# Register user (only 1 allowed)
POST /api/signup
{
  "username": "admin",
  "password": "yourpassword"
}

# Login and get token
POST /api/login
{
  "username": "admin",
  "password": "yourpassword"
}
```

**Features:**
- Bcrypt password hashing (cost: 10)
- Token generation
- Single-user enforcement
- Database-backed authentication

### Enhanced Media API (v0.1)

**New Endpoints:**
```bash
# Player information
GET /api/v0.1/player/info
GET /api/v0.1/player/info/dbus
GET /api/v0.1/player/info/dbus?player=org.mpris.MediaPlayer2.spotify
GET /api/v0.1/player/info/windows

# Player lists
GET /api/v0.1/player/list
GET /api/v0.1/player/mpris/list
GET /api/v0.1/player/windows/list

# Player controls
POST /api/v0.1/player/play-pause
POST /api/v0.1/player/play
POST /api/v0.1/player/pause
POST /api/v0.1/player/next
POST /api/v0.1/player/previous
```

### System Information API

**New Endpoints:**
```bash
GET /api/v0.1/system/wifi
GET /api/v0.1/system/bluetooth
```

### Static Asset Serving

**Endpoints:**
```bash
GET /assets/css/style.css
GET /assets/js/app.js
GET /assets/images/logo.png
```

### Startup Banners

**7 Customizable Variants:**
- Variant1: Modern detailed (default)
- Variant2: Gradient with stats
- Variant3: Bold block letters
- Variant4: Compact modern
- Variant5: Double border
- Variant6: Minimal elegant
- Variant7: Retro ASCII

**Usage:**
```go
banner.Show()        // Default (Variant1)
banner.Variant2()    // Gradient style
banner.Variant3()    // Bold letters
// ... etc
```

### Centralized Routing

**New Package:** `internal/api/router.go`

All routes registered in one place:
```go
api.SetupRoutes()  // Registers all 18+ endpoints
```

---

## 🔄 API Changes

### Breaking Changes

**Endpoints Moved to v0.1:**
```bash
# Old (v0.0.1.2)        # New (v0.0.1.3)
POST /player/play      → POST /api/v0.1/player/play
GET /player/info       → GET /api/v0.1/player/info
GET /system/wifi       → GET /api/v0.1/system/wifi
```

**New Endpoints:**
```bash
POST /api/signup                              # User registration
POST /api/login                               # User login
GET /api/v0.1/player/info/dbus               # D-Bus player info
GET /api/v0.1/player/info/windows            # Windows player info
GET /api/v0.1/player/mpris/list              # List MPRIS players
GET /api/v0.1/player/windows/list            # List Windows players
GET /assets/*                                 # Static files
```

---

## 🔧 Technical Improvements

### Code Quality
- ✅ Go standard project layout
- ✅ Snake case file naming (`media_info.go`)
- ✅ Clean import paths (`Quazaar/internal/*`)
- ✅ Platform-specific build tags
- ✅ Proper error handling
- ✅ Comprehensive logging

### Platform Support
- ✅ Linux D-Bus/MPRIS integration
- ✅ Windows PowerShell/SMTC support (experimental)
- ✅ Build tags for platform-specific code
- ✅ Stub functions for unsupported platforms

### Performance
| Metric        | v0.0.1.2 | v0.0.1.3 | Change                |
|---------------|----------|----------|-----------------------|
| Binary Size   | 10 MB    | 13 MB    | +30% (SQLite + auth)  |
| Memory Usage  | ~50 MB   | ~70 MB   | +40% (database)       |
| Startup Time  | ~100ms   | ~150ms   | +50% (DB init)        |
| API Endpoints | 8        | 18+      | +125%                 |
| Packages      | 1        | 10       | Modular architecture  |

---

## 🐛 Bug Fixes

1. ✅ Fixed missing `SpawnProcess` imports in system package
2. ✅ Fixed duplicate imports in media files
3. ✅ Fixed Windows build tags
4. ✅ Fixed error message capitalization (linter)
5. ✅ Fixed route registration organization
6. ✅ Fixed file naming conventions

---

## 🔐 Security Features

### Current
- ✅ Bcrypt password hashing (cost: 10)
- ✅ Unique username constraint
- ✅ Unique token constraint
- ✅ Single-user enforcement (CHECK constraint)
- ✅ No plaintext password storage

### Planned (v0.0.1.4)
- [ ] Token-based endpoint protection
- [ ] Authentication middleware
- [ ] Token expiration enforcement
- [ ] Rate limiting
- [ ] CORS configuration

---

## 📚 Documentation

### New Documentation
- `docs/changelog/` - Version changelogs
- `docs/beta/RELEASE_v0.0.1.3.md` - Complete release notes
- Updated project structure documentation
- API testing guides
- Authentication system documentation

---

## ⚠️ Known Limitations

- **Authentication Incomplete** - Endpoints not yet protected
- **Single User Only** - Database limited to 1 user
- **No Token Expiry** - Tokens don't expire yet
- **No Rate Limiting** - Unlimited API requests
- **No TLS** - HTTP only (localhost recommended)
- **Windows Experimental** - Primary platform is Linux

---

## 🚀 Migration Guide

### From v0.0.1.2

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

4. **Update API calls:**
   - Add `/api/v0.1/` prefix to all player and system endpoints
   - Use new authentication endpoints if needed

---

## 📦 Build Information

### Binary Details
- **Filename:** `quazaar_v0.0.1.3_linux_x64`
- **Size:** 13 MB
- **Platform:** Linux x86-64
- **Go Version:** 1.21+
- **Build Command:** `go build ./cmd/server`

### Dependencies
- SQLite3 (embedded)
- D-Bus (for MPRIS)
- bcrypt (for password hashing)
- gorilla/websocket
- godotenv

---

## 🎯 What's Next (v0.0.1.4)

### Planned Features
- [ ] Complete authentication middleware
- [ ] Protected endpoints
- [ ] Token validation
- [ ] Password change endpoint
- [ ] Token refresh mechanism
- [ ] User profile endpoint
- [ ] Enhanced error handling
- [ ] Rate limiting

**ETA:** 1-2 weeks

---

## 🤝 Contributing

Repository: https://github.com/codershubinc/Quazaar

### Development Setup
```bash
git clone https://github.com/codershubinc/Quazaar.git
cd Quazaar
git checkout beta
go mod download
go build ./cmd/server
```

---

## 📄 License

**MIT License**

Copyright © 2025 Swapnil Ingle

---

## 🔗 Links

- **Repository:** https://github.com/codershubinc/Quazaar
- **Issues:** https://github.com/codershubinc/Quazaar/issues
- **Documentation:** https://github.com/codershubinc/Quazaar/tree/beta/docs
- **License:** https://github.com/codershubinc/Quazaar/blob/beta/LICENSE

---

**Contributors:** Swapnil Ingle  
**Project Status:** 🟡 Beta - Major Restructure  
**Release Date:** November 17, 2025
