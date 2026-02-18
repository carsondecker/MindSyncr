package ws

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
)

type WSHandler struct {
	cfg *sutils.Config
}

func NewWSHandler(cfg *sutils.Config) *WSHandler {
	return &WSHandler{
		cfg,
	}
}

func (h *WSHandler) GetConfig() *sutils.Config {
	return h.cfg
}

func (h *WSHandler) HandleGetWSTicket(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
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
