package spotifyAuth

import (
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	tk, err := spotifyTokens.GetSpotifyAccessToken()
	if err != nil {
		http.Error(w, "Failed to get Spotify access token", http.StatusInternalServerError)
		return
	}

	user, err := GetUserProfile(tk)
	if err != nil {
		fmt.Println("Err", err)
		http.Error(w, "Failed to get Spotify user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
