package v1

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/acakp/dumbchat/config"
	"github.com/acakp/dumbchat/internal/adapter/templates"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/internal/domain"
)

type Handler struct {
	Cfg   config.Config
	DB    *sql.DB
	Hub   *ws.Hub
	URLs  domain.URLs
	Tmpls *templates.ParsedTemplates
}

func createURLs(cfg config.Config) domain.URLs {
	base := strings.TrimRight(cfg.BasePath, "/")

	return domain.URLs{
		Base:        base,
		Post:        base + "/messages",
		DeleteRoute: base + "/messages/{messageID}",
		Delete: func(id int) string {
			return fmt.Sprintf("%s/messages/%d", base, id)
		},
		WS:      base + "/ws",
		Message: base + "/message",
	}
}

func New(cfg config.Config, db *sql.DB, hub *ws.Hub, tmpls *templates.ParsedTemplates) *Handler {
	return &Handler{
		Cfg:   cfg,
		DB:    db,
		Hub:   hub,
		URLs:  createURLs(cfg),
		Tmpls: tmpls,
	}
}

// func NewURLs(base string) URLs {
// 	base = strings.TrimRight(base, "/")

// 	return URLs{
// 		Base:        base,
// 		Post:        base + "/messages",
// 		DeleteRoute: base + "/messages/{messageID}",
// 		Delete: func(id int) string {
// 			return fmt.Sprintf("%s/messages/%d", base, id)
// 		},
// 		WS:      base + "/ws",
// 		Message: base + "/message",
// 	}
// }
