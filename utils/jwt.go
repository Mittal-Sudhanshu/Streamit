package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET []byte = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET)
}
