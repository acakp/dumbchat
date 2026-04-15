package chat

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"

	"github.com/acakp/dumbchat/config"
	"github.com/acakp/dumbchat/internal/chat"
	"github.com/go-chi/chi/v5"
)

type App struct {
	handler *chat.Handler
	// other unexported fields
}

func (a *App) AttachTemplates(t *template.Template) error {
	tmpls := chat.ParseTemplates(t)
	if tmpls.Err != nil {
		return tmpls.Err
	}
	a.handler.Tmpls = &tmpls

	return nil
}

func (a *App) Execute(wr io.Writer) {
	a.handler.Tmpls.ChatTmpl.Execute(wr, a.handler.URLs)
}

func New(db *sql.DB) (*App, error) {
	cfg, err := config.Init()
	if err != nil {
		return &App{}, fmt.Errorf("Error initializing config for new app (chat.go): %v\n", err)
	}
	urls := chat.NewURLs(cfg.BasePath)

	h := &chat.Handler{
		DB:   db,
		URLs: urls,
	}

	return &App{handler: h}, nil
}

func (a *App) RegisterRoutes(r chi.Router) {
	a.handler.RegisterRoutes(r)
}

func (a *App) CreateTables(db *sql.DB) {
	chat.CreateTables(db)
}
