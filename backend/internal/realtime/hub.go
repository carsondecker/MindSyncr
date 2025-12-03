package realtime

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Event
	Clients    map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Event),
		Clients:    make(map[string]*Client),
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

func receivesEvent(c *Client, e Event) bool {
	return c.RoomID == e.RoomID
}
