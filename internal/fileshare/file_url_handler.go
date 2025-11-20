package fileShare

import (
	"Quazaar/internal/db"
	"Quazaar/pkg/helpers"
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
	storedDeviceId, err := db.GetFileShareDeviceToken(token)

	if err != nil {
		return false
	}
	if storedDeviceId != deviceId {
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
