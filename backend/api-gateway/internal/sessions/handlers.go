package sessions

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/utils"
)

type SessionsHandler struct {
	cfg *utils.Config
}

func NewSessionsHandler(cfg *utils.Config) *SessionsHandler {
	return &SessionsHandler{
		cfg,
	}
}

func (h *SessionsHandler) GetConfig() *utils.Config {
	return h.cfg
}

func (h *SessionsHandler) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithBodyAndClaims(h, w, r,
		http.StatusCreated,
		func(data CreateSessionRequest, claims *utils.Claims) (CreateSessionResponse, *utils.ServiceError) {
			joinCode, sErr := utils.GetPathValue(r, "join_code")
			if sErr != nil {
				return CreateSessionResponse{}, sErr
			}

			res, sErr := h.createSessionService(r.Context(), claims.UserId, joinCode, data.Name)
			if sErr != nil {
				return CreateSessionResponse{}, sErr
			}

			return res, nil
		},
	)
}
