package spotifyTokens

import (
	"Quazaar/internal/db"
	"time"
)

var spotifyClientRefreshToken string
var spotifyClientAccessToken string
var spotifyAccessTokenExpiry time.Time

func GetSpotifyRefreshToken() (string, error) {
	// Return from cache if available
	if spotifyClientRefreshToken != "" {
		return spotifyClientRefreshToken, nil
	}

	// Fetch from database and cache it
	token, err := db.GetToken("spotifyClientRefreshToken")
	if err != nil {
		return "", err
	}
	spotifyClientRefreshToken = token
	return spotifyClientRefreshToken, nil
}

// GetSpotifyAccessToken retrieves a valid Spotify access token
// Returns cached token if still valid, otherwise refreshes it using the refresh token
func GetSpotifyAccessToken() (token string, err error) {
	// Return cached token if still valid
	if spotifyClientAccessToken != "" && time.Now().Before(spotifyAccessTokenExpiry) {
		return spotifyClientAccessToken, nil
	}

	// Get refresh token
	refreshToken, err := GetSpotifyRefreshToken()
	if err != nil {
		return "", err
	}

	tokenResponse, err := RefreshAccessToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Cache the new token and expiry time
	spotifyClientAccessToken = tokenResponse.AccessToken
	spotifyAccessTokenExpiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	return spotifyClientAccessToken, nil
}

// SetSpotifyRefreshToken stores a new refresh token in the database and cache
func SetSpotifyRefreshToken(newRefreshToken string, expiry int) error {
	err := db.StoreToken("spotifyClientRefreshToken", "spotify", newRefreshToken, expiry)
	if err != nil {
		return err
	}
	spotifyClientRefreshToken = newRefreshToken
	return nil
}
