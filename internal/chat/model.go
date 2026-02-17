package chat

import (
	"database/sql"
	"html/template"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type Handler struct {
	DB    *sql.DB
	Hub   *Hub
	URLs  URLs
	Tmpls *parsedTemplates
}

type parsedTemplates struct {
	Err         error
	ChatTmpl    *template.Template
	MessageTmpl *template.Template
	LoginTmpl   *template.Template
}

type URLs struct {
	Base        string
	Post        string
	Poll        string
	DeleteRoute string
	Delete      func(id int) string
	WS          string
	Message     string
}

type ChatView struct {
	Messages []MessageView
	IsAdmin  bool
	URLs     URLs
}

type MessageView struct {
	URLs    URLs
	Msg     Message
	IsAdmin bool
}

type Message struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m Message) FormattedTime() string {
	return m.CreatedAt.Format("15:04 02.01.06")
}

type Hub struct {
	IpCounts   sync.Map
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
