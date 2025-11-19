# ⚡ Quazaar

Real-time music player integration.
With WebSocket remote control for Linux.

`As it in developmental phase. Expect bugs and breaking changes. Also refer the beta branch for active changes`

---

`Android app update , Beta version is out check `[Quazaar App Download](https://github.com/codershubinc/QuazaarApp/releases/download/v0.1.0-beta/Quazaar_v0.1.0-beta.apk) `

`for using the app with server refer 
`[docs/beta/README.md](docs/beta/README.md)

---

## 🎯 Features

- **Remote Command Execution**: Control your PC from any device on your network
- **Real-time Music Display**: Shows currently playing track with album artwork using `playerctl`
- **WebSocket Communication**: Fast, bidirectional communication between devices
- **Secure Command Allowlist**: Only pre-approved commands can be executed
- **Modern Web Interface**: Clean, responsive UI that works on desktop and mobile
- **Auto-updating Music Info**: Track information refreshes every 1 seconds

## 🚀 Quick Start

### Prerequisites

- Go 1.16 or higher
- `playerctl` (for music integration)

  ```bash
  # Arch Linux
  sudo pacman -S playerctl

  # Ubuntu/Debian
  sudo apt install playerctl
  ```

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/codershubinc/Quazaar.git
   cd Quazaar
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Build the server**

   ```bash
   go build -o quazaar
   ```

4. **Run the server**
   ```bash
   ./quazaar
   ```

The server will start on `ws://0.0.0.0:8765/ws`

### Usage

1. **Start the server** on your PC

   ```bash
   ./quazaar
   ```

# ⚡ Quazaar

Quazaar is a lightweight WebSocket-based remote control server for Linux with optional music/player integration (Spotify, MPRIS).

This repository contains the server, WebSocket handlers, and example web client.

For full documentation, configuration details, and guides see the `docs/` folder.

## Quick Start

Build and run:

```bash
git clone https://github.com/codershubinc/Quazaar.git
cd Quazaar
go mod download
go build -o quazaar ./cmd/server
./quazaar
```

Default server address: `127.0.0.1:8765` (change via environment variables or `cmd/server`).

## Where to look next

- `docs/` — Full integration guides, API reference, and troubleshooting
- `internal/spotify/` — Spotify integration and token management
- `temp/web/` — Example web client for manual testing

## Contributing

Contributions and issues are welcome — please open a GitHub issue or PR. See `docs/CONTRIBUTING.md` if present.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.

---

**Made with ❤️ by [Swapnil Ingle](https://github.com/codershubinc) • [@codershubinc](https://github.com/codershubinc)**

All rights reserved.
