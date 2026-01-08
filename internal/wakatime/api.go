package wakatime

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleGetWakaTimeStats handles GET /api/v0.1/system/wakatime
func HandleGetWakaTimeStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := GetWakaTimeStats()
	if err != nil {
		log.Printf("❌ Failed to get wakatime stats: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable) // Or 500
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to get wakatime stats",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"stats":   stats,
	})
}
