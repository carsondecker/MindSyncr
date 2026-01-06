package utils

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Sec-WebSocket-Protocol")
		if authHeader == "" {
			Error(w, http.StatusUnauthorized, ErrMissingAccessToken, "no ws access token found")
			return
		}

		parts := strings.Split(authHeader, ",")
		if len(parts) < 2 {
			Error(w, http.StatusUnauthorized, ErrMissingAccessToken, "missing ws access token")
			return
		}

		protocol := strings.TrimSpace(parts[0])
		ticket := strings.TrimSpace(parts[1])

		if protocol == "" || ticket == "" {
			Error(w, http.StatusUnauthorized, ErrMissingAccessToken, "missing ws access token")
			return
		}

		claims, err := GetClaimsFromToken(ticket)
		if err != nil {
			Error(w, http.StatusUnauthorized, ErrInvalidAccessToken, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
