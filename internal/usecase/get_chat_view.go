package usecase

import (
	"fmt"

	"github.com/acakp/dumbchat/internal/adapter/postgres"
	"github.com/acakp/dumbchat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetChatView(db *pgxpool.Pool, isAdmin bool, urls domain.URLs) (domain.ChatView, error) {
	msgs, err := postgres.GetMessages(db)
	if err != nil {
		return domain.ChatView{}, fmt.Errorf("GetChatView: %w", err)
	}

	views := make([]domain.MessageView, 0, len(msgs))
	for _, msg := range msgs {
		views = append(views, domain.MessageView{
			URLs:    urls,
			Msg:     msg,
			IsAdmin: isAdmin,
		})
	}

	return domain.ChatView{
		Messages: views,
		IsAdmin:  isAdmin,
		URLs:     urls,
	}, nil
}
