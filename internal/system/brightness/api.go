package systemBrightness

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleGetBrightness handles GET /api/v0.1/system/brightness
func HandleGetBrightness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	br, err := GetCurrent()
	if err != nil {
		log.Printf("❌ Failed to get brightness: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to get brightness",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"brightness": br,
	})
}
