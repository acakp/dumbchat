package v1

import (
	"errors"
	"github.com/acakp/dumbchat/internal/adapter"
	"github.com/acakp/dumbchat/internal/usecase"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (h *Handler) AdminPost(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	// extract form value
	pwd := r.FormValue("password")

	//compare hash and password
	sessionID, err := adapter.CheckAdminPassword(h.DB, pwd, h.Cfg.AdminHash)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			http.Error(w, "Authentication Error", http.StatusUnauthorized)
		} else {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	usecase.IssueAdminSession(w, sessionID)
	http.Redirect(w, r, h.URLs.Base, http.StatusSeeOther)
}
