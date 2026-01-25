package chat

import (
	"embed"
	"html/template"
	"os"

	"github.com/acakp/dumbchat/internal/chat"
	"github.com/go-chi/chi/v5"
)

var templatesFS embed.FS

func AttachTemplates(t *template.Template) error {
	chatTmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return err
	}

	for _, tmpl := range chatTmpl.Templates() {
		if tmpl.Name() == "" {
			continue
		}
		_, err := t.AddParseTree(tmpl.Name(), tmpl.Tree)
		if err != nil {
			return err
		}
	}

	return nil
}

func ChatInit(r chi.Router) error {
	db, errdb := chat.OpenDB()
	if errdb != nil {
		return errdb
	}
	defer db.Close()

	if err := chat.CreateTables(db); err != nil {
		return err
	}

	chatURLs := chat.NewURLs(os.Getenv("CHAT_BASE_PATH"))
	handler := chat.Handler{
		DB:   db,
		URLs: chatURLs,
		// Tmpls: &ts,
	}

	chat.RegisterRoutes(r, handler)

	return nil
}
