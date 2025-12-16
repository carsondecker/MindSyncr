package rooms

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (h *RoomsHandler) createRoomService(ctx context.Context, userId uuid.UUID, name, description string) (Room, *utils.ServiceError) {
	joinCode, err := createUniqueJoinCode(ctx, h.cfg.Queries)
	if err != nil {
		return Room{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrCreateJoinCodeFail,
			Message:    err.Error(),
		}
	}

	row, err := h.cfg.Queries.InsertRoom(ctx, sqlc.InsertRoomParams{
		OwnerID:     userId,
		Name:        name,
		Description: description,
		JoinCode:    joinCode,
	})
	if err != nil {
		return Room{}, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := Room{
		Id:          row.ID,
		OwnerId:     row.OwnerID,
		Name:        row.Name,
		Description: row.Description,
		JoinCode:    row.JoinCode,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}

	return res, nil
}

func (h *RoomsHandler) getOwnedRoomsService(ctx context.Context, userId uuid.UUID) ([]Room, *utils.ServiceError) {
	rows, err := h.cfg.Queries.GetRoomsByUser(ctx, userId)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	rooms := make([]Room, 0)
	for _, row := range rows {
		rooms = append(rooms, Room{
			Id:          row.ID,
			OwnerId:     row.OwnerID,
			Name:        row.Name,
			Description: row.Description,
			JoinCode:    row.JoinCode,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return rooms, nil
}

func (h *RoomsHandler) getJoinedRoomsService(ctx context.Context, userId uuid.UUID) ([]Room, *utils.ServiceError) {
	rows, err := h.cfg.Queries.GetRoomsByMembership(ctx, userId)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	rooms := make([]Room, 0)
	for _, row := range rows {
		rooms = append(rooms, Room{
			Id:          row.ID,
			OwnerId:     row.OwnerID,
			Name:        row.Name,
			Description: row.Description,
			JoinCode:    row.JoinCode,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return rooms, nil
}

func (h *RoomsHandler) updateRoomsService(ctx context.Context, userId uuid.UUID, joinCode string, data PatchRoomRequest) (Room, *utils.ServiceError) {
	row, err := h.cfg.Queries.UpdateRoom(ctx, sqlc.UpdateRoomParams{
		OwnerID:     userId,
		JoinCode:    joinCode,
		Name:        NewNullString(data.Name),
		Description: NewNullString(data.Description),
	})
	if err != nil {
		return Room{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := Room{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		JoinCode:    row.JoinCode,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}

	return res, nil
}

func (h *RoomsHandler) deleteRoomService(ctx context.Context, userId uuid.UUID, joinCode string) *utils.ServiceError {
	err := h.cfg.Queries.DeleteRoom(ctx, sqlc.DeleteRoomParams{
		OwnerID:  userId,
		JoinCode: joinCode,
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

// TODO: make it so you can't join your own session
func (h *RoomsHandler) joinRoomService(ctx context.Context, userId uuid.UUID, joinCode string) *utils.ServiceError {
	err := h.cfg.Queries.JoinRoom(ctx, sqlc.JoinRoomParams{
		UserID:   userId,
		JoinCode: joinCode,
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

func (h *RoomsHandler) leaveRoomService(ctx context.Context, userId uuid.UUID, joinCode string) *utils.ServiceError {
	err := h.cfg.Queries.LeaveRoom(ctx, sqlc.LeaveRoomParams{
		UserID:   userId,
		JoinCode: joinCode,
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
