package usecase

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/acakp/dumbchat/internal/domain"
)

func IssueAdminSession(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   36000, // 10h
	}
	http.SetCookie(w, cookie)
}

func IsAdminSession(db *sql.DB, cookie *http.Cookie) error {
	rows, err := db.Query(`
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
