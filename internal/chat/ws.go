package chat

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			fmt.Println("h.register!!!")
			h.clients[c] = true
		case c := <-h.unregister:
			fmt.Println("h.UNregister!!")
			delete(h.clients, c)
			close(c.send)
		case msg := <-h.broadcast:
			for c := range h.clients {
				fmt.Println("recieved new msg!:", msg)
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
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("broadcast!", msg)
		c.hub.broadcast <- msg
	}
}
