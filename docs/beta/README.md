# Quazaar Beta Releases - Linux x64

> **Note:** These are beta releases for Linux systems only. Windows releases are distributed separately.

## 📦 Available Releases

### Latest Beta - Quazaar (Current Development)

- **File:** `quazaar`
- **Version:** Development Build
- **Size:** 13 MB
- **Date:** November 16, 2025
- **Platform:** Linux x86-64
- **Status:** 🟢 Active Development

**Features:**

- Go standard project layout (cmd/, internal/, pkg/)
- SQLite database with authentication
- WebSocket support for real-time updates
- Media player integration (D-Bus/MPRIS)
- System information (WiFi, Bluetooth)
- Static asset serving
- Token-based authentication (in progress)

---

## Blitz is renamed to Quazaar. Previous Blitz beta builds are listed below for reference.

### Blitz v0.0.1.2

- **File:** `blitz_v0.0.1.2`
- **Version:** 0.0.1.2
- **Size:** 10 MB
- **Date:** November 13, 2025
- **Platform:** Linux x86-64
- **Status:** 🟡 Beta

**Changes:**

- Improved media polling
- Enhanced WebSocket stability
- Bug fixes and optimizations

---

### Blitz v0.0.1.1

- **File:** `blitz_v0.0.1.1`
- **Version:** 0.0.1.1
- **Size:** 8.8 MB
- **Date:** November 3, 2025
- **Platform:** Linux x86-64
- **Status:** 🟡 Beta

**Changes:**

- WebSocket improvements
- Media info enhancements
- System utilities added

---

### Blitz v0.0.1.0

- **File:** `blitz_v0.0.1.0`
- **Version:** 0.0.1.0 (Initial Beta)
- **Size:** 8.8 MB
- **Date:** November 2, 2025
- **Platform:** Linux x86-64
- **Status:** 🟡 Beta

**Features:**

- Initial beta release
- Basic media player control
- WebSocket communication
- System information retrieval

---

### Additional Builds

#### Blitz (Latest Stable Build)

- **File:** `blitz`
- **Size:** 9.9 MB
- **Date:** November 14, 2025
- **Status:** 🟢 Stable

#### Blitz Backup

- **File:** `blitz_backup`
- **Size:** 10 MB
- **Date:** November 14, 2025
- **Status:** 📦 Backup Build

---

## 🚀 Installation

### Download

```bash
# Navigate to release directory
cd release/

# Make executable
chmod +x quazaar  # or blitz_v0.0.1.x

# Run
./quazaar
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/codershubinc/Blitz.git
cd Blitz

# Checkout beta branch
git checkout beta

# Build
go build -o quazaar ./cmd/server

# Run
./quazaar
```

---

## 📋 Requirements

### System Requirements

- **OS:** Linux (any distribution with D-Bus support)
- **Architecture:** x86-64 (AMD64)
- **Libraries:**
  - glibc 2.27+ (for GNU/Linux 4.4.0+)
  - D-Bus (for media player integration)
  - SQLite3

### Optional Dependencies

- `nmcli` - WiFi information
- `bluetoothctl` - Bluetooth device info
- `iw` - Advanced WiFi statistics

---

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
LOCAL_HOST_IP=127.0.0.1
LOCAL_HOST_PORT=8765
```

### Database Location

Default: `~/.quazaar/quazaar.db`

---

## 🌐 API Endpoints

### Authentication

- `POST /api/signup` - Create new user
- `POST /api/login` - User login

### WebSocket

- `GET /ws` - WebSocket connection for real-time updates

### Static Assets

- `GET /assets/css/*` - CSS files
- `GET /assets/js/*` - JavaScript files
- `GET /assets/images/*` - Images

---

## 📝 Version History

| Version        | Date         | Size   | Key Features                  |
| -------------- | ------------ | ------ | ----------------------------- |
| Quazaar (dev)  | Nov 16, 2025 | 13 MB  | Modern structure, auth system |
| Blitz v0.0.1.2 | Nov 13, 2025 | 10 MB  | Stability improvements        |
| Blitz v0.0.1.1 | Nov 3, 2025  | 8.8 MB | WebSocket enhancements        |
| Blitz v0.0.1.0 | Nov 2, 2025  | 8.8 MB | Initial beta release          |

---

## ⚠️ Known Limitations

- **Linux Only:** These builds are compiled for Linux x86-64
- **Beta Software:** May contain bugs and incomplete features
- **No Windows Support:** Windows builds distributed separately
- **Authentication:** Token-based auth is currently in development

---

## 🐛 Bug Reports

Report issues on GitHub:

- **Repository:** [codershubinc/Quazaar](https://github.com/codershubinc/Quazaar)
- **Branch:** beta
- **Issues:** Create an issue with detailed description

---

## 📄 License

MIT License - See LICENSE file for details

---

## 🔗 Links

- **Main Repository:** https://github.com/codershubinc/Quazaar
- **Documentation:** `/docs/`
- **Project Structure:** `/docs/PROJECT_STRUCTURE.md`

---

**Last Updated:** November 17, 2025
