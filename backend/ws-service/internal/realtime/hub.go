package realtime

import (
	"github.com/carsondecker/MindSyncr-WS/internal/utils"
	"github.com/google/uuid"
)

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan utils.Event
	Clients    map[uuid.UUID]*Client
}

// TODO: make clients in a map based on session ID instead
func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan utils.Event),
		Clients:    make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
			go client.ReadPump()
			go client.WritePump()

		case client := <-h.Unregister:
			delete(h.Clients, client.ID)

		case event := <-h.Broadcast:
			for _, client := range h.Clients {
				if receivesEvent(client, event) {
					client.SendChan <- event
				}
			}
		}
	}
}

func receivesEvent(c *Client, e utils.Event) bool {
	return c.SessionID == e.SessionID
}
