package main

import (
	"Quazaar/internal/api"
	"Quazaar/internal/auth"
	"Quazaar/internal/banner"
	"Quazaar/internal/db"
	"Quazaar/internal/poller"
	"Quazaar/internal/spotify"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// Display startup banner
	// Options: banner.Variant1(), banner.Variant2(), banner.Variant3(), banner.Variant4()
	banner.Show() // Uses Variant1 by default

	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatal("❌ Failed to initialize database:", err)
	}
	defer db.CloseDB()
	// Initialize auth package with DB
	auth.InitDB(db.DB)

	// Check if user is logged in

	isLoggedIn := auth.IsLoggedIn()
	if !isLoggedIn {
		log.Println("⚠️  No user found. Please sign up first.")
	} else {
		log.Println("✅ User found. Proceeding to start the server.")
	}

	// Setup all API routes
	api.SetupRoutes()

	// Setup Spotify routes
	spotify.SetupRoutes()

	// Start media poller
	go poller.Handle()

	// Initialize Spotify integration
	spotify.Init()

	// Start the server
	fmt.Println("")
	fmt.Println("📡 Starting server...")
	fmt.Println("")
	localAddr := os.Getenv("LOCAL_HOST_IP") + ":" + os.Getenv("LOCAL_HOST_PORT")
	if localAddr == ":" {
		localAddr = "127.0.0.1:8765"
	}
	log.Println("✅ Server listening at http://" + localAddr)
	if err := http.ListenAndServe(localAddr, nil); err != nil {
		log.Fatal("❌ Server error:", err)
	}
}
