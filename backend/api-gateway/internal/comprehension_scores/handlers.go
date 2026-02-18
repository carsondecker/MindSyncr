package comprehensionscores

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
)

type ComprehensionScoresHandler struct {
	cfg *sutils.Config
}

func NewComprehensionScoresHandler(cfg *sutils.Config) *ComprehensionScoresHandler {
	return &ComprehensionScoresHandler{
		cfg,
	}
}

func (h *ComprehensionScoresHandler) GetConfig() *sutils.Config {
	return h.cfg
}

func (h *ComprehensionScoresHandler) HandleCreateComprehensionScore(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusCreated,
		func(data CreateComprehensionScoreRequest, claims *utils.Claims) (ComprehensionScore, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return ComprehensionScore{}, sErr
			}

			res, sErr := h.createComprehensionScoreService(r.Context(), claims.UserId, sessionId, data.Score)
			if sErr != nil {
				return ComprehensionScore{}, sErr
			}

			return res, nil
		},
	)
}

func (h *ComprehensionScoresHandler) HandleGetComprehensionScores(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFunc(h, w, r,
		http.StatusOK,
		func() ([]ComprehensionScore, *utils.ServiceError) {
			session_id, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return nil, sErr
			}

			res, sErr := h.getComprehensionScoresService(r.Context(), session_id)
			if sErr != nil {
				return nil, sErr
			}

			return res, nil
		},
	)
}
