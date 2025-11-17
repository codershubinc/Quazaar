package spotifyAuth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"Quazaar/internal/spotify"
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
	req, err := http.NewRequest("GET", spotify.SpotifyAPIBaseURL+"/me", nil)
	fmt.Println("Req url" + spotify.SpotifyAPIBaseURL + "/me")
	if err != nil {
		return any(nil), err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return any(nil), err
	}
	defer resp.Body.Close()
	var user models.SpotifyUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return any(nil), err
	}
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
		"scope":         {"user-read-private user-read-email user-read-currently-playing user-read-playback-state user-read-recently-played user-follow-read user-library-read user-modify-playback-state"},
		"redirect_uri":  {getRedirectURI()},
		"state":         {state},
	}.Encode()

	// Redirect to Spotify authorization page
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// storeRefreshToken stores the refresh token in the database
func storeRefreshToken(refreshToken string, expiresIn int) error {
	if err := spotifyTokens.SetSpotifyRefreshToken(refreshToken, expiresIn); err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
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
	if err := storeRefreshToken(tokenResponse.RefreshToken, tokenResponse.ExpiresIn); err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to store refresh token: %v", err)
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
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
