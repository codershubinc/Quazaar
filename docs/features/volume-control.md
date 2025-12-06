# System Volume Control

The system volume control module allows managing the system's audio volume and mute state via WebSocket commands. It uses `pactl` (PulseAudio Control) under the hood.

## WebSocket API

All volume commands should be sent with the following base structure:

```json
{
  "type": "system",
  "msg_of": "volume",
  "action": "ACTION_NAME",
  ...params
}
```

### Actions

#### 1. Get Volume

Retrieves the current system volume.

**Request:**

```json
{
  "type": "system",
  "msg_of": "volume",
  "action": "get"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_volume": 45,
    "action": "get",
    "success": true
  }
}
```

#### 2. Increase Volume

Increases the system volume by 5%.

**Request:**

```json
{
  "type": "system",
  "msg_of": "volume",
  "action": "inc"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_volume": 50,
    "action": "inc",
    "success": true
  }
}
```

#### 3. Decrease Volume

Decreases the system volume by 5%.

**Request:**

```json
{
  "type": "system",
  "msg_of": "volume",
  "action": "dec"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_volume": 45,
    "action": "dec",
    "success": true
  }
}
```

#### 4. Set Volume

Sets the system volume to a specific percentage (0-100).

**Request:**

```json
{
  "type": "system",
  "msg_of": "volume",
  "action": "set",
  "set_to_vol": 75
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_volume": 75,
    "action": "set",
    "success": true
  }
}
```

#### 5. Toggle Mute

Toggles the system mute state.

**Request:**

```json
{
  "type": "system",
  "msg_of": "volume",
  "action": "mute"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "action": "mute",
    "success": true,
    "meta": {
      "is_muted": true
    }
  }
}
```

## Implementation Details

- **Backend**: `internal/system/volume/`
  - `volume.go`: Core logic using `pactl` commands.
  - `handle_volume.go`: WebSocket message handler.
- **Frontend**: `temp/web/index.html` (Test Client)

## Dependencies

- `pactl` (PulseAudio Utils) must be installed on the host system.
