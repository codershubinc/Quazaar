# Quazaar Android Integration Guide (Kotlin)

This guide provides instructions and code snippets for integrating Quazaar's latest features (v0.1.6) into an Android application using Kotlin.

## 📚 Prerequisites

- **Retrofit 2**: For REST API calls.
- **OkHttp 3/4**: For WebSocket connections.
- **Gson/Moshi**: For JSON parsing.
- **Coroutines**: For asynchronous operations.

## 1. Authentication

### Login Endpoint

- **URL**: `POST /api/v0.1/login`
- **Body**: `{"username": "...", "password": "..."}`
- **Response**: `{"token": "...", "tokenType": "deviceId", ...}`

### Kotlin Implementation

**Data Classes:**

```kotlin
data class LoginRequest(
    val username: String,
    val password: String
)

data class LoginResponse(
    val success: Boolean,
    val message: String,
    val token: String,
    val tokenType: String,
    val username: String
)
```

**Retrofit Interface:**

```kotlin
interface QuazaarApi {
    @POST("api/v0.1/login")
    suspend fun login(@Body request: LoginRequest): Response<LoginResponse>
}
```

**Usage:**

```kotlin
val response = api.login(LoginRequest("user", "pass"))
if (response.isSuccessful) {
    val token = response.body()?.token
    // Save token to SharedPreferences/DataStore
}
```

## 2. WebSocket Connection

### Connection URL

- **URL**: `ws://<host>:8765/ws?deviceId=<token>`
- **Note**: The `deviceId` query parameter is **mandatory** for authentication.

### OkHttp Implementation

```kotlin
val client = OkHttpClient.Builder().build()

fun connectWebSocket(host: String, token: String) {
    val request = Request.Builder()
        .url("ws://$host:8765/ws?deviceId=$token")
        .build()

    val listener = object : WebSocketListener() {
        override fun onOpen(webSocket: WebSocket, response: Response) {
            Log.d("Quazaar", "Connected!")
        }

        override fun onMessage(webSocket: WebSocket, text: String) {
            Log.d("Quazaar", "Received: $text")
            // Parse JSON message here
        }

        override fun onClosing(webSocket: WebSocket, code: Int, reason: String) {
            webSocket.close(1000, null)
        }

        override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
            Log.e("Quazaar", "Error: ${t.message}")
        }
    }

    client.newWebSocket(request, listener)
}
```

## 3. System Controls (Volume & Brightness)

Send these JSON payloads over the active WebSocket connection.

### Data Models

```kotlin
data class SystemCommand(
    val type: String = "system",
    val msg_of: String, // "volume" or "brightness"
    val action: String, // "inc", "dec", "set", "mute", "get"
    val set_to_vol: Int? = null, // For volume set
    val set_to: Int? = null      // For brightness set
)
```

### Volume Commands

**Increase Volume:**

```kotlin
val cmd = SystemCommand(msg_of = "volume", action = "inc")
webSocket.send(gson.toJson(cmd))
```

**Decrease Volume:**

```kotlin
val cmd = SystemCommand(msg_of = "volume", action = "dec")
webSocket.send(gson.toJson(cmd))
```

**Set Volume (0-100):**

```kotlin
val cmd = SystemCommand(msg_of = "volume", action = "set", set_to_vol = 50)
webSocket.send(gson.toJson(cmd))
```

**Mute Toggle:**

```kotlin
val cmd = SystemCommand(msg_of = "volume", action = "mute")
webSocket.send(gson.toJson(cmd))
```

### Brightness Commands

**Increase Brightness:**

```kotlin
val cmd = SystemCommand(msg_of = "brightness", action = "inc")
webSocket.send(gson.toJson(cmd))
```

**Decrease Brightness:**

```kotlin
val cmd = SystemCommand(msg_of = "brightness", action = "dec")
webSocket.send(gson.toJson(cmd))
```

**Set Brightness (0-100):**

```kotlin
val cmd = SystemCommand(msg_of = "brightness", action = "set", set_to = 75)
webSocket.send(gson.toJson(cmd))
```

## 4. Media Controls

Send these JSON payloads for media playback.

### Data Model

```kotlin
data class MediaCommand(
    val command: String // "play", "pause", "next", "previous"
)
```

### Usage

**Play:**

```kotlin
webSocket.send(gson.toJson(MediaCommand("play")))
```

**Pause:**

```kotlin
webSocket.send(gson.toJson(MediaCommand("pause")))
```

**Next Track:**

```kotlin
webSocket.send(gson.toJson(MediaCommand("next")))
```

**Previous Track:**

```kotlin
webSocket.send(gson.toJson(MediaCommand("previous")))
```

## 5. Handling Responses

The server sends JSON messages asynchronously. You should parse the `type` or `message` field to determine how to handle the data.

### Response Structure

```kotlin
data class WebSocketResponse(
    val type: String?,    // "system", etc.
    val message: String?, // "media_info", "spotify_artist", etc.
    val status: String?,  // "success", "error"
    val data: JsonElement? // Dynamic payload
)
```

### Example Handler

```kotlin
fun handleMessage(json: String) {
    val response = gson.fromJson(json, WebSocketResponse::class.java)

    when {
        response.type == "system" -> {
            // Handle volume/brightness updates
            // e.g., update SeekBars
        }
        response.message == "media_info" -> {
            // Update Now Playing UI (Title, Artist, Album Art)
        }
        response.message == "spotify_artist_get" -> {
            // Show Artist Info
        }
    }
}
```

## 6. AI Assistant Prompts

Use these prompts to quickly generate implementation code with GitHub Copilot or other AI assistants.

### Implement System Controls (Volume & Brightness Only)

Copy and paste this prompt to generate the helper class and SeekBar integration. This assumes you already have WebSocket connectivity set up.

```text
I have an existing Android app (Kotlin) connected to a Quazaar server via WebSocket.
I need to add ONLY the System Controls for Volume and Brightness.

Please provide:
1. A data class `SystemCommand` for the specific payload structure.
2. A helper function to send these commands using my existing WebSocket.
3. Code to wire up two SeekBars (one for Volume, one for Brightness) to send "set" commands.

API Specs:
- Payload: { "type": "system", "msg_of": "volume"|"brightness", "action": "inc"|"dec"|"set"|"mute", "set_to_vol": INT, "set_to": INT }
- Volume Set: msg_of="volume", action="set", set_to_vol=0-100
- Brightness Set: msg_of="brightness", action="set", set_to=0-100
```
