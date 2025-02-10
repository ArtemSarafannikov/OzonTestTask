package utils

import (
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

func GenerateNewID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), err
}

// GenerateJWT This function maded as a variable for tests
var GenerateJWT = func(userID string) (string, error) {
	JWTSecret, exist := os.LookupEnv("JWT_SECRET")
	if !exist {
		return "", cstErrors.NoJWTSecretError
	}
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

//func GenerateJWT(userID string) (string, error) {
//	JWTSecret, exist := os.LookupEnv("JWT_SECRET")
//	if !exist {
//		return "", cstErrors.NoJWTSecretError
//	}
//	claims := &jwt.RegisteredClaims{
//		Subject:   userID,
//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return token.SignedString([]byte(JWTSecret))
//}
