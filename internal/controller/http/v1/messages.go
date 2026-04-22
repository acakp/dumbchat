package v1

import (
	"encoding/json"
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/internal/usecase"
	"github.com/acakp/dumbchat/pkg/render"
)

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg, err := usecase.ParseMessage(r)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "Error parsing form, content field may be empty")
		return
	}

	// check nickname for banned words (e.g. "admin")
	c, err := r.Cookie("admin_session")
	isAdmin := false
	if err == nil {
		if err = usecase.IsAdminSession(h.DB, c); err == nil {
			isAdmin = true
		}
	}
	if isAdmin == false {
		if err = usecase.ValidateNickname(msg, h.Cfg.BannedNicknames); err != nil {
			render.Error(w, err, http.StatusBadRequest, "Nickname contains prohibited words")
			return
		}
	}

	// process the form data
	msg.TruncateMessageContent()
	msg.ID, err = adapter.InsertMessage(h.DB, msg)
	if err != nil {
		render.Error(w, err, http.StatusInternalServerError, "Failed to save message")
		return
	}

	// notify websocket hub about new message
	event := ws.Event{
		Type: "new_message",
		Data: msg,
	}
	jsonData, _ := json.Marshal(event)
	h.Hub.Broadcast <- jsonData
}
