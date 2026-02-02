package chat

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Event
}

type Client struct {
	conn *websocket.Conn
	send chan Event
}

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			delete(h.clients, c)
			close(c.send)
		case msg := <-h.broadcast:
			for c := range h.clients {
				c.send <- msg
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		c.conn.WriteJSON(msg)
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
