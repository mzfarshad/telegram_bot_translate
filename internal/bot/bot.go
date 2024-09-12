package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/key"
	"github.com/mzfarshad/tlg_bot/internal/setting"
	"github.com/mzfarshad/tlg_bot/internal/storange"
)

// Bot represents the Telegram bot with its API and managers.

type Bot struct {
	API            *tgbotapi.BotAPI // API instance to interact with Telegram API
	HandlerManager *HandlerManager  // Manager for handling different commands and interactions
	MenuManager    *MenuManager     // Manager for handling menu logic
}

// NewBot creates a new instance of Bot, initializes API, handlers, and database connection.

func NewBot(tkn string) (*Bot, error) {

	// Create a new bot API instance with the given token
	botApi, err := tgbotapi.NewBotAPI(tkn)
	if err != nil {
		return nil, fmt.Errorf("error connect to telegram by token: %v", err)
	}

	// Initialize the Bot struct with the API instance
	bot := &Bot{
		API: botApi,
	}

	// Create handler and menu managers for the bot
	bot.HandlerManager = CreateHandlerManager(bot)
	bot.MenuManager = CreateMenuManager(bot)

	// Initialize the database connection
	if err := storange.InitDB(); err != nil {
		log.Fatal(err)
	}

	return bot, nil
}

// Start begins polling for updates from Telegram and processes each update.

func (b *Bot) Start() {

	// Create an Update configuration with a timeout of 60 seconds
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get a channel of updates from the Telegram API
	updates := b.API.GetUpdatesChan(u)

	// Process each update received
	for update := range updates {
		b.handleUpdate(update)
	}
}

// handleUpdate processes incoming updates from Telegram, including messages and callback queries.

func (b *Bot) handleUpdate(update tgbotapi.Update) {

	var chatID int64
	var msg *tgbotapi.Message
	var lang key.Language
	var err error

	if update.Message != nil {

		// Handle incoming message updates
		msg = update.Message
		chatID = update.Message.Chat.ID

		// Retrieve the language setting for the bot based on the user ID and chat ID
		lang, err = setting.BotLanguage(int(update.Message.From.ID), chatID)
		if err != nil {
			log.Println(err)
			return
		}
		// Process commands in the message
		if msg.IsCommand() {
			cmd := msg.Command()
			cmd = strings.ToLower(cmd)
			if cmd == string(key.StartHandler) {

				// Handle the /start command
				startHandler := &StartHandler{bot: b}
				startHandler.Handle(chatID, update.Message, lang)
			} else {
				selectLangPairs := &SelectLanguagePairs{bot: b}
				selectLangPairs.Handle(chatID, msg, lang)
			}
		}
	} else if update.CallbackQuery != nil {
		// Handle callback queries
		if update.CallbackQuery.Message != nil {
			chatID = update.CallbackQuery.Message.Chat.ID

			// Retrieve the language setting for the bot based on the user ID and chat ID
			lang, err = setting.BotLanguage(int(update.CallbackQuery.From.ID), chatID)
			if err != nil {
				log.Println(err)
				return
			}

			// Handle the callback interaction based on the callback data
			handlerName := update.CallbackQuery.Data
			b.HandlerManager.handlInteraction(handlerName, update.CallbackQuery, chatID, lang)

		} else {

			log.Println("CallbackQuery has no associated message")
		}

	} else if update.InlineQuery != nil {
		userID := update.InlineQuery.From.ID

		go b.inlineQueryHandle(int(userID), update.InlineQuery)

	} else {
		log.Println("Update has neither message, callback query nor inline query")
	}
}
