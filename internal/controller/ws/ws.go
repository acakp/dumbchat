package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type Hub struct {
	mu         sync.Mutex
	IpCounts   map[string]int
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

type Client struct {
	ip   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	rate *rate.Limiter
}

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func New() *Hub {
	return &Hub{
		IpCounts:   make(map[string]int),
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) releaseConnection(ip string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	count, ok := h.IpCounts[ip]
	if !ok {
		return
	}
	count = count - 1
	if count <= 0 {
		delete(h.IpCounts, ip)
	} else {
		h.IpCounts[ip] = count
	}
}

// allows 6 websocket connections per IP
// if there is more than 6 conns on current IP,
// returns error, nil otherwise
func (h *Hub) trackConnection(ip string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	count, ok := h.IpCounts[ip]
	if !ok {
		h.IpCounts[ip] = 0
	}
	if count > 5 {
		return fmt.Errorf("too many connections from this IP")
	}
	h.IpCounts[ip] = count + 1
	return nil
}

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
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(h.Clients, c)
				}
			}
		}
	}
}

func (c *Client) writePump(h *Hub) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		h.releaseConnection(c.ip)
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
		h.releaseConnection(c.ip)
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
