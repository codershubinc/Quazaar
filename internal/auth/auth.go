package auth

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB(db *sql.DB) {
	DB = db
}

func addUser(username, password, name string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	// Schema: id, name, pass
	_, err = DB.Exec(`INSERT INTO users (name, pass, username) VALUES (?, ?, ?)`,
		name, hashedPassword, username)
	if err != nil {
		log.Printf("❌ Error inserting user: %v", err)
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func loginUser(username, password string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("database not initialized")
	}

	var hashedPassword string
	// Schema: name is username column
	err := DB.QueryRow(`SELECT pass FROM users WHERE name = ?`, username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("❌ User not found: %v", err)
		return false, fmt.Errorf("user not found")
	}

	if !checkPasswordHash(password, hashedPassword) {
		return false, fmt.Errorf("invalid password")
	}

	return true, nil
}

func generateAuthToken(username string) (string, error) {
	if DB == nil {
		return "", fmt.Errorf("database not initialized")
	}

	// Generate token using username + timestamp
	tokenData := username + time.Now().String()
	hashedToken, err := hashPassword(tokenData)
	if err != nil {
		return "", err
	}

	// Schema: tokenOf, tokenType, token, expiry
	_, err = DB.Exec(`INSERT INTO tokens (tokenOf, tokenType, token, expiry) VALUES (?, ?, ?, ?)`,
		username, "deviceId", hashedToken, 0)
	if err != nil {
		log.Printf("❌ Error creating token: %v", err)
		return "", fmt.Errorf("failed to create token: %v", err)
	}
	return hashedToken, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// validateToken checks if token exists and returns the username
func validateToken(token string) (string, error) {
	if DB == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var username string
	err := DB.QueryRow(`SELECT tokenOf FROM tokens WHERE token = ?`, token).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	return username, nil
}

// invalidateToken removes a token from the database
func invalidateToken(token string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	result, err := DB.Exec(`DELETE FROM tokens WHERE token = ?`, token)
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

// updatePassword updates user's password
func updatePassword(username, newPassword string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`UPDATE users SET pass = ? WHERE name = ?`, hashedPassword, username)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

// getUserInfo retrieves user information (excluding password)
func getUserInfo(username string) (map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var id int
	var name string
	err := DB.QueryRow(`SELECT id, name FROM users WHERE name = ?`, username).Scan(&id, &name)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return map[string]interface{}{
		"id":       id,
		"username": name,
	}, nil
}

// getUserTokens retrieves all active tokens for a user
func getUserTokens(username string) ([]map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := DB.Query(`SELECT id, tokenType, expiry FROM tokens WHERE tokenOf = ?`, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %v", err)
	}
	defer rows.Close()

	var tokens []map[string]interface{}
	for rows.Next() {
		var id int
		var tokenType string
		var expiry int
		if err := rows.Scan(&id, &tokenType, &expiry); err != nil {
			continue
		}

		tokens = append(tokens, map[string]interface{}{
			"id":        id,
			"tokenType": tokenType,
			"expiry":    expiry,
		})
	}

	return tokens, nil
}
