package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter"
	"github.com/acakp/dumbchat/internal/domain"
	"github.com/acakp/dumbchat/internal/usecase"
)

func (h *Handler) RenderMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := usecase.ExtractMessageID(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	msg, err := adapter.GetMessage(h.DB, messageID)
	if err != nil {
		if errors.Is(err, domain.ErrMessageNotFound) {
			http.Error(w, "Message not found", http.StatusNotFound)
		} else {
			log.Fatal("Error rendering message (Handler.renderMessage): ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	msv := domain.MessageView{
		Msg:     msg,
		IsAdmin: false,
		URLs:    h.URLs,
	}
	w.Header().Set("Content-Type", "text/html")
	h.Tmpls.MessageTmpl.ExecuteTemplate(w, "msg", msv)
}
