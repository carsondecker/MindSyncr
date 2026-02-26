package questionlikes

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
)

type QuestionLikesHandler struct {
	cfg *sutils.Config
}

func NewQuestionLikesHandler(cfg *sutils.Config) *QuestionLikesHandler {
	return &QuestionLikesHandler{
		cfg,
	}
}

func (h *QuestionLikesHandler) GetConfig() *sutils.Config {
	return h.cfg
}

func (h *QuestionLikesHandler) HandleCreateQuestionLike(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusCreated,
		func(claims *utils.Claims) (QuestionLike, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return QuestionLike{}, sErr
			}

			questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
			if sErr != nil {
				return QuestionLike{}, sErr
			}

			res, sErr := h.createQuestionLikeService(r.Context(), claims.UserId, sessionId, questionId)
			if sErr != nil {
				return QuestionLike{}, sErr
			}

			return res, nil
		},
	)
}

func (h *QuestionLikesHandler) HandleGetQuestionLikes(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFunc(h, w, r,
		http.StatusOK,
		func() ([]QuestionLike, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return nil, sErr
			}

			res, sErr := h.getQuestionLikesService(r.Context(), sessionId)
			if sErr != nil {
				return nil, sErr
			}

			return res, nil
		},
	)
}

func (h *QuestionLikesHandler) HandleDeleteQuestionLike(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return struct{}{}, sErr
			}

			questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
			if sErr != nil {
				return struct{}{}, sErr
			}

			sErr = h.deleteQuestionLikeService(r.Context(), sessionId, claims.UserId, questionId)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}
