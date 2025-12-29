package spotifyTrack

import (
	"Quazaar/pkg/helpers"
	"log"
	"net/http"
)

func CurrentlyPlayingTrackApi(w http.ResponseWriter, r *http.Request) {
	data, err := CurrentlyPlayingTrack()
	if err != nil {
		log.Println("got err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.SendJsonDataToClient(w, 200, data)
}

func AddToLibraryApi(w http.ResponseWriter, r *http.Request) {
	err := AddToLibrary(nil)
	if err != nil {
		http.Error(w, "Failed to add track to library :: "+err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.SendJsonDataToClient(w, 200, map[string]string{"status": "Track added to library"})
}

func CheckUserHasTrackInLibrary(w http.ResponseWriter, r *http.Request) {
	is, err := CheckTrackInLibrary(nil)
	if err != nil {
		http.Error(w, "Failed to check track in library :: "+err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.SendJsonDataToClient(w, 200, map[string]bool{"is_in_library": is})
}
