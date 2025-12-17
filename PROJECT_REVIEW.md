# Project Review: Quazaar

## 🎯 Features

Quazaar is a comprehensive remote control and media integration server for Linux.

*   **Remote Media Control**:
    *   **Linux**: Integrates with `playerctl` (MPRIS) to control media players (Play, Pause, Next, Previous).
    *   **Windows**: Uses a dedicated C# sidecar (`QuazaarMedia.exe`) embedded in the binary to interface with the Windows Media Session Manager (SMTC). This allows control of Spotify, Chrome, generic media apps, and system volume on Windows.
    *   Real-time retrieval of metadata (Title, Artist, Album, Art) on both platforms.
    *   Specific support for Spotify integration (Authorization, Artist info).
*   **System Control (Linux Only)**:
    *   View WiFi status and Bluetooth devices (relies on `nmcli` and `bluetoothctl`).
    *   Volume control implementation (referenced in Todo).
*   **File Sharing**:
    *   Temporary file upload mechanism allowing remote devices to send files to the host.
    *   QR-code/Token-based authentication for file transfers.
*   **Connectivity**:
    *   **WebSocket Server**: Enables real-time, bidirectional communication for remote control apps.
    *   **HTTP API**: RESTful endpoints for authentication, player info, and system status.
*   **Security**:
    *   User Authentication (Signup/Login) with SQLite database.
    *   Token-based session management.
*   **Deployment**:
    *   Single-binary distribution using Go `embed` for static assets (web interface).

## ✅ Pros

*   **Modular Architecture**: The project follows the Standard Go Project Layout (using `cmd/`, `internal/`, `pkg/`), keeping concerns separated (e.g., `auth`, `player`, `spotify`, `websocket`).
*   **Single Binary Deployment**: Usage of `embed` package allows bundling the web frontend and assets directly into the Go binary, simplifying distribution.
*   **Standard Library Usage**: Heavy reliance on the standard library (`net/http`, `image`, `os`) reduces external dependency bloat.
*   **Cross-Platform Architecture**:
    *   **Linux**: Native implementation via DBus/MPRIS.
    *   **Windows**: Creative "Sidecar" pattern where a helper .exe is embedded in the Go binary, extracted at runtime, and communicates via Stdin/Stdout JSON streams. This solves the difficulty of calling Windows Runtime APIs from Go.
*   **Real-Time Capabilities**: The poller architecture ensures connected clients receive up-to-date media information without manual refreshing.

## ❌ Cons

*   **Feature Parity**:
    *   System information (WiFi, Bluetooth) is currently implemented only for Linux. The Windows implementation lacks these system-level features.

*   **Observability & Debugging**:
    *   Logging is inconsistent, mixing `fmt.Println` and `log.Println`. This makes parsing logs in production difficult.
    *   Lack of a centralized configuration struct; environment variables are accessed in various places.
*   **Reliability**:
    *   **No Graceful Shutdown**: The server kills active connections and background tasks immediately upon process termination.
    *   **Memory Management**: File uploads load the entire file into memory (or temp disk) before processing, which is risky for large files.
*   **Testing**:
    *   There is a noticeable lack of unit tests (`_test.go` files). Crucial logic like authentication and file handling is untested.
*   **Code Hygiene**:
    *   Presence of "scratchpad" comments and commented-out code in the codebase.
    *   Hardcoded paths (e.g., `Downloads/Quazaar`) might not work on all environments or Docker containers.

## 💡 Additional Observations

*   **Security Posture**: The "Command Allowlist" mentioned in the README is a critical feature. Care must be taken to ensure that the `playerctl` or system commands cannot be injected with malicious arguments.
*   **Spotify Integration**: The project seems to be moving towards a hybrid approach—using local MPRIS control for generic players and the Spotify Web API for richer data. This is a good strategy but adds complexity in token management.
*   **Concurrency**: The use of global variables in some packages (e.g., database instance, configuration) could lead to race conditions as the application scales. A dependency injection approach would be cleaner.
*   **Frontend Integration**: The project includes a web client. Ensuring the API remains backwards compatible is crucial as the Android app (Beta) and Web client evolve.
