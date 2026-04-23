package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

func CheckAdminPassword(db *pgxpool.Pool, pwd, pwdHash string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd))
	if err != nil {
		return "", fmt.Errorf("error comparing hash and password: %w", err)
	}
	sessionID, _ := newSessionID()
	// add session id to db
	query := `
		INSERT INTO admin_sessions (id, expires_at)
		VALUES ($1, $2);
		`
	_, err = db.Exec(context.Background(),
		query,
		sessionID,
		time.Now().Add(10*time.Hour),
	)
	if err != nil {
		return "", fmt.Errorf("error saving admin session id to db: %w", err)
	}
	return sessionID, nil
}

func newSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), fmt.Errorf("error generating session id: %w", err)
}
