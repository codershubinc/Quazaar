package fileShare

import (
	"Quazaar/pkg/helpers"
	"log"
	"net/http"
)

func RequestTempFileShareAccept(w http.ResponseWriter, r *http.Request) {
	log.Println("🔔 Temp File Share Accept Request Handler Invoked")
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
	log.Println("🔔 Temp File Share Accept Handler Invoked")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceId := r.URL.Query().Get("deviceId")
	token := r.URL.Query().Get("token")
	log.Printf("🔑 Temp File Share Accept Request - DeviceID: %s, Token: %s", deviceId, token)

	// Clean up the used token
	defer DelFileAcceptTempUris(token)

	if !ValidateFileAcceptTempUri(deviceId, token) {
		log.Println("❌ Invalid or expired token")
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(100 << 20) // 100 MB max memory
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	err = helpers.StoreFile(header.Filename, file)
	if err != nil {
		http.Error(w, "Failed to store file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
