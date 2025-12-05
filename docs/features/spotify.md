# Spotify Integration

## Architecture

```mermaid
graph TD
    User[User] -->|1. Auth Request| Server[Quazaar Server]
    Server -->|2. Redirect| SpotifyAuth[Spotify Auth Page]
    SpotifyAuth -->|3. Callback Code| Server
    Server -->|4. Exchange Code| SpotifyAPI[Spotify API]
    SpotifyAPI -->|5. Access/Refresh Token| Server

    subgraph "Runtime"
        Poller[Media Poller] -->|Poll| SpotifyAPI
        SpotifyAPI -->|Track Data| Poller
        Poller -->|Broadcast| WS[WebSocket]
        WS -->|Update UI| User

        User -->|Command (Play/Pause)| WS
        WS -->|API Call| SpotifyAPI
    end
```

## Getting Started

1. **Create App**: Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications) and create an app.
2. **Configure**: Set Redirect URI to `http://localhost:8765/spotify/callback`.
3. **Environment**: Add credentials to your `.env` file:

```bash
export SPOTIFY_CLIENT_ID="your_client_id"
export SPOTIFY_CLIENT_SECRET="your_client_secret"
export SPOTIFY_REDIRECT_URI="http://localhost:8765/spotify/callback"
```

## Authentication

Navigate to `http://localhost:8765/spotify/auth` to login and grant permissions.

## Features & Commands

The integration supports the following WebSocket commands:

| Command | Description | Payload |
|---------|-------------|---------|
| `spotify_auth_status` | Check if user is authenticated | None |
| `spotify_current_track` | Get currently playing track info | None |
| `spotify_play` | Resume playback | `{"device_id": "..."}` (optional) |
| `spotify_pause` | Pause playback | None |
| `spotify_next` | Skip to next track | None |
| `spotify_previous` | Skip to previous track | None |
| `spotify_volume` | Set volume (0-100) | `{"volume": 50}` |
| `spotify_playlists` | List user playlists | None |

## Scopes

The application requests permissions for:
- Playback control & state reading
- Library & Playlist access
- User top/recent tracks
