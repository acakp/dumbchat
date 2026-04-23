package postgres

import (
	"context"
	"fmt"

	"github.com/acakp/dumbchat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetMessage(db *pgxpool.Pool, messageID int) (domain.Message, error) {
	rows, err := db.Query(context.Background(), `
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
