package utils

import (
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = os.Getenv("SECRET_KEY") // For test purposes

type Claims struct {
	UserID string `json:"_id"`
	jwt.StandardClaims
}

func ValidateJWT(tokenString string) (*Claims, error) {
	if jwtSecret == "" {
		log.Fatal("Secret key not set")
	}
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
