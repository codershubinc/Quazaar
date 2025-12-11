//go:build windows

package player

import (
	"Quazaar/internal/sidecar"
	"Quazaar/pkg/models"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// windowsBackend manages the persistent connection to QuazaarMedia.exe
type windowsBackend struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	reader *bufio.Scanner
	mu     sync.Mutex
}

// Global instance
var wb = &windowsBackend{}

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
		if err := wb.startSidecar(); err != nil {
			log.Printf("❌ Windows Sidecar failed: %v", err)
		} else {
			log.Println("✅ Windows Sidecar Connected")
		}
	}()
}

// --- Connection Management ---

func (w *windowsBackend) startSidecar() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Extract the embedded sidecar
	sidecarPath, err := sidecar.ExtractSidecar()
	if err != nil {
		w.cmd = nil
		return fmt.Errorf("failed to extract sidecar: %w", err)
	}

	w.cmd = exec.Command(sidecarPath)

	w.stdin, err = w.cmd.StdinPipe()
	if err != nil {
		w.cmd = nil
		return err
	}

	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		w.cmd = nil
		return err
	}

	if err := w.cmd.Start(); err != nil {
		w.cmd = nil
		return err
	}

	w.reader = bufio.NewScanner(stdout)

	// Wait for handshake with 2s timeout
	readyCh := make(chan bool)
	go func() {
		if w.reader.Scan() {
			line := w.reader.Text()
			log.Printf("Sidecar initial response: %s", line)
			if strings.Contains(line, "ready") {
				readyCh <- true
			}
		}
	}()

	select {
	case <-readyCh:
		return nil
	case <-time.After(2 * time.Second):
		w.cmd.Process.Kill()
		w.cmd = nil
		return fmt.Errorf("handshake timeout")
	}
}

// Helper to send command and wait for "status":"ok"
func (w *windowsBackend) send(action string, args map[string]interface{}) (bool, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cmd == nil {
		return false, fmt.Errorf("sidecar offline")
	}

	payload := map[string]interface{}{"Action": action}
	for k, v := range args {
		payload[k] = v
	}

	jsonBytes, _ := json.Marshal(payload)
	// Write with newline
	w.stdin.Write(append(jsonBytes, '\n'))

	if w.reader.Scan() {
		resp := w.reader.Text()
		// Loose check handles "status": "ok" or "status":"ok"
		if strings.Contains(resp, "ok") {
			return true, nil
		}
		return false, fmt.Errorf("sidecar error: %s", resp)
	}
	return false, fmt.Errorf("stream closed")
}

// --- Implementation Functions ---

func playPauseWindows() (bool, error) {
	return wb.send("play_pause", nil)
}

func nextWindows() (bool, error) {
	return wb.send("next", nil)
}

func prevWindows() (bool, error) {
	return wb.send("prev", nil)
}

func seekBackwardWindows() (bool, error) {
	// C# sidecar needs to handle "seek_backward" or logic implies -10s
	return wb.send("seek_backward", nil)
}

func seekForwardWindows() (bool, error) {
	return wb.send("seek_forward", nil)
}

func seekToWindows(position int64) (bool, error) {
	return wb.send("seek", map[string]interface{}{"Position": position})
}

func setVolumeWindows(volume int) (bool, error) {
	return wb.send("volume", map[string]interface{}{"Level": volume})
}

func getCurrentPlayerMetadataWindows() (models.MediaInfo, error) {
	wb.mu.Lock()
	defer wb.mu.Unlock()

	if wb.cmd == nil {
		return models.MediaInfo{}, fmt.Errorf("offline")
	}

	wb.stdin.Write([]byte(`{"Action": "info"}` + "\n"))

	if wb.reader.Scan() {
		resp := wb.reader.Text()

		// Define a temporary struct to match the C# JSON output (which uses numbers)
		type windowsResponse struct {
			Title      string  `json:"Title"`
			Artist     string  `json:"Artist"`
			Album      string  `json:"Album"`
			Status     string  `json:"Status"`
			Position   float64 `json:"Position"`
			Duration   float64 `json:"Duration"`
			App        string  `json:"App"`
			ArtworkUri string  `json:"ArtworkUri"`
		}

		var wInfo windowsResponse
		if err := json.Unmarshal([]byte(resp), &wInfo); err != nil {
			return models.MediaInfo{}, err
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
			Artwork:  tempArtworkUriToBytesHandler(wInfo.ArtworkUri),
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
		return ""
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
