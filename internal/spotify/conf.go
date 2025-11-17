package spotify

import (
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"Quazaar/pkg/helpers"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SpotifyAPIBaseURL = os.Getenv("SPOTIFY_API_BASE_URL")

func Init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	SpotifyAPIBaseURL = os.Getenv("SPOTIFY_API_BASE_URL")

	fmt.Println("Checking for spotify tokens .../.../.../...")
	_, err := spotifyTokens.GetSpotifyRefreshToken()
	if err != nil {
		fmt.Println("❌ Spotify tokens not found. Please authenticate with Spotify.")
		fmt.Println("Refresh token not found in DB redirecting to auth spotify ../.../../..")
		helpers.LogMessage(
			helpers.WARNING,
			"Spotify refresh token not found in DB, redirecting to auth spotify",
			nil,
		)
	} else {
		fmt.Println("✅ Spotify tokens found.")
	}
}
