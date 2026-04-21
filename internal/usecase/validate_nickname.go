package usecase

import (
	"fmt"
	"strings"

	"github.com/acakp/dumbchat/internal/domain"
)

func ValidateNickname(msg domain.Message, bannedNicknames []string) error {
	if len(bannedNicknames) == 0 {
		return nil
	}
	for _, banned := range bannedNicknames {
		if strings.Contains(msg.Nickname, banned) {
			return fmt.Errorf("prohibited nickname")
		}
	}

	return nil
}
