package api

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/auth"
	comprehensionscores "github.com/carsondecker/MindSyncr/internal/comprehension_scores"
	"github.com/carsondecker/MindSyncr/internal/rooms"
	"github.com/carsondecker/MindSyncr/internal/sessions"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/carsondecker/MindSyncr/internal/ws"
)

func GetRouter(cfg *utils.Config) *http.ServeMux {
	baseRouter := http.NewServeMux()
	router := http.NewServeMux()

	baseRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	middlewareHandler := utils.NewMiddlewareHandler(cfg)

	router.HandleFunc("GET /healthz", healthzHandler)

	authRouter := http.NewServeMux()
	authHandler := auth.NewAuthHandler(cfg)
	authRouter.HandleFunc("POST /register", authHandler.HandleRegister)
	authRouter.HandleFunc("POST /login", authHandler.HandleLogin)
	authRouter.HandleFunc("POST /refresh", authHandler.HandleRefresh)
	authRouter.Handle("POST /logout", utils.AuthMiddleware(http.HandlerFunc(authHandler.HandleLogout)))
	authRouter.Handle("GET /me", utils.AuthMiddleware(http.HandlerFunc(authHandler.HandleGetUser)))

	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	roomsRouter := http.NewServeMux()
	roomsHandler := rooms.NewRoomsHandler(cfg)

	roomsRouter.Handle("POST /", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleCreateRoom)))
	roomsRouter.Handle("GET /", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleGetRooms)))
	roomsRouter.Handle("GET /{room_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipByRoomId(http.HandlerFunc(roomsHandler.HandleGetRoom)),
	))
	roomsRouter.Handle("PATCH /{room_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipByRoomId(http.HandlerFunc(roomsHandler.HandleUpdateRoom)),
	))
	roomsRouter.Handle("DELETE /{room_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipByRoomId(http.HandlerFunc(roomsHandler.HandleDeleteRoom)),
	))
	roomsRouter.Handle("POST /{join_code}/join", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleJoinRoom)))
	roomsRouter.Handle("POST /{room_id}/leave", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleLeaveRoom)))

	sessionsHandler := sessions.NewSessionsHandler(cfg)

	roomsRouter.Handle("POST /{room_id}/sessions", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipByRoomId(http.HandlerFunc(sessionsHandler.HandleCreateSession)),
	))
	roomsRouter.Handle("GET /{room_id}/sessions", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipByRoomId(http.HandlerFunc(sessionsHandler.HandleGetSessions)),
	))

	router.Handle("/rooms/", http.StripPrefix("/rooms", roomsRouter))

	sessionsRouter := http.NewServeMux()

	sessionsRouter.Handle("GET /{session_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipBySessionId(http.HandlerFunc(sessionsHandler.HandleGetSession)),
	))
	sessionsRouter.Handle("DELETE /{session_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipBySessionId(http.HandlerFunc(sessionsHandler.HandleDeleteSession)),
	))
	sessionsRouter.Handle("POST /{session_id}/end", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipBySessionId(http.HandlerFunc(sessionsHandler.HandleEndSession)),
	))
	sessionsRouter.Handle("POST /{session_id}/join", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipBySessionId(http.HandlerFunc(sessionsHandler.HandleJoinSession)),
	))
	sessionsRouter.Handle("POST /{session_id}/leave", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipBySessionId(
			middlewareHandler.CheckSessionActive(http.HandlerFunc(sessionsHandler.HandleLeaveSession)),
		),
	))

	comprehensionScoresHandler := comprehensionscores.NewComprehensionScoresHandler(cfg)

	sessionsRouter.Handle("POST /{session_id}/comprehension-scores", utils.AuthMiddleware(
		middlewareHandler.CheckSessionMembershipOnly(
			middlewareHandler.CheckSessionActive(http.HandlerFunc(comprehensionScoresHandler.HandleCreateComprehensionScore)),
		),
	))
	sessionsRouter.Handle("GET /{session_id}/comprehension-scores", utils.AuthMiddleware(
		middlewareHandler.CheckSessionMembershipOnly(http.HandlerFunc(comprehensionScoresHandler.HandleGetComprehensionScores)),
	))

	router.Handle("/sessions/", http.StripPrefix("/sessions", sessionsRouter))

	wsRouter := http.NewServeMux()
	wsHandler := ws.NewWSHandler(cfg)

	wsRouter.Handle("GET /{session_id}", utils.AuthMiddleware(
		middlewareHandler.CheckSessionMembership(http.HandlerFunc(wsHandler.HandleGetWSTicket)),
	))

	router.Handle("/ws/", http.StripPrefix("/ws", wsRouter))

	return baseRouter
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, 200, struct {
		Status string `json:"status"`
	}{
		Status: "healthy",
	})
}
