package sutils

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/utils"
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
			utils.Error(w, http.StatusUnauthorized, utils.ErrMissingAccessToken, "no access token cookie found")
			return
		}

		claims, err := utils.GetClaimsFromToken(cookie.Value)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrInvalidAccessToken, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *MiddlewareHandler) CheckRoomMembershipByRoomId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckRoomMembershipByRoomId(ctx, sqlc.CheckRoomMembershipByRoomIdParams{
			ID:     roomId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be a member of this room to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckRoomOwnershipByRoomId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckRoomOwnershipByRoomId(ctx, sqlc.CheckRoomOwnershipByRoomIdParams{
			ID:      roomId,
			OwnerID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the owner of this room to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckRoomMembershipBySessionId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckRoomMembershipBySessionId(ctx, sqlc.CheckRoomMembershipBySessionIdParams{
			ID:     sessionId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the owner of this session's room to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckRoomOwnershipBySessionId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckRoomOwnershipBySessionId(ctx, sqlc.CheckRoomOwnershipBySessionIdParams{
			ID:      sessionId,
			OwnerID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the owner of this session's room to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckSessionMembershipOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckSessionMembershipOnly(ctx, sqlc.CheckSessionMembershipOnlyParams{
			ID:     sessionId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must only be a member of this session to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckSessionMembership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		_, err := h.cfg.Queries.CheckSessionMembership(ctx, sqlc.CheckSessionMembershipParams{
			ID:      sessionId,
			OwnerID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be a member of this session to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckSessionActive(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckSessionActive(r.Context(), sessionId)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "the session must be active to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckQuestionBelongsToSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckQuestionBelongsToSession(r.Context(), sqlc.CheckQuestionBelongsToSessionParams{
			ID:        questionId,
			SessionID: sessionId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "this question does not belong to the provided session")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckCanDeleteQuestion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckCanDeleteQuestion(r.Context(), sqlc.CheckCanDeleteQuestionParams{
			ID:     questionId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the owner of this session or the writer of this question to delete it")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckCanDeleteQuestionLike(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckCanDeleteQuestionLike(r.Context(), sqlc.CheckCanDeleteQuestionLikeParams{
			UserID:     claims.UserId,
			QuestionID: questionId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the user who liked this question to delete it")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckOwnsQuestion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckOwnsQuestion(r.Context(), sqlc.CheckOwnsQuestionParams{
			ID:     questionId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the owner of this question to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckDoesNotOwnQuestion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckOwnsQuestion(r.Context(), sqlc.CheckOwnsQuestionParams{
			ID:     questionId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "the owner of this question cannot perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckReplyBelongsToSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		replyId, sErr := utils.GetUUIDPathValue(r, "reply_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckReplyBelongsToSession(r.Context(), sqlc.CheckReplyBelongsToSessionParams{
			ID:        replyId,
			SessionID: sessionId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "this reply does not belong to the provided session")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckOwnsReply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		replyId, sErr := utils.GetUUIDPathValue(r, "reply_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckOwnsReply(r.Context(), sqlc.CheckOwnsReplyParams{
			ID:     replyId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the owner of this reply to perform this action")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) CheckCanDeleteReply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		raw := ctx.Value(UserContextKey)
		claims, ok := raw.(*utils.Claims)
		if !ok || claims == nil {
			utils.Error(w, http.StatusUnauthorized, utils.ErrGetUserDataFail, "failed to get user claims from context")
			return
		}

		replyId, sErr := utils.GetUUIDPathValue(r, "reply_id")
		if sErr != nil {
			utils.SError(w, sErr)
			return
		}

		ok, err := h.cfg.Queries.CheckCanDeleteReply(r.Context(), sqlc.CheckCanDeleteReplyParams{
			ID:     replyId,
			UserID: claims.UserId,
		})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, utils.ErrDbtxFail, err.Error())
			return
		}
		if !ok {
			utils.Error(w, http.StatusUnauthorized, utils.ErrDbtxFail, "you must be the owner of this session or the writer of this reply to delete it")
			return
		}

		next.ServeHTTP(w, r)
	})
}
