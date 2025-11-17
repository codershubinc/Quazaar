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

var spotifyClientRefreshToken string
var spotifyClientAccessToken string
var spotifyAccessTokenExpiry time.Time

func GetSpotifyRefreshToken() (string, error) {
	if spotifyClientRefreshToken != "" {
		return spotifyClientRefreshToken, nil
	}
	token, err := db.GetToken("spotifyClientRefreshToken")
	if err != nil {
		return "", err
	}
	spotifyClientRefreshToken = token
	return spotifyClientRefreshToken, nil
}

func GetSpotifyAccessToken() (token string, err error) {
	if spotifyClientAccessToken != "" && time.Now().Before(spotifyAccessTokenExpiry) {
		return spotifyClientAccessToken, nil
	}
	refreshToken, err := GetSpotifyRefreshToken()
	if err != nil {
		return "", err
	}

	token, err = refreshSpotifyAccessToken(refreshToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

func refreshSpotifyAccessToken(refreshToken string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// Basic Auth: Base64(client_id:client_secret)
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
	spotifyAccessTokenExpiry = time.Now().Add(time.Duration(result["expires_in"].(float64)) * time.Second)
	spotifyClientAccessToken = result["access_token"].(string)

	return spotifyClientAccessToken, nil
}

func SetSpotifyRefreshToken(newRefreshToken string, expiry int) error {
	err := db.StoreToken("spotifyClientRefreshToken", "spotify", newRefreshToken, expiry)
	if err != nil {
		return err
	}
	spotifyClientRefreshToken = newRefreshToken
	return nil
}
