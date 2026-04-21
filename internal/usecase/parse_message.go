package usecase

import (
	"fmt"
	"net/http"

	"github.com/acakp/dumbchat/internal/domain"
)

func ParseMessage(r *http.Request) (domain.Message, error) {
	err := r.ParseForm()
	if err != nil {
		return domain.Message{}, fmt.Errorf("Error parsing form: ", err)
	}

	msg := ExtractMessageFormValues(r)
	if msg.Content == "" {
		return domain.Message{}, fmt.Errorf("Content field is required: ", err)
	}
	return msg, nil
}
