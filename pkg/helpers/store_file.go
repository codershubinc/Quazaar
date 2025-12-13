package helpers

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func StoreFile(filePath string, data multipart.File) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error getting user home directory:", err)
		return err
	}

	fileDir := filepath.Join(homeDir, "Downloads", "Quazaar")

	err = os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		log.Println("Error while creating directory ", fileDir, ":", err)
		return err
	}

	fullPath := filepath.Join(fileDir, filePath)
	dest, err := os.Create(fullPath)
	log.Println("Storing file at:", fullPath)
	if err != nil {
		log.Println("Error while creating file ", fullPath, ":", err)
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, data)
	if err != nil {
		return err
	}

	return nil
}
