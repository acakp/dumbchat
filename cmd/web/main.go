package main

import (
	"fmt"
	"log"
	"net/http"

	ch "github.com/acakp/dumbchat/internal/chat"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", fs)

	ts := ch.ParseTemplates()
	if ts.Err != nil {
		log.Fatal(ts.Err)
	}

	db, errdb := ch.OpenDB()
	if errdb != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := ch.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	chatURLs := ch.NewURLs("/chat")
	handler := ch.Handler{
		DB:    db,
		URLs:  chatURLs,
		Tmpls: &ts,
	}

	ch.RegisterRoutes(r, handler)

	// r.Get(chatURLs.Base, ch.ChatHandler(ts.ChatTmpl))
	// r.Post(chatURLs.Poll, ch.MessagesHandler(db, ts.MessageTmpl))
	// r.Delete(chatURLs.Delete, ch.RequireAdmin(db, ch.DeleteMessageHandler(db)))
	// r.Get(chatURLs.Poll, ch.PollHandler(db, ts.MessageTmpl))
	// r.Get("/admin/login", ch.AdminGetHandler(ts.LoginTmpl))
	// r.Post("/admin/login", ch.AdminPostHandler(db))
	fmt.Println("starting on :8888...")
	http.ListenAndServe(":8888", r)
}
