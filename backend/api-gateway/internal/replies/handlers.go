package replies

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
)

type RepliesHandler struct {
	cfg *sutils.Config
}

func NewRepliesHandler(cfg *sutils.Config) *RepliesHandler {
	return &RepliesHandler{
		cfg,
	}
}

func (h *RepliesHandler) GetConfig() *sutils.Config {
	return h.cfg
}

func (h *RepliesHandler) HandleCreateReply(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusCreated,
		func(data CreateReplyRequest, claims *utils.Claims) (Reply, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return Reply{}, sErr
			}

			questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
			if sErr != nil {
				return Reply{}, sErr
			}

			res, sErr := h.createReplyService(r.Context(), claims.UserId, sessionId, questionId, data)
			if sErr != nil {
				return Reply{}, sErr
			}

			return res, nil
		},
	)
}

func (h *RepliesHandler) HandleGetReplies(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFunc(h, w, r,
		http.StatusOK,
		func() ([]Reply, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return nil, sErr
			}

			res, sErr := h.getRepliesService(r.Context(), sessionId)
			if sErr != nil {
				return nil, sErr
			}

			return res, nil
		},
	)
}

func (h *RepliesHandler) HandleDeleteReply(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return struct{}{}, sErr
			}

			replyId, sErr := utils.GetUUIDPathValue(r, "reply_id")
			if sErr != nil {
				return struct{}{}, sErr
			}

			sErr = h.deleteReplyService(r.Context(), sessionId, claims.UserId, replyId)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}

func (h *RepliesHandler) HandleUpdateReply(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusOK,
		func(data PatchReplyRequest, claims *utils.Claims) (Reply, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return Reply{}, sErr
			}

			replyId, sErr := utils.GetUUIDPathValue(r, "reply_id")
			if sErr != nil {
				return Reply{}, sErr
			}

			res, sErr := h.updateReplyService(r.Context(), sessionId, claims.UserId, replyId, data)
			if sErr != nil {
				return Reply{}, sErr
			}

			return res, nil
		},
	)
}
