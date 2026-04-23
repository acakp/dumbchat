package postgres

import (
	"database/sql"
	"fmt"

	"github.com/acakp/dumbchat/internal/domain"
)

func DeleteMessage(db *sql.DB, messageID int) error {
	query := "DELETE FROM messages WHERE id = $1;"
	res, err := db.Exec(query, messageID)
	if err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return domain.ErrMessageNotFound
	}
	return nil
}
