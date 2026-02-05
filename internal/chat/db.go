package chat

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"os"
	"time"

	// au "github.com/acakp/dumbchat/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func insertMessage(db *sql.DB, msg Message) (int64, error) {
	query := `
	INSERT INTO messages (nickname, content, created_at)
	VALUES (?, ?, ?)
	`
	res, err := db.Exec(
		query,
		msg.Nickname,
		msg.Content,
		msg.CreatedAt,
	)
	if err != nil {
		return -1, err
	}
	msgid, err := res.LastInsertId()
	return msgid, err
}

var ErrMessageNotFound = errors.New("Message with given ID not found")

func deleteMessage(db *sql.DB, messageID int) error {
	query := "DELETE FROM messages WHERE id = ?"
	res, err := db.Exec(query, messageID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrMessageNotFound
	}
	return nil
}

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./chat.db")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS messages(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nickname TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS admin_sessions(
			id TEXT PRIMARY KEY,
			expires_at DATETIME NOT NULL
		);
	`
	_, err := db.Exec(query)
	return err
}

func getMessage(db *sql.DB, messageID int) (Message, error) {
	rows, err := db.Query(`
		SELECT id, nickname, content, created_at
		FROM messages WHERE id=?
	`, messageID)
	if err != nil {
		return Message{}, err
	}
	defer rows.Close()

	var msg Message
	for rows.Next() {
		if err := rows.Scan(&msg.ID, &msg.Nickname, &msg.Content, &msg.CreatedAt); err != nil {
			return Message{}, err
		}
	}
	return msg, nil
}

func getMessages(db *sql.DB) ([]Message, error) {
	rows, err := db.Query(`
		SELECT id, nickname, content, created_at
		FROM messages
	`)
	if err != nil {
		return []Message{}, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Nickname, &m.Content, &m.CreatedAt); err != nil {
			return []Message{}, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func newSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}

func checkAdminPassword(db *sql.DB, pwd string) (string, error) {
	pwdHash := os.Getenv("ADMIN_PASSWORD_HASH")
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd))
	if err != nil {
		return "", err
	}
	sessionID, _ := newSessionID()
	// add session id to db
	query := `
		INSERT INTO admin_sessions (id, expires_at)
		VALUES (?, ?)
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
