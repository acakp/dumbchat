package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel      string `env:"LOGGER_LEVEL" envDefault:"error"`
	PrettyConsole bool   `env:"LOGGER_PRETTY_CONSOLE" envDefault:"false"`
}

func Init(cfg Config) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	log.Logger = log.With().Timestamp().Logger()

	if cfg.PrettyConsole {
		// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Logger initialized")
}
