package sessions

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

// TODO: make it so sessions automatically start on creation and that the previous session must be ended before creation
func (h *SessionsHandler) createSessionService(ctx context.Context, userId, roomId uuid.UUID, name string) (Session, *utils.ServiceError) {
	row, err := h.cfg.Queries.InsertSession(ctx, sqlc.InsertSessionParams{
		OwnerID: userId,
		RoomID:  roomId,
		Name:    name,
	})
	if err != nil {
		return Session{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := Session{
		Id:        row.ID,
		RoomId:    row.RoomID,
		OwnerID:   row.OwnerID,
		Name:      row.Name,
		IsActive:  row.IsActive,
		StartedAt: row.StartedAt,
		EndedAt:   utils.NewNullTime(row.EndedAt),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	return res, nil
}

func (h *SessionsHandler) getSessionsService(ctx context.Context, roomId uuid.UUID) ([]Session, *utils.ServiceError) {
	rows, err := h.cfg.Queries.GetSessionsByRoomId(ctx, roomId)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	rooms := make([]Session, 0)
	for _, row := range rows {
		rooms = append(rooms, Session{
			Id:        row.ID,
			RoomId:    row.RoomID,
			OwnerID:   row.OwnerID,
			Name:      row.Name,
			IsActive:  row.IsActive,
			StartedAt: row.StartedAt,
			EndedAt:   utils.NewNullTime(row.EndedAt),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}

	return rooms, nil
}

func (h *SessionsHandler) getSessionService(ctx context.Context, sessionId uuid.UUID) (Session, *utils.ServiceError) {
	row, err := h.cfg.Queries.GetSessionById(ctx, sessionId)
	if err != nil {
		return Session{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := Session{
		Id:        row.ID,
		RoomId:    row.RoomID,
		OwnerID:   row.OwnerID,
		Name:      row.Name,
		IsActive:  row.IsActive,
		StartedAt: row.StartedAt,
		EndedAt:   utils.NewNullTime(row.EndedAt),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	return res, nil
}
