package questions

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

func (h *QuestionsHandler) createQuestionService(ctx context.Context, userId, sessionId uuid.UUID, text string) (Question, *utils.ServiceError) {
	row, err := h.cfg.Queries.InsertQuestion(ctx, sqlc.InsertQuestionParams{
		UserID:    userId,
		SessionID: sessionId,
		Text:      text,
	})
	if err != nil {
		return Question{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := Question{
		Id:         row.ID,
		UserId:     row.UserID,
		SessionId:  row.SessionID,
		Text:       row.Text,
		IsAnswered: row.IsAnswered,
		AnsweredAt: utils.NewNullTime(row.AnsweredAt),
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	sErr := h.cfg.RedisClient.Broadcast("questions", "created", sessionId, userId, res.Id, res)
	if sErr != nil {
		return Question{}, sErr
	}

	return res, nil
}
