package rooms

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/utils"
)

type RoomsHandler struct {
	cfg *utils.Config
}

func NewRoomsHandler(cfg *utils.Config) *RoomsHandler {
	return &RoomsHandler{
		cfg,
	}
}

func (h RoomsHandler) GetConfig() *utils.Config {
	return h.cfg
}

/*
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
*/

func (h *RoomsHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerWithClaimsFunc(
		h,
		w,
		r,
		201,
		func(data CreateRoomRequest, claims *utils.Claims) (CreateRoomResponse, *utils.ServiceError) {
			return h.createRoomService(r.Context(), claims.UserId, data.Name, data.Description)
		},
	)
}
