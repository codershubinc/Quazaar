# Architecture & Project Structure

Quazaar follows the standard Go project layout, organizing code into logical layers to ensure scalability and maintainability.

## 📂 Directory Structure

```
Quazaar/
├── cmd/                    # Application entry points
│   └── server/             # Main server application
│       └── main.go         # Entry point
│
├── internal/               # Private application code (not importable by other projects)
│   ├── api/                # HTTP API router and handlers
│   ├── auth/               # Authentication logic (signup, login, tokens)
│   ├── db/                 # Database connection and schema management
│   ├── fileshare/          # File sharing functionality
│   ├── media/              # Media player integration (playerctl, Windows)
│   ├── middleware/         # HTTP middleware (auth, logging)
│   ├── player/             # Player control logic and handlers
│   ├── poller/             # Background polling for media updates
│   ├── spotify/            # Spotify API integration
│   ├── system/             # System info (WiFi, Bluetooth, App Launcher)
│   └── websocket/          # WebSocket connection handling
│
├── pkg/                    # Public library code (can be used by other projects)
│   ├── helpers/            # Utility functions (error handling, random strings)
│   └── models/             # Shared data structures (JSON responses, DB models)
│
├── docs/                   # Documentation files
├── assets/                 # Static assets (images, etc.)
└── temp/                   # Temporary files (uploads, cache)
```

## 🏗️ Core Components

### 1. Server (`cmd/server`)
The entry point of the application. It initializes the database, sets up the HTTP router, starts the WebSocket handler, and launches the background poller.

### 2. Internal Packages (`internal/`)
- **`api`**: Defines the HTTP routes and links them to handlers.
- **`auth`**: Manages user registration, login, and JWT/token generation.
- **`db`**: Handles SQLite database connections and migrations.
- **`media`**: Interacts with the OS media controls (MPRIS on Linux, Win32 APIs on Windows).
- **`websocket`**: Manages real-time bidirectional communication with clients.
- **`poller`**: Periodically checks for media changes and broadcasts updates to connected WebSocket clients.

### 3. Shared Packages (`pkg/`)
- **`models`**: Defines the structs used for API responses and database records, ensuring consistency across the app.
- **`helpers`**: Provides common utility functions to reduce code duplication.

## 🔄 Data Flow

1. **Client Request**: A client (web or mobile) sends an HTTP request or WebSocket message.
2. **Router/Handler**: The request is routed to the appropriate handler in `internal/`.
3. **Service Layer**: The handler calls logic in `internal/` packages (e.g., `player`, `spotify`).
4. **System Interaction**: The service interacts with the OS (e.g., `playerctl`) or external APIs (Spotify).
5. **Response**: The result is formatted using `pkg/models` and sent back to the client.

## 🗄️ Database Schema

Quazaar uses SQLite with the following main tables:
- **`users`**: Stores user credentials (hashed passwords).
- **`tokens`**: Manages authentication tokens for devices.
- **`file_share_device_tokens`**: Manages temporary tokens for file sharing.

## 🔌 External Integrations

- **Playerctl**: Used on Linux for controlling media players.
- **Spotify Web API**: Used for advanced Spotify control and metadata.
- **D-Bus**: Used for low-level Linux system communication.
