package auth

import (
	"Quazaar/pkg/models"
	"encoding/json"
	"log"
	"net/http"
)

// HandleSignup handles HTTP POST request for user registration
func HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
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
	err := Signup(req.Username, req.Password, req.Username)
	if err != nil {
		log.Printf("❌ Signup failed: %v", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "User registered successfully",
		"username": req.Username,
	})

	log.Printf("✅ User registered: %s", req.Username)
}

// HandleLogin handles HTTP POST request for user login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := Login(req.Username, req.Password)
	if err != nil || token == "" {
		log.Printf("❌ Login failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Login successful",
		"username":  req.Username,
		"token":     token,
		"tokenType": "deviceId",
	})

	log.Printf("✅ User logged in: %s", req.Username)
}

// Signup wrapper function
func Signup(username, password, name string) error {
	return addUser(username, password, name)
}

// Login wrapper function
func Login(username, password string) (string, error) {
	isLogin, err := loginUser(username, password)
	if err != nil {
		return "", err
	}
	if !isLogin {
		return "", nil
	}
	return generateAuthToken(username)
}

func IsLoggedIn() bool {
	// Check if any user exists in the database
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		log.Printf("❌ Database error: %v", err)
		return false
	}
	return count > 0
}

// HandleRefreshToken handles HTTP POST request for token refresh
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate old token and get username
	username, err := validateToken(req.Token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Invalidate old token
	err = invalidateToken(req.Token)
	if err != nil {
		log.Printf("⚠️ Failed to invalidate old token: %v", err)
	}

	// Generate new token
	newToken, err := generateAuthToken(username)
	if err != nil {
		http.Error(w, "Failed to generate new token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Token refreshed successfully",
		"token":     newToken,
		"tokenType": "deviceId",
	})

	log.Printf("✅ Token refreshed for user: %s", username)
}

// HandleChangePassword handles HTTP POST request for password change
func HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate token and get username
	username, err := validateToken(req.Token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Validate new password
	if len(req.NewPassword) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Verify old password
	isValid, err := loginUser(username, req.OldPassword)
	if err != nil || !isValid {
		http.Error(w, "Invalid old password", http.StatusUnauthorized)
		return
	}

	// Update password
	err = updatePassword(username, req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})

	log.Printf("✅ Password changed for user: %s", username)
}

// HandleGetUserInfo handles HTTP GET request for user info
func HandleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from request
	token := extractTokenFromRequest(r)
	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	// Validate token and get username
	username, err := validateToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Get user info
	userInfo, err := getUserInfo(username)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    userInfo,
	})
}

// HandleLogout handles HTTP POST request for user logout
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate token
	username, err := validateToken(req.Token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Invalidate token
	err = invalidateToken(req.Token)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})

	log.Printf("✅ User logged out: %s", username)
}

// HandleGetTokens handles HTTP GET request to list user's active tokens
func HandleGetTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from request
	token := extractTokenFromRequest(r)
	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	// Validate token and get username
	username, err := validateToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Get all user tokens
	tokens, err := getUserTokens(username)
	if err != nil {
		http.Error(w, "Failed to get tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tokens":  tokens,
	})
}

// Helper function to extract token from request
func extractTokenFromRequest(r *http.Request) string {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Format: Bearer <token>
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:]
		}
	}

	// Check query parameter
	return r.URL.Query().Get("token")
}
