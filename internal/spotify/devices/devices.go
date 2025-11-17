package spotifyDevices

import (
	"Quazaar/internal/spotify"
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"encoding/json"
	"net/http"
)

func getUserDevices() any {
	endpoint := spotify.SpotifyAPIBaseURL + "/me/player/devices"

	req, _ := http.NewRequest("GET", endpoint, nil)
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

	// parse response
	var devices any
	json.NewDecoder(resp.Body).Decode(&devices)

	return devices
}
