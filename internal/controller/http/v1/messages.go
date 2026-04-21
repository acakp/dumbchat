package v1

import (
	"encoding/json"
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/internal/usecase"
)

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg, err := usecase.ParseMessage(r)
	if err != nil {
		http.Error(w, "Error parsing form, content field may be empty", http.StatusBadRequest)
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
			http.Error(w, "Nickname contains prohibited words", http.StatusBadRequest)
			return
		}
	}

	// process the form data
	msg.TruncateMessageContent()
	msg.ID, err = adapter.InsertMessage(h.DB, msg)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
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
