package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	ID        int64
	Nickname  string
	Content   string
	CreatedAt time.Time
}

func (m Message) FormattedTime() string {
	return m.CreatedAt.Format("15:04 02.01.06")
}

type Messages struct {
	Msgs []Message
}

var (
	messageTmpl *template.Template
)

func chatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("html/layout.html", "html/message.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func messagesHandler(db *sql.DB) http.HandlerFunc {
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
		msg.ID, err = processMessage(db, msg)
		if err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		messageTmpl.ExecuteTemplate(w, "msg", msg)
	}
}

func pollHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("html/message.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// lastIDStr := r.URL.Query().Get("after_id")
		// lastid := 0
		// if lastIDStr != "" {
		// 	lastid, _ = strconv.Atoi(lastIDStr)
		// }
		// messages, err := getMessagesAfter(db, lastid)
		// if err != nil {
		// 	http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		// 	return
		// }

		// load all messages every time to make deleted messages disappear w/o refreshing the page
		messages, err := getMessages(db)
		if err != nil {
			http.Error(w, "Failed to load messages", http.StatusInternalServerError)
			return
		}

		// check for admin
		c, err := r.Cookie("admin_session")
		isAdmin := false
		if err == nil {
			if erra := isAdminSession(db, c); erra == nil {
				isAdmin = true
			}
		}
		// show msgs
		for _, msg := range messages.Msgs {
			nice := struct {
				Msg     Message
				IsAdmin bool
			}{
				msg,
				isAdmin,
			}
			_ = tmpl.ExecuteTemplate(w, "msg", nice)
		}
	}
}

func processMessage(db *sql.DB, msg Message) (int64, error) {
	query := `
	INSERT INTO messages (nickname, content, created_at)
	VALUES (?, ?, ?)
	`
	res, err := db.Exec(
		query,
		msg.Nickname,
		msg.Content,
		msg.CreatedAt,
	)
	if err != nil {
		return -1, err
	}
	msgid, err := res.LastInsertId()
	return msgid, err
}

func opendb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS messages(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS admin_sessions(
			id TEXT PRIMARY KEY,
			expires_at DATETIME NOT NULL
		);
	`
	_, err := db.Exec(query)
	return err
}

// func getMessagesAfter(db *sql.DB, lastID int) (Messages, error) {
// 	rows, err := db.Query(`
// 		SELECT id, nickname, content, created_at
// 		FROM messages
// 		WHERE id > ?
// 	`, lastID)
// 	if err != nil {
// 		return Messages{}, err
// 	}
// 	defer rows.Close()

//		var messages Messages
//		for rows.Next() {
//			var m Message
//			if err := rows.Scan(&m.ID, &m.Nickname, &m.Content, &m.CreatedAt); err != nil {
//				return Messages{}, err
//			}
//			messages.Msgs = append(messages.Msgs, m)
//		}
//		return messages, nil
//	}
func getMessages(db *sql.DB) (Messages, error) {
	rows, err := db.Query(`
		SELECT id, nickname, content, created_at
		FROM messages
	`)
	if err != nil {
		return Messages{}, err
	}
	defer rows.Close()

	var messages Messages
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Nickname, &m.Content, &m.CreatedAt); err != nil {
			return Messages{}, err
		}
		messages.Msgs = append(messages.Msgs, m)
	}
	return messages, nil
}

func adminGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("html/login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func adminPostHandler(db *sql.DB) http.HandlerFunc {
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
			sessionID, _ := newSessionID()
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
			issueAdminSession(w, sessionID)

			fmt.Println("admin auth: success")
			return
		} else {
			fmt.Println("admin auth: failure.", err)
			http.Error(w, "Authentication Error", http.StatusUnauthorized)
			return
		}
	}
}

func newSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}

func issueAdminSession(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true is HTTPS only
		SameSite: http.SameSiteStrictMode,
		MaxAge:   36000, // 10h
	}
	http.SetCookie(w, cookie)
}

func requireAdmin(db *sql.DB, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Valid() != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}
		if err = isAdminSession(db, cookie); err != nil {
			http.Error(w, "Unauthorized", 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAdminSession(db *sql.DB, cookie *http.Cookie) error {
	rows, err := db.Query(`
			SELECT id, expires_at
			FROM admin_sessions
			WHERE id = ?
			AND expires_at > CURRENT_TIMESTAMP
		`, cookie.Value)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func deleteMessageHandler(db *sql.DB) http.HandlerFunc {
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	// parse the message template once
	messageTmpl, err = template.ParseFiles("html/message.html")
	if err != nil {
		log.Fatal(err)
	}

	db, errdb := opendb()
	if errdb != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	r.Get("/chat", chatHandler())
	r.Post("/messages", messagesHandler(db))
	r.Delete("/messages/{messageID}", requireAdmin(db, deleteMessageHandler(db)))
	r.Get("/poll", pollHandler(db))
	r.Get("/admin/login", adminGetHandler())
	r.Post("/admin/login", adminPostHandler(db))
	fmt.Println("starting on :8888...")
	http.ListenAndServe(":8888", r)
}
