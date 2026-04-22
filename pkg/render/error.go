package render

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Error(w http.ResponseWriter, err error, status int, message string) {
	log.Error().Err(err).Msg(message)
	http.Error(w, message, status)
}
