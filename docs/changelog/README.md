# Quazaar Changelog

Complete version history and release notes for all Quazaar releases.

---

## Version Index

- [v0.0.1.3](#v0013---november-17-2025) - Current Release (Major Restructure)
- [v0.0.1.2](#v0012---november-16-2025) - Feature Enhancement
- [v0.0.1.1](#v0011---mid-november-2025) - Stability Update
- [v0.0.1.0](#v0010---early-november-2025) - Initial Release

---

## [v0.0.1.3] - November 17, 2025

**Project Renamed:** Blitz → Quazaar  
**Status:** 🟡 Beta - Major Restructure

### Highlights

- Complete project restructure to Go standard layout
- SQLite database integration
- Authentication system foundation
- API versioning (v0.1)
- 7 customizable startup banners
- Centralized routing
- Cross-platform support improvements

### Major Changes

- Restructured to `cmd/`, `internal/`, `pkg/` architecture
- All files renamed to snake_case
- New database at `~/.quazaar/quazaar.db`
- User authentication with bcrypt
- Token generation system
- D-Bus/MPRIS direct integration
- Static asset serving

### New Endpoints

```
POST /api/signup
POST /api/login
GET /api/v0.1/player/info/dbus
GET /api/v0.1/player/info/windows
GET /api/v0.1/player/mpris/list
GET /api/v0.1/player/windows/list
GET /api/v0.1/system/wifi
GET /api/v0.1/system/bluetooth
GET /assets/*
```

[Full Changelog](./CHANGELOG_v0.0.1.3.md)

---

## [v0.0.1.2] - November 16, 2025

**Project Name:** Blitz  
**Status:** 🟡 Beta - Feature Enhancement

### Highlights

- Player control commands
- System information APIs
- Enhanced media handling
- Utils directory organization

### New Features

- Player controls (play/pause/next/previous)
- WiFi information endpoint
- Bluetooth devices endpoint
- Active players list
- Better code organization with utils/

### New Endpoints

```
POST /player/play-pause
POST /player/play
POST /player/pause
POST /player/next
POST /player/previous
GET /player/info
GET /player/list
GET /system/wifi
GET /system/bluetooth
```

[Full Changelog](./CHANGELOG_v0.0.1.2.md)

---

## [v0.0.1.1] - Mid November 2025

**Project Name:** Blitz  
**Status:** 🟡 Beta - Stability Update

### Highlights

- Improved media detection
- Enhanced WebSocket stability
- Better error handling
- Code cleanup and refactoring

### Improvements

- Better playerctl integration
- More reliable WebSocket updates
- Enhanced metadata parsing
- Improved artwork retrieval
- Better handling of edge cases

### Bug Fixes

- WebSocket disconnection issues
- Player detection edge cases
- Missing metadata handling
- Artwork retrieval failures

[Full Changelog](./CHANGELOG_v0.0.1.1.md)

---

## [v0.0.1.0] - Early November 2025

**Project Name:** Blitz  
**Status:** 🟢 Beta - Initial Release

### Highlights

- First beta release
- WebSocket server implementation
- Basic media player integration
- System information retrieval

### Core Features

- Real-time WebSocket communication
- Media player integration via playerctl
- WiFi information
- Bluetooth device enumeration
- Basic HTTP server

### Endpoints

```
GET /ws
GET /
```

[Full Changelog](./CHANGELOG_v0.0.1.0.md)

---

## Version Timeline

```
v0.0.1.0 (Initial)
    │
    ├─ Basic WebSocket
    ├─ Playerctl integration
    └─ Simple architecture

v0.0.1.1 (Stability)
    │
    ├─ Better media detection
    ├─ WebSocket improvements
    └─ Bug fixes

v0.0.1.2 (Features)
    │
    ├─ Player controls
    ├─ System APIs
    ├─ Utils organization
    └─ More endpoints

v0.0.1.3 (Architecture) ← Current
    │
    ├─ Project renamed (Blitz → Quazaar)
    ├─ Go standard layout
    ├─ Database integration
    ├─ Authentication system
    ├─ API versioning
    ├─ Startup banners
    └─ Centralized routing
```

---

## Statistics

### Release Progression

| Version  | Binary Size | Memory | Endpoints | Packages | LOC   |
| -------- | ----------- | ------ | --------- | -------- | ----- |
| v0.0.1.0 | ~8 MB       | ~40 MB | 2         | 1        | ~1200 |
| v0.0.1.1 | ~8.5 MB     | ~45 MB | 2         | 1        | ~1400 |
| v0.0.1.2 | ~10 MB      | ~50 MB | 8         | 1+utils  | ~1800 |
| v0.0.1.3 | ~13 MB      | ~70 MB | 18+       | 10       | ~2500 |

### Feature Growth

| Feature           | v0.0.1.0 | v0.0.1.1 | v0.0.1.2 | v0.0.1.3 |
| ----------------- | -------- | -------- | -------- | -------- |
| WebSocket         | ✅       | ✅       | ✅       | ✅       |
| Media Info        | ✅       | ✅       | ✅       | ✅       |
| Player Controls   | ❌       | ❌       | ✅       | ✅       |
| System Info       | ⚠️       | ⚠️       | ✅       | ✅       |
| Authentication    | ❌       | ❌       | ❌       | ✅       |
| Database          | ❌       | ❌       | ❌       | ✅       |
| API Versioning    | ❌       | ❌       | ❌       | ✅       |
| Startup Banners   | ❌       | ❌       | ❌       | ✅       |
| Static Assets     | ❌       | ❌       | ❌       | ✅       |
| D-Bus Integration | ⚠️       | ⚠️       | ⚠️       | ✅       |
| Cross-Platform    | ❌       | ❌       | ❌       | ⚠️       |

✅ Fully Implemented | ⚠️ Partial | ❌ Not Available

---

## Breaking Changes

### v0.0.1.3

- Project renamed: `blitz` → `quazaar`
- Binary renamed: `blitz_v0.0.1.2` → `quazaar_v0.0.1.3_linux_x64`
- API endpoints versioned: `/player/*` → `/api/v0.1/player/*`
- Build command changed: `go build` → `go build ./cmd/server`
- Import paths restructured

### v0.0.1.2

- No breaking changes (backward compatible)

### v0.0.1.1

- No breaking changes (backward compatible)

---

## Migration Guides

### From v0.0.1.2 to v0.0.1.3

**API Changes:**

```bash
# Old                           # New
POST /player/play       →       POST /api/v0.1/player/play
GET /player/info        →       GET /api/v0.1/player/info
GET /system/wifi        →       GET /api/v0.1/system/wifi
```

**Binary Changes:**

```bash
# Old
./blitz_v0.0.1.2

# New
./quazaar_v0.0.1.3_linux_x64
```

**New Features to Consider:**

- User authentication (signup/login)
- Database storage
- Startup banners
- Static asset serving

### From v0.0.1.1 to v0.0.1.2

**No breaking changes** - Backward compatible

- New endpoints added
- Enhanced functionality

### From v0.0.1.0 to v0.0.1.1

**No breaking changes** - Backward compatible

- Stability improvements
- Bug fixes

---

## Deprecation Notices

### v0.0.1.3

- Old endpoint paths without `/api/v0.1/` prefix are **deprecated**
- Will be removed in v0.1.0

### v0.0.1.2

- `utils/` directory structure deprecated
- Replaced with proper Go packages in v0.0.1.3

---

## Known Issues by Version

### v0.0.1.3

- Authentication endpoints not yet protected
- Token expiry not enforced
- Windows support experimental
- No rate limiting
- No TLS support

### v0.0.1.2

- No authentication
- No database
- File naming not following Go conventions
- No API versioning

### v0.0.1.1

- Limited API endpoints
- Single file architecture
- No authentication

### v0.0.1.0

- Basic functionality only
- Limited error handling
- No API versioning
- Platform-specific code not separated

---

## Roadmap

### v0.0.1.4 (Planned)

- Complete authentication middleware
- Protected endpoints
- Token validation
- Password change functionality
- Token refresh mechanism

### v0.1.0 (Future)

- Stable API v1
- Full authentication system
- Rate limiting
- Enhanced error handling
- Comprehensive testing

### v0.2.0 (Future)

- Windows native support
- macOS support
- Mobile app integration
- TLS/SSL support

### v1.0.0 (Future)

- Production ready
- Stable API
- Full documentation
- Security audit
- Plugin system

---

## License

**MIT License**

Copyright © 2025 Swapnil Ingle

All versions are licensed under the MIT License.

---

## Links

- **Repository:** https://github.com/codershubinc/Quazaar
- **Issues:** https://github.com/codershubinc/Quazaar/issues
- **Releases:** https://github.com/codershubinc/Quazaar/releases
- **Documentation:** https://github.com/codershubinc/Quazaar/tree/beta/docs

---

**Last Updated:** November 17, 2025  
**Current Version:** v0.0.1.3  
**Project Status:** 🟡 Beta - Active Development
