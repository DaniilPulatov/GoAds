package utils

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

// Put it into env variable or config file in production
var SecretKey = []byte("your_secret_key") // Replace it with your actual secret key

func GenerateToken(userID string, duration int) (string, error) {
	expAt := time.Duration(duration) * time.Minute
	// Create claims with user data
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),            // iat
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(expAt)), // exp
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		log.Printf("failed to sign token: %v", err)
		return "", errors.New("failed to sign token")
	}

	return signedToken, nil
}

func IsValidPhone(phone string) bool {
	// Регулярное выражение для проверки +998 и 9 цифр после
	pattern := `^\+998\d{9}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return false
	}
	return matched
}
