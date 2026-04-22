package main

import (
	"log"

	"github.com/acakp/dumbchat/config"
	"github.com/acakp/dumbchat/internal/app"
	"github.com/acakp/dumbchat/pkg/logger"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("error initializing config: ", err)
	}

	logger.Init(cfg.Logger)

	err = app.Run(cfg)
	if err != nil {
		log.Fatal("app.Run: ", err)
	}
}
