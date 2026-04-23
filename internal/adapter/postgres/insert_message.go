package postgres

import (
	"context"
	"fmt"

	"github.com/acakp/dumbchat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertMessage(db *pgxpool.Pool, msg domain.Message) (int64, error) {
	query := `
	INSERT INTO messages (nickname, content, created_at)
	VALUES ($1, $2, $3) RETURNING id;
	`
	var msgID int
	err := db.QueryRow(
		context.Background(),
		query,
		msg.Nickname,
		msg.Content,
		msg.CreatedAt,
	).Scan(&msgID)
	if err != nil {
		return -1, fmt.Errorf("error inserting messages to db: %w", err)
	}
	return int64(msgID), nil
}
