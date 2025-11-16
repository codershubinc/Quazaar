package auth

import (
	"Quazaar/models"
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
		"tokenType": "auth",
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
