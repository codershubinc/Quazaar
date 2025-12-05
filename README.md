# ⚡ Quazaar

Real-time music player integration with WebSocket remote control for Linux.

> **Note**: This project is in a developmental phase. Expect bugs and breaking changes. Refer to the beta branch for active changes.

---

📱 **Android App Update**: Beta version is out! [Download Quazaar App v0.1.0-beta](https://github.com/codershubinc/QuazaarApp/releases/download/v0.1.0-beta/Quazaar_v0.1.0-beta.apk)
For using the app with the server, refer to [docs/beta/README.md](docs/beta/README.md).

---

## 📚 Documentation

For the complete project documentation, including architecture, development journey, and detailed guides, please visit:

👉 **[Project Documentation & Journey](docs/PROJECT_DOCUMENTATION.md)** 👈

---

## 🎯 Features

- **Remote Command Execution**: Control your PC from any device on your network.
- **Real-time Music Display**: Shows currently playing track with album artwork using `playerctl`.
- **WebSocket Communication**: Fast, bidirectional communication between devices.
- **Secure Command Allowlist**: Only pre-approved commands can be executed.
- **Modern Web Interface**: Clean, responsive UI that works on desktop and mobile.
- **Auto-updating Music Info**: Track information refreshes every 1 second.

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

### Installation & Run

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
   go build -o quazaar ./cmd/server
   ```

4. **Run the server**

   ```bash
   ./quazaar
   ```

The server will start on `ws://0.0.0.0:8765/ws` (default).

## 📂 Where to look next

- **[docs/PROJECT_DOCUMENTATION.md](docs/PROJECT_DOCUMENTATION.md)** — **Start Here!** Full project overview and journey.
- `docs/` — Full integration guides, API reference, and troubleshooting.
- `internal/spotify/` — Spotify integration and token management.
- `temp/web/` — Example web client for manual testing.

## 🤝 Contributing

Contributions and issues are welcome — please open a GitHub issue or PR. See `CONTRIBUTING.md` for details.

## 📄 License

This project is licensed under the MIT License - see the `LICENSE` file for details.

---

**Made with ❤️ by [Swapnil Ingle](https://github.com/codershubinc) • [@codershubinc](https://github.com/codershubinc)**

All rights reserved.
