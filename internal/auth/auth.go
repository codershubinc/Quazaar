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
		username, "auth", hashedToken, 0)
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
