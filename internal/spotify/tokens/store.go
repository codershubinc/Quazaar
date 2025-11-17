package spotifyTokens

import (
	"Quazaar/internal/db"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// In-memory cache for tokens to avoid repeated database queries
var spotifyClientRefreshToken string
var spotifyClientAccessToken string
var spotifyAccessTokenExpiry time.Time

// GetSpotifyRefreshToken retrieves the Spotify refresh token from cache or database
// Refresh tokens are long-lived and used to obtain new access tokens
//
// Returns:
//   - string: The refresh token
//   - error: If token is not found in cache or database
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
// Automatically refreshes the token if it has expired
//
// Returns:
//   - string: A valid access token
//   - error: If refresh token is unavailable or refresh fails
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

	// Refresh the access token
	token, err = refreshSpotifyAccessToken(refreshToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

// refreshSpotifyAccessToken exchanges a refresh token for a new access token
// This is called automatically by GetSpotifyAccessToken when the token expires
//
// # Access tokens expire after 1 hour and must be refreshed
//
// Parameters:
//   - refreshToken: The refresh token obtained during OAuth authorization
//
// Returns:
//   - string: New access token
//   - error: If the refresh request fails
func refreshSpotifyAccessToken(refreshToken string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// Build Basic Auth header with client credentials
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Cache the new token and expiry time
	spotifyAccessTokenExpiry = time.Now().Add(time.Duration(result["expires_in"].(float64)) * time.Second)
	spotifyClientAccessToken = result["access_token"].(string)

	return spotifyClientAccessToken, nil
}

// SetSpotifyRefreshToken stores a new refresh token in the database and cache
// This is typically called after successful OAuth authorization
//
// Parameters:
//   - newRefreshToken: The refresh token received from Spotify OAuth
//   - expiry: Token expiry time in seconds (usually 3600 for access tokens)
//
// Returns:
//   - error: If database storage fails
func SetSpotifyRefreshToken(newRefreshToken string, expiry int) error {
	err := db.StoreToken("spotifyClientRefreshToken", "spotify", newRefreshToken, expiry)
	if err != nil {
		return err
	}
	spotifyClientRefreshToken = newRefreshToken
	return nil
}
