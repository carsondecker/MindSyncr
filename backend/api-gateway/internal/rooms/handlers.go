package rooms

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
)

type RoomsHandler struct {
	cfg *sutils.Config
}

func NewRoomsHandler(cfg *sutils.Config) *RoomsHandler {
	return &RoomsHandler{
		cfg,
	}
}

func (h RoomsHandler) GetConfig() *sutils.Config {
	return h.cfg
}

func (h *RoomsHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusCreated,
		func(data CreateRoomRequest, claims *utils.Claims) (Room, *utils.ServiceError) {
			return h.createRoomService(r.Context(), claims.UserId, data.Name, data.Description)
		},
	)
}

func (h *RoomsHandler) HandleGetRooms(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")

	switch role {
	case "owner":
		h.HandleGetOwnedRooms(w, r)
	case "member":
		h.HandleGetJoinedRooms(w, r)
	default:
		utils.Error(w, http.StatusBadRequest, utils.ErrBadRequest, "endpoint requires role query parameter of either \"owner\" or \"member\"")
	}
}

func (h *RoomsHandler) HandleGetOwnedRooms(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) ([]Room, *utils.ServiceError) {
			return h.getOwnedRoomsService(r.Context(), claims.UserId)
		},
	)
}

func (h *RoomsHandler) HandleGetJoinedRooms(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) ([]Room, *utils.ServiceError) {
			return h.getJoinedRoomsService(r.Context(), claims.UserId)
		},
	)
}

func (h *RoomsHandler) HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFunc(h, w, r,
		http.StatusOK,
		func() (Room, *utils.ServiceError) {
			roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
			if sErr != nil {
				return Room{}, sErr
			}

			res, sErr := h.getRoomService(r.Context(), roomId)
			if sErr != nil {
				return Room{}, sErr
			}

			return res, nil
		},
	)
}

func (h *RoomsHandler) HandleUpdateRoom(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusOK,
		func(data PatchRoomRequest, claims *utils.Claims) (Room, *utils.ServiceError) {
			roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
			if sErr != nil {
				return Room{}, sErr
			}

			res, sErr := h.updateRoomsService(r.Context(), claims.UserId, roomId, data)
			if sErr != nil {
				return Room{}, sErr
			}

			return res, nil
		},
	)
}

func (h *RoomsHandler) HandleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
			if sErr != nil {
				return struct{}{}, sErr
			}

			sErr = h.deleteRoomService(r.Context(), claims.UserId, roomId)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}

// TODO: stop users from joining a room they own
func (h *RoomsHandler) HandleJoinRoom(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			joinCode, sErr := utils.GetPathValue(r, "join_code")
			if sErr != nil {
				return struct{}{}, sErr
			}

			sErr = h.joinRoomService(r.Context(), claims.UserId, joinCode)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}

// TODO: stop users from leaving a room they are not a member of
func (h *RoomsHandler) HandleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
			if sErr != nil {
				return struct{}{}, sErr
			}

			sErr = h.leaveRoomService(r.Context(), claims.UserId, roomId)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}
