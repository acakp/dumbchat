package usecase

import (
	"database/sql"

	"github.com/acakp/dumbchat/internal/adapter"
	"github.com/acakp/dumbchat/internal/domain"
)

func GetChatView(db *sql.DB, isAdmin bool, urls domain.URLs) (domain.ChatView, error) {
	msgs, err := adapter.GetMessages(db)
	if err != nil {
		return domain.ChatView{}, err
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
