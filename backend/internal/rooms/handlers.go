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

func (h *RoomsHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusCreated,
		func(data CreateRoomRequest, claims *utils.Claims) (CreateRoomResponse, *utils.ServiceError) {
			return h.createRoomService(r.Context(), claims.UserId, data.Name, data.Description)
		},
	)
}

func (h *RoomsHandler) HandleGetRooms(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) ([]Room, *utils.ServiceError) {
			return h.getRoomsService(r.Context(), claims.UserId)
		},
	)
}

func (h *RoomsHandler) HandleUpdateRoom(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusOK,
		func(data PatchRoomRequest, claims *utils.Claims) (Room, *utils.ServiceError) {
			joinCode := r.PathValue("join_code")
			if joinCode == "" || len(joinCode) != JoinCodeLength {
				return Room{}, &utils.ServiceError{
					StatusCode: http.StatusInternalServerError,
					Code:       utils.ErrValidationFailed,
				}
			}

			res, sErr := h.updateRoomsService(r.Context(), claims.UserId, joinCode, data)
			if sErr != nil {
				return Room{}, sErr
			}

			return res, nil
		},
	)
}
