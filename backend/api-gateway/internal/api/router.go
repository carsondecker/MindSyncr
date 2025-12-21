package api

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/auth"
	"github.com/carsondecker/MindSyncr/internal/rooms"
	"github.com/carsondecker/MindSyncr/internal/sessions"
	"github.com/carsondecker/MindSyncr/internal/utils"
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

	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	roomsRouter := http.NewServeMux()
	roomsHandler := rooms.NewRoomsHandler(cfg)

	roomsRouter.Handle("POST /", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembership(http.HandlerFunc(roomsHandler.HandleCreateRoom)),
	))
	roomsRouter.Handle("GET /", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleGetRooms)))
	roomsRouter.Handle("GET /{room_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembership(http.HandlerFunc(roomsHandler.HandleGetRoom)),
	))
	roomsRouter.Handle("PATCH /{room_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnership(http.HandlerFunc(roomsHandler.HandleUpdateRoom)),
	))
	roomsRouter.Handle("DELETE /{room_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnership(http.HandlerFunc(roomsHandler.HandleDeleteRoom)),
	))
	roomsRouter.Handle("POST /{join_code}/join", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleJoinRoom)))
	roomsRouter.Handle("POST /{room_id}/leave", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleLeaveRoom)))

	router.Handle("/rooms/", http.StripPrefix("/rooms", roomsRouter))

	dependentSessionsRouter := http.NewServeMux()
	sessionsHandler := sessions.NewSessionsHandler(cfg)

	dependentSessionsRouter.Handle("POST /", utils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnership(http.HandlerFunc(sessionsHandler.HandleCreateSession)),
	))
	dependentSessionsRouter.Handle("GET /", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembership(http.HandlerFunc(sessionsHandler.HandleGetSessions)),
	))

	roomsRouter.Handle("/{room_id}/sessions/", http.StripPrefix("/{room_id}/sessions", dependentSessionsRouter))

	sessionsRouter := http.NewServeMux()

	sessionsRouter.Handle("GET /{session_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembership(http.HandlerFunc(sessionsHandler.HandleGetSession)),
	))
	sessionsRouter.Handle("DELETE /{session_id}", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembership(http.HandlerFunc(sessionsHandler.HandleEndSession)),
	))
	sessionsRouter.Handle("POST /{session_id}/end", utils.AuthMiddleware(
		middlewareHandler.CheckRoomMembership(http.HandlerFunc(sessionsHandler.HandleEndSession)),
	))

	router.Handle("/sessions/", http.StripPrefix("/sessions", sessionsRouter))

	return baseRouter
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, 200, struct {
		Status string `json:"status"`
	}{
		Status: "healthy",
	})
}
