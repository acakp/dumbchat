package usecase

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ExtractMessageID(r *http.Request) (int, error) {
	id := chi.URLParam(r, "messageID")
	messageID, err := strconv.Atoi(id)
	if err != nil {
		return -1, fmt.Errorf("error extracting message id: %w", err)
	}
	return messageID, nil
}
