package api

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/auth"
	"github.com/carsondecker/MindSyncr/internal/config"
	"github.com/carsondecker/MindSyncr/internal/rooms"
	"github.com/carsondecker/MindSyncr/internal/utils"
)

func GetRouter(cfg *config.Config) *http.ServeMux {
	baseRouter := http.NewServeMux()
	router := http.NewServeMux()

	baseRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

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
	roomsRouter.Handle("POST /", utils.AuthMiddleware(http.HandlerFunc(roomsHandler.HandleCreateRoom)))

	router.Handle("/rooms/", http.StripPrefix("/rooms", roomsRouter))

	return baseRouter
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, 200, struct {
		Status string `json:"status"`
	}{
		Status: "healthy",
	})
}
