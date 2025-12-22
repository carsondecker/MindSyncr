package sessions

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (h *SessionsHandler) endSessionService(ctx context.Context, userId, sessionId uuid.UUID) *utils.ServiceError {
	err := h.cfg.Queries.EndSession(ctx, sqlc.EndSessionParams{
		OwnerID: userId,
		ID:      sessionId,
	})
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	return nil
}

func (h *SessionsHandler) deleteSessionService(ctx context.Context, userId, sessionId uuid.UUID) *utils.ServiceError {
	err := h.cfg.Queries.DeleteSession(ctx, sqlc.DeleteSessionParams{
		OwnerID: userId,
		ID:      sessionId,
	})
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	return nil
}

func (h *SessionsHandler) joinSessionService(ctx context.Context, userId, sessionId uuid.UUID) *utils.ServiceError {
	ownerId, err := h.cfg.Queries.GetSessionOwnerById(ctx, sessionId)
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}
	if userId == ownerId {
		return &utils.ServiceError{
			StatusCode: http.StatusForbidden,
			Code:       utils.ErrForbidden,
			Message:    "users cannot join a room they own",
		}
	}

	err = h.cfg.Queries.JoinSession(ctx, sqlc.JoinSessionParams{
		UserID:    userId,
		SessionID: sessionId,
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return &utils.ServiceError{
					StatusCode: http.StatusBadRequest,
					Code:       utils.ErrUserAlreadyExists,
					Message:    "this user has already joined this room",
				}
			}
		}
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	return nil
}

func (h *SessionsHandler) leaveSessionService(ctx context.Context, userId, sessionId uuid.UUID) *utils.ServiceError {
	err := h.cfg.Queries.LeaveSession(ctx, sqlc.LeaveSessionParams{
		UserID:    userId,
		SessionID: sessionId,
	})
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	return nil
}
