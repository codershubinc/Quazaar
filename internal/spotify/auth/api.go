package spotifyAuth

import (
	spotifyTokens "Quazaar/internal/spotify/tokens"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUser handles HTTP requests to fetch the authenticated Spotify user's profile
// Requires a valid Spotify access token to be available
//
// Response: JSON representation of the Spotify user profile
// Status Codes:
//   - 200: User profile retrieved successfully
//   - 500: Failed to get access token or user profile
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get a valid access token (automatically refreshes if expired)
	tk, err := spotifyTokens.GetSpotifyAccessToken()
	if err != nil {
		http.Error(w, "Failed to get Spotify access token", http.StatusInternalServerError)
		return
	}

	// Fetch user profile from Spotify API
	user, err := GetUserProfile(tk)
	if err != nil {
		fmt.Println("Err ", err)
		http.Error(w, "Failed to get Spotify user profile", http.StatusInternalServerError)
		return
	}

	// Return user profile as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
