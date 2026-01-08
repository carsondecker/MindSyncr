package app

import (
	"net/http"

	"github.com/carsondecker/MindSyncr-WS/internal/utils"
)

func (a *App) registerRoutes() {
	a.Router.Handle("/ws", utils.AuthMiddleware(http.HandlerFunc(a.Hub.WebSocketHandler)))
}
