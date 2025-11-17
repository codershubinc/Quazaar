package spotifyDevices

import (
	"encoding/json"
	"net/http"
)

func GetUserDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	devices := getUserDevices()
	if devices == nil {
		http.Error(w, "Failed to get Spotify user devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
