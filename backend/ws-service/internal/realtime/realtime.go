package realtime

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *Hub) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.Header.Get("UserId")
	sessionIdStr := r.Header.Get("SessionId")

	if userIdStr == "" || sessionIdStr == "" {
		http.Error(w, "Missing required headers", http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil || userId == uuid.Nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	sessionId, err := uuid.Parse(sessionIdStr)
	if err != nil || userId == uuid.Nil {
		http.Error(w, "Invalid session id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("New WebSocket connection established")

	client := NewClient(userId, sessionId, conn, h)

	h.Register <- client
}
