package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/acakp/dumbchat/config"
	ch "github.com/acakp/dumbchat/internal/chat"
	"github.com/acakp/dumbchat/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := config.Init()
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

	db, errdb := ch.OpenDB(cfg)
	if errdb != nil {
		log.Fatal(errdb)
	}
	defer db.Close()

	if err := ch.CreateTables(db, cfg.DBDriver); err != nil {
		log.Fatal(err)
	}

	hub := ch.Hub{
		IpCounts:   make(map[string]int),
		Clients:    make(map[*ch.Client]bool),
		Register:   make(chan *ch.Client),
		Unregister: make(chan *ch.Client),
		Broadcast:  make(chan []byte),
	}
	go hub.Run()

	chatURLs := ch.NewURLs(cfg.BasePath)
	handler := ch.Handler{
		Cfg:   cfg,
		DB:    db,
		Hub:   &hub,
		URLs:  chatURLs,
		Tmpls: &ts,
	}

	r.Route(cfg.BasePath, func(r chi.Router) {
		handler.RegisterRoutes(r)
	})

	fmt.Println("starting on :8888...")
	http.ListenAndServe(":8888", r)
}
