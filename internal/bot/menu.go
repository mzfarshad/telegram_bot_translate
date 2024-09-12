package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/help"
	"github.com/mzfarshad/tlg_bot/internal/key"
)

// CreateMenuManager initializes and returns a new MenuManager with registered menus.

func CreateMenuManager(bot *Bot) *MenuManager {
	menus := newMenuManager()

	// Register menus for different parts of the bot
	menus.rigesterMenu(string(key.MenuMain), &MainMenu{bot: bot})
	menus.rigesterMenu(string(key.MenuSetting), &SettingMenu{bot: bot})
	menus.rigesterMenu(string(key.MenuSettingLanguage), &SettingLanguageMenu{bot: bot})
	menus.rigesterMenu(string(key.MenuTranslation), &TranslationMenu{bot: bot})
	menus.rigesterMenu(string(key.MenuTranslationLanguagePairs), &TranslateSelectLanguagePairsMenu{bot: bot})
	menus.rigesterMenu(string(key.MenuFinishTranslateSetup), &TranslationFinishMenu{bot: bot})
	menus.rigesterMenu(string(key.MenuResetTranslate), &TranslationResetTranslateMenu{bot: bot})
	menus.rigesterMenu(string(key.MenuHelp), &HelpMenu{bot: bot})

	return menus
}

// MainMenu represents the main menu of the bot.

type MainMenu struct {
	bot *Bot
}

// ShowMenu displays the main menu to the user.

func (m *MainMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	// Push the current menu state to stack
	if err := pushState(userID, chatID, string(key.MenuMain)); err != nil {
		log.Fatal(err)
	}

	// Create the inline keyboard with options
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(key.GetKey(lang, key.KeyTranslaion),
				string(key.KeyTranslaion)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(key.GetKey(lang, key.KeySettings),
				string(key.KeySettings)),
			tgbotapi.NewInlineKeyboardButtonData(key.GetKey(lang, key.KeyContactUs),
				string(key.KeyContactUs)),
			tgbotapi.NewInlineKeyboardButtonData(key.GetKey(lang, key.KeyHelp),
				string(key.KeyHelp)),
		),
	)

	// Create and send the message with the main menu
	message := tgbotapi.NewMessage(chatID, key.GetMenuMessage(lang, key.MainMessage))
	message.ReplyMarkup = keyboard
	m.bot.API.Send(message)
}

// SettingMenu represents the settings menu of the bot.

type SettingMenu struct {
	bot *Bot
}

// ShowMenu displays the settings menu to the user.

func (m *SettingMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	// Push the current menu state to stack
	if err := pushState(userID, chatID, string(key.MenuSetting)); err != nil {
		log.Fatal(err)
	}

	buttons := []key.TextButton{key.KeySettingsLanguage}
	keyboard := createMenuKeyboard(lang, buttons)

	// Create and send the message with the settings menu
	message := tgbotapi.NewMessage(chatID, key.GetMenuMessage(lang, key.SettingMessage))
	message.ReplyMarkup = keyboard
	_, err := m.bot.API.Send(message)
	if err != nil {
		log.Printf("error sending message from main menu: %v", err)
	}
}

// SettingLanguageMenu represents the language settings menu of the bot.

type SettingLanguageMenu struct {
	bot *Bot
}

// ShowMenu displays the language settings menu to the user.

func (m *SettingLanguageMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	// Push the current menu state to stack
	if err := pushState(userID, chatID, string(key.MenuSettingLanguage)); err != nil {
		log.Fatalf("error push settings/lanuage menu state: %v", err)
	}

	// Create the inline keyboard with language options
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("English", string(key.LangEN)),
			tgbotapi.NewInlineKeyboardButtonData("فارسی", string(key.LangFA)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(key.GetKey(lang, key.KeyBack), string(key.KeyBack)),
		),
	)

	// Create and send the message with the language settings menu
	message := tgbotapi.NewMessage(chatID, key.GetMenuMessage(lang, key.SettingLanguageMessage))
	message.ReplyMarkup = keyboard
	m.bot.API.Send(message)
}

type TranslationMenu struct {
	bot *Bot
}

func (m *TranslationMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	if err := pushState(userID, chatID, string(key.MenuTranslation)); err != nil {
		log.Fatalf("error push translation menu state: %v", err)
	}

	buttons := []key.TextButton{
		key.KeyTranslateSentMessage,
		key.KeyResetTranslationSetting,
	}

	keyboard := createMenuKeyboard(lang, buttons)

	msg := key.GetMenuMessage(lang, key.TranslationMenuMessage)
	message := tgbotapi.NewMessage(chatID, msg)
	message.ReplyMarkup = keyboard

	m.bot.API.Send(message)
}

type TranslateSelectLanguagePairsMenu struct {
	bot *Bot
}

func (m *TranslateSelectLanguagePairsMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	if err := pushState(userID, chatID, string(key.MenuTranslationLanguagePairs)); err != nil {
		log.Fatalf("error push translatin language pairs menu state: %v", err)
	}

	buttons := []key.TextButton{}

	keyboard := createMenuKeyboard(lang, buttons)

	var path string
	if lang == key.LangEN {
		path = "internal/key/select_lang_pairs_en.txt"
	} else {
		path = "internal/key/select_lang_pairs_fa.txt"
	}

	msg, err := key.ReadSelectLanguagePairsMsg(path)
	if err != nil {
		log.Println(err)
	}

	message := tgbotapi.NewMessage(chatID, msg)
	message.ReplyMarkup = keyboard

	m.bot.API.Send(message)
}

type TranslationFinishMenu struct {
	bot *Bot
}

func (m *TranslationFinishMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	if err := pushState(userID, chatID, string(key.MenuFinishTranslateSetup)); err != nil {
		log.Printf("error push state menu in sent language pairs: %v", err)
	}

	buttons := []key.TextButton{key.KeyFinishSetup}
	keyboard := createMenuKeyboard(lang, buttons)

	msg := key.GetMenuMessage(lang, key.TranslateFinishMessage)
	message := tgbotapi.NewMessage(chatID, msg)
	message.ReplyMarkup = keyboard

	m.bot.API.Send(message)
}

type TranslationResetTranslateMenu struct {
	bot *Bot
}

func (m *TranslationResetTranslateMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	if err := pushState(userID, chatID, string(key.MenuResetTranslate)); err != nil {
		log.Printf("error push state menu in translate reset: %v", err)
	}

	button := []key.TextButton{key.KeyResetTranslateYes}
	keyboard := createMenuKeyboard(lang, button)

	message := tgbotapi.NewMessage(chatID, key.GetMenuMessage(lang, key.ResetTranslateSettingMessage))
	message.ReplyMarkup = keyboard

	m.bot.API.Send(message)
}

type HelpMenu struct {
	bot *Bot
}

func (m *HelpMenu) ShowMenu(userID int, chatID int64, lang key.Language) {

	var helpFilePath string

	if lang == "en" {
		helpFilePath = "internal/help/help_en.txt"
	} else {
		helpFilePath = "internal/help/help_fa.txt"
	}

	helpText, err := help.ReadHelpFile(helpFilePath)
	if err != nil {
		log.Println(err)
		helpFilePath = "Sorry, we couldn't load the help content at the moment."
	}

	message := tgbotapi.NewMessage(chatID, helpText)

	m.bot.API.Send(message)

}
