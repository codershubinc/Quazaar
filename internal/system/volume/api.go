package systemVolume

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleGetVolume handles GET /api/v0.1/system/volume
func HandleGetVolume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vol, err := CurrentSystemVolume()
	if err != nil {
		log.Printf("❌ Failed to get volume: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to get volume",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"volume":  vol,
	})
}
