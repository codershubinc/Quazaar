package spotifyArtist

import (
	"Quazaar/pkg/models"
	"errors"
)

// HandleSpotifyArtistWsMessage handles Spotify artist WebSocket messages
// Returns ServerResponse to send back to client
func HandleSpotifyArtistWsMessage(msg string, action string, artistId string) (models.ServerResponse, error) {
	if msg != "spotify_artist" {
		return models.ServerResponse{
			Status:  "error",
			Message: "invalid_message_type",
			Data:    map[string]string{"error": "invalid message type"},
		}, errors.New("invalid message type")
	}

	switch action {
	case "get":
		return getArtistInfoWs(artistId)
	case "follow":
		return followArtistWs(artistId)
	default:
		return models.ServerResponse{
			Status:  "error",
			Message: "invalid_action",
			Data:    map[string]string{"error": "invalid action"},
		}, errors.New("invalid action")
	}
}

func followArtistWs(artistId string) (models.ServerResponse, error) {
	_, err := followArtist(artistId)
	if err != nil {
		return models.ServerResponse{
			Status:  "error",
			Message: "spotify_artist_follow",
			Data:    map[string]string{"error": err.Error()},
		}, err
	}

	return models.ServerResponse{
		Status:  "success",
		Message: "spotify_artist_follow",
		Data:    nil,
	}, nil
}

func getArtistInfoWs(artistID string) (models.ServerResponse, error) {
	data, _, err := getArtistInfo(artistID)
	if err != nil {
		return models.ServerResponse{
			Status:  "error",
			Message: "spotify_artist_get",
			Data:    map[string]string{"error": err.Error()},
		}, err
	}

	return models.ServerResponse{
		Status:  "success",
		Message: "spotify_artist_get",
		Data:    data,
	}, nil
}
