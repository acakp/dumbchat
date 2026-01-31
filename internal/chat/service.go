package chat

import (
	"fmt"
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
		Poll:        base + "/poll",
		DeleteRoute: base + "/messages/{messageID}",
		Delete: func(id int) string {
			return fmt.Sprintf("%s/messages/%d", base, id)
		},
	}
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
	bannedNicknames := strings.SplitSeq(os.Getenv("BANNED_NICKNAMES"), ",")
	for banned := range bannedNicknames {
		if strings.Contains(msg.Nickname, banned) {
			return fmt.Errorf("prohibited nickname")
		}
	}

	return nil
}
