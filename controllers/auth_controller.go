package controllers

import (
	"encoding/json"
	"expense-app-backend/models"
	"expense-app-backend/utils"
	"net/http"

	"gorm.io/gorm"
)

func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request body",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		defer r.Body.Close()

		_, err := user.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		if err := db.Create(&user).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to create user",
			}

			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusCreated)
		user.Password = ""

		type UserResponse struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}

		response := struct {
			StatusCode int          `json:"status_code"`
			Message    string       `json:"message"`
			Data       UserResponse `json:"data"`
		}{
			StatusCode: http.StatusCreated,
			Message:    "User created successfully",
			Data: UserResponse{
				ID:        user.ID.String(),
				Name:      user.Name,
				Email:     user.Email,
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		}

		json.NewEncoder(w).Encode(response)
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var requestUser models.User
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request body",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		defer r.Body.Close()

		if err := db.Where("email = ?", requestUser.Email).First(&user).Error; err != nil {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		if err := user.CheckPassword(requestUser.Password); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid password",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse := struct {
				StatusCode int    `json:"status_code"`
				Message    string `json:"message"`
			}{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to generate token",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := struct {
			StatusCode int    `json:"status_code"`
			Message    string `json:"message"`
			Token      string `json:"token"`
		}{
			StatusCode: http.StatusOK,
			Message:    "Login successful",
			Token:      token,
		}
		json.NewEncoder(w).Encode(response)
	}
}
