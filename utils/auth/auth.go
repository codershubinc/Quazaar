package auth

import (
	"Quazaar/utils/db"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenRevoked       = errors.New("token has been revoked")
	ErrUserNotFound       = errors.New("user not found")
)

// User represents the single user in the system
type User struct {
	ID           int
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

// Token represents a service token
type Token struct {
	ID        int
	Name      string     // Token name (e.g., "Mobile App", "Web Client")
	Token     string     // The actual token string
	Service   string     // Service name (e.g., "websocket", "api")
	ExpiresAt *time.Time // When it expires (nil = never)
	CreatedAt time.Time
	LastUsed  *time.Time
	Active    bool
}

// GenerateToken creates a cryptographically secure random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword verifies a password against a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RegisterUser creates the single user (call this once on first setup)
func RegisterUser(username, password string) (*User, error) {
	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Try to insert user with id=1 (only one allowed due to CHECK constraint)
	result, err := db.DB.Exec(
		"INSERT INTO users (id, username, password_hash) VALUES (1, ?, ?)",
		username, hashedPassword,
	)
	if err != nil {
		// User already exists
		if err.Error() == "UNIQUE constraint failed: users.username" || err.Error() == "UNIQUE constraint failed: users.id" {
			return nil, errors.New("user already registered - use login instead")
		}
		return nil, err
	}

	// Get user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	log.Printf("✅ User registered: %s", username)

	return &User{
		ID:           int(userID),
		Username:     username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}, nil
}

// LoginUser verifies username/password and returns user
func LoginUser(username, password string) (*User, error) {
	var user User
	err := db.DB.QueryRow(
		"SELECT id, username, password_hash, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	log.Printf("✅ User logged in: %s", username)
	return &user, nil
}

// CreateToken creates a new service token
func CreateToken(name, service string, duration *time.Duration) (*Token, error) {
	// Generate random token
	tokenString, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	if duration != nil {
		expTime := time.Now().Add(*duration)
		expiresAt = &expTime
	}

	// Insert token into database
	result, err := db.DB.Exec(
		"INSERT INTO tokens (name, token, service, expires_at, active) VALUES (?, ?, ?, ?, TRUE)",
		name, tokenString, service, expiresAt,
	)
	if err != nil {
		return nil, err
	}

	tokenID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	expireMsg := "never"
	if expiresAt != nil {
		expireMsg = expiresAt.Format("2025-01-02 15:04:05")
	}
	log.Printf("🔑 Token created: %s (service: %s, expires: %s)", name, service, expireMsg)

	return &Token{
		ID:        int(tokenID),
		Name:      name,
		Token:     tokenString,
		Service:   service,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Active:    true,
	}, nil
}

// ValidateToken checks if a token is valid
// Returns true if token exists, is active, and not expired
func ValidateToken(tokenString string) (bool, error) {
	var token Token
	var expiresAt *time.Time

	err := db.DB.QueryRow(
		"SELECT id, name, token, service, expires_at, active FROM tokens WHERE token = ?",
		tokenString,
	).Scan(&token.ID, &token.Name, &token.Token, &token.Service, &expiresAt, &token.Active)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, ErrInvalidToken
		}
		return false, err
	}

	// Check if token is active
	if !token.Active {
		return false, ErrTokenRevoked
	}

	// Check if token is expired
	if expiresAt != nil && time.Now().After(*expiresAt) {
		return false, ErrTokenExpired
	}

	// Update last used timestamp
	go func() {
		db.DB.Exec("UPDATE tokens SET last_used = datetime('now') WHERE token = ?", tokenString)
	}()

	return true, nil
}

// RevokeToken deactivates a token
func RevokeToken(tokenString string) error {
	result, err := db.DB.Exec(
		"UPDATE tokens SET active = FALSE WHERE token = ?",
		tokenString,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrInvalidToken
	}

	log.Printf("🔒 Token revoked: %s...", tokenString[:10])
	return nil
}

// RevokeAllTokens revokes all tokens (except those you want to keep)
func RevokeAllTokens() error {
	result, err := db.DB.Exec(
		"UPDATE tokens SET active = FALSE WHERE active = TRUE",
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("🔒 Revoked all %d tokens", rowsAffected)
	return nil
}

// GetAllTokens returns all tokens for the user
func GetAllTokens() ([]Token, error) {
	rows, err := db.DB.Query(
		"SELECT id, name, token, service, expires_at, created_at, last_used, active FROM tokens ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []Token
	for rows.Next() {
		var token Token
		var expiresAt, lastUsed *time.Time

		if err := rows.Scan(&token.ID, &token.Name, &token.Token, &token.Service, &expiresAt, &token.CreatedAt, &lastUsed, &token.Active); err != nil {
			return nil, err
		}

		token.ExpiresAt = expiresAt
		token.LastUsed = lastUsed
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// CleanExpiredTokens removes expired tokens from the database
func CleanExpiredTokens() (int64, error) {
	result, err := db.DB.Exec(
		"DELETE FROM tokens WHERE expires_at IS NOT NULL AND expires_at < datetime('now')",
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected > 0 {
		log.Printf("🧹 Cleaned %d expired tokens", rowsAffected)
	}
	return rowsAffected, nil
}
