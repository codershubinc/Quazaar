# System Brightness Control

The system brightness control module allows managing the system's screen brightness via WebSocket commands. It uses `brightnessctl` under the hood.

## WebSocket API

All brightness commands should be sent with the following base structure:

```json
{
  "type": "system",
  "msg_of": "brightness",
  "action": "ACTION_NAME",
  ...params
}
```

### Actions

#### 1. Get Brightness

Retrieves the current system brightness percentage.

**Request:**

```json
{
  "type": "system",
  "msg_of": "brightness",
  "action": "get"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_brightness": 60,
    "action": "get",
    "success": true
  }
}
```

#### 2. Increase Brightness

Increases the system brightness by 5%.

**Request:**

```json
{
  "type": "system",
  "msg_of": "brightness",
  "action": "inc"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_brightness": 65,
    "action": "inc",
    "success": true
  }
}
```

#### 3. Decrease Brightness

Decreases the system brightness by 5%.

**Request:**

```json
{
  "type": "system",
  "msg_of": "brightness",
  "action": "dec"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_brightness": 60,
    "action": "dec",
    "success": true
  }
}
```

#### 4. Set Brightness

Sets the system brightness to a specific percentage (0-100).

**Request:**

```json
{
  "type": "system",
  "msg_of": "brightness",
  "action": "set",
  "set_to": 80
}
```

**Response:**

```json
{
  "status": "success",
  "message": "system",
  "data": {
    "current_brightness": 80,
    "action": "set",
    "success": true
  }
}
```

## Implementation Details

- **Backend**: `internal/system/brightness/`
  - `brightness.go`: Core logic using `brightnessctl` commands.
  - `handle_brightness.go`: WebSocket message handler.
- **Frontend**: `temp/web/index.html` (Test Client)

## Dependencies

- `brightnessctl` must be installed on the host system.
