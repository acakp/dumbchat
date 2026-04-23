package postgres

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func OpenDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", getPostgresURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database (postgres): %w", err)
	}
	return db, nil
}

func getPostgresURL(cfg Config) string {
	url := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", cfg.PGHost, cfg.PGPort, cfg.PGDBName, cfg.PGUser, cfg.PGPassword)
	return url
}
