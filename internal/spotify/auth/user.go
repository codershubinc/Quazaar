package spotifyAuth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	spotifyConfig "Quazaar/internal/spotify/config"
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"Quazaar/pkg/helpers"
	"Quazaar/pkg/models"
)

// TokenResponse represents the response from Spotify's token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func GetUserProfile(accessToken string) (any, error) {
	req, err := http.NewRequest("GET", spotifyConfig.SpotifyAPIBaseURL+"/me", nil)
	if err != nil {
		return any(nil), err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Err at GetUserProfile :: ", err)
		return any(nil), err
	}
	defer resp.Body.Close()

	// Read the raw response body first
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err reading response body :: ", err)
		return any(nil), err
	}

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Spotify API returned status %d for /me endpoint\n", resp.StatusCode)
		fmt.Printf("Response body: %s\n", string(bodyBytes))

		// Try to parse as JSON error response
		var errorResp map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
			fmt.Printf("Spotify error response: %+v\n", errorResp)
		}

		return any(nil), fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode JSON response directly into user struct
	var user models.SpotifyUser
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		fmt.Println("Err at GetUserProfile Decode :: ", err)
		fmt.Printf("Response body: %s\n", string(bodyBytes))
		return any(nil), err
	}

	fmt.Printf("Successfully retrieved user profile: %s\n", user.ID)
	return user, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	state := generateRandomString(16)

	helpers.LogMessage(
		helpers.WARNING,
		"Spotify Login Request Initiated - User being redirected to OAuth authorization | State: %s",
		state,
	)

	authURL := "https://accounts.spotify.com/authorize?" + url.Values{
		"response_type": {"code"},
		"client_id":     {getClientID()},
		"scope":         {"user-read-private user-read-email user-read-currently-playing user-read-playback-state user-read-recently-played user-follow-read user-library-read user-modify-playback-state user-follow-modify user-library-modify"},
		"redirect_uri":  {getRedirectURI()},
		"state":         {state},
	}.Encode()

	// Redirect to Spotify authorization page
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state == "" {
		helpers.LogMessage(helpers.ERROR, "Spotify OAuth callback - state mismatch")
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	if code == "" {
		helpers.LogMessage(helpers.ERROR, "Spotify OAuth callback - no authorization code")
		http.Error(w, "No authorization code", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for access token
	tokenResponse, err := spotifyTokens.ExchangeCodeForToken(code, getRedirectURI())
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Store the refresh token
	if err := spotifyTokens.SetSpotifyRefreshToken(tokenResponse.RefreshToken, 0); err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to store refresh token: %v", err)
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	helpers.LogMessage(helpers.INFO, "Spotify authentication successful - Token expires in %d seconds", tokenResponse.ExpiresIn)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message":    "Spotify authentication successful",
		"expires_in": tokenResponse.ExpiresIn,
	})
}

// Helper function to generate random string for state parameter
func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// Helper functions to get environment variables
func getClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

func getRedirectURI() string {
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	if redirectURI == "" {
		return "http://127.0.0.1:8765/api/v0.1/spotify/callback"
	}
	return redirectURI
}
