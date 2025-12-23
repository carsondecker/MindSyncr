package realtime

import (
	"log"

	"github.com/carsondecker/MindSyncr-WS/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID        uuid.UUID
	SessionID uuid.UUID
	Conn      *websocket.Conn
	SendChan  chan utils.Event
	Close     chan struct{}
	Hub       *Hub
}

func NewClient(id, sessionId uuid.UUID, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:        id,
		SessionID: sessionId,
		Conn:      conn,
		SendChan:  make(chan utils.Event),
		Close:     make(chan struct{}),
		Hub:       hub,
	}
}

// TODO: redo close with context?
func (c *Client) close() {
	select {
	case <-c.Close:
		return
	default:
		close(c.Close)
		c.Hub.Unregister <- c
		c.Conn.Close()
		log.Println("Websocket connection closed.")
	}
}

func (c *Client) ReadPump() {
	defer c.close()

	for {
		select {
		case <-c.Close:
			return
		default:
			_, _, err := c.Conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) WritePump() {
	defer c.close()

	for {
		select {
		case <-c.Close:
			return
		case event := <-c.SendChan:
			c.Conn.WriteJSON(event)
		}
	}
}
