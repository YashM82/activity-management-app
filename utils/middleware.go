package utils

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *Claims {
	user, ok := ctx.Value(UserContextKey).(*Claims)
	if !ok {
		return nil
	}
	return user
}
