package chat

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/acakp/dumbchat/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.chat)
	r.Post("/messages", h.messages)
	r.Delete("/messages/{messageID}", requireAdmin(h.DB, http.HandlerFunc(h.deleteMessage)))
	r.Get("/admin/login", h.adminGet)
	r.Post("/admin/login", h.adminPost)
	r.Get("/ws", handleWS(h.Hub))
	r.Get("/message/{messageID}", h.renderMessage)
}

func handleWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Error upgrading to websockets", http.StatusUpgradeRequired)
			return
		}

		clientIp := clientIP(r)
		err = hub.trackConnection(clientIp)
		if err != nil {
			http.Error(w, "Too many connections", http.StatusTooManyRequests)
			return
		}
		client := &Client{
			ip:   clientIp,
			hub:  hub,
			conn: conn,
			send: make(chan []byte),
			rate: rate.NewLimiter(1, 5),
		}
		hub.Register <- client

		go client.writePump(hub)
		go client.readPump(hub)
	}
}

func (h *Handler) chat(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("admin_session")
	isAdmin := false
	if err == nil {
		if erra := auth.IsAdminSession(h.DB, c); erra == nil {
			isAdmin = true
		}
	}

	chatView, err := getChatView(h.DB, isAdmin, h.URLs)
	if err != nil {
		http.Error(w, "Failed to load chat", http.StatusInternalServerError)
		return
	}

	err = h.Tmpls.ChatTmpl.Execute(w, chatView)
	if err != nil {
		http.Error(w, "Failed to load chat template", http.StatusInternalServerError)
	}
}

func (h *Handler) messages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg, err := parseMessage(r)
	if err != nil {
		http.Error(w, "Error parsing form, content field may be empty", http.StatusBadRequest)
		return
	}

	// check nickname for banned words (e.g. "admin")
	c, err := r.Cookie("admin_session")
	isAdmin := false
	if err == nil {
		if err = auth.IsAdminSession(h.DB, c); err == nil {
			isAdmin = true
		}
	}
	if isAdmin == false {
		if err = validateNickname(msg); err != nil {
			http.Error(w, "Nickname contains prohibited words", http.StatusBadRequest)
			return
		}
	}

	// process the form data
	msg.truncateMessageContent()
	msg.ID, err = insertMessage(h.DB, msg)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// notify websocket hub about new message
	event := Event{
		Type: "new_message",
		Data: msg,
	}
	h.Hub.Broadcast <- event.ToJSON()
}

func (h *Handler) adminGet(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpls.LoginTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) adminPost(w http.ResponseWriter, r *http.Request) {
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
	auth.IssueAdminSession(w, sessionID)
	http.Redirect(w, r, h.URLs.Base, http.StatusSeeOther)
}

func (h *Handler) renderMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := extractMessageID(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	msg, err := getMessage(h.DB, messageID)
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			http.Error(w, "Message not found", http.StatusNotFound)
		} else {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
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

func (h *Handler) deleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := extractMessageID(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	msg, err := getMessage(h.DB, messageID)
	err = deleteMessage(h.DB, messageID)
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			http.Error(w, "Message not found", http.StatusNotFound)
		} else {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	// notify websocket hub about deleting a  message
	event := Event{
		Type: "delete_message",
		Data: msg,
	}
	h.Hub.Broadcast <- event.ToJSON()
	w.WriteHeader(http.StatusOK)
}

func requireAdmin(db *sql.DB, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Valid() != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err = auth.IsAdminSession(db, cookie); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
