package chat

import (
	"database/sql"
	"html/template"
	"io"
	"os"

	"github.com/acakp/dumbchat/internal/chat"
	"github.com/go-chi/chi/v5"
)

type App struct {
	handler *chat.Handler
	// other unexported fields
}

type Config struct {
	DB       *sql.DB
	BasePath string
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

func New(cfg Config) (*App, error) {
	urls := chat.NewURLs(os.Getenv("CHAT_BASE_PATH"))

	h := &chat.Handler{
		DB:   cfg.DB,
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
