package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/key"
	"github.com/mzfarshad/tlg_bot/internal/storange"
)

// MenuHandler interface defines the contract for handling menu interactions.
// Any type that implements this interface can be used to show different menus.

type MenuHandler interface {
	ShowMenu(userID int, chatID int64, lang key.Language)
}

// MenuManager manages a collection of MenuHandlers, each associated with a unique path.

type MenuManager struct {
	Menu map[string]MenuHandler // Map to store menu handlers by their path
}

// newMenuManager creates and returns a new instance of MenuManager with an initialized menu map.

func newMenuManager() *MenuManager {

	return &MenuManager{
		Menu: make(map[string]MenuHandler),
	}
}

// registerMenu adds a new MenuHandler to the MenuManager under the specified path.

func (mm *MenuManager) rigesterMenu(path string, menuHandler MenuHandler) {

	mm.Menu[path] = menuHandler
}

// menuInteraction processes a menu interaction using the appropriate MenuHandler based on the path.
// If the menu handler does not exist, it logs a fatal error message.

func (mm *MenuManager) menuInteraction(userID int, chatID int64, path string, lang key.Language) {

	if menu, exist := mm.Menu[path]; exist {
		menu.ShowMenu(userID, chatID, lang)
	} else {
		log.Fatalf("menu is not exist in menu manager: %s", path)
	}
}

// pushState saves the new state to the menu state stack for a user and chat.
// If the new state is different from the last one, it appends the new state.

func pushState(userID int, chatID int64, newState string) error {

	// Retrieve the current menu state from storage
	state, err := storange.GetMenuState(chatID, userID)
	if err != nil {
		return err
	}

	// Append the new state if it is different from the last state
	if len(state) == 0 || state[len(state)-1] != newState {
		state = append(state, newState)
		if err := storange.SaveMenuState(userID, chatID, state); err != nil {
			return err
		}
	}
	return nil
}

// popState removes the last state from the menu state stack for a user and chat.
// It returns the new top state or an empty string if the stack is empty.

func popState(userID int, chatID int64) string {

	// Retrieve the current menu state from storag.
	state, _ := storange.GetMenuState(chatID, userID)

	// Remove the last state if there are more than one state
	if len(state) > 1 {
		state = state[:len(state)-1]
	}

	// Return the new top state or an empty string if the stack is empty
	storange.SaveMenuState(userID, chatID, state)
	if len(state) > 0 {
		return state[len(state)-1]
	}
	return ""
}

func createMenuKeyboard(lang key.Language, buttonKeys []key.TextButton) tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for i, buttonKey := range buttonKeys {
		buttonText := key.GetKey(lang, buttonKey)
		button := tgbotapi.NewInlineKeyboardButtonData(buttonText, string(buttonKey))

		if i%2 == 0 {
			keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{button})
		} else {
			keyboardRows[len(keyboardRows)-1] = append(keyboardRows[len(keyboardRows)-1], button)
		}
	}

	backButtonRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(key.GetKey(lang, key.KeyBack), string(key.KeyBack)),
	)
	keyboardRows = append(keyboardRows, backButtonRow)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	return inlineKeyboard
}
