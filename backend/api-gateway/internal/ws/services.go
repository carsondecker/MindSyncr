package ws

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

func (h *WSHandler) getWSTicketService(userId, sessionId uuid.UUID) (WSTicketResponse, *utils.ServiceError) {
	wsTicket, err := utils.CreateWSJWT(userId, sessionId)
	if err != nil {
		return WSTicketResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	res := WSTicketResponse{
		Ticket: wsTicket,
	}

	return res, nil
}
