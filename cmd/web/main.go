package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	ch "github.com/acakp/dumbchat/internal/chat"
	"github.com/acakp/dumbchat/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.FileServer(http.FS(web.StaticFS))
	r.Handle("/static/*", fs)

	ts := ch.ParseTemplatesCmd()
	if ts.Err != nil {
		log.Fatal(ts.Err)
	}

	db, errdb := ch.OpenDB()
	if errdb != nil {
		log.Fatal(errdb)
	}
	defer db.Close()

	if err := ch.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	hub := ch.Hub{
		IpCounts:   sync.Map{},
		Clients:    make(map[*ch.Client]bool),
		Register:   make(chan *ch.Client),
		Unregister: make(chan *ch.Client),
		Broadcast:  make(chan []byte),
	}
	go hub.Run()

	basePath := os.Getenv("CHAT_BASE_PATH")
	chatURLs := ch.NewURLs(basePath)
	handler := ch.Handler{
		DB:    db,
		Hub:   &hub,
		URLs:  chatURLs,
		Tmpls: &ts,
	}

	r.Route(basePath, func(r chi.Router) {
		handler.RegisterRoutes(r)
	})

	fmt.Println("starting on :8888...")
	http.ListenAndServe(":8888", r)
}
