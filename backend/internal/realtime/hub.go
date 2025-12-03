package realtime

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Event
	Clients    map[string]*Client
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
		
		case client := <-h.Unregister:
			delete(h.Clients, client.ID)

		case event := <-h.Broadcast:
			for 
		} 
	}
}

func recievesEvent(c *Client, e Event) {

}