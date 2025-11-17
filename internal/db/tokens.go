package db

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
