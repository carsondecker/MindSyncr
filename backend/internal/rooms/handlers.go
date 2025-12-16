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
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) ([]Room, *utils.ServiceError) {
			return h.getOwnedRoomsService(r.Context(), claims.UserId)
		},
	)
}

func (h *RoomsHandler) HandleGetJoinedRooms(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) ([]Room, *utils.ServiceError) {
			return h.getJoinedRoomsService(r.Context(), claims.UserId)
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
					StatusCode: http.StatusUnprocessableEntity,
					Code:       utils.ErrValidationFailed,
					Message:    "failed to get join code",
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

func (h *RoomsHandler) HandleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			joinCode := r.PathValue("join_code")
			if joinCode == "" || len(joinCode) != JoinCodeLength {
				return struct{}{}, &utils.ServiceError{
					StatusCode: http.StatusUnprocessableEntity,
					Code:       utils.ErrValidationFailed,
					Message:    "failed to get join code",
				}
			}

			sErr := h.deleteRoomService(r.Context(), claims.UserId, joinCode)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}

func (h *RoomsHandler) HandleJoinRoom(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			joinCode := r.PathValue("join_code")
			if joinCode == "" || len(joinCode) != JoinCodeLength {
				return struct{}{}, &utils.ServiceError{
					StatusCode: http.StatusUnprocessableEntity,
					Code:       utils.ErrValidationFailed,
					Message:    "failed to get join code",
				}
			}

			sErr := h.joinRoomService(r.Context(), claims.UserId, joinCode)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}

func (h *RoomsHandler) HandleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			joinCode := r.PathValue("join_code")
			if joinCode == "" || len(joinCode) != JoinCodeLength {
				return struct{}{}, &utils.ServiceError{
					StatusCode: http.StatusUnprocessableEntity,
					Code:       utils.ErrValidationFailed,
					Message:    "failed to get join code",
				}
			}

			sErr := h.leaveRoomService(r.Context(), claims.UserId, joinCode)
			if sErr != nil {
				return struct{}{}, sErr
			}

			return struct{}{}, nil
		},
	)
}
