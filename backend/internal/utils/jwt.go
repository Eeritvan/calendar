package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId string) (string, error) {
	JWTkey := os.Getenv("JWT_KEY")
	secretKey := []byte(JWTkey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	returnToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT")
	}
	return returnToken, nil
}
