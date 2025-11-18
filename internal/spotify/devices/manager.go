package spotifyDevices

import (
	spotifyConfig "Quazaar/internal/spotify/config"
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"encoding/json"
	"net/http"
)

// getUserDevices fetches the list of available Spotify playback devices for the authenticated user
// This is an internal function used by the GetUserDevices HTTP handler
//
// Returns:
//   - any: JSON response from Spotify containing device information
//   - nil: If the request fails at any stage
//
// The response includes device details such as:
//   - Device ID, name, and type (Computer, Smartphone, Speaker, etc.)
//   - Active status and volume level
//   - Playback restrictions and capabilities
func getUserDevices() any {
	endpoint := spotifyConfig.SpotifyAPIBaseURL + "/me/player/devices"

	req, _ := http.NewRequest("GET", endpoint, nil)

	// Get valid access token (auto-refreshes if expired)
	accessToken, err := spotifyTokens.GetSpotifyAccessToken()
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Parse and return the raw device list response
	var devices any
	json.NewDecoder(resp.Body).Decode(&devices)

	return devices
}
