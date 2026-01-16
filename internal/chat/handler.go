package chat

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	au "github.com/acakp/dumbchat/internal/auth"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func ChatHandler(chatTmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := chatTmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func MessagesHandler(db *sql.DB, msgTmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// extract form values
		msg := Message{
			Nickname:  r.FormValue("nickname"),
			Content:   r.FormValue("content"),
			CreatedAt: time.Now(),
		}

		if msg.Content == "" {
			http.Error(w, "Content field is required", http.StatusBadRequest)
			return
		}
		if msg.Nickname == "" {
			msg.Nickname = "anonymous"
		}

		// process the form data
		msg.ID, err = insertMessage(db, msg)
		if err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		msgTmpl.ExecuteTemplate(w, "msg", msg)
	}
}

func PollHandler(db *sql.DB, msgTmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages, err := getMessages(db)
		if err != nil {
			http.Error(w, "Failed to load messages", http.StatusInternalServerError)
			return
		}

		// check for admin
		c, err := r.Cookie("admin_session")
		isAdmin := false
		if err == nil {
			if erra := au.IsAdminSession(db, c); erra == nil {
				isAdmin = true
			}
		}
		// show msgs
		for _, msg := range messages {
			nice := struct {
				Msg     Message
				IsAdmin bool
			}{
				msg,
				isAdmin,
			}
			_ = msgTmpl.ExecuteTemplate(w, "msg", nice)
		}
	}
}
func AdminGetHandler(loginTmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := loginTmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func AdminPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// extract form value
		pwd := r.FormValue("password")

		//compare hash and password
		pwdHash := os.Getenv("ADMIN_PASSWORD_HASH")
		if err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd)); err == nil {
			sessionID, _ := NewSessionID()
			// add session id to db
			query := `
				INSERT INTO admin_sessions (id, expires_at)
				VALUES (?, ?)
				`
			_, err := db.Exec(
				query,
				sessionID,
				time.Now().Add(10*time.Hour),
			)
			if err != nil {
				log.Fatal(err)
			}
			// set cookie
			au.IssueAdminSession(w, sessionID)

			fmt.Println("admin auth: success")
			return
		} else {
			fmt.Println("admin auth: failure.", err)
			http.Error(w, "Authentication Error", http.StatusUnauthorized)
			return
		}
	}
}

func DeleteMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "messageID")
		messageID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}

		query := "DELETE FROM messages WHERE id = ?"
		res, err := db.Exec(query, messageID)
		rows, _ := res.RowsAffected()
		if rows == 0 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func RequireAdmin(db *sql.DB, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Valid() != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}
		if err = au.IsAdminSession(db, cookie); err != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
