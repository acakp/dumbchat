package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func chatHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	tmpl, err := template.ParseFiles("html/layout.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, messages)
}

var (
	messages Messages
	mu       sync.RWMutex
	lastID   int
)

func messagesHandler(w http.ResponseWriter, r *http.Request) {
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
	processMessage(msg)

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func processMessage(msg Message) {
	mu.Lock()
	defer mu.Unlock()
	lastID++
	msg.ID = lastID
	messages.Msgs = append(messages.Msgs, msg)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	r.Get("/chat", chatHandler)
	r.Post("/messages", messagesHandler)
	fmt.Println("starting on :8888...")
	http.ListenAndServe(":8888", r)
}
