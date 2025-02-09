package middlewares

import (
	"context"
	"fmt"
	cstErrors "github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get("Authorization")
		if auth != "" {
			tokenString := strings.TrimPrefix(auth, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				jwtSecret, exists := os.LookupEnv("JWT_SECRET")
				if !exists {
					return nil, cstErrors.NoJWTSecretError
				}
				return []byte(jwtSecret), nil
			})

			if err == nil && token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				sub := claims["sub"].(string)
				ctx = context.WithValue(ctx, utils.UserIdCtxKey, sub)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
