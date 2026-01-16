package utils

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
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
func (h *MiddlewareHandler) CheckRoomMembershipByRoomId(next http.Handler) http.Handler {
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

		_, err := h.cfg.Queries.CheckRoomMembershipByRoomId(ctx, sqlc.CheckRoomMembershipByRoomIdParams{
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
func (h *MiddlewareHandler) CheckRoomOwnershipByRoomId(next http.Handler) http.Handler {
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

		_, err := h.cfg.Queries.CheckRoomOwnershipByRoomId(ctx, sqlc.CheckRoomOwnershipByRoomIdParams{
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

// TODO: update error handling to give a different error if no rows are returned
func (h *MiddlewareHandler) CheckRoomMembershipBySessionId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*Claims)
		if !ok || claims == nil {
			Error(w, http.StatusUnauthorized, ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			SError(w, sErr)
			return
		}

		_, err := h.cfg.Queries.CheckRoomMembershipBySessionId(ctx, sqlc.CheckRoomMembershipBySessionIdParams{
			ID:     sessionId,
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
func (h *MiddlewareHandler) CheckRoomOwnershipBySessionId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*Claims)
		if !ok || claims == nil {
			Error(w, http.StatusUnauthorized, ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			SError(w, sErr)
			return
		}

		_, err := h.cfg.Queries.CheckRoomOwnershipBySessionId(ctx, sqlc.CheckRoomOwnershipBySessionIdParams{
			ID:      sessionId,
			OwnerID: claims.UserId,
		})
		if err != nil {
			Error(w, http.StatusInternalServerError, ErrDbtxFail, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TODO: update error handling to give a different error if no rows are returned
func (h *MiddlewareHandler) CheckSessionMembershipOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*Claims)
		if !ok || claims == nil {
			Error(w, http.StatusUnauthorized, ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			SError(w, sErr)
			return
		}

		_, err := h.cfg.Queries.CheckSessionMembershipOnly(ctx, sqlc.CheckSessionMembershipOnlyParams{
			ID:     sessionId,
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
func (h *MiddlewareHandler) CheckSessionMembership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*Claims)
		if !ok || claims == nil {
			Error(w, http.StatusUnauthorized, ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			SError(w, sErr)
			return
		}

		_, err := h.cfg.Queries.CheckSessionMembership(ctx, sqlc.CheckSessionMembershipParams{
			ID:      sessionId,
			OwnerID: claims.UserId,
		})
		if err != nil {
			Error(w, http.StatusInternalServerError, ErrDbtxFail, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckSessionActive(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, sErr := GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			SError(w, sErr)
			return
		}

		_, err := h.cfg.Queries.CheckSessionActive(r.Context(), sessionId)
		if err != nil {
			Error(w, http.StatusInternalServerError, ErrDbtxFail, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
