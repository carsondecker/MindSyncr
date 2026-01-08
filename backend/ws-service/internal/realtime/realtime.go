package realtime

import (
	"log"
	"net/http"

	"github.com/carsondecker/MindSyncr-WS/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
	Subprotocols:    []string{"mindsyncr-ws"},
}

func (h *Hub) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	utils.RegisterCustomValidations(validate)

	claims, sErr := utils.GetClaims(r, validate)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("New WebSocket connection established")

	client := NewClient(claims.UserId, claims.SessionId, conn, h)

	h.Register <- client

	<-client.Close
}
