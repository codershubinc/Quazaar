package spotifyArtist

import (
	"Quazaar/pkg/helpers"
	"fmt"
	"net/http"
)

func GetArtistInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Log :: GetArtistInfo request got")
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		http.Error(w, "Artist ID is required", http.StatusBadRequest)
		return
	}

	data, res, err := getArtistInfo(artistID)
	if err != nil {
		http.Error(w, "Failed to get artist info", http.StatusInternalServerError)
		return
	}
	helpers.SendJsonDataToClient(w, res, data)

	defer r.Body.Close()
}

func FollowArtist(w http.ResponseWriter, r *http.Request) {
	artistId := r.URL.Query().Get("id")
	if artistId == "" {
		http.Error(w, "Artist ID is required", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	req, err := followArtist(artistId)
	if err != nil {
		http.Error(w, "Failed to follow artist", http.StatusInternalServerError)
		return
	}

	if req != nil && req.StatusCode != http.StatusNoContent {
		http.Error(w, "Failed to follow artist", req.StatusCode)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
