package rooms

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/config"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

type RoomsHandler struct {
	cfg *config.Config
}

func NewRoomsHandler(cfg *config.Config) *RoomsHandler {
	return &RoomsHandler{
		cfg,
	}
}

func (h *RoomsHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var createRoomRequest CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&createRoomRequest); err != nil {
		utils.Error(w, http.StatusBadRequest, utils.ErrBadRequest, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	err := h.cfg.Validator.Struct(createRoomRequest)
	if err != nil {
		utils.Error(w, http.StatusUnprocessableEntity, utils.ErrValidationFailed, err.Error())
		return
	}

	ctx := r.Context()

	userId := ctx.Value(utils.UserContextKey).(*utils.Claims).UserId
	if userId == uuid.Nil {
		utils.Error(w, http.StatusInternalServerError, utils.ErrGetUserDataFail, "failed to get user id from access token")
	}

	res, sErr := h.createRoomService(ctx, userId, createRoomRequest.Name, createRoomRequest.Description)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	utils.Success(w, http.StatusCreated, res)
}
