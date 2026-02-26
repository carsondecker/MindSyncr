package questionlikes

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
)

func (h *QuestionLikesHandler) createQuestionLikeService(ctx context.Context, userId, sessionId, questionId uuid.UUID) (QuestionLike, *utils.ServiceError) {
	row, err := h.cfg.Queries.InsertQuestionLike(ctx, sqlc.InsertQuestionLikeParams{
		UserID:     userId,
		QuestionID: questionId,
	})
	if err != nil {
		return QuestionLike{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := QuestionLike{
		Id:         row.ID,
		UserId:     row.UserID,
		QuestionId: row.QuestionID,
		CreatedAt:  row.CreatedAt,
	}

	sErr := h.cfg.RedisClient.Broadcast("question_likes", "created", sessionId, userId, res.Id, res)
	if sErr != nil {
		return QuestionLike{}, sErr
	}

	return res, nil
}

func (h *QuestionLikesHandler) getQuestionLikesService(ctx context.Context, sessionId uuid.UUID) ([]QuestionLike, *utils.ServiceError) {
	rows, err := h.cfg.Queries.GetQuestionLikesBySession(ctx, sessionId)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	questionLikes := make([]QuestionLike, 0)
	for _, row := range rows {
		questionLikes = append(questionLikes, QuestionLike{
			Id:         row.ID,
			UserId:     row.UserID,
			QuestionId: row.QuestionID,
			CreatedAt:  row.CreatedAt,
		})
	}

	return questionLikes, nil
}

func (h *QuestionLikesHandler) deleteQuestionLikeService(ctx context.Context, sessionId, userId, questionId uuid.UUID) *utils.ServiceError {
	id, err := h.cfg.Queries.DeleteQuestionLike(ctx, sqlc.DeleteQuestionLikeParams{
		UserID:     userId,
		QuestionID: questionId,
		SessionID:  sessionId,
	})
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	sErr := h.cfg.RedisClient.Broadcast("question_likes", "deleted", sessionId, userId, id, QuestionLike{})
	if sErr != nil {
		return sErr
	}

	return nil
}
