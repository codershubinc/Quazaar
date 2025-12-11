package sidecar

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:embed QuazaarMedia.exe
var sidecarFS embed.FS

// ExtractSidecar extracts the embedded QuazaarMedia.exe to a temporary location
// and returns the path to the executable.
func ExtractSidecar() (string, error) {
	// Create a temp directory or use a specific one
	tempDir := os.TempDir()
	exePath := filepath.Join(tempDir, "QuazaarMedia.exe")

	// Check if file exists and is up to date?
	// For simplicity, we overwrite it every time or check size.
	// To avoid permission issues if it's running, we might need a unique name or handle errors.

	// Try to open the embedded file
	src, err := sidecarFS.Open("QuazaarMedia.exe")
	if err != nil {
		return "", fmt.Errorf("failed to open embedded sidecar: %w", err)
	}
	defer src.Close()

	// Create the destination file
	// If it exists and is running, this might fail on Windows.
	// We can try to remove it first.
	os.Remove(exePath)

	dst, err := os.OpenFile(exePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		// If we can't write to temp, try local directory?
		// Or maybe it's already running.
		return "", fmt.Errorf("failed to create sidecar file at %s: %w", exePath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy sidecar: %w", err)
	}

	return exePath, nil
}
