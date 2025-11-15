package main

import (
	"Quazaar/utils/auth"
	"Quazaar/utils/db"
	"Quazaar/utils/poller"
	"Quazaar/utils/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
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

	// Setup HTTP routes
	http.HandleFunc("/api/register", auth.HandleRegister)
	http.HandleFunc("/api/login", auth.HandleLogin)
	http.HandleFunc("/api/tokens/create", auth.HandleCreateToken)
	http.HandleFunc("/api/tokens/list", auth.HandleListTokens)
	http.HandleFunc("/api/tokens/revoke", auth.HandleRevokeToken)
	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/ws", websocket.Handle)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/auth", serveAuthPage)

	// Start media poller
	go poller.Handle()

	// Start cleanup goroutine (runs every hour)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			auth.CleanExpiredTokens()
		}
	}()

	// Start the server
	fmt.Println("")
	fmt.Println("📡 Starting server...")
	fmt.Println("")
	localAddr := os.Getenv("LOCAL_HOST_IP") + ":" + os.Getenv("LOCAL_HOST_PORT")
	if localAddr == ":" {
		localAddr = "127.0.0.1:8765"
	}
	fmt.Printf("🌐 Server running at http://%s\n", localAddr)
	fmt.Printf("📝 API Endpoints:\n")
	fmt.Printf("   POST   /api/register          - Register user\n")
	fmt.Printf("   POST   /api/login             - Login user\n")
	fmt.Printf("   POST   /api/tokens/create     - Create service token\n")
	fmt.Printf("   GET    /api/tokens/list       - List all tokens\n")
	fmt.Printf("   POST   /api/tokens/revoke     - Revoke a token\n")
	fmt.Printf("   WS     /ws                    - WebSocket endpoint\n")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")
	fmt.Println("")

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

func serveAuthPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "temp/web/auth.html")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","service":"quazaar","version":"1.0"}`)
	log.Println("✅ Health check")
}
