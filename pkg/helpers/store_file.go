package helpers

import (
	"io"
	"log"
	"mime/multipart"
	"os"
)

func StoreFile(filePath string, data multipart.File) error {
	fileDir := "/home/swap/Downloads/Quazaar/"

	err := os.Mkdir(fileDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal("Error while creating directory ", filePath, ":", err)
		return err
	}

	dest, err := os.Create(fileDir + filePath)
	log.Println("Storing file at:", fileDir+filePath)
	if err != nil {
		log.Print("Error while creating path ", fileDir+filePath, ":", err)
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, data)
	if err != nil {
		return err
	}

	return nil
}
