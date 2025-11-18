package spotifyDevices

import (
	"encoding/json"
	"net/http"
)

// GetUserDevices handles HTTP requests to fetch the user's available Spotify devices
// This includes all devices where the user has an active Spotify session
//
// Response: JSON array of Spotify devices with their status and capabilities
// Status Codes:
//   - 200: Devices retrieved successfully
//   - 405: Method not allowed (only GET is supported)
//   - 500: Failed to retrieve devices from Spotify API
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
