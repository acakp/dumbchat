package usecase

import (
	"github.com/acakp/dumbchat/internal/domain"
	"net/http"
	"time"
)

func ExtractMessageFormValues(r *http.Request) domain.Message {
	msg := domain.Message{
		Nickname:  r.FormValue("nickname"),
		Content:   r.FormValue("content"),
		CreatedAt: time.Now(),
	}
	if msg.Nickname == "" {
		msg.Nickname = "anonymous"
	}
	return msg
}
