package spotifyConfig

import (
	"os"
)

var SpotifyAPIBaseURL = os.Getenv("SPOTIFY_API_BASE_URL")
