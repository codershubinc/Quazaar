package api

import (
	"Quazaar/internal/auth"
	"Quazaar/internal/player"
	"Quazaar/internal/system"
	"Quazaar/internal/websocket"
	"net/http"
)

// SetupRoutes configures all HTTP routes for the API
func SetupRoutes() {
	// Root
	http.HandleFunc("/", serveHome)

	// WebSocket
	http.HandleFunc("/ws", websocket.Handle)

	// Static assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	// API v0.1 - Authentication
	http.HandleFunc("/api/v0.1/signup", auth.HandleSignup)
	http.HandleFunc("/api/v0.1/login", auth.HandleLogin)

	// API v0.1 - Player Info
	http.HandleFunc("/api/v0.1/player/info", player.HandleGetPlayerInfo)
	http.HandleFunc("/api/v0.1/player/info/dbus", player.HandleGetPlayerInfoDBus)
	http.HandleFunc("/api/v0.1/player/info/windows", player.HandleGetPlayerInfoWindows)

	// API v0.1 - Player Lists
	http.HandleFunc("/api/v0.1/player/list", player.HandleGetActivePlayers)
	http.HandleFunc("/api/v0.1/player/mpris/list", player.HandleGetMPRISPlayers)
	http.HandleFunc("/api/v0.1/player/windows/list", player.HandleGetWindowsActivePlayers)

	// API v0.1 - Player Controls
	http.HandleFunc("/api/v0.1/player/play-pause", player.HandlePlayPause)
	http.HandleFunc("/api/v0.1/player/play", player.HandlePlay)
	http.HandleFunc("/api/v0.1/player/pause", player.HandlePause)
	http.HandleFunc("/api/v0.1/player/next", player.HandleNext)
	http.HandleFunc("/api/v0.1/player/previous", player.HandlePrevious)

	// API v0.1 - System Info
	http.HandleFunc("/api/v0.1/system/wifi", system.HandleGetWiFiInfo)
	http.HandleFunc("/api/v0.1/system/bluetooth", system.HandleGetBluetoothDevices)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "temp/web/index.html")
}
