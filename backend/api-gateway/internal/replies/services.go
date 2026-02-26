package replies

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
)

func (h *RepliesHandler) createReplyService(ctx context.Context, userId, sessionId, questionId uuid.UUID, data CreateReplyRequest) (Reply, *utils.ServiceError) {
	row, err := h.cfg.Queries.InsertReply(ctx, sqlc.InsertReplyParams{
		UserID:     userId,
		QuestionID: questionId,
		ParentID:   data.ParentId,
		Text:       data.Text,
	})
	if err != nil {
		return Reply{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := Reply{
		Id:         row.ID,
		UserId:     row.UserID,
		QuestionId: row.QuestionID,
		ParentId:   row.ParentID,
		Text:       row.Text,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	sErr := h.cfg.RedisClient.Broadcast("replies", "created", sessionId, userId, res.Id, res)
	if sErr != nil {
		return Reply{}, sErr
	}

	return res, nil
}

func (h *RepliesHandler) getRepliesService(ctx context.Context, sessionId uuid.UUID) ([]Reply, *utils.ServiceError) {
	rows, err := h.cfg.Queries.GetRepliesBySession(ctx, sessionId)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	replies := make([]Reply, 0)
	for _, row := range rows {
		replies = append(replies, Reply{
			Id:         row.ID,
			UserId:     row.UserID,
			QuestionId: row.QuestionID,
			ParentId:   row.ParentID,
			Text:       row.Text,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
		})
	}

	return replies, nil
}

func (h *RepliesHandler) deleteReplyService(ctx context.Context, sessionId, userId, replyId uuid.UUID) *utils.ServiceError {
	id, err := h.cfg.Queries.DeleteReply(ctx, sqlc.DeleteReplyParams{
		UserID:    userId,
		ID:        replyId,
		SessionID: sessionId,
	})
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	sErr := h.cfg.RedisClient.Broadcast("replies", "deleted", sessionId, userId, id, Reply{})
	if sErr != nil {
		return sErr
	}

	return nil
}

func (h *RepliesHandler) updateReplyService(ctx context.Context, sessionId, userId, replyId uuid.UUID, data PatchReplyRequest) (Reply, *utils.ServiceError) {
	row, err := h.cfg.Queries.UpdateReply(ctx, sqlc.UpdateReplyParams{
		UserID: userId,
		ID:     replyId,
		Text:   utils.NewNullString(data.Text),
	})
	if err != nil {
		return Reply{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := Reply{
		Id:         row.ID,
		UserId:     row.UserID,
		QuestionId: row.QuestionID,
		ParentId:   row.ParentID,
		Text:       row.Text,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	sErr := h.cfg.RedisClient.Broadcast("replies", "updated", sessionId, userId, res.Id, res)
	if sErr != nil {
		return Reply{}, sErr
	}

	return res, nil
}
