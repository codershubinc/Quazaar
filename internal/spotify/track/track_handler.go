package spotifyTrack

import (
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"Quazaar/pkg/helpers"
	"fmt"
	"log"
)

func CurrentlyPlayingTrack() (data any, err error) {
	tk, err := spotifyTokens.GetSpotifyAccessToken()

	res, err := helpers.SendRequest("https://api.spotify.com/v1/me/player/currently-playing", "GET", map[string]string{"Authorization": "Bearer " + tk}, nil)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	parsedData, err := helpers.GetParsedRequestData(res)
	if err != nil {
		return data, err
	}

	return parsedData, nil
}

func AddToLibrary(currentTrack any) error {
	tk, err := spotifyTokens.GetSpotifyAccessToken()
	if err != nil {
		return err
	}

	if currentTrack == nil {
		currentTrack, err = CurrentlyPlayingTrack()
	}
	if err != nil {
		return err
	}

	trackMap := currentTrack.(map[string]any)
	itemMap := trackMap["item"].(map[string]any)
	trackID := itemMap["id"].(string)

	// Use ISO 8601 format with UTC timezone (e.g., 2023-01-15T14:30:00Z)
	// Build body matching sample:
	// {
	//   "ids": ["id"],
	//   "timestamped_ids": [{"id":"id","added_at":"timestamp"}]
	// }

	type timestampedID struct {
		ID      string `json:"id"`
		AddedAt string `json:"added_at"`
	}

	body := map[string]any{
		"ids": []string{trackID},
	}

	req, err := helpers.SendRequest("https://api.spotify.com/v1/me/tracks", "PUT", map[string]string{"Authorization": "Bearer " + tk}, body)
	if err != nil {
		log.Println("err ", err)
		return err
	}
	defer req.Body.Close()
	reqData, _ := helpers.GetParsedRequestData(req)

	log.Println("\n\n\n=====+++++++++AddToLibrary response data:", reqData)
	if reqData == nil {
		return nil
	}

	if req.StatusCode != 200 {
		// try to extract "message" from parsed response, otherwise return generic status error
		if m, ok := reqData.(map[string]any); ok {
			if msg, ok := m["message"].(string); ok && msg != "" {
				return fmt.Errorf("spotify API error: %s", msg)
			}
		}
		return fmt.Errorf("spotify API returned status %d", req.StatusCode)
	}
	return nil

}

func CheckTrackInLibrary(currentTrack any) (bool, error) {
	var err error
	if currentTrack == nil {
		currentTrack, err = CurrentlyPlayingTrack()
	}
	if err != nil {
		return false, err
	}

	trackMap := currentTrack.(map[string]any)
	itemMap := trackMap["item"].(map[string]any)
	trackID := itemMap["id"].(string)
	tk, err := spotifyTokens.GetSpotifyAccessToken()
	if err != nil {
		return false, err
	}

	res, err := helpers.SendRequest("https://api.spotify.com/v1/me/tracks/contains?ids="+trackID, "GET", map[string]string{"Authorization": "Bearer " + tk}, nil)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	parsedData, err := helpers.GetParsedRequestData(res)
	if err != nil {
		return false, err
	}

	dataArray, ok := parsedData.([]any)
	if !ok || len(dataArray) == 0 {
		return false, fmt.Errorf("unexpected response format from Spotify API")
	}

	isInLibrary, ok := dataArray[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected data type for track presence in library")
	}

	if isInLibrary {
		return true, nil
	}
	return false, nil
}
