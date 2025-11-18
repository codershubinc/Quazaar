package spotify

import (
	"Quazaar/internal/middleware"
	spotifyAuth "Quazaar/internal/spotify/auth"
	spotifyDevices "Quazaar/internal/spotify/devices"
	"net/http"
)

// SetupRoutes configures all Spotify-specific HTTP routes
// This includes OAuth flow, user profile, and device management endpoints
func SetupRoutes() {
	// Spotify OAuth Authentication Flow
	http.HandleFunc("/api/v0.1/spotify/login", middleware.AuthenticationMiddleware(spotifyAuth.Login))
	http.HandleFunc("/api/v0.1/spotify/callback", spotifyAuth.Callback)

	// Spotify User & Device Management
	http.HandleFunc("/api/v0.1/spotify/user", middleware.AuthenticationMiddleware(spotifyAuth.GetUser))
	http.HandleFunc("/api/v0.1/spotify/devices", middleware.AuthenticationMiddleware(spotifyDevices.GetUserDevices))
}
