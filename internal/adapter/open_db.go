package adapter

import (
	"database/sql"
	"fmt"
	"github.com/acakp/dumbchat/config"

	_ "modernc.org/sqlite"
)

func OpenDB(cfg config.Config) (*sql.DB, error) {
	switch cfg.DBDriver {
	case "sqlite":
		db, err := sql.Open("sqlite", "./chat.db")
		if err != nil {
			return nil, err
		}
		db.SetMaxOpenConns(1)
		return db, nil
	case "postgres", "postgresql":
		db, err := sql.Open("pgx", getPostgresURL(cfg))
		if err != nil {
			return nil, fmt.Errorf("Unable to connect to database (postgres): %v\n", err)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("Unable to connect to database: environment variable 'DB' is invalid")
	}
}

func getPostgresURL(cfg config.Config) string {
	url := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", cfg.PGHost, cfg.PGPort, cfg.PGDBName, cfg.PGUser, cfg.PGPassword)
	return url
}
