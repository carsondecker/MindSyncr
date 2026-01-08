package realtime

import (
	"log"
	"sync"
	"time"

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
	closeOnce sync.Once
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
	c.closeOnce.Do(func() {
		close(c.Close)
		c.Hub.Unregister <- c
		c.Conn.Close()
		log.Println("WebSocket connection closed.")
	})
}

func (c *Client) ReadPump() {
	defer c.close()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-c.Close:
			return
		default:
			_, _, err := c.Conn.ReadMessage()
			if err != nil {
				log.Println("ReadPump error:", err)
				return
			}
		}
	}
}

func (c *Client) WritePump() {
	defer c.close()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Close:
			return
		case event := <-c.SendChan:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteJSON(event); err != nil {
				log.Println("WriteJSON error:", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
