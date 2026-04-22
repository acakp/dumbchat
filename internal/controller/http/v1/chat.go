package v1

import (
	"net/http"

	"github.com/acakp/dumbchat/internal/usecase"
	"github.com/acakp/dumbchat/pkg/render"
)

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("admin_session")
	isAdmin := false
	if err == nil {
		if erra := usecase.IsAdminSession(h.DB, c); erra == nil {
			isAdmin = true
		}
	}

	chatView, err := usecase.GetChatView(h.DB, isAdmin, h.URLs)
	if err != nil {
		render.Error(w, err, http.StatusInternalServerError, "Failed to load chat")
		return
	}

	err = h.Tmpls.ChatTmpl.Execute(w, chatView)
	if err != nil {
		render.Error(w, err, http.StatusInternalServerError, "Failed to load chat template")
		return
	}
}
