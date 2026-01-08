# Player API

## `/api/v0.1/player/info` [GET]

**Description:** Gets current player info via D-Bus with fallback (default method).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "player": {
    "Title": "Song Title",
    "Artist": "Artist Name",
    "Album": "Album Name",
    "Artwork": "http://...",
    "Length": "3:45",
    "Status": "Playing",
    "Player": "spotify"
  }
}
```

## `/api/v0.1/player/info/dbus` [GET]

**Description:** Gets player info strictly via D-Bus (Linux).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Optional Query Params:**
- `player`: `org.mpris.MediaPlayer2.spotify` (specific player service name)

**Response Data (JSON):**
```json
{
  "success": true,
  "player": {
    "Title": "Song Title",
    "Artist": "Artist Name",
    "Player": "spotify"
  }
}
```

## `/api/v0.1/player/info/windows` [GET]

**Description:** Gets player info using Windows APIs (Windows Only).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "player": { "..." }
}
```

## `/api/v0.1/player/list` [GET]

**Description:** Lists all active media players detected by the system.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "players": ["spotify", "vlc", "firefox"],
  "count": 3
}
```

## `/api/v0.1/player/mpris/list` [GET]

**Description:** Lists active MPRIS players via D-Bus.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "players": ["org.mpris.MediaPlayer2.spotify"],
  "count": 1,
  "source": "D-Bus MPRIS"
}
```

## `/api/v0.1/player/windows/list` [GET]

**Description:** Lists active Windows media players.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "players": ["Spotify.exe"],
  "count": 1,
  "source": "Windows APIs"
}
```

## `/api/v0.1/player/play-pause` [POST]

**Description:** Toggles play/pause on the current active player.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "message": "Play/Pause toggled"
}
```

## `/api/v0.1/player/play` [POST]

**Description:** Sends 'play' command to the current player (via playerctl).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "message": "Playing"
}
```

## `/api/v0.1/player/pause` [POST]

**Description:** Sends 'pause' command to the current player (via playerctl).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "message": "Paused"
}
```

## `/api/v0.1/player/next` [POST]

**Description:** Skips to the next track.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "message": "Skipped to next track"
}
```

## `/api/v0.1/player/previous` [POST]

**Description:** Skips to the previous track.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "message": "Skipped to previous track"
}
```
