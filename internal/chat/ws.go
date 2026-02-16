package chat

import (
	"time"

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
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if ok {
				c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				c.conn.WriteMessage(websocket.TextMessage, msg)
			} else {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
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
