package ws

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/utils"
)

type WSHandler struct {
	cfg *utils.Config
}

func NewWSHandler(cfg *utils.Config) *WSHandler {
	return &WSHandler{
		cfg,
	}
}

func (h *WSHandler) GetConfig() *utils.Config {
	return h.cfg
}

func (h *WSHandler) HandleGetWSTicket(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (WSTicketResponse, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return WSTicketResponse{}, sErr
			}

			res, sErr := h.getWSTicketService(claims.UserId, sessionId)
			if sErr != nil {
				return WSTicketResponse{}, sErr
			}

			return res, nil
		},
	)
}
