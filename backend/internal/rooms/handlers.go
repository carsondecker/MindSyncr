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
