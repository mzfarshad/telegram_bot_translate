package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/key"
)

// CommandHandler interface defines the contract for handling commands and callback queries.
// Any type that implements this interface can be used as a handler in the HandlerManager.

type CommandHadler interface {
	Handle(chaID int64, callback *tgbotapi.CallbackQuery, lang key.Language)
}

// HandlerManager manages a collection of CommandHandlers, each associated with a unique handler name.

type HandlerManager struct {
	handler map[string]CommandHadler // Map to store handler instances by their names
}

// newHandlerManager creates and returns a new instance of HandlerManager with an initialized handler map.

func newHandlerManager() *HandlerManager {
	return &HandlerManager{
		handler: make(map[string]CommandHadler),
	}
}

// registerHandler adds a new CommandHandler to the HandlerManager under the specified handler name.

func (hm *HandlerManager) rigesterHandler(handlerName string, handler CommandHadler) {
	hm.handler[handlerName] = handler
}

// handleInteraction processes a callback query using the appropriate CommandHandler based on the handler name.
// If the handler does not exist, it logs an error message.

func (hm *HandlerManager) handlInteraction(
	handlerName string, callback *tgbotapi.CallbackQuery, chatID int64, lang key.Language) {
	if handler, exist := hm.handler[handlerName]; exist {
		handler.Handle(chatID, callback, lang)
	} else {
		log.Printf("handler is not exist in command manager: %s", handlerName)
	}
}
