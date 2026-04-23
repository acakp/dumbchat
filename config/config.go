package config

import (
	"flag"
	"fmt"

	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/pkg/logger"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger          logger.Config
	DBConfig        postgres.Config
	HttpPort        string   `env:"HTTP_PORT" envDefault:"8080"`
	AdminHash       string   `env:"ADMIN_PASSWORD_HASH,required"`
	BasePath        string   `env:"CHAT_BASE_PATH" envDefault:"/chat"`
	BannedNicknames []string `env:"BANNED_NICKNAMES"`
}

func Init() (Config, error) {
	envPath := flag.String("e", ".env", "path to the env file")
	flag.Parse()

	err := godotenv.Load(*envPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error loading env file (godotenv): %v\n", err)
	}

	var config Config
	err = env.Parse(&config)
	if err != nil {
		return Config{}, fmt.Errorf("Error loading env file (carlos0/env): %v\n", err)
	}

	return config, nil
}
