package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte("your_secret_key") // Use a secure key here

type JWTClaim struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID, email string) (string, error) {
	claims := &JWTClaim{
		UserID: userID,
		Email:  email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
