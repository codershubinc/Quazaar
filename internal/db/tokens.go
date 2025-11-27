package db

import (
	"log"
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
	log.Println("🔑 Storing File Share Device Token - DeviceID:", deviceId, "Token:", token, expiry)
	log.Println("🔑 StoreFileShareDeviceToken Exec Error:", err)
	return err
}

func GetFileShareDeviceToken(token string) (string, string, error) {
	var deviceId string
	var expiry int
	err := DB.QueryRow("SELECT deviceId , expiry FROM file_share_device_tokens WHERE token = ? LIMIT 1", token).Scan(&deviceId, &expiry)
	log.Println("🔑 Retrieved File Share Device Token - DeviceID:", deviceId, "Token:", token, "Expiry:", expiry)
	if err != nil {
		log.Println("🔑 GetFileShareDeviceToken Query Error:", err)
		return "", "", err
	}
	return deviceId, token, nil
}

func DeleteFileShareDeviceToken(token string) error {
	_, err := DB.Exec("DELETE FROM file_share_device_tokens WHERE token = ?", token)
	return err
}
