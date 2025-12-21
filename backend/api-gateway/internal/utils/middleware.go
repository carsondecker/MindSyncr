package utils

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/google/uuid"
)

type contextKey string

const UserContextKey contextKey = "user"

type MiddlewareHandler struct {
	cfg *Config
}

func NewMiddlewareHandler(cfg *Config) *MiddlewareHandler {
	return &MiddlewareHandler{
		cfg,
	}
}

func (h *MiddlewareHandler) GetConfig() *Config {
	return h.cfg
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			Error(w, http.StatusUnauthorized, ErrMissingAccessToken, "no access token cookie found")
			return
		}

		claims, err := GetClaimsFromToken(cookie.Value)
		if err != nil {
			Error(w, http.StatusUnauthorized, ErrInvalidAccessToken, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TODO: update error handling to give a different error if no rows are returned
func (h *MiddlewareHandler) CheckRoomMembership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*Claims)
		if !ok || claims == nil {
			Error(w, http.StatusUnauthorized, ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		roomIdStr, sErr := GetPathValue(r, "room_id")
		if sErr != nil {
			SError(w, sErr)
			return
		}
		roomId, err := uuid.FromBytes([]byte(roomIdStr))
		if err != nil {
			Error(w, http.StatusBadRequest, ErrBadRequest, "invalid room id")
			return
		}

		_, err = h.cfg.Queries.CheckRoomMembership(ctx, sqlc.CheckRoomMembershipParams{
			ID:     roomId,
			UserID: claims.UserId,
		})
		if err != nil {
			Error(w, http.StatusInternalServerError, ErrDbtxFail, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TODO: update error handling to give a different error if no rows are returned
func (h *MiddlewareHandler) CheckRoomOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*Claims)
		if !ok || claims == nil {
			Error(w, http.StatusUnauthorized, ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		roomId, sErr := GetUUIDPathValue(r, "room_id")
		if sErr != nil {
			SError(w, sErr)
			return
		}

		_, err := h.cfg.Queries.CheckRoomOwnership(ctx, sqlc.CheckRoomOwnershipParams{
			ID:      roomId,
			OwnerID: claims.UserId,
		})
		if err != nil {
			Error(w, http.StatusInternalServerError, ErrDbtxFail, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
