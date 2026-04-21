package usecase

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ExtractMessageID(r *http.Request) (int, error) {
	id := chi.URLParam(r, "messageID")
	messageID, err := strconv.Atoi(id)
	if err != nil {
		// http.Error(w, "Bad request", http.StatusBadRequest)
		return -1, err
	}
	return messageID, err
}
