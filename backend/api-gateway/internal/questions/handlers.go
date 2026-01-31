package questions

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/utils"
)

type QuestionsHandler struct {
	cfg *utils.Config
}

func NewQuestionsHandler(cfg *utils.Config) *QuestionsHandler {
	return &QuestionsHandler{
		cfg,
	}
}

func (h *QuestionsHandler) GetConfig() *utils.Config {
	return h.cfg
}

func (h *QuestionsHandler) HandleCreateQuestion(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusCreated,
		func(data CreateQuestionRequest, claims *utils.Claims) (Question, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return Question{}, sErr
			}

			return h.createQuestionService(r.Context(), claims.UserId, sessionId, data.Text)
		},
	)
}
