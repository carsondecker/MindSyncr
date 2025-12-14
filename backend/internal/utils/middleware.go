package utils

import (
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			Error(w, http.StatusUnauthorized, ErrMissingAccessToken, "no access token cookie found")
			return
		}

		claims, err := GetClaims(cookie.Value)
		if err != nil {
			Error(w, http.StatusUnauthorized, ErrInvalidAccessToken, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
