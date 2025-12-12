//go:build windows

package player

import (
	"Quazaar/internal/sidecar"
	"Quazaar/pkg/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Global instance from sidecar package
var ws = sidecar.WS

func init() {
	// 1. Initialize the handler with our functions (even if sidecar isn't ready yet)
	// These functions will check if 'wb.cmd' is nil internally.
	WindowsPlayerHandler = models.PlayerFunctions{
		PlayPause:                playPauseWindows,
		Next:                     nextWindows,
		Prev:                     prevWindows,
		SeekBackward:             seekBackwardWindows,
		SeekForward:              seekForwardWindows,
		SeekTo:                   seekToWindows,
		SetVolume:                setVolumeWindows,
		GetCurrentPlayerMetadata: getCurrentPlayerMetadataWindows,
		GetAllPlayers:            getAllPlayersWindows,
	}

	// 2. Assign to global service variable immediately
	CurrentPlayerFunc = WindowsPlayerHandler

	// 3. Start Sidecar in background
	go func() {
		if err := sidecar.WS.StartSidecar(); err != nil {
			log.Printf("❌ Windows Sidecar failed: %v", err)
		} else {
			log.Println("✅ Windows Sidecar Connected")
		}
	}()
}

func playPauseWindows() (bool, error) {
	return ws.Send("play_pause", nil)
}

func nextWindows() (bool, error) {
	return ws.Send("next", nil)
}

func prevWindows() (bool, error) {

	return ws.Send("prev", nil)
}

func seekBackwardWindows() (bool, error) {
	return ws.Send("seek_backward", nil)
}

func seekForwardWindows() (bool, error) {
	return ws.Send("seek_forward", nil)
}

func seekToWindows(position int64) (bool, error) {
	return ws.Send("seek", map[string]interface{}{"Position": position})
}

func setVolumeWindows(volume int) (bool, error) {
	return ws.Send("volume", map[string]interface{}{"Level": volume})
}

func getCurrentPlayerMetadataWindows() (models.MediaInfo, error) {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()

	if ws.Cmd == nil {
		return models.MediaInfo{}, fmt.Errorf("offline")
	}

	ws.Stdin.Write([]byte(`{"Action": "info"}` + "\n"))

	if ws.Reader.Scan() {
		resp := ws.Reader.Text()

		// Define a temporary struct to match the C# JSON output (which uses numbers)
		type windowsResponse struct {
			Title         string  `json:"Title"`
			Artist        string  `json:"Artist"`
			Album         string  `json:"Album"`
			Status        string  `json:"Status"`
			Position      float64 `json:"Position"`
			Duration      float64 `json:"Duration"`
			App           string  `json:"App"`
			ArtworkUri    string  `json:"ArtworkUri"`
			ArtworkBase64 string  `json:"ArtworkBase64"`
			Message       string  `json:"message"`
		}

		var wInfo windowsResponse
		if err := json.Unmarshal([]byte(resp), &wInfo); err != nil {
			return models.MediaInfo{}, err
		}

		if wInfo.Status == "error" {
			return models.MediaInfo{}, fmt.Errorf("sidecar error: %s  	 ", wInfo.Message, wInfo)
		}

		artwork := ""
		if wInfo.ArtworkBase64 != "" {
			artwork = "data:image/jpeg;base64," + wInfo.ArtworkBase64
		} else {
			artwork = tempArtworkUriToBytesHandler(wInfo.ArtworkUri)
		}

		// Convert to models.MediaInfo (which expects strings for Position/Length)
		// We convert ms (from Windows) to us (microseconds) to match playerctl behavior
		return models.MediaInfo{
			Title:    wInfo.Title,
			Artist:   wInfo.Artist,
			Album:    wInfo.Album,
			Status:   wInfo.Status,
			Position: fmt.Sprintf("%.0f", wInfo.Position*1000),
			Length:   fmt.Sprintf("%.0f", wInfo.Duration*1000),
			Player:   wInfo.App,
			Artwork:  artwork,
		}, nil
	}
	return models.MediaInfo{}, fmt.Errorf("stream closed")
}

func getAllPlayersWindows() ([]string, error) {
	// Windows Media Session Manager auto-selects the active source.
	// We return a static list for now as multi-session selection
	// requires advanced C# logic not present in basic sidecar.
	return []string{"System Media Transport"}, nil
}

// Exported variable populated by init()
var WindowsPlayerHandler models.PlayerFunctions

func initializeDefaultPlayer() models.PlayerFunctions {
	return WindowsPlayerHandler
}

func tempArtworkUriToBytesHandler(uri string) string {
	// On Windows, the sidecar provides a full file path for artwork.
	// We read the file and return its bytes.
	// fmt.Println("Got uri", uri)
	if uri == "" {
		return ""
	}
	// The uri is likely a full path like "C:\Users\..."
	bytes, err := os.ReadFile(uri)
	if err != nil {
		fallbackImage, _ := os.ReadFile(filepath.Join("assets", "img", "artwork-fallback.jpg"))

		return "data:image/" + "jpg" + ";base64," + base64.StdEncoding.EncodeToString(fallbackImage)
	}

	// Determine extension
	ext := strings.ToLower(filepath.Ext(uri))
	imageExtension := "jpeg" // default
	switch ext {
	case ".png":
		imageExtension = "png"
	case ".jpg", ".jpeg":
		imageExtension = "jpeg"
	case ".gif":
		imageExtension = "gif"
	case ".webp":
		imageExtension = "webp"
	case ".bmp":
		imageExtension = "bmp"
	}

	return "data:image/" + imageExtension + ";base64," + base64.StdEncoding.EncodeToString(bytes)
}
