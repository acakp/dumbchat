package adapter

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

func CheckAdminPassword(db *sql.DB, pwd, pwdHash string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd))
	if err != nil {
		return "", err
	}
	sessionID, _ := newSessionID()
	// add session id to db
	query := `
		INSERT INTO admin_sessions (id, expires_at)
		VALUES ($1, $2);
		`
	_, err = db.Exec(
		query,
		sessionID,
		time.Now().Add(10*time.Hour),
	)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func newSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}
