package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	ID        int
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

func chatHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("html/layout.html", "html/message.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		messages, err := getMessages(db, 500)
		if err != nil {
			http.Error(w, "Failed to load messages", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, messages)
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
			ID:        1,
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
		if err = processMessage(db, msg); err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		messageTmpl.ExecuteTemplate(w, "msg", msg)
	}
}

func processMessage(db *sql.DB, msg Message) error {
	query := `
	INSERT INTO messages (nickname, content, created_at)
	VALUES (?, ?, ?)
	`
	_, err := db.Exec(
		query,
		msg.Nickname,
		msg.Content,
		msg.CreatedAt,
	)
	return err
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
	`
	_, err := db.Exec(query)
	return err
}

func getMessages(db *sql.DB, limit int) (Messages, error) {
	rows, err := db.Query(`
		SELECT id, nickname, content, created_at
		FROM messages
		LIMIT ?
	`, limit)
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

	r.Get("/chat", chatHandler(db))
	r.Post("/messages", messagesHandler(db))
	fmt.Println("starting on :8888...")
	http.ListenAndServe(":8888", r)
}
