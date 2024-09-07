package storange

import (
	"database/sql"
	"fmt"
	"strings"
)

// SaveBotLanguage saves or updates the bot's language setting for a user and chat.

func SaveBotLanguage(chatID int64, userID int, language string) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO bot_setting (user_id, chat_id, bot_language)
		VALUES (?, ?, ?)`, userID, chatID, language)
	if err != nil {
		return fmt.Errorf("failed to saving bot language in db:%v", err)
	}
	return nil
}

// GetBotLanguage retrieves the bot's language setting for a user and chat.

func GetBotLanguage(userID int, chatID int64) (string, error) {
	var lang string
	err := db.QueryRow(`SELECT bot_language FROM bot_setting WHERE user_id = ? AND chat_id = ?`,
		userID, chatID).Scan(&lang)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get bot language in db: %v", err)
	}
	return lang, nil
}

// SaveMenuState saves or updates the menu state for a user and chat.

func SaveMenuState(userID int, chatID int64, state []string) error {
	stateStr := strings.Join(state, "/")
	_, err := db.Exec(`INSERT OR REPLACE INTO menu_state (user_id, chat_id, state)
	 VALUES (?, ?, ?)`, userID, chatID, stateStr)
	if err != nil {
		return fmt.Errorf("failed to save state menu in db: %v", err)
	}
	return nil
}

// GetMenuState retrieves the menu state for a user and chat.

func GetMenuState(chatID int64, userID int) ([]string, error) {
	var stateStr string
	err := db.QueryRow(`SELECT state FROM menu_state WHERE user_id = ? AND chat_id = ?`,
		userID, chatID).Scan(&stateStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get state menu in db: %v", err)
	}
	state := strings.Split(stateStr, "/")
	return state, nil
}

// ClearMenuState removes the menu state for a user and chat.

func ClearMenuState(chatID int64, userID int) error {
	_, err := db.Exec(`DELETE FROM menu_state WHERE chat_id = ? AND user_id = ?`, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to clear menu state: %v", err)
	}
	return nil
}
