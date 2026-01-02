# Project Review: Quazaar

## 🎯 Features

Quazaar is a comprehensive remote control and media integration server for Linux and Windows.

- **Remote Media Control**:
  - **Linux**: Native integration with DBus/MPRIS to control media players.
    - **Smart Player Detection**: Dynamically identifies the active player with a fallback priority (Playing > Paused > Last Active > Spotify > First Available).
    - **Artwork Caching**: Automatically downloads, caches, and serves artwork from remote URLs (e.g., Spotify CDN) to ensure fast loading and offline availability.
  - **Windows**: Uses a dedicated C# sidecar (`QuazaarMedia.exe`) embedded in the binary to interface with the Windows Media Session Manager (SMTC).
  - **Universal Metadata**: Real-time retrieval of Title, Artist, Album, and Artwork across platforms.
- **System Control**:
  - **Linux**: View WiFi status and Bluetooth devices (via `nmcli` and `bluetoothctl`).
  - **Volume Control**: Supported on both platforms (via DBus on Linux, Sidecar on Windows).
- **File Sharing**:
  - Temporary file upload mechanism allowing remote devices to send files to the host.
  - QR-code/Token-based authentication for file transfers.
- **Connectivity**:
  - **WebSocket Server**: Enables real-time, bidirectional communication for remote control apps.
  - **HTTP API**: RESTful endpoints for authentication, player info, and system status.
- **Security**:
  - User Authentication (Signup/Login) with SQLite database.
  - Token-based session management.
- **Deployment**:
  - Single-binary distribution using Go `embed` for static assets (web interface).

## ✅ Pros

- **Modular Architecture**: Follows the Standard Go Project Layout (`cmd/`, `internal/`, `pkg/`), keeping concerns separated (e.g., `auth`, `player`, `spotify`, `websocket`).
- **Robust Linux Integration**: The recent updates to `provider_linux.go` provide a resilient media control experience that handles multiple players gracefully.
- **Cross-Platform Strategy**:
  - **Linux**: Direct DBus/MPRIS implementation.
  - **Windows**: "Sidecar" pattern (embedded .exe) effectively bridges the gap between Go and Windows Runtime APIs.
- **Performance**:
  - **Artwork Caching**: The `internal/media` package implements a caching layer for artwork, reducing network requests and improving responsiveness.
  - **Unbuffered Channels**: The WebSocket implementation uses unbuffered channels for fresh messages, preventing lag from stale data.
- **Single Binary**: Bundling the web frontend and Windows helper executable directly into the Go binary simplifies distribution.

## ❌ Cons

- **Hardcoded Paths**:
  - `internal/media/artwork.go` contains a hardcoded fallback path: `/home/swap/Downloads/vector-music-note-icon.jpg`. This will fail on other machines.
  - Cache directories (`temp/spotify`, `temp/artwork`) are created relative to the execution path, which might not be ideal for all deployments.
- **Feature Parity**:
  - System information (WiFi, Bluetooth) is currently limited to Linux.
- **Observability**:
  - Logging is inconsistent (mix of `fmt` and `log`).
  - No centralized configuration struct; environment variables are accessed ad-hoc.
- **Reliability**:
  - **No Graceful Shutdown**: Active connections and background tasks are terminated immediately upon process exit.
  - **Error Handling**: Some DBus connection errors could be handled more gracefully with automatic reconnection logic.
- **Testing**:
  - Limited unit test coverage (`_test.go` files are sparse).
  - Integration tests exist as shell scripts (`tests/`), but Go-native tests are missing for core logic.

## 💡 Recommendations & Next Steps

### 1. Configuration & Environment

- **Centralized Config**: Create a `pkg/config` package. Instead of reading `os.Getenv` in multiple places, load all configuration into a struct at startup.
  - _Suggestion_: Use libraries like `kelseyhightower/envconfig` or `viper` to handle defaults and `.env` files.
- **Path Management**: Remove hardcoded paths (e.g., `/home/swap/...`).
  - Use `os.UserCacheDir()` or `os.UserConfigDir()` for storing temporary files and databases.
  - Embed a default "unknown album art" image using `//go:embed` instead of relying on a local file path.

### 2. Code Quality & Observability

- **Structured Logging**: Replace `fmt.Println` and standard `log` with Go's structured logger `log/slog` (available in Go 1.21+).
  - This allows for log levels (Info, Error, Debug) and machine-readable output (JSON) for production.
- **Linting**: Integrate `golangci-lint` into the workflow to catch potential bugs, unused variables, and style issues automatically.

### 3. Architecture & Parity

- **Refactor `internal/system`**: The current `internal/system` package directly calls Linux commands (`nmcli`, `bluetoothctl`).
  - **Define Interfaces**: Create `WiFiProvider`, `BluetoothProvider`, and `VolumeProvider` interfaces.
  - **Platform Implementations**:
    - Move existing logic to `linux_wifi.go`, `linux_bluetooth.go`, etc.
    - Create `windows_wifi.go` (stubbed or implemented via Sidecar) to satisfy the interface.
  - **Factory Pattern**: Use build tags (`//go:build linux`) or a factory function to initialize the correct provider at runtime.
- **Windows Sidecar Protocol**: Version the JSON protocol used between Go and the C# sidecar. If the sidecar binary is updated but the Go binary isn't (or vice versa), a version mismatch warning should be logged.

### 4. Reliability & Security

- **Graceful Shutdown**: Implement `context.Context` propagation for cancellation. Listen for `os.Interrupt` signals to close database connections, WebSocket listeners, and clean up temp files before exiting.
- **Input Validation**:
  - **File Uploads**: Validate file magic numbers (MIME types) server-side, not just file extensions. Enforce strict file size limits to prevent DoS.
  - **Command Injection**: Ensure that any shell commands (like `nmcli` or `playerctl` fallbacks) use `exec.Command` with separate arguments, never `bash -c` with concatenated user input.

### 5. Testing Strategy

- **Unit Tests**: Prioritize writing tests for:
  - `internal/media`: Parsing logic for metadata.
  - `internal/auth`: Token generation and validation.
- **Mocking**: Use interfaces for the DBus and File System interactions to allow testing the business logic without a running desktop environment.

### 6. File Structure & Naming Conventions

- **Test Co-location**: Adopt the standard Go practice of placing unit tests (`_test.go`) directly alongside the source code (e.g., `internal/auth/auth_test.go`). Keep the `tests/` directory reserved for end-to-end or integration test scripts.
- **Asset Organization**:
  - **Web Assets**: Move frontend assets (HTML, JS, CSS) to a top-level `web/` directory (e.g., `web/static`, `web/templates`). This aligns with standard Go project layouts and makes it clear where the frontend code lives.
  - **Project Assets**: Use a root-level `assets/` directory for repository-related images (logos, diagrams for README) that are not embedded in the binary.
  - **Avoid Redundancy**: The current path `internal/assets/assets/` is redundant and deep. Flattening this improves readability.
- **Package Naming**:
  - Continue using `snake_case` for Go filenames and directories.
  - Avoid generic package names like `helpers`. Try to scope them by domain (e.g., `osutil`, `netutil`) to avoid a "junk drawer" package.
- **Sidecar Isolation**: Since `cmd/windows-helper` contains a C# project, consider moving the C# source code to a top-level `sidecars/windows` directory. This clearly separates the Go application entry points (`cmd/`) from external component sources.
