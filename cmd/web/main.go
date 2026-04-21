package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/acakp/dumbchat/config"
	"github.com/acakp/dumbchat/internal/adapter"
	httpctrl "github.com/acakp/dumbchat/internal/controller/http"
	v1 "github.com/acakp/dumbchat/internal/controller/http/v1"
	"github.com/acakp/dumbchat/internal/controller/ws"

	// ch "github.com/acakp/dumbchat/internal/chat"
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

	ts := adapter.ParseTemplatesCmd()
	if ts.Err != nil {
		log.Fatal(ts.Err)
	}

	db, errdb := adapter.OpenDB(cfg)
	if errdb != nil {
		log.Fatal(errdb)
	}
	defer db.Close()

	if err := adapter.CreateTables(db, cfg.DBDriver); err != nil {
		log.Fatal(err)
	}

	hub := ws.New()
	go hub.Run()

	handler := v1.New(cfg, db, hub, &ts)

	r.Route(cfg.BasePath, func(r chi.Router) {
		httpctrl.RegisterRoutes(r, handler)
	})

	fmt.Printf("starting on :%s...", cfg.HttpPort)
	http.ListenAndServe(net.JoinHostPort("", cfg.HttpPort), r)
}
