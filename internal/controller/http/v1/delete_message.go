package v1

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/internal/domain"
	"github.com/acakp/dumbchat/internal/usecase"
)

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := usecase.ExtractMessageID(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	msg, err := adapter.GetMessage(h.DB, messageID)
	err = adapter.DeleteMessage(h.DB, messageID)
	if err != nil {
		if errors.Is(err, domain.ErrMessageNotFound) {
			http.Error(w, "Message not found", http.StatusNotFound)
		} else {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	// notify websocket hub about deleting a  message
	event := ws.Event{
		Type: "delete_message",
		Data: msg,
	}
	jsonData, _ := json.Marshal(event)
	h.Hub.Broadcast <- jsonData
	w.WriteHeader(http.StatusOK)
}
