package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
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
			fmt.Println("Method not allowed")
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

		lastIDStr := r.URL.Query().Get("after_id")
		lastid := 0
		if lastIDStr != "" {
			lastid, _ = strconv.Atoi(lastIDStr)
		}
		messages, err := getMessagesAfter(db, lastid)
		if err != nil {
			http.Error(w, "Failed to load messages", http.StatusInternalServerError)
			return
		}

		for _, msg := range messages.Msgs {
			_ = tmpl.ExecuteTemplate(w, "msg", msg)
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

func getMessagesAfter(db *sql.DB, lastID int) (Messages, error) {
	rows, err := db.Query(`
		SELECT id, nickname, content, created_at
		FROM messages
		WHERE id > ?
	`, lastID)
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

func adminHandler() http.HandlerFunc {
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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	// parse the message template once
	var err error
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
	r.Get("/poll", pollHandler(db))
	r.Get("/admin/login", adminHandler())
	fmt.Println("starting on :8888...")
	http.ListenAndServe(":8888", r)
}
