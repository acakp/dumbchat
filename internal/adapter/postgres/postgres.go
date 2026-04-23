package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	PGHost     string `env:"PGHOST"`
	PGPort     string `env:"PGPORT"`
	PGDBName   string `env:"PGDBNAME"`
	PGUser     string `env:"PGUSER"`
	PGPassword string `env:"PGPASSWORD"`
}

func New(cfg Config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), getPostgresURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool (postgres): %w", err)
	}
	return dbpool, nil
}

func getPostgresURL(cfg Config) string {
	url := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", cfg.PGHost, cfg.PGPort, cfg.PGDBName, cfg.PGUser, cfg.PGPassword)
	return url
}
