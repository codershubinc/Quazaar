package spotifyTokens

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"Quazaar/pkg/helpers"
)

// TokenExchangeResponse represents the response from Spotify's token endpoint
type TokenExchangeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// getting refresh token
func ExchangeCodeForToken(code, redirectURI string) (*TokenExchangeResponse, error) {

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {redirectURI},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to create token exchange request: %v", err)
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to exchange code for token: %v", err)
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to read token response: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		helpers.LogMessage(helpers.ERROR, "Token exchange failed with status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("token exchange failed: status %d", resp.StatusCode)
	}

	var tokenResponse TokenExchangeResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to parse token response: %v", err)
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	helpers.LogMessage(helpers.INFO, "Successfully exchanged authorization code for tokens | Expires in: %d seconds", tokenResponse.ExpiresIn)

	return &tokenResponse, nil
}

// RefreshAccessToken refreshes an expired access token using the refresh token
// Returns a new access token and updates the expiry time
func RefreshAccessToken(refreshToken string) (*TokenExchangeResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to create token refresh request: %v", err)
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to refresh access token: %v", err)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to read refresh response: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		helpers.LogMessage(helpers.ERROR, "Token refresh failed with status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("token refresh failed: status %d", resp.StatusCode)
	}

	var tokenResponse TokenExchangeResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to parse refresh response: %v", err)
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	helpers.LogMessage(helpers.INFO, "Successfully refreshed access token | Expires in: %d seconds", tokenResponse.ExpiresIn)

	return &tokenResponse, nil
}

// ValidateToken checks if a token is valid by making a test request to Spotify API
func ValidateToken(accessToken string) (bool, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create validation request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		helpers.LogMessage(helpers.INFO, "Token validation successful")
		return true, nil
	}

	helpers.LogMessage(helpers.WARNING, "Token validation failed with status: %d", resp.StatusCode)
	return false, nil
}

// RevokeToken revokes a Spotify access or refresh token
func RevokeToken(token string) error {
	data := url.Values{
		"token": {token},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token/revoke", strings.NewReader(data.Encode()))
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to create revoke request: %v", err)
		return fmt.Errorf("failed to create revoke request: %w", err)
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		helpers.LogMessage(helpers.ERROR, "Failed to revoke token: %v", err)
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		helpers.LogMessage(helpers.ERROR, "Token revocation failed with status %d: %s", resp.StatusCode, string(body))
		return fmt.Errorf("token revocation failed: status %d", resp.StatusCode)
	}

	helpers.LogMessage(helpers.INFO, "Token successfully revoked")
	return nil
}
