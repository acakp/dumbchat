package v1

import (
	"errors"
	"net/http"

	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/internal/usecase"
	"github.com/acakp/dumbchat/pkg/render"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) AdminPost(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "Error parsing form")
		return
	}
	// extract form value
	pwd := r.FormValue("password")

	// compare hash and password
	sessionID, err := postgres.CheckAdminPassword(h.DB, pwd, h.Cfg.AdminHash)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			render.Error(w, err, http.StatusUnauthorized, "Authentication Error")
		} else {
			render.Error(w, err, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}
	usecase.IssueAdminSession(w, sessionID)
	http.Redirect(w, r, h.URLs.Base, http.StatusSeeOther)
}
