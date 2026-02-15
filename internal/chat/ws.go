package chat

import (
	"github.com/gorilla/websocket"
)

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true
		case c := <-h.Unregister:
			delete(h.Clients, c)
			close(c.send)
		case msg := <-h.Broadcast:
			for c := range h.Clients {
				c.send <- msg
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.conn.Close()
	}()

	for {
		if !c.rate.Allow() {
			c.conn.Close()
			return
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.Broadcast <- msg
	}
}
