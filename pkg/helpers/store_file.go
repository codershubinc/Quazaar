package helpers

import (
	"io"
	"mime/multipart"
	"os"
)

func StoreFile(filePath string, data multipart.File) error {
	dest, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, data)
	if err != nil {
		return err
	}

	return nil
}
