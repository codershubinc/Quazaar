package sidecar

import (
	"bufio"
	"embed"
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

// Embed the Windows helper executable into the Go binary.
// This allows the application to be distributed as a single file,
// with the helper being extracted at runtime.
//
//go:embed QuazaarMedia.exe
var sidecarFS embed.FS

// windowsSidecar manages the lifecycle of the external helper process.
// It handles starting the process, communicating via stdin/stdout,
// and ensuring thread-safe access.
type WindowsSidecar struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Reader *bufio.Scanner
	Mu     sync.Mutex
}

var WS = &WindowsSidecar{}

// Globle Sidecar Config
var SidecarConfig = struct {
	ExecutablePath string
}{}

// ExtractSidecar extracts the embedded QuazaarMedia.exe to a temporary location
// and returns the path to the executable.
// This is necessary because Windows cannot execute a binary directly from memory/embedding;
// it must exist as a file on disk.
func ExtractSidecar() (string, error) {
	// Determine the destination path for the executable.
	// Defaults to the system temp directory, but can be overridden via config.
	tempDir := os.TempDir()
	exePath := filepath.Join(tempDir, "QuazaarMedia.exe")
	if SidecarConfig.ExecutablePath != "" {
		exePath = SidecarConfig.ExecutablePath
	}

	// Open the embedded executable file from the virtual filesystem.
	src, err := sidecarFS.Open("QuazaarMedia.exe")
	if err != nil {
		return "", fmt.Errorf("failed to open embedded sidecar: %w", err)
	}
	defer src.Close()

	// Attempt to remove any existing instance of the file.
	// This ensures we are always running the version embedded in this build.
	// Note: On Windows, this will fail if the process is currently running.
	os.Remove(exePath)

	// Create the destination file on the host filesystem.
	dst, err := os.OpenFile(exePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create sidecar file at %s: %w", exePath, err)
	}
	defer dst.Close()

	// Copy the binary data from the embedded FS to the disk file.
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy sidecar: %w", err)
	}

	return exePath, nil
}

// StartSidecar initializes and starts the external helper process.
// It establishes stdin/stdout pipes for communication and performs a handshake
// to ensure the process is ready to receive commands.
func (w *WindowsSidecar) StartSidecar() error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	// Step 1: Extract the binary to disk
	sidecarPath, err := ExtractSidecar()
	if err != nil {
		w.Cmd = nil
		errReport := map[string]string{
			"component": "sidecar",
			"stage":     "extraction",
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		}
		// Report the crash/error for diagnostics
		defer HandleCrash(errReport)
		return fmt.Errorf("failed to extract sidecar: %w", err)
	}

	// Step 2: Prepare the command
	w.Cmd = exec.Command(sidecarPath)

	// Step 3: Set up pipes for IPC (Inter-Process Communication)
	w.Stdin, err = w.Cmd.StdinPipe()
	if err != nil {
		w.Cmd = nil
		errReport := map[string]string{
			"component": "sidecar",
			"stage":     "stdin_pipe",
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		}
		defer HandleCrash(errReport)
		return err
	}

	stdout, err := w.Cmd.StdoutPipe()
	if err != nil {
		w.Cmd = nil
		errReport := map[string]string{
			"component": "sidecar",
			"stage":     "stdout_pipe",
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		}
		defer HandleCrash(errReport)
		return err
	}

	// Step 4: Start the process
	if err := w.Cmd.Start(); err != nil {
		w.Cmd = nil
		errReport := map[string]string{
			"component": "sidecar",
			"stage":     "start_process",
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		}
		defer HandleCrash(errReport)
		return err
	}

	w.Reader = bufio.NewScanner(stdout)

	// Step 5: Perform Handshake
	// Wait for the sidecar to print "ready" to confirm it initialized successfully.
	// This prevents sending commands to a process that isn't fully started.
	readyCh := make(chan bool)
	go func() {
		if w.Reader.Scan() {
			line := w.Reader.Text()
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
		// If handshake fails, kill the process to avoid zombie processes
		w.Cmd.Process.Kill()
		w.Cmd = nil
		errReport := map[string]string{
			"component": "sidecar",
			"stage":     "handshake_timeout",
			"error":     "handshake timeout",
			"timestamp": time.Now().Format(time.RFC3339),
		}
		// handler will try to start the new sidecar process
		defer HandleCrash(errReport)
		return fmt.Errorf("handshake timeout")
	}
}

func (w *WindowsSidecar) Send(action string, args map[string]interface{}) (bool, error) {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	if w.Cmd == nil {
		return false, fmt.Errorf("sidecar offline")
	}

	payload := map[string]interface{}{"Action": action}
	for k, v := range args {
		payload[k] = v
	}

	jsonBytes, _ := json.Marshal(payload)
	// Write with newline
	w.Stdin.Write(append(jsonBytes, '\n'))

	if w.Reader.Scan() {
		resp := w.Reader.Text()
		// Loose check handles "status": "ok"
		if strings.Contains(resp, "ok") {
			return true, nil
		}
		return false, fmt.Errorf("sidecar error: %s", resp)
	}
	return false, fmt.Errorf("stream closed")
}
