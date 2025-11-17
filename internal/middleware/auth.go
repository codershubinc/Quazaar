package middleware

import (
	"Quazaar/internal/db"
	"fmt"
	"net/http"
)

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("URL ==>", r.URL.String())
		_, authToken := ExtractToken(r)
		valid, err := ValidateDeviceId(authToken)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	}
}

func ExtractToken(r *http.Request) (typeOfToken, token string) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return "Bearer", authHeader[7:]
	}
	fmt.Println("got deviceId ", r.URL.Query().Get("deviceId"))
	return "fromQuery", r.URL.Query().Get("deviceId")
}

func ValidateDeviceId(tokenString string) (bool, error) {
	DB := db.GetDB()
	err := DB.QueryRow("SELECT COUNT(1) FROM tokens WHERE token = ? AND tokenType = 'deviceId'", tokenString).Scan(new(int))
	if err != nil {
		return false, err
	}
	var count int
	err = DB.QueryRow("SELECT COUNT(1) FROM tokens WHERE token = ? AND tokenType = 'deviceId'", tokenString).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
