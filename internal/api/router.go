package api

import (
	"Quazaar/internal/auth"
	fileShare "Quazaar/internal/fileshare"
	"Quazaar/internal/middleware"
	"Quazaar/internal/player"
	spotifyArtist "Quazaar/internal/spotify/artist"
	spotifyAuth "Quazaar/internal/spotify/auth"
	"Quazaar/internal/system"
	"Quazaar/internal/websocket"
	"net/http"
)

// SetupRoutes configures all HTTP routes for the API
func SetupRoutes() {
	// Root
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/filesharetest", serveFileShareTestPage)

	// WebSocket
	http.HandleFunc("/ws", middleware.AuthenticationMiddleware(websocket.Handle))

	// Static assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	// API v0.1 - Authentication
	http.HandleFunc("/api/v0.1/signup", auth.HandleSignup)
	http.HandleFunc("/api/v0.1/login", auth.HandleLogin)
	http.HandleFunc("/api/v0.1/auth/refresh", auth.HandleRefreshToken)
	http.HandleFunc("/api/v0.1/auth/change-password", auth.HandleChangePassword)
	http.HandleFunc("/api/v0.1/auth/user", auth.HandleGetUserInfo)
	http.HandleFunc("/api/v0.1/auth/logout", auth.HandleLogout)
	http.HandleFunc("/api/v0.1/auth/tokens", auth.HandleGetTokens)

	// API v0.1 - Player Info
	http.HandleFunc("/api/v0.1/player/info", middleware.AuthenticationMiddleware(player.HandleGetPlayerInfo))
	http.HandleFunc("/api/v0.1/player/info/dbus", middleware.AuthenticationMiddleware(player.HandleGetPlayerInfoDBus))
	http.HandleFunc("/api/v0.1/player/info/windows", middleware.AuthenticationMiddleware(player.HandleGetPlayerInfoWindows))

	// API v0.1 - Player Lists
	http.HandleFunc("/api/v0.1/player/list", middleware.AuthenticationMiddleware(player.HandleGetActivePlayers))
	http.HandleFunc("/api/v0.1/player/mpris/list", middleware.AuthenticationMiddleware(player.HandleGetMPRISPlayers))
	http.HandleFunc("/api/v0.1/player/windows/list", middleware.AuthenticationMiddleware(player.HandleGetWindowsActivePlayers))

	// API v0.1 - Player Controls
	http.HandleFunc("/api/v0.1/player/play-pause", middleware.AuthenticationMiddleware(player.HandlePlayPause))
	http.HandleFunc("/api/v0.1/player/play", middleware.AuthenticationMiddleware(player.HandlePlay))
	http.HandleFunc("/api/v0.1/player/pause", middleware.AuthenticationMiddleware(player.HandlePause))
	http.HandleFunc("/api/v0.1/player/next", middleware.AuthenticationMiddleware(player.HandleNext))
	http.HandleFunc("/api/v0.1/player/previous", middleware.AuthenticationMiddleware(player.HandlePrevious))
	// API v0.1 - System Info
	http.HandleFunc("/api/v0.1/system/wifi", middleware.AuthenticationMiddleware(system.HandleGetWiFiInfo))
	http.HandleFunc("/api/v0.1/system/bluetooth", middleware.AuthenticationMiddleware(system.HandleGetBluetoothDevices))

	http.HandleFunc("/api/v0.1/spotify/artist", middleware.AuthenticationMiddleware(spotifyArtist.GetArtistInfo))
	http.HandleFunc("/api/v0.1/spotify/artist/follow", middleware.AuthenticationMiddleware(spotifyArtist.FollowArtist))
	http.HandleFunc("/api/v0.1/spotify/me", middleware.AuthenticationMiddleware(spotifyAuth.GetUser))

	// API v0.1 - File Share
	http.HandleFunc("/api/v0.1/fileshare/create-accept-uri", fileShare.RequestTempFileShareAccept)
	http.HandleFunc("/api/v0.1/fileshare/acceptfile", fileShare.HandleTempFileShareAccept)
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

func serveFileShareTestPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filesharetest" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "temp/fileshare_test.html")
}
