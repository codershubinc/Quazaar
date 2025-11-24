package fileShare

import (
	"Quazaar/internal/db"
	"Quazaar/pkg/helpers"
	"log"
)

func CreateFileAcceptTempUri() (string, error) {
	dId, err := helpers.GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	token, err := helpers.GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	expiryTime := 3600 // seconds

	err = db.StoreFileShareDeviceToken(token, dId, expiryTime)
	if err != nil {
		return "", err
	}
	return "/api/v0.1/fileshare/acceptfile?deviceId=" + dId + "&token=" + token, nil
}

func ValidateFileAcceptTempUri(deviceId, token string) bool {
	storedDeviceId, storedToken, err := db.GetFileShareDeviceToken(token)
	log.Println("🔑 Validating Token - DeviceID:", deviceId, "Token:", token)
	log.Println("🔑 Stored DeviceID from DB:", storedDeviceId)

	if err != nil {
		log.Println("❌ Error validating token:", err)
		return false
	}
	if storedDeviceId != deviceId || storedToken != token {
		return false
	}
	return true
}

func DelFileAcceptTempUris(token string) bool {
	err := db.DeleteFileShareDeviceToken(token)
	if err != nil {
		return err == nil
	}
	return true
}
