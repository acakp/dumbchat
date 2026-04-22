package adapter

import (
	"database/sql"
	"fmt"

	"github.com/acakp/dumbchat/internal/domain"
)

func GetMessages(db *sql.DB) ([]domain.Message, error) {
	rows, err := db.Query(`
		SELECT id, nickname, content, created_at
		FROM messages;
	`)
	if err != nil {
		return []domain.Message{}, fmt.Errorf("error getting MESSAGES from db: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.Nickname, &m.Content, &m.CreatedAt); err != nil {
			return []domain.Message{}, fmt.Errorf("error scanning MESSAGES from db: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}
