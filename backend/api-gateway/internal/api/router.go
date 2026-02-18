package api

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/auth"
	comprehensionscores "github.com/carsondecker/MindSyncr/internal/comprehension_scores"
	"github.com/carsondecker/MindSyncr/internal/questions"
	"github.com/carsondecker/MindSyncr/internal/rooms"
	"github.com/carsondecker/MindSyncr/internal/sessions"
	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/internal/ws"
	"github.com/carsondecker/MindSyncr/utils"
)

func GetRouter(cfg *sutils.Config) *http.ServeMux {
	baseRouter := http.NewServeMux()
	router := http.NewServeMux()

	baseRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	middlewareHandler := sutils.NewMiddlewareHandler(cfg)

	router.HandleFunc("GET /healthz", healthzHandler)

	authRouter := http.NewServeMux()
	authHandler := auth.NewAuthHandler(cfg)
	authRouter.HandleFunc("POST /register", authHandler.HandleRegister)
	authRouter.HandleFunc("POST /login", authHandler.HandleLogin)
	authRouter.HandleFunc("POST /refresh", authHandler.HandleRefresh)
	authRouter.Handle("POST /logout", sutils.AuthMiddleware(http.HandlerFunc(authHandler.HandleLogout)))
	authRouter.Handle("GET /me", sutils.AuthMiddleware(http.HandlerFunc(authHandler.HandleGetUser)))

	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	roomsRouter := http.NewServeMux()
	roomsHandler := rooms.NewRoomsHandler(cfg)

	roomsRouter.Handle("POST /", sutils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleCreateRoom)))
	roomsRouter.Handle("GET /", sutils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleGetRooms)))
	roomsRouter.Handle("GET /{room_id}", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipByRoomId(http.HandlerFunc(roomsHandler.HandleGetRoom)),
	))
	roomsRouter.Handle("PATCH /{room_id}", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipByRoomId(http.HandlerFunc(roomsHandler.HandleUpdateRoom)),
	))
	roomsRouter.Handle("DELETE /{room_id}", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipByRoomId(http.HandlerFunc(roomsHandler.HandleDeleteRoom)),
	))
	roomsRouter.Handle("POST /{join_code}/join", sutils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleJoinRoom)))
	roomsRouter.Handle("POST /{room_id}/leave", sutils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleLeaveRoom)))

	sessionsHandler := sessions.NewSessionsHandler(cfg)

	roomsRouter.Handle("POST /{room_id}/sessions", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipByRoomId(http.HandlerFunc(sessionsHandler.HandleCreateSession)),
	))
	roomsRouter.Handle("GET /{room_id}/sessions", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipByRoomId(http.HandlerFunc(sessionsHandler.HandleGetSessions)),
	))

	router.Handle("/rooms/", http.StripPrefix("/rooms", roomsRouter))

	sessionsRouter := http.NewServeMux()

	sessionsRouter.Handle("GET /{session_id}", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipBySessionId(http.HandlerFunc(sessionsHandler.HandleGetSession)),
	))
	sessionsRouter.Handle("DELETE /{session_id}", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipBySessionId(http.HandlerFunc(sessionsHandler.HandleDeleteSession)),
	))
	sessionsRouter.Handle("POST /{session_id}/end", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomOwnershipBySessionId(http.HandlerFunc(sessionsHandler.HandleEndSession)),
	))
	sessionsRouter.Handle("POST /{session_id}/join", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipBySessionId(http.HandlerFunc(sessionsHandler.HandleJoinSession)),
	))
	sessionsRouter.Handle("POST /{session_id}/leave", sutils.AuthMiddleware(
		middlewareHandler.CheckRoomMembershipBySessionId(
			middlewareHandler.CheckSessionActive(http.HandlerFunc(sessionsHandler.HandleLeaveSession)),
		),
	))

	comprehensionScoresHandler := comprehensionscores.NewComprehensionScoresHandler(cfg)

	sessionsRouter.Handle("POST /{session_id}/comprehension-scores", sutils.AuthMiddleware(
		middlewareHandler.CheckSessionMembershipOnly(
			middlewareHandler.CheckSessionActive(http.HandlerFunc(comprehensionScoresHandler.HandleCreateComprehensionScore)),
		),
	))
	sessionsRouter.Handle("GET /{session_id}/comprehension-scores", sutils.AuthMiddleware(
		middlewareHandler.CheckSessionMembership(http.HandlerFunc(comprehensionScoresHandler.HandleGetComprehensionScores)),
	))

	questionsHandler := questions.NewQuestionsHandler(cfg)

	sessionsRouter.Handle("POST /{session_id}/questions", sutils.AuthMiddleware(
		middlewareHandler.CheckSessionMembershipOnly(
			middlewareHandler.CheckSessionActive(http.HandlerFunc(questionsHandler.HandleCreateQuestion)),
		),
	))
	sessionsRouter.Handle("GET /{session_id}/questions", sutils.AuthMiddleware(
		middlewareHandler.CheckSessionMembership(http.HandlerFunc(questionsHandler.HandleGetQuestions)),
	))

	router.Handle("/sessions/", http.StripPrefix("/sessions", sessionsRouter))

	questionsRouter := http.NewServeMux()

	questionsRouter.Handle("DELETE /{question_id}", sutils.AuthMiddleware(
		middlewareHandler.CheckSessionMembershipOnly(
			middlewareHandler.CheckSessionActive(http.HandlerFunc(questionsHandler.HandleDeleteQuestion)),
		),
	))

	router.Handle("/questions/", http.StripPrefix("/questions", questionsRouter))

	wsRouter := http.NewServeMux()
	wsHandler := ws.NewWSHandler(cfg)

	wsRouter.Handle("GET /{session_id}", sutils.AuthMiddleware(
		middlewareHandler.CheckSessionMembership(
			middlewareHandler.CheckSessionActive(http.HandlerFunc(wsHandler.HandleGetWSTicket)),
		),
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
