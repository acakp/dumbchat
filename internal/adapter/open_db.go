package adapter

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func OpenDB(cfg Config) (*sql.DB, error) {
	switch cfg.DBDriver {
	case "sqlite":
		db, err := sql.Open("sqlite", "./chat.db")
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database (sqlite): %w", err)
		}
		db.SetMaxOpenConns(1)
		return db, nil
	case "postgres", "postgresql":
		db, err := sql.Open("pgx", getPostgresURL(cfg))
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database (postgres): %w", err)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unable to connect to database: environment variable 'DB' is invalid")
	}
}

func getPostgresURL(cfg Config) string {
	url := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", cfg.PGHost, cfg.PGPort, cfg.PGDBName, cfg.PGUser, cfg.PGPassword)
	return url
}
