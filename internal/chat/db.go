package chat

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
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

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chat.db")
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

func NewSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}
