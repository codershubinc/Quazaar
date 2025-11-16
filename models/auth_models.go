package models

import "time"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is the payload for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenRequest is the payload for creating a token
type TokenRequest struct {
	Name     string `json:"name"`           // Token name (e.g., "Mobile App")
	Service  string `json:"service"`        // Service name (e.g., "websocket")
	Duration int    `json:"duration_hours"` // Duration in hours (0 = never expires)
}

// TokenResponse is the response when creating/listing tokens
type TokenResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Token     string     `json:"token"`
	Service   string     `json:"service"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	LastUsed  *time.Time `json:"last_used"`
	Active    bool       `json:"active"`
}
