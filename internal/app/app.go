package app

import (
	"fmt"
	"net"
	"net/http"

	"github.com/acakp/dumbchat/config"
	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/internal/adapter/templates"
	httpctrl "github.com/acakp/dumbchat/internal/controller/http"
	v1 "github.com/acakp/dumbchat/internal/controller/http/v1"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/web"
	"github.com/go-chi/chi/v5"
)

func Run(cfg config.Config) error {
	r := chi.NewRouter()

	fs := http.FileServer(http.FS(web.StaticFS))
	r.Handle("/static/*", fs)

	ts := templates.ParseTemplatesCmd()
	if ts.Err != nil {
		return fmt.Errorf("templates.ParseTemplatesCmd: %w", ts.Err)
	}

	dbpool, err := postgres.New(cfg.DBConfig)
	if err != nil {
		return fmt.Errorf("templates.OpenDB: %w", err)
	}
	defer dbpool.Close()

	if err = postgres.CreateTables(dbpool); err != nil {
		return fmt.Errorf("templates.CreateTables: %w", err)
	}

	hub := ws.New()
	go hub.Run()

	handler := v1.New(cfg, dbpool, hub, &ts)

	r.Route(cfg.BasePath, func(r chi.Router) {
		httpctrl.RegisterRoutes(r, handler)
	})

	fmt.Printf("starting on :%s...\n", cfg.HttpPort)
	http.ListenAndServe(net.JoinHostPort("", cfg.HttpPort), r)

	return nil
}
