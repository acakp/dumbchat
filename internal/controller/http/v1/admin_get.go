package v1

import (
	"net/http"

	"github.com/acakp/dumbchat/pkg/render"
)

func (h *Handler) AdminGet(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpls.LoginTmpl.Execute(w, nil)
	if err != nil {
		render.Error(w, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
