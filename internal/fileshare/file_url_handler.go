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

	err = db.StoreFileShareDeviceToken(token, dId, 0)
	if err != nil {
		return "", err
	}
	return "/filesharetest?deviceId=" + dId + "&token=" + token, nil
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
