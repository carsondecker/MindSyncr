package app

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/utils"
)

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrMissingAccessToken, "no access token cookie found")
			return
		}

		claims, err := utils.GetClaims(cookie.Value)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrInvalidAccessToken, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
