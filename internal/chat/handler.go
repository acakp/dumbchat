package chat

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	au "github.com/acakp/dumbchat/internal/auth"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(r chi.Router, h Handler) {
	r.Get(h.URLs.Base, h.Chat)
	r.Post(h.URLs.Post, h.Messages)
	r.Delete(h.URLs.DeleteRoute, h.DeleteMessage)
	r.Get(h.URLs.Poll, h.Poll)
	r.Get("/admin/login", h.AdminGet)
	r.Post("/admin/login", h.AdminPost)
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpls.ChatTmpl.Execute(w, h.URLs)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg, err := parseMessage(r)
	if err != nil {
		http.Error(w, "Error parsing form, content field may be empty", http.StatusBadRequest)
		return
	}

	// process the form data
	msg.ID, err = insertMessage(h.DB, msg)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	msv := MessageView{
		Msg:     msg,
		IsAdmin: false,
		URLs:    h.URLs,
	}

	w.Header().Set("Content-Type", "text/html")
	h.Tmpls.MessageTmpl.ExecuteTemplate(w, "msg", msv)
}

func (h *Handler) Poll(w http.ResponseWriter, r *http.Request) {
	// check for admin
	c, err := r.Cookie("admin_session")
	isAdmin := false
	if err == nil {
		if erra := au.IsAdminSession(h.DB, c); erra == nil {
			isAdmin = true
		}
	}
	msv := MessageView{
		IsAdmin: isAdmin,
		URLs:    h.URLs,
	}
	// show msgs
	err = showAllMessages(w, h.DB, h.Tmpls.MessageTmpl, msv)
	if err != nil {
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AdminGet(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpls.LoginTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

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
	sessionID, err := checkAdminPassword(h.DB, pwd)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			http.Error(w, "Authentication Error", http.StatusUnauthorized)
		} else {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	au.IssueAdminSession(w, sessionID)
}

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := extractMessageID(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	err = deleteMessage(h.DB, messageID)
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			http.Error(w, "Message not found", http.StatusNotFound)
		} else {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func RequireAdmin(db *sql.DB, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Valid() != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err = au.IsAdminSession(db, cookie); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
