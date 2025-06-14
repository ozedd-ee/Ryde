package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = os.Getenv("SECRET_KEY") // For test purposes

type Claims struct {
	UserID string `json:"_id"`
	Email  string `json:"email"`
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

// To be validated in payment service
func GeneratePaymentJWT(userID, email string) (string, error) {
	paymentSecret := os.Getenv("PAYMENT_SECRET")
	if paymentSecret == "" {
		log.Fatal("payment shared secret not set")
	}
	expirationTime := time.Now().Add(time.Minute * 30)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(paymentSecret)
}
