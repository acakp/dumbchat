package domain

import "time"

type Message struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m Message) FormattedTime() string {
	return m.CreatedAt.Format("15:04 02.01.06")
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

func (m *Message) TruncateMessageContent() {
	maxLen := 4000
	if len(m.Content) > maxLen {
		m.Content = m.Content[:maxLen]
	}
}
