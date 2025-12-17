package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LocalHostPortDev string
	LocalHostIP      string
	SpotifyClientID  string
	SpotifySecret    string
}

var GlobalConfig Config

func Load() {
	// Load .env file, but don't fail if it doesn't exist (e.g. production env vars)
	_ = godotenv.Load()

	GlobalConfig = Config{
		LocalHostPortDev: getEnv("LOCAL_HOST_PORT_DEV", "8765"),
		LocalHostIP:      getEnv("LOCAL_HOST_IP", "0.0.0.0"),
		SpotifyClientID:  getEnv("SPOTIFY_CLIENT_ID", ""),
		SpotifySecret:    getEnv("SPOTIFY_CLIENT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
