package websocket

import (
	"Quazaar/internal/player"
	spotifyArtist "Quazaar/internal/spotify/artist"
	"Quazaar/internal/system"
	"Quazaar/pkg/models"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Handle(res http.ResponseWriter, req *http.Request) {
	conn, err := CreateWebSocketConnection(res, req)
	if err != nil {
		http.Error(res, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Create unique client with unbuffered channel (fresh messages only)
	client := &Client{
		Conn: conn,
		Send: make(chan models.ServerResponse), // Unbuffered - fresh messages only
		ID:   fmt.Sprintf("%s-%d", req.RemoteAddr, time.Now().UnixNano()),
	}

	// Register client
	RegisterClient(client)
	defer UnregisterClient(client.ID)

	// No read deadline - connection stays open indefinitely
	// Clients won't timeout due to inactivity

	// Send welcome message
	msg := models.ServerResponse{
		Message: "Welcome to the WebSocket server!",
	}
	if err := SendWebSocketMessage(msg, conn); err != nil {
		log.Printf("Failed to send welcome message to %s", client.ID)
		return
	}

	// Writer goroutine - sends messages to this specific client
	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)
		for response := range client.Send {
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("Error writing to client %s: %v", client.ID, err)
				// Stop reading on write error
				return
			}
		}
	}()

	// Reader goroutine - receives messages from client
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Client %s disconnected: %v", client.ID, err)
			break
		}

		log.Printf("📨 Received from %s: %+v", client.ID, msg)
		// type of accept msg should
		// type = [any of system , spotify etc add more category as needed]
		//
		typeOfMsg := msg["type"]
		switch typeOfMsg {
		case "system":
			res, err := system.HandleWebSocket(msg)
			if err != nil {
				continue
			}
			if err := conn.WriteJSON(res); err != nil {
				log.Printf("Error sending system response to client %s: %v", client.ID, err)
			}
			continue

		}
		// Handle Spotify artist messages
		if msgType, ok := msg["message"].(string); ok && msgType == "spotify_artist" {
			log.Printf("🎵 Processing Spotify artist message")

			action, _ := msg["action"].(string)
			artistId, _ := msg["artistId"].(string)

			response, err := spotifyArtist.HandleSpotifyArtistWsMessage(msgType, action, artistId)
			if err != nil {
				log.Printf("⚠️  Spotify artist request failed: %v", err)
			}
			// Send response (success or error)
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("Error sending Spotify artist response to client %s: %v", client.ID, err)
			}
			continue // Skip other handlers after processing Spotify message
		}

		// Handle player commands
		if commandVal, ok := msg["command"]; ok {
			command, ok := commandVal.(string)
			if !ok {
				log.Printf("⚠️  Invalid command format: %v", commandVal)
				continue
			}

			log.Printf("🎮 Processing command: %v", command)

			var err error
			var success bool = true

			switch command {
			case "player_toggle", "play-pause":
				success, err = player.LinuxDBusPlayer.PlayPause()
			case "next":
				success, err = player.LinuxDBusPlayer.Next()
			case "prev", "previous":
				success, err = player.LinuxDBusPlayer.Prev()
			case "seek_forward":
				success, err = player.LinuxDBusPlayer.SeekForward()
			case "seek_backward":
				success, err = player.LinuxDBusPlayer.SeekBackward()
			case "play":
				// Check status first to avoid toggling if already playing
				info, e := player.LinuxDBusPlayer.GetCurrentPlayerMetadata()
				if e == nil {
					if info.Status != "Playing" {
						success, err = player.LinuxDBusPlayer.PlayPause()
					}
				} else {
					err = e
				}
			case "pause":
				// Check status first to avoid toggling if already paused
				info, e := player.LinuxDBusPlayer.GetCurrentPlayerMetadata()
				if e == nil {
					if info.Status == "Playing" {
						success, err = player.LinuxDBusPlayer.PlayPause()
					}
				} else {
					err = e
				}
			default:
				err = fmt.Errorf("unknown command: %s", command)
			}

			if err != nil || !success {
				log.Printf("⚠️  Command failed: %v", err)
				// Send error response to client
				errorMsg := models.ServerResponse{
					Status:  "error",
					Message: "command_failed",
					Data: map[string]string{
						"error": fmt.Sprintf("%v", err),
					},
				}
				conn.WriteJSON(errorMsg)
			} else {
				// Send success response to client
				successMsg := models.ServerResponse{
					Status:  "success",
					Message: "command_executed",
					Data: map[string]string{
						"command": command,
					},
				}
				log.Printf("✅ Command executed successfully: %v", command)
				conn.WriteJSON(successMsg)
			}
		}
	}
	<-writerDone
}
