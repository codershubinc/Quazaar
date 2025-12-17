package fileShare

import (
	"Quazaar/internal/logger"
	"Quazaar/pkg/helpers"
	"io"
	"net/http"
)

func RequestTempFileShareAccept(w http.ResponseWriter, r *http.Request) {
	logger.Info("🔔 Temp File Share Accept Request Handler Invoked")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tempUri, err := CreateFileAcceptTempUri()
	if err != nil {
		http.Error(w, "Failed to create temp URI", http.StatusInternalServerError)
		return
	}
	helpers.SendJsonDataToClient(w, 200, map[string]any{
		"acceptUri": tempUri,
		"message":   "Temporary file share URI created successfully",
		"expiry":    3600,
		"time":      http.TimeFormat,
	})
}

func HandleTempFileShareAccept(w http.ResponseWriter, r *http.Request) {
	logger.Info("🔔 Temp File Share Accept Handler Invoked")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceId := r.URL.Query().Get("deviceId")
	token := r.URL.Query().Get("token")
	logger.Info("🔑 Temp File Share Accept Request", "deviceId", deviceId, "token", token)

	// Clean up the used token
	defer DelFileAcceptTempUris(token)

	if !ValidateFileAcceptTempUri(deviceId, token) {
		logger.Warn("❌ Invalid or expired token")
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Failed to get multipart reader", http.StatusBadRequest)
		return
	}

	fileFound := false
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to read part", http.StatusInternalServerError)
			return
		}

		if part.FormName() == "file" {
			filename := part.FileName()
			if filename == "" {
				continue
			}

			err = helpers.StoreFile(filename, part)
			if err != nil {
				logger.Error("Failed to store file", "error", err)
				http.Error(w, "Failed to store file", http.StatusInternalServerError)
				return
			}
			part.Close()
			fileFound = true
		}
	}

	if !fileFound {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
