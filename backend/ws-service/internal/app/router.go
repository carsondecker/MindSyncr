package app

func (a *App) registerRoutes() {
	a.Router.HandleFunc("/ws", a.Hub.WebSocketHandler).Methods("GET")
}
