# Features

Quazaar is a modern media controller and file sharing server designed for local networks.

## 🎵 Media Control

Control your media playback from any device on your network.

- **Universal Control**: Works with any MPRIS-compatible media player (Spotify, VLC, Rhythmbox, etc.).
- **Real-time Updates**: See what's playing, artist, album, and progress in real-time.
- **Playback Controls**: Play, pause, next, previous, stop, and volume control.
- **Album Art**: Automatically fetches and displays high-quality album artwork.
- **Lyrics**: (Coming Soon) Synchronized lyrics display.

## 🔐 Secure Authentication

Robust security model to keep your server safe.

- **Single-User Owner**: Only one admin account ensures full control.
- **Token-Based Access**: Generate unique tokens for each device (phone, tablet, laptop).
- **Revocable Access**: Lost a device? Revoke its token instantly without affecting others.
- **Granular Permissions**: (Planned) Restrict tokens to specific features (e.g., "read-only" for guests).

## 📂 File Sharing

Share files easily across your local network.

- **Upload & Download**: Simple drag-and-drop interface for uploading files.
- **Temporary Links**: Generate shareable links that expire automatically.
- **Secure Storage**: Files are stored securely on the server.
- **Management**: List, delete, and manage shared files from the dashboard.

## 📡 WebSocket API

Real-time communication for developers and custom integrations.

- **Event-Driven**: Receive instant updates when media changes.
- **Low Latency**: Fast and responsive control via persistent WebSocket connections.
- **Simple JSON Protocol**: Easy to understand and implement in any language.

## 🖥️ Cross-Platform

Run Quazaar on your preferred operating system.

- **Linux**: Native support with `playerctl` integration.
- **Windows**: (Beta) Support via Windows Media APIs.
- **macOS**: (Planned) Future support for macOS media controls.
- **Docker**: Run in a container for easy deployment and isolation.

## 🛠️ Developer Friendly

Built for hackers and tinkerers.

- **Open Source**: MIT licensed, free to use and modify.
- **Go Backend**: High-performance, concurrent backend written in Go.
- **Clean Architecture**: Modular design makes it easy to add new features.
- **Comprehensive Docs**: Detailed guides and API references.
