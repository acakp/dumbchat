package chat

import (
	"database/sql"
	"html/template"
	"time"
)

type Handler struct {
	DB    *sql.DB
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
}

type MessageView struct {
	URLs    URLs
	Msg     Message
	IsAdmin bool
}

type Message struct {
	ID        int64
	Nickname  string
	Content   string
	CreatedAt time.Time
}

func (m Message) FormattedTime() string {
	return m.CreatedAt.Format("15:04 02.01.06")
}
