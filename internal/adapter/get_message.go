package adapter

import (
	"database/sql"
	"fmt"

	"github.com/acakp/dumbchat/internal/domain"
)

func GetMessage(db *sql.DB, messageID int) (domain.Message, error) {
	rows, err := db.Query(`
		SELECT id, nickname, content, created_at
		FROM messages WHERE id=$1;
	`, messageID)
	if err != nil {
		return domain.Message{}, fmt.Errorf("error getting message from db %w", err)
	}
	defer rows.Close()

	var msg domain.Message
	for rows.Next() {
		if err := rows.Scan(&msg.ID, &msg.Nickname, &msg.Content, &msg.CreatedAt); err != nil {
			return domain.Message{}, fmt.Errorf("error scanning message: %w", err)
		}
	}
	return msg, nil
}
