// Copyright 2025 Swapnil Ingle
//
// Licensed under the MIT License. See LICENSE file for details.

package main

import (
	"Quazaar/internal/auth"
	"Quazaar/internal/db"
	"Quazaar/internal/poller"
	"Quazaar/internal/websocket"
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

	fmt.Println("🚀 Hello Quazaar Server ...")

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
	// Setup HTTP routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", websocket.Handle)
	http.HandleFunc("/api/signup", auth.HandleSignup)
	http.HandleFunc("/api/login", auth.HandleLogin)

	// Serve static assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	// Start media poller
	go poller.Handle()

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
