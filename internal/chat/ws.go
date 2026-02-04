package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
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

func (e Event) ToJSON() []byte {
	jsonData, _ := json.Marshal(e)
	return jsonData
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			fmt.Println("h.register!!!")
			h.Clients[c] = true
		case c := <-h.Unregister:
			fmt.Println("h.UNregister!!")
			delete(h.Clients, c)
			close(c.send)
		case msg := <-h.Broadcast:
			for c := range h.Clients {
				fmt.Println("recieved new msg!:", string(msg))
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
		h.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("broadcast!", string(msg))
		c.hub.Broadcast <- msg
	}
}
