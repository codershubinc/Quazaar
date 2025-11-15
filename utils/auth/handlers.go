package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// RegisterRequest is the payload for user registration
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

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// HandleRegister handles user registration (first time setup)
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if len(req.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Create user
	user, err := RegisterUser(req.Username, req.Password)
	if err != nil {
		log.Printf("❌ Registration failed: %v", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "User registered successfully",
		"user_id":  user.ID,
		"username": user.Username,
	})

	log.Printf("✅ User registered via API: %s", user.Username)
}

// HandleLogin handles user login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := LoginUser(req.Username, req.Password)
	if err != nil {
		log.Printf("❌ Login failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create an auth token for this session
	duration := 24 * 7 * time.Hour // 7 days
	authToken, err := CreateToken("auth-"+req.Username, "auth", &duration)
	if err != nil {
		log.Printf("❌ Failed to create auth token: %v", err)
		http.Error(w, "Login succeeded but failed to create auth token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"message":    "Login successful",
		"user_id":    user.ID,
		"username":   user.Username,
		"token":      authToken.Token,
		"auth_token": authToken.Token,
	})

	log.Printf("✅ User logged in via API: %s (auth token created)", user.Username)
}

// HandleCreateToken handles token creation for services
func HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	token := r.Header.Get("Authorization")
	if token == "" {
		// Log request details for debugging
		log.Printf("❌ Create token attempt - No Authorization header")
		log.Printf("   Headers: %v", r.Header)
		http.Error(w, "Unauthorized - Authorization header required", http.StatusUnauthorized)
		return
	}

	log.Printf("🔍 Create token attempt - Auth header: %s...", token[:min(20, len(token))])

	// Validate the token
	valid, err := ValidateToken(token)
	if !valid || err != nil {
		log.Printf("❌ Token validation failed: %v (valid: %v)", err, valid)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" {
		http.Error(w, "Token name is required", http.StatusBadRequest)
		return
	}

	// Parse duration
	var duration *time.Duration
	if req.Duration > 0 {
		d := time.Duration(req.Duration) * time.Hour
		duration = &d
	}

	// Create token
	newToken, err := CreateToken(req.Name, req.Service, duration)
	if err != nil {
		log.Printf("❌ Failed to create token: %v", err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TokenResponse{
		ID:        newToken.ID,
		Name:      newToken.Name,
		Token:     newToken.Token,
		Service:   newToken.Service,
		ExpiresAt: newToken.ExpiresAt,
		CreatedAt: newToken.CreatedAt,
		Active:    newToken.Active,
	})

	log.Printf("✅ Token created via API: %s", req.Name)
}

// HandleListTokens lists all active tokens
func HandleListTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate the token
	valid, err := ValidateToken(token)
	if !valid || err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Get all tokens
	tokens, err := GetAllTokens()
	if err != nil {
		log.Printf("❌ Failed to list tokens: %v", err)
		http.Error(w, "Failed to list tokens", http.StatusInternalServerError)
		return
	}

	var response []TokenResponse
	for _, t := range tokens {
		response = append(response, TokenResponse{
			ID:        t.ID,
			Name:      t.Name,
			Token:     t.Token,
			Service:   t.Service,
			ExpiresAt: t.ExpiresAt,
			CreatedAt: t.CreatedAt,
			LastUsed:  t.LastUsed,
			Active:    t.Active,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tokens":  response,
		"count":   len(response),
	})

	log.Printf("✅ Listed %d tokens via API", len(response))
}

// HandleRevokeToken revokes a specific token
func HandleRevokeToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate the token
	valid, err := ValidateToken(token)
	if !valid || err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Revoke token
	if err := RevokeToken(req.Token); err != nil {
		http.Error(w, "Failed to revoke token", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token revoked successfully",
	})

	log.Printf("✅ Token revoked via API")
}
