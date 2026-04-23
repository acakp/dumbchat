package postgres

import (
	"context"
	"fmt"
	"net/http"

	"github.com/acakp/dumbchat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func IsAdminSession(db *pgxpool.Pool, cookie *http.Cookie) error {
	rows, err := db.Query(context.Background(), `
			SELECT id, expires_at
			FROM admin_sessions
			WHERE id = $1
			AND expires_at > CURRENT_TIMESTAMP;
		`, cookie.Value)
	if err != nil {
		return fmt.Errorf("error checking admin session in db: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("error checking admin session: %w", domain.ErrNotFound)
	}

	return nil
}
