package api

import (
	"Quazaar/internal/auth"
	fileShare "Quazaar/internal/fileshare"
	"Quazaar/internal/middleware"
	"Quazaar/internal/player"
	spotifyArtist "Quazaar/internal/spotify/artist"
	spotifyAuth "Quazaar/internal/spotify/auth"
	spotifyTrack "Quazaar/internal/spotify/track"
	"Quazaar/internal/system"
	systemBattery "Quazaar/internal/system/battery"
	systemBluetooth "Quazaar/internal/system/bluetooth"
	systemBrightness "Quazaar/internal/system/brightness"
	systemSound "Quazaar/internal/system/sound"
	systemUsage "Quazaar/internal/system/usage"
	systemVolume "Quazaar/internal/system/volume"
	systemWakaTime "Quazaar/internal/wakatime"
	"Quazaar/internal/websocket"
	"bytes"
	"embed"
	iofs "io/fs"
	"log"
	"net/http"
	"time"
)

var fs embed.FS

// SetupRoutes configures all HTTP routes for the API
func SetupRoutes(embedFS embed.FS) {
	fs = embedFS

	// Root
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/filesharetest", serveFileShareTestPage)

	// WebSocket
	http.HandleFunc("/ws", middleware.AuthenticationMiddleware(websocket.Handle))

	// Static assets
	assetsFS, err := iofs.Sub(fs, "assets")
	if err != nil {
		log.Println("Error creating assetsFS:", err)
	} else {
		http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))
	}

	// Web static files
	webFS, err := iofs.Sub(fs, "statics/web")
	if err != nil {
		log.Println("Error creating webFS:", err)
	} else {
		http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.FS(webFS))))
	}

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
	http.HandleFunc("/api/v0.1/system/info", middleware.AuthenticationMiddleware(system.HandleGetSystemInfo))
	http.HandleFunc("/api/v0.1/system/wifi", middleware.AuthenticationMiddleware(system.HandleGetWiFiInfo))
	http.HandleFunc("/api/v0.1/system/bluetooth", middleware.AuthenticationMiddleware(systemBluetooth.HandleGetBluetoothDevices))
	http.HandleFunc("/api/v0.1/system/volume", middleware.AuthenticationMiddleware(systemVolume.HandleGetVolume))
	http.HandleFunc("/api/v0.1/system/brightness", middleware.AuthenticationMiddleware(systemBrightness.HandleGetBrightness))
	http.HandleFunc("/api/v0.1/system/sound/devices", middleware.AuthenticationMiddleware(systemSound.HandleListDevices))
	http.HandleFunc("/api/v0.1/system/sound/device", middleware.AuthenticationMiddleware(systemSound.HandleSetDevice))

	// API v0.1 - Usage
	http.HandleFunc("/api/v0.1/system/usage", middleware.AuthenticationMiddleware(systemUsage.HandleGetSystemUsage))
	http.HandleFunc("/api/v0.1/system/storage", middleware.AuthenticationMiddleware(systemUsage.HandleGetStorageUsage))

	// API v0.1 - WakaTime
	http.HandleFunc("/api/v0.1/system/wakatime", middleware.AuthenticationMiddleware(systemWakaTime.HandleGetWakaTimeStats))

	// API v0.1 - Spotify
	http.HandleFunc("/api/v0.1/spotify/track/currently-playing", middleware.AuthenticationMiddleware(spotifyTrack.CurrentlyPlayingTrackApi))
	http.HandleFunc("/api/v0.1/spotify/track/add-to-library", middleware.AuthenticationMiddleware(spotifyTrack.AddToLibraryApi))
	http.HandleFunc("/api/v0.1/spotify/track/check-in-library", middleware.AuthenticationMiddleware(spotifyTrack.CheckUserHasTrackInLibrary))
	http.HandleFunc("/api/v0.1/spotify/artist", middleware.AuthenticationMiddleware(spotifyArtist.GetArtistInfo))
	http.HandleFunc("/api/v0.1/spotify/artist/follow", middleware.AuthenticationMiddleware(spotifyArtist.FollowArtist))
	http.HandleFunc("/api/v0.1/spotify/me", middleware.AuthenticationMiddleware(spotifyAuth.GetUser))
	// API v0.1 - File Share
	http.HandleFunc("/api/v0.1/fileshare/create-accept-uri", fileShare.RequestTempFileShareAccept)
	http.HandleFunc("/api/v0.1/fileshare/acceptfile", fileShare.HandleTempFileShareAccept)

	// API v0.1 - System  TODO: Move to systemBattery package
	http.HandleFunc("/api/v0.1/system/battery", systemBattery.GetBatteryInfoApi)
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
	data, err := fs.ReadFile("statics/web/index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	reader := bytes.NewReader(data)
	http.ServeContent(w, r, "index.html", time.Time{}, reader)
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
	data, err := fs.ReadFile("statics/web/fileshare_test.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	reader := bytes.NewReader(data)
	http.ServeContent(w, r, "fileshare_test.html", time.Time{}, reader)
}
