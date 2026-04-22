package usecase

import (
	"fmt"
	"net/http"

	"github.com/acakp/dumbchat/internal/domain"
)

func ParseMessage(r *http.Request) (domain.Message, error) {
	err := r.ParseForm()
	if err != nil {
		return domain.Message{}, fmt.Errorf("error parsing form: %w", err)
	}

	msg := ExtractMessageFormValues(r)
	if msg.Content == "" {
		return domain.Message{}, fmt.Errorf("content field is required: %w", err)
	}
	return msg, nil
}
