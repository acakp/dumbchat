package auth

import (
	"database/sql"
	"errors"
	"net/http"
)

func IssueAdminSession(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true is HTTPS only
		SameSite: http.SameSiteStrictMode,
		MaxAge:   36000, // 10h
	}
	http.SetCookie(w, cookie)
}

func IsAdminSession(db *sql.DB, cookie *http.Cookie) error {
	rows, err := db.Query(`
			SELECT id, expires_at
			FROM admin_sessions
			WHERE id = ?
			AND expires_at > CURRENT_TIMESTAMP
		`, cookie.Value)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("no such session")
	}

	return nil
}
