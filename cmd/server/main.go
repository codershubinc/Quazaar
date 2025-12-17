package main

import (
	"Quazaar/internal/api"
	"Quazaar/internal/assets"
	"Quazaar/internal/auth"
	"Quazaar/internal/banner"
	"Quazaar/internal/config"
	"Quazaar/internal/db"
	"Quazaar/internal/logger"
	"Quazaar/internal/poller"
	"Quazaar/internal/spotify"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	config.Load()

	//temp:
	logger.Info("OS type", "os", runtime.GOOS)
	// Display startup banner
	// Options: banner.Variant1(), banner.Variant2(), banner.Variant3(), banner.Variant4()
	banner.Show()

	// Initialize database
	if err := db.Init(); err != nil {
		logger.Error("❌ Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.CloseDB()
	// Initialize auth package with DB
	auth.InitDB(db.DB)

	// Check if user is logged in

	isLoggedIn := auth.IsLoggedIn()
	if !isLoggedIn {
		logger.Warn("⚠️  No user found. Please sign up first.")
	} else {
		logger.Info("✅ User found. Proceeding to start the server.")
	}

	// Setup all API routes
	api.SetupRoutes(assets.Assets)

	// Setup Spotify routes
	spotify.SetupRoutes()

	// Create a context that listens for signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start media poller
	go poller.Handle(ctx)

	// Initialize Spotify integration
	spotify.Init()

	port := config.GlobalConfig.LocalHostPortDev
	ip := config.GlobalConfig.LocalHostIP

	localAddr := ip + ":" + port
	srv := &http.Server{
		Addr:    localAddr,
		Handler: nil, // uses DefaultServeMux
	}

	// Start the server in a goroutine
	go func() {
		fmt.Println("")
		fmt.Println("📡 Starting server...")
		fmt.Println("")
		logger.Info("✅ Server listening", "url", "http://"+localAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("❌ Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("🛑 Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Cancel the main context to stop the poller and other background tasks
	cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("❌ Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("👋 Server exiting")
}
