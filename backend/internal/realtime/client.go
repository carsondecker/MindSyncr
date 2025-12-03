package realtime

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	RoomID   string
	Conn     *websocket.Conn
	SendChan chan Event
	Close    chan struct{}
	Hub      *Hub
}

func NewClient(id, roomId string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:       id,
		RoomID:   roomId,
		Conn:     conn,
		SendChan: make(chan Event),
		Close:    make(chan struct{}),
		Hub:      hub,
	}
}

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
