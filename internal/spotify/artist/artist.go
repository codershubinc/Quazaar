package spotifyArtist

import (
	spotifyConfig "Quazaar/internal/spotify/config"
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"Quazaar/pkg/helpers"
	"fmt"
	"net/http"
)

func getArtistInfo(artistID string) (any, *http.Response, error) {
	base := spotifyConfig.SpotifyAPIBaseURL
	if base == "" {
		base = "https://api.spotify.com/v1"
	}

	url := base + "/artists/" + artistID
	fmt.Println("Spotify req uri ::", url)

	accessToken, err := spotifyTokens.GetSpotifyAccessToken()

	if err != nil {
		return nil, nil, err
	}

	req, err := helpers.SendRequest(url, "GET", map[string]string{
		"Authorization": "Bearer " + accessToken,
	}, nil)
	if err != nil {
		return nil, nil, err
	}
	defer req.Body.Close()

	data, err := helpers.GetParsedRequestData(req)
	if err != nil {
		return nil, nil, err
	}

	return data, req, nil
}

func followArtist(artistId string) (*http.Response, error) {

	if artistId == "" {
		return nil, nil
	}
	url := spotifyConfig.SpotifyAPIBaseURL + "/me/following?type=artist&ids=" + artistId

	accessToken, err := spotifyTokens.GetSpotifyAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := helpers.SendRequest(url, "PUT", map[string]string{
		"Authorization": "Bearer " + accessToken,
	}, nil)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if req.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to follow artist with status code: %d", req.StatusCode)
	}

	return req, nil
}
