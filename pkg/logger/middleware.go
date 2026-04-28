package logger

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func Middleware(next http.Handler) http.Handler {
	return hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		l := hlog.FromRequest(r)

		if status == 0 {
			status = 200
		}

		event := l.Info()
		if status >= 500 {
			event = log.Error()
		} else if status >= 400 {
			event = log.Warn()
		}
		event.
			Dur("duration", duration).
			Int("size", size).
			Str("ip", r.RemoteAddr).
			Int("status", status).
			Str("path", r.URL.Path).
			Str("method", r.Method).
			Msg("")
	})(next)
}
