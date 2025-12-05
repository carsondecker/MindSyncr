package app

func (a *App) registerRoutes() {
	a.Router.HandleFunc("GET /ws", a.Hub.WebSocketHandler)
}
