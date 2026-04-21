package v1

import (
	"net/http"
)

func (h *Handler) AdminGet(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpls.LoginTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
