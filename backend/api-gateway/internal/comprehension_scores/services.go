package comprehensionscores

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
)

func (h *ComprehensionScoresHandler) createComprehensionScoreService(ctx context.Context, userId, sessionId uuid.UUID, score int16) (ComprehensionScore, *utils.ServiceError) {
	row, err := h.cfg.Queries.InsertComprehensionScore(ctx, sqlc.InsertComprehensionScoreParams{
		UserID:    userId,
		SessionID: sessionId,
		Score:     score,
	})
	if err != nil {
		return ComprehensionScore{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := ComprehensionScore{
		Id:        row.ID,
		SessionId: row.SessionID,
		UserId:    row.UserID,
		Score:     row.Score,
		CreatedAt: row.CreatedAt,
	}

	sErr := h.cfg.RedisClient.Broadcast("comprehension_scores", "created", sessionId, userId, res.Id, res)
	if sErr != nil {
		return ComprehensionScore{}, sErr
	}

	return res, nil
}

func (h *ComprehensionScoresHandler) getComprehensionScoresService(ctx context.Context, sessionId uuid.UUID) ([]ComprehensionScore, *utils.ServiceError) {
	rows, err := h.cfg.Queries.GetComprehensionScoresBySession(ctx, sessionId)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	scores := make([]ComprehensionScore, 0)
	for _, row := range rows {
		scores = append(scores, ComprehensionScore{
			Id:        row.ID,
			SessionId: row.SessionID,
			UserId:    row.UserID,
			Score:     row.Score,
			CreatedAt: row.CreatedAt,
		})
	}

	return scores, nil
}
