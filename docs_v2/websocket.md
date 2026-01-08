# WebSocket API

// endpoint /ws [GET]
// Description: Upgrades the connection to a WebSocket connection.
// Required Headers:
//   Authorization: Bearer <token>
//   Connection: Upgrade
//   Upgrade: websocket
// Response:
//   On success, connection is upgraded.
//   Server sends: {"message": "Welcome to the WebSocket server!"}

// Message Type: Player Commands
// Description: Control media playback.
// Direction: Client -> Server
// JSON Payload:
// {
//   "command": "play"  // Options: "play", "pause", "player_toggle", "next", "prev", "seek_forward", "seek_backward"
// }
// Response (JSON):
// Success:
// {
//   "status": "success",
//   "message": "command_executed",
//   "data": { "command": "play" }
// }
// Error:
// {
//   "status": "error",
//   "message": "command_failed",
//   "data": { "error": "description" }
// }

// Message Type: System - Volume
// Description: Control system volume.
// Direction: Client -> Server
// JSON Payload:
// {
//   "type": "system",
//   "msg_of": "volume",
//   "action": "inc" // Options: "inc", "dec", "get", "mute", "set"
//   // "set_to_vol": 50 (Required for action="set")
// }
// Response (JSON):
// {
//   "status": "success",
//   "message": "system",
//   "data": {
//     "current_volume": 55,
//     "previous_volume": 50,
//     "action": "inc",
//     "success": true
//   }
// }

// Message Type: System - Brightness
// Description: Control screen brightness.
// Direction: Client -> Server
// JSON Payload:
// {
//   "type": "system",
//   "msg_of": "brightness",
//   "action": "inc" // Options: "inc", "dec", "get", "set"
//   // "set_to": 70 (Required for action="set")
// }
// Response (JSON):
// {
//   "status": "success",
//   "message": "system",
//   "data": {
//     "current_brightness": 75,
//     "action": "inc",
//     "success": true
//   }
// }

// Message Type: Spotify Artist
// Description: Get artist info or follow artist.
// Direction: Client -> Server
// JSON Payload:
// {
//   "message": "spotify_artist",
//   "action": "get",      // Options: "get", "follow"
//   "artistId": "spotify_artist_id"
// }
// Response (JSON):
// {
//   "status": "success",
//   "message": "spotify_artist_get",
//   "data": { ...artist info... }
// }
