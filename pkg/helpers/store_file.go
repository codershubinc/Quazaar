package helpers

import (
	"Quazaar/internal/logger"
	"io"
	"os"
	"path/filepath"
)

func StoreFile(filePath string, data io.Reader) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Error getting user home directory", "error", err)
		return err
	}

	fileDir := filepath.Join(homeDir, "Downloads", "Quazaar")

	err = os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		logger.Error("Error while creating directory", "dir", fileDir, "error", err)
		return err
	}

	fullPath := filepath.Join(fileDir, filePath)
	dest, err := os.Create(fullPath)
	logger.Info("Storing file at", "path", fullPath)
	if err != nil {
		logger.Error("Error while creating file", "path", fullPath, "error", err)
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, data)
	if err != nil {
		return err
	}

	return nil
}
