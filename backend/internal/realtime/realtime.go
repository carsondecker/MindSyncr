package realtime

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *Hub) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserId")
	roomID := r.Header.Get("RoomId")

	if userID == "" || roomID == "" {
		http.Error(w, "Missing required headers", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("New WebSocket connection established")

	client := NewClient(userID, roomID, conn, h)

	h.Register <- client
}
