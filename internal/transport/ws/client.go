package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	UserID int64
	Conn   *websocket.Conn
	Send   chan Message
	mu     sync.Mutex
}

func (c *Client) WriteMessage(msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(msg)
}

func (c *Client) Close() {
	close(c.Send)
	_ = c.Conn.Close()
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteJSON(message); err != nil {
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
