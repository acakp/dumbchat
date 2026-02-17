package chat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewURLs(base string) URLs {
	base = strings.TrimRight(base, "/")

	return URLs{
		Base:        base,
		Post:        base + "/messages",
		DeleteRoute: base + "/messages/{messageID}",
		Delete: func(id int) string {
			return fmt.Sprintf("%s/messages/%d", base, id)
		},
		WS:      base + "/ws",
		Message: base + "/message",
	}
}

func getChatView(db *sql.DB, isAdmin bool, urls URLs) (ChatView, error) {
	msgs, err := getMessages(db)
	if err != nil {
		return ChatView{}, err
	}

	views := make([]MessageView, 0, len(msgs))
	for _, msg := range msgs {
		views = append(views, MessageView{
			URLs:    urls,
			Msg:     msg,
			IsAdmin: isAdmin,
		})
	}

	return ChatView{
		Messages: views,
		IsAdmin:  isAdmin,
		URLs:     urls,
	}, nil
}

// customize json marshaling to include formatted time
func (m Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		*Alias
		FormattedTime string `json:"formattedTime"`
	}{
		Alias:         (*Alias)(&m),
		FormattedTime: m.FormattedTime(),
	})
}

func (m *Message) truncateMessageContent() {
	maxLen := 4000
	if len(m.Content) > maxLen {
		m.Content = m.Content[:maxLen]
	}
}

// allows 6 websocket connections per IP
// if there is more than 6 conns on current IP,
// returns error, nil otherwise
func (h *Hub) trackConnection(ip string) error {
	count, _ := h.IpCounts.LoadOrStore(ip, 0)
	if count.(int) > 5 {
		return fmt.Errorf("too many connections from this IP")
	}
	h.IpCounts.Store(ip, count.(int)+1)
	return nil
}

func (h *Hub) releaseConnection(ip string) {
	v, ok := h.IpCounts.Load(ip)
	if !ok {
		return
	}
	count := v.(int) - 1
	if count <= 0 {
		h.IpCounts.Delete(ip)
	} else {
		h.IpCounts.Store(ip, count)
	}
}

func clientIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	if ip == "127.0.0.1" {
		if real := r.Header.Get("X-Real-IP"); real != "" {
			return real
		}
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			return strings.Split(xff, ",")[0]
		}
	}

	return ip
}

func parseMessage(r *http.Request) (Message, error) {
	err := r.ParseForm()
	if err != nil {
		return Message{}, err // error parsing form
	}

	msg := extractMessageFormValues(r)
	if msg.Content == "" {
		return Message{}, err // content field is required
	}
	return msg, nil
}

func extractMessageFormValues(r *http.Request) Message {
	msg := Message{
		Nickname:  r.FormValue("nickname"),
		Content:   r.FormValue("content"),
		CreatedAt: time.Now(),
	}
	if msg.Nickname == "" {
		msg.Nickname = "anonymous"
	}
	return msg
}

func extractMessageID(r *http.Request) (int, error) {
	id := chi.URLParam(r, "messageID")
	messageID, err := strconv.Atoi(id)
	if err != nil {
		// http.Error(w, "Bad request", http.StatusBadRequest)
		return -1, err
	}
	return messageID, err
}

func validateNickname(msg Message) error {
	bn := os.Getenv("BANNED_NICKNAMES")
	if bn == "" {
		return nil
	}
	bannedNicknames := strings.SplitSeq(bn, ",")
	for banned := range bannedNicknames {
		if strings.Contains(msg.Nickname, banned) {
			return fmt.Errorf("prohibited nickname")
		}
	}

	return nil
}

func (e Event) ToJSON() []byte {
	jsonData, _ := json.Marshal(e)
	return jsonData
}
