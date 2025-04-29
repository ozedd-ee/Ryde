package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = os.Getenv("USER_SECRET_KEY")

type Claims struct {
	UserID string `json:"_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID string) (string, error) {
	if jwtSecret == "" {
		log.Fatal("Secret key not set")
	}
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	if jwtSecret == "" {
		log.Fatal("Secret key not set")
	}
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func (token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
