package sessions

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

func (h *SessionsHandler) createSessionService(ctx context.Context, userId uuid.UUID, joinCode, name string) (CreateSessionResponse, *utils.ServiceError) {
	ownerId, err := h.cfg.Queries.GetRoomOwnerIdByJoinCode(ctx, joinCode)
	if err != nil {
		return CreateSessionResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}
	if userId != ownerId {
		return CreateSessionResponse{}, &utils.ServiceError{
			StatusCode: http.StatusForbidden,
			Code:       utils.ErrForbidden,
			Message:    "users cannot create a session on a room they don't own",
		}
	}

	row, err := h.cfg.Queries.InsertSession(ctx, sqlc.InsertSessionParams{
		OwnerID:  userId,
		JoinCode: joinCode,
		Name:     name,
	})
	if err != nil {
		return CreateSessionResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := CreateSessionResponse{
		Id:        row.ID,
		RoomId:    row.RoomID,
		OwnerID:   row.OwnerID,
		Name:      row.Name,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	return res, nil
}
