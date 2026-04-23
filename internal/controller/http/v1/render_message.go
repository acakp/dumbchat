package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/internal/domain"
	"github.com/acakp/dumbchat/internal/usecase"
	"github.com/acakp/dumbchat/pkg/render"
)

func (h *Handler) RenderMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := usecase.ExtractMessageID(r)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "Bad request")
		return
	}

	msg, err := postgres.GetMessage(h.DBPool, messageID)
	if err != nil {
		if errors.Is(err, domain.ErrMessageNotFound) {
			render.Error(w, err, http.StatusNotFound, "Message not found")
		} else {
			err = fmt.Errorf("Error rendering message (Handler.renderMessage): %w", err)
			render.Error(w, err, http.StatusInternalServerError, "Internal Server Error")
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
