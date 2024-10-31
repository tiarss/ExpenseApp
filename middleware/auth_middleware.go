package middleware

import (
	"encoding/json"
	"expense-app-backend/utils"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		token := strings.Split(authHeader, "Bearer ")[1]
		if _, err := utils.ValidateToken(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		next.ServeHTTP(w, r)
	})
}
