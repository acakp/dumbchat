package app

import (
	"fmt"
	"net"
	"net/http"

	"github.com/acakp/dumbchat/config"
	"github.com/acakp/dumbchat/internal/adapter"
	httpctrl "github.com/acakp/dumbchat/internal/controller/http"
	v1 "github.com/acakp/dumbchat/internal/controller/http/v1"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(cfg config.Config) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.FileServer(http.FS(web.StaticFS))
	r.Handle("/static/*", fs)

	ts := adapter.ParseTemplatesCmd()
	if ts.Err != nil {
		return fmt.Errorf("adapter.ParseTemplatesCmd: %w", ts.Err)
	}

	db, err := adapter.OpenDB(cfg)
	if err != nil {
		return fmt.Errorf("adapter.OpenDB: %w", err)
	}
	defer db.Close()

	if err = adapter.CreateTables(db, cfg.DBDriver); err != nil {
		return fmt.Errorf("adapter.CreateTables: %w", err)
	}

	hub := ws.New()
	go hub.Run()

	handler := v1.New(cfg, db, hub, &ts)

	r.Route(cfg.BasePath, func(r chi.Router) {
		httpctrl.RegisterRoutes(r, handler)
	})

	fmt.Printf("starting on :%s...", cfg.HttpPort)
	http.ListenAndServe(net.JoinHostPort("", cfg.HttpPort), r)

	return nil
}
