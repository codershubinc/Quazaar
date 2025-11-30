# API Reference

Quazaar provides a REST API for authentication and token management, and a WebSocket API for real-time communication and player control.

## REST API

The REST API is used primarily for authentication and managing access tokens.

### Base URL

```text
http://<server-ip>:8765/api
```

### Endpoints

#### 1. Register User

Creates the initial owner account. This endpoint only works if no user is registered.

- **Endpoint**: `/register`
- **Method**: `POST`
- **Content-Type**: `application/json`

**Request Body**:

```json
{
  "username": "admin",
  "password": "secure_password"
}
```

**Response (201 Created)**:

```json
{
  "success": true,
  "message": "User registered successfully",
  "user_id": 1,
  "username": "admin"
}
```

#### 2. Login

Authenticates the user and returns a session token.

- **Endpoint**: `/login`
- **Method**: `POST`
- **Content-Type**: `application/json`

**Request Body**:

```json
{
  "username": "admin",
  "password": "secure_password"
}
```

**Response (200 OK)**:

```json
{
  "success": true,
  "message": "Login successful",
  "user_id": 1,
  "username": "admin",
  "token": "session_token_string..."
}
```

#### 3. Create Token

Generates a new persistent token for a device or service.

- **Endpoint**: `/tokens/create`
- **Method**: `POST`
- **Headers**: `Authorization: <token>`

**Request Body**:

```json
{
  "name": "My Device",
  "service": "mobile",
  "duration_hours": 720
}
```

**Response (201 Created)**:

```json
{
  "id": 1,
  "name": "My Device",
  "token": "new_token_string...",
  "service": "mobile",
  "expires_at": "2025-12-31T23:59:59Z",
  "active": true
}
```

#### 4. List Tokens

Retrieves all active and inactive tokens.

- **Endpoint**: `/tokens/list`
- **Method**: `GET`
- **Headers**: `Authorization: <token>`

**Response (200 OK)**:

```json
{
  "success": true,
  "tokens": [
    {
      "id": 1,
      "name": "My Device",
      "active": true
      // ... other fields
    }
  ]
}
```

#### 5. Revoke Token

Invalidates a specific token.

- **Endpoint**: `/tokens/revoke`
- **Method**: `POST`
- **Headers**: `Authorization: <token>`

**Request Body**:

```json
{
  "token": "token_to_revoke"
}
```

## WebSocket API

The WebSocket API is used for real-time media updates and controlling the player.

### Connection

- **URL**: `ws://<server-ip>:8765/ws`
- **Authentication**: Pass `token` in query string or `Authorization` header.

### Player Commands

Send these JSON messages to control the media player.

#### Play

```json
{
  "command": "play"
}
```

#### Pause

```json
{
  "command": "pause"
}
```

#### Toggle Play/Pause

```json
{
  "command": "player_toggle"
}
```

#### Next Track

```json
{
  "command": "next"
}
```

#### Previous Track

```json
{
  "command": "prev"
}
```

#### Volume Control

```json
{
  "command": "volume_up"
}
```

```json
{
  "command": "volume_down"
}
```

#### Stop

```json
{
  "command": "stop"
}
```

### Server Responses

The server sends JSON messages to connected clients.

#### Media Update

Sent when track changes or playback status updates.

```json
{
  "type": "media_update",
  "data": {
    "title": "Song Title",
    "artist": "Artist Name",
    "album": "Album Name",
    "status": "Playing",
    "position": 120,
    "duration": 300,
    "artUrl": "http://..."
  }
}
```

#### Command Response

Sent after executing a command.

```json
{
  "status": "success",
  "message": "command_executed",
  "data": {
    "command": "player_toggle"
  }
}
```

## Error Handling

Standard HTTP status codes are used for the REST API:

- **200 OK**: Success
- **201 Created**: Resource created
- **400 Bad Request**: Invalid input
- **401 Unauthorized**: Missing or invalid token
- **500 Internal Server Error**: Server-side issue
