package questions

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
)

type QuestionsHandler struct {
	cfg *sutils.Config
}

func NewQuestionsHandler(cfg *sutils.Config) *QuestionsHandler {
	return &QuestionsHandler{
		cfg,
	}
}

func (h *QuestionsHandler) GetConfig() *sutils.Config {
	return h.cfg
}

func (h *QuestionsHandler) HandleCreateQuestion(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusCreated,
		func(data CreateQuestionRequest, claims *utils.Claims) (Question, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return Question{}, sErr
			}

			res, sErr := h.createQuestionService(r.Context(), claims.UserId, sessionId, data.Text)
			if sErr != nil {
				return Question{}, sErr
			}

			return res, nil
		},
	)
}

func (h *QuestionsHandler) HandleGetQuestions(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFunc(h, w, r,
		http.StatusOK,
		func() ([]Question, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return nil, sErr
			}

			res, sErr := h.getQuestionsService(r.Context(), sessionId)
			if sErr != nil {
				return nil, sErr
			}

			return res, nil
		},
	)
}

func (h *QuestionsHandler) HandleDeleteQuestion(w http.ResponseWriter, r *http.Request) {
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

			sErr = h.deleteQuestionService(r.Context(), sessionId, claims.UserId, questionId)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}

func (h *QuestionsHandler) HandleUpdateQuestion(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusOK,
		func(data PatchQuestionRequest, claims *utils.Claims) (Question, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return Question{}, sErr
			}

			questionId, sErr := utils.GetUUIDPathValue(r, "question_id")
			if sErr != nil {
				return Question{}, sErr
			}

			res, sErr := h.updateQuestionService(r.Context(), sessionId, claims.UserId, questionId, data)
			if sErr != nil {
				return Question{}, sErr
			}

			return res, nil
		},
	)
}
