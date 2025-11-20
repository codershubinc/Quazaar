package db

import (
	"fmt"
	"time"
)

func GetToken(tokenType string) (string, error) {
	var token string
	err := DB.QueryRow("SELECT token FROM tokens WHERE tokenType = ? ORDER BY id DESC LIMIT 1", tokenType).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func StoreToken(tokenType, tokenOf, token string, expiry int) error {
	_, err := DB.Exec("INSERT INTO tokens (tokenType,tokenOf, token,expiry) VALUES (?, ? , ?, ?)", tokenType, tokenOf, token, expiry)
	return err
}
func StoreFileShareDeviceToken(token, deviceId string, expiry int) error {
	_, err := DB.Exec("INSERT INTO file_share_device_tokens (token, deviceId, expiry) VALUES (?, ? , ?)", token, deviceId, expiry)
	return err
}

func GetFileShareDeviceToken(token string) (string, error) {
	var deviceId string
	var expiry int
	var createdAt time.Time
	err := DB.QueryRow("SELECT deviceId , expiry, createdAt FROM file_share_device_tokens WHERE token = ? LIMIT 1", token).Scan(&deviceId, &expiry, &createdAt)
	if err != nil {
		return "", err
	}
	if expiry > 0 {
		timeStored := time.Now().Add(-time.Duration(expiry) * time.Second)
		if timeStored.After(createdAt) {
			return "", fmt.Errorf("token expired")
		}
	}
	return deviceId, nil
}

func DeleteFileShareDeviceToken(token string) error {
	_, err := DB.Exec("DELETE FROM file_share_device_tokens WHERE token = ?", token)
	return err
}
