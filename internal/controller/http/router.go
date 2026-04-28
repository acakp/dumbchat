package http

import (
	"net/http"

	v1 "github.com/acakp/dumbchat/internal/controller/http/v1"
	"github.com/acakp/dumbchat/internal/controller/ws"
	"github.com/acakp/dumbchat/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func RegisterRoutes(r chi.Router, h *v1.Handler) {
	r.Use(hlog.NewHandler(log.Logger))
	r.Use(logger.Middleware)

	r.Get("/", h.Chat)
	r.Post("/messages", h.Messages)
	r.Delete("/messages/{messageID}", v1.RequireAdmin(h.DBPool, http.HandlerFunc(h.DeleteMessage)))
	r.Get("/admin/login", h.AdminGet)
	r.Post("/admin/login", h.AdminPost)
	r.Get("/ws", ws.HandleWS(h.Hub))
	r.Get("/message/{messageID}", h.RenderMessage)
}
