package chat

import "time"

type URLs struct {
	Base   string
	Post   string
	Poll   string
	Delete string
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
