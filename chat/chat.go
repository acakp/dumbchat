package chat

import (
	"fmt"
	"html/template"
	"io"

	"github.com/acakp/dumbchat/config"
	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/internal/adapter/templates"
	httpctrl "github.com/acakp/dumbchat/internal/controller/http"
	v1 "github.com/acakp/dumbchat/internal/controller/http/v1"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	handler *v1.Handler
	// other unexported fields
}

func (a *App) AttachTemplates(t *template.Template) error {
	tmpls := templates.ParseTemplates(t)
	if tmpls.Err != nil {
		return tmpls.Err
	}
	a.handler.Tmpls = &tmpls

	return nil
}

func (a *App) Execute(wr io.Writer) {
	a.handler.Tmpls.ChatTmpl.Execute(wr, a.handler.URLs)
}

func New(dbpool *pgxpool.Pool) (*App, error) {
	cfg, err := config.Init()
	if err != nil {
		return &App{}, fmt.Errorf("Error initializing config for new app (chat.go): %v\n", err)
	}

	hub := ws.New()
	h := v1.New(cfg, dbpool, hub, nil)

	return &App{handler: h}, nil
}

func (a *App) RegisterRoutes(r chi.Router) {
	httpctrl.RegisterRoutes(r, a.handler)
}

func (a *App) CreateTables(dbpool *pgxpool.Pool) {
	postgres.CreateTables(dbpool)
}
