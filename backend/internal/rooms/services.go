package rooms

import (
	"context"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

func (h *RoomsHandler) createRoomService(ctx context.Context, userId uuid.UUID, name, description string) (CreateRoomResponse, *utils.ServiceError) {
	joinCode, err := createUniqueJoinCode(ctx, h.cfg.Queries)
	if err != nil {
		return CreateRoomResponse{}, &utils.ServiceError{
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
		return CreateRoomResponse{}, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := CreateRoomResponse{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		JoinCode:    row.JoinCode,
		CreatedAt:   row.CreatedAt,
	}

	return res, nil
}
