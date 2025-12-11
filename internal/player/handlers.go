package player

import (
	"Quazaar/internal/media"
	"Quazaar/pkg/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// HandleGetPlayerInfo handles GET /api/v0.1/player/info - gets current player info via D-Bus with fallback
func HandleGetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use platform player functions for Metadata
	pf := GetCurrentPlayerFunctions()
	info, err := pf.GetCurrentPlayerMetadata()
	if err != nil {
		log.Printf("❌ Failed to get player info: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "No player running or D-Bus unavailable",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"player":  info,
	})

	log.Printf("✅ Player info retrieved: %s", info.Title)
}

// HandleGetPlayerInfoDBus handles GET /api/v0.1/player/info/dbus - gets player info via D-Bus only
// Query parameter: player (optional) - specific MPRIS player service name
// Example: /api/v0.1/player/info/dbus?player=org.mpris.MediaPlayer2.spotify
func HandleGetPlayerInfoDBus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for player query parameter
	playerParam := r.URL.Query().Get("player")

	var info media.DBusMediaInfo
	var err error

	if playerParam != "" {
		// Get info for specific player
		info, err = media.GetPlayerInfoViaDBusForPlayer(playerParam)
		if err != nil {
			log.Printf("❌ D-Bus player info failed for %s: %v", playerParam, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Player not found or unavailable: %s", playerParam),
				"player":  playerParam,
			})
			return
		}
	} else {
		// Get info for first available player (original behavior)
		info, err = media.GetPlayerInfoViaDBus()
		if err != nil {
			log.Printf("❌ D-Bus player info failed: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "D-Bus not available or no player running",
			})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"player":  info,
	})

	log.Printf("✅ D-Bus player info retrieved for %s: %s", info.Player, info.Title)
}

// HandleGetActivePlayers handles GET /api/v0.1/player/list - lists all active players
func HandleGetActivePlayers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	players, err := media.GetAllActivePlayers()
	if err != nil {
		log.Printf("❌ Failed to get active players: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve active players",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"players": players,
		"count":   len(players),
	})

	log.Printf("✅ Listed %d active players", len(players))
}

// HandleGetMPRISPlayers handles GET /api/v0.1/player/mpris/list - lists MPRIS players via D-Bus
func HandleGetMPRISPlayers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	players, err := media.GetMPRISPlayers()
	if err != nil {
		log.Printf("❌ Failed to get MPRIS players: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve MPRIS players via D-Bus",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"players": players,
		"count":   len(players),
		"source":  "D-Bus MPRIS",
	})

	log.Printf("✅ Listed %d MPRIS players via D-Bus", len(players))
}

// HandlePlayPause handles POST /api/v0.1/player/play-pause - toggle play/pause
func HandlePlayPause(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use platform player functions for Play/Pause
	pf := GetCurrentPlayerFunctions()
	success, err := pf.PlayPause()
	if err != nil || !success {
		log.Printf("❌ Play/Pause failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to toggle play/pause",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Play/Pause toggled",
	})

	log.Printf("✅ Play/Pause toggled")
}

// HandleNext handles POST /api/v0.1/player/next - skip to next track
func HandleNext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use platform player functions for Next
	pf := GetCurrentPlayerFunctions()
	success, err := pf.Next()
	if err != nil || !success {
		log.Printf("❌ Next track failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to skip to next track",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Skipped to next track",
	})

	log.Printf("✅ Skipped to next track")
}

// HandlePrevious handles POST /api/v0.1/player/previous - skip to previous track
func HandlePrevious(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use platform player functions for Previous
	pf := GetCurrentPlayerFunctions()
	success, err := pf.Prev()
	if err != nil || !success {
		log.Printf("❌ Previous track failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to skip to previous track",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Skipped to previous track",
	})

	log.Printf("✅ Skipped to previous track")
}

// HandlePlay handles POST /api/v0.1/player/play - play
func HandlePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := helpers.SpawnProcess("playerctl", []string{"play"})
	if err != nil {
		log.Printf("❌ Play failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to play",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Playing",
	})

	log.Printf("✅ Playing")
}

// HandlePause handles POST /api/v0.1/player/pause - pause
func HandlePause(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := helpers.SpawnProcess("playerctl", []string{"pause"})
	if err != nil {
		log.Printf("❌ Pause failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to pause",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Paused",
	})

	log.Printf("✅ Paused")
}

// HandleGetPlayerInfoWindows handles GET /api/v0.1/player/info/windows - gets player info via Windows APIs
func HandleGetPlayerInfoWindows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info, err := media.GetPlayerInfoViaWindowsWithFallback()
	if err != nil {
		log.Printf("❌ Failed to get Windows player info: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "No player running or Windows APIs unavailable",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"player":  info,
	})

	log.Printf("✅ Windows player info retrieved: %s", info.Title)
}

// HandleGetWindowsActivePlayers handles GET /api/v0.1/player/windows/list - lists active Windows media players
func HandleGetWindowsActivePlayers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	players, err := media.GetWindowsActivePlayers()
	if err != nil {
		log.Printf("❌ Failed to get Windows active players: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve Windows active players",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"players": players,
		"count":   len(players),
		"source":  "Windows APIs",
	})

	log.Printf("✅ Listed %d Windows active players", len(players))
}

// HandleGetPlayerInfoByPlayer handles GET /api/v0.1/player/info?player=<player_name> - gets info for specific player
func HandleGetPlayerInfoByPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	playerParam := r.URL.Query().Get("player")
	if playerParam == "" {
		http.Error(w, "Missing player parameter", http.StatusBadRequest)
		return
	}

	// Try Windows-specific method first (for Windows builds)
	if runtime.GOOS == "windows" {
		info, err := media.GetWindowsMediaInfoByPlayer(playerParam)
		if err == nil && info.Title != "" {
			mediaInfo := media.ConvertWindowsMediaToMediaInfo(info)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"player":  mediaInfo,
			})
			log.Printf("✅ Player info retrieved for %s: %s", playerParam, mediaInfo.Title)
			return
		}
	}

	// Fallback to D-Bus method (for Linux builds)
	dbusInfo, err := media.GetPlayerInfoViaDBus()
	if err == nil && dbusInfo.Title != "" {
		mediaInfo := media.ConvertDBusMediaToMediaInfo(dbusInfo)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"player":  mediaInfo,
		})
		log.Printf("✅ Player info retrieved for %s: %s", playerParam, mediaInfo.Title)
		return
	}

	// If no player found
	log.Printf("❌ Player %s not found or not playing", playerParam)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   fmt.Sprintf("Player %s not found or not playing media", playerParam),
	})
}
