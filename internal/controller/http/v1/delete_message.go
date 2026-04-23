package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/internal/domain"
	"github.com/acakp/dumbchat/internal/usecase"
	"github.com/acakp/dumbchat/pkg/render"
)

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := usecase.ExtractMessageID(r)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "Bad request")
		return
	}
	msg, err := postgres.GetMessage(h.DB, messageID)
	err = postgres.DeleteMessage(h.DB, messageID)
	if err != nil {
		if errors.Is(err, domain.ErrMessageNotFound) {
			render.Error(w, err, http.StatusNotFound, "Message not found")
		} else {
			render.Error(w, err, http.StatusInternalServerError, "Internal Server Error")
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
