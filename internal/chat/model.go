package chat

import (
	"database/sql"
	"html/template"
	"time"
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
