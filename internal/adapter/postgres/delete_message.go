package postgres

import (
	"context"
	"fmt"

	"github.com/acakp/dumbchat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DeleteMessage(db *pgxpool.Pool, messageID int) error {
	query := "DELETE FROM messages WHERE id = $1;"
	res, err := db.Exec(context.Background(), query, messageID)
	if err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}
	if rows := res.RowsAffected(); rows == 0 {
		return domain.ErrMessageNotFound
	}
	return nil
}
