package setting

import (
	"fmt"

	"github.com/mzfarshad/tlg_bot/internal/key"
	"github.com/mzfarshad/tlg_bot/internal/storange"
)

// BotLanguage retrieves the bot's language setting for a specific user and chat.
// If no language is set in the database, it defaults to English.

func BotLanguage(userID int, chatID int64) (key.Language, error) {
	lang, err := storange.GetBotLanguage(userID, chatID)
	if err != nil {
		return "", fmt.Errorf("failed to get language from db: %v", err)
	}
	if lang == "" {
		return key.LangEN, nil
	}
	return key.Language(lang), nil
}
