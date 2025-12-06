package spotify

import (
	spotifyConfig "Quazaar/internal/spotify/config"
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"Quazaar/pkg/helpers"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	spotifyConfig.SpotifyAPIBaseURL = os.Getenv("SPOTIFY_API_BASE_URL")

	fmt.Println("Checking for spotify tokens .../.../.../...")
	rt, err := spotifyTokens.GetSpotifyRefreshToken()
	if err != nil {
		fmt.Println("❌ Spotify tokens not found. Please authenticate with Spotify.")
		fmt.Println("Refresh token not found in DB redirecting to auth spotify ../.../../..")
		helpers.LogMessage(
			helpers.WARNING,
			"Spotify refresh token not found in DB, redirecting to auth spotify",
		)
	}
	fmt.Println("✅ Spotify tokens found.", strings.Split(rt, "")[0]+"...")
	helpers.LogMessage(
		helpers.INFO,
		"Spotify refresh token found in DB",
	)
	helpers.LogMessage(
		helpers.INFO,
		"Spotify :: Generating access token from refresh token",
	)
	_, err = spotifyTokens.GetSpotifyAccessToken()
	if err != nil {
		helpers.LogMessage(
			helpers.ERROR,
			"Spotify :: Failed to generate access token from refresh token: %s",
			err.Error(),
		)
	}
	helpers.LogMessage(
		helpers.INFO,
		"Spotify :: Access token generated successfully",
	)
}
