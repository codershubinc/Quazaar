// Test program to simulate Windows media info on Linux
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("🖥️  Windows Media Info Test (Running on", runtime.GOOS, ")")
	fmt.Println("")

	// Simulate Windows media detection
	fmt.Println("🎵 Simulating Windows media player detection...")

	// Mock Windows media players
	players := []string{
		"Spotify.exe",
		"Music.UI.exe",
		"wmplayer.exe",
		"chrome.exe",
	}

	fmt.Println("📋 Available Windows media players:")
	for i, player := range players {
		fmt.Printf("   %d. %s\n", i+1, player)
	}

	fmt.Println("")
	fmt.Println("🎯 Testing media info retrieval...")

	// Mock media info
	mockInfo := map[string]interface{}{
		"title":    "Mock Song Title",
		"artist":   "Mock Artist Name",
		"album":    "Mock Album Name",
		"position": "2:34",
		"length":   "4:12",
		"status":   "Playing",
		"player":   "Spotify.exe",
	}

	fmt.Println("📊 Current track info:")
	for key, value := range mockInfo {
		fmt.Printf("   %-10s: %v\n", key, value)
	}

	fmt.Println("")
	fmt.Println("🎮 Testing playback controls...")
	controls := []string{"Play", "Pause", "Next", "Previous", "Stop"}

	for _, control := range controls {
		fmt.Printf("   ✅ %s command sent\n", control)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("")
	fmt.Println("✅ Windows media info test completed!")
	fmt.Println("💡 To test on real Windows: copy quazaar-windows.exe to a Windows machine")
}
