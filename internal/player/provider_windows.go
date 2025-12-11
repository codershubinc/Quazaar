//go:build windows

package player

import (
	"Quazaar/pkg/models"
	"bufio"
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

	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)
	// Fix for 'go run' which puts binary in temp folders
	if strings.Contains(exPath, "go-build") {
		exPath = "."
	}

	sidecarPath := filepath.Join(exPath, "QuazaarMedia.exe")

	// Check if file exists to give better error message
	if _, err := os.Stat(sidecarPath); os.IsNotExist(err) {
		return fmt.Errorf("QuazaarMedia.exe not found at %s", sidecarPath)
	}

	w.cmd = exec.Command(sidecarPath)

	w.stdin, err = w.cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := w.cmd.Start(); err != nil {
		return err
	}

	w.reader = bufio.NewScanner(stdout)

	// Wait for handshake with 2s timeout
	readyCh := make(chan bool)
	go func() {
		if w.reader.Scan() && strings.Contains(w.reader.Text(), "ready") {
			readyCh <- true
		}
	}()

	select {
	case <-readyCh:
		return nil
	case <-time.After(2 * time.Second):
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
		var info models.MediaInfo
		if err := json.Unmarshal([]byte(resp), &info); err != nil {
			return models.MediaInfo{}, err
		}
		return info, nil
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
