# Spotify API

// endpoint /api/v0.1/spotify/track/currently-playing [GET]
// Description: Gets the currently playing track from Spotify API.
// Required Headers:
//   Authorization: Bearer <token>
// Response Data (JSON):
// {
//   "item": {
//     "id": "track_id",
//     "name": "Track Name",
//     "artists": [{"name": "Artist"}],
//     "album": {"name": "Album", "images": [...]},
//     ...
//   },
//   "is_playing": true,
//   "progress_ms": 12000
// }

// endpoint /api/v0.1/spotify/track/add-to-library [GET/PUT]
// Description: Adds the currently playing track (or implicit track) to the user's Spotify library.
// Note: While typically a PUT action, the implementation may support GET or body-less PUT using current context.
// Required Headers:
//   Authorization: Bearer <token>
// Response Data (JSON):
// {
//   "status": "Track added to library"
// }

// endpoint /api/v0.1/spotify/track/check-in-library [GET]
// Description: Checks if the currently playing track is already in the user's Spotify library.
// Required Headers:
//   Authorization: Bearer <token>
// Response Data (JSON):
// {
//   "is_in_library": true
// }

// endpoint /api/v0.1/spotify/artist [GET]
// Description: Retrieves information about a specific artist.
// Required Headers:
//   Authorization: Bearer <token>
// Query Params:
//   id: "spotify_artist_id"
// Response Data (JSON):
// {
//   "id": "artist_id",
//   "name": "Artist Name",
//   "followers": { "total": 1000000 },
//   "genres": ["pop", "rock"],
//   "images": [...]
// }

// endpoint /api/v0.1/spotify/artist/follow [GET/PUT]
// Description: Follows a specific artist on Spotify.
// Required Headers:
//   Authorization: Bearer <token>
// Query Params:
//   id: "spotify_artist_id"
// Response Data:
// Status 204 No Content (on success)

// endpoint /api/v0.1/spotify/me [GET]
// Description: Retrieves the Spotify profile of the authenticated user.
// Required Headers:
//   Authorization: Bearer <token>
// Response Data (JSON):
// {
//   "id": "spotify_user_id",
//   "display_name": "John Doe",
//   "email": "john@example.com",
//   "product": "premium",
//   "images": [...]
// }
