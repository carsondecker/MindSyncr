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
		func(data CreateSessionRequest, claims *utils.Claims) (Session, *utils.ServiceError) {
			roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
			if sErr != nil {
				return Session{}, sErr
			}

			res, sErr := h.createSessionService(r.Context(), claims.UserId, roomId, data.Name)
			if sErr != nil {
				return Session{}, sErr
			}

			return res, nil
		},
	)
}

func (h *SessionsHandler) HandleGetSessions(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFunc(h, w, r,
		http.StatusOK,
		func() ([]Session, *utils.ServiceError) {
			roomId, sErr := utils.GetUUIDPathValue(r, "room_id")
			if sErr != nil {
				return nil, sErr
			}

			res, sErr := h.getSessionsService(r.Context(), roomId)
			if sErr != nil {
				return nil, sErr
			}

			return res, nil
		},
	)
}

func (h *SessionsHandler) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFunc(h, w, r,
		http.StatusOK,
		func() (Session, *utils.ServiceError) {
			sessionId, sErr := utils.GetUUIDPathValue(r, "session_id")
			if sErr != nil {
				return Session{}, sErr
			}

			res, sErr := h.getSessionService(r.Context(), sessionId)
			if sErr != nil {
				return Session{}, sErr
			}

			return res, nil
		},
	)
}
