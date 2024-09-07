package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/key"
	"github.com/mzfarshad/tlg_bot/internal/setting"
	"github.com/mzfarshad/tlg_bot/internal/storange"
)

// CreateHandlerManager initializes and returns a new HandlerManager with registered handlers.

func CreateHandlerManager(bot *Bot) *HandlerManager {
	hm := newHandlerManager()

	// Register handlers for different actions
	hm.rigesterHandler(string(key.KeySettings), &SettingsHandler{bot: bot})
	hm.rigesterHandler(string(key.KeyBack), &BackHandler{bot: bot})
	hm.rigesterHandler(string(key.KeySettingsLanguage), &SettingLanguageHandler{bot: bot})
	hm.rigesterHandler(string(key.LangEN), &SettingSelectLanguageHandler{bot: bot})
	hm.rigesterHandler(string(key.LangFA), &SettingSelectLanguageHandler{bot: bot})
	hm.rigesterHandler(string(key.KeyTranslaion), &TranslationHandler{bot: bot})
	hm.rigesterHandler(string(key.KeyTranslateSentMessage), &TranslationSentMessagesHandler{bot: bot})
	hm.rigesterHandler(string(key.KeyFinishSetup), &TranslationFinishSetup{bot: bot})
	hm.rigesterHandler(string(key.KeyResetTranslationSetting), &TranslateResetSetting{bot: bot})
	hm.rigesterHandler(string(key.KeyResetTranslateYes), &TranslationResetSettingYes{bot: bot})
	hm.rigesterHandler(string(key.KeyHelp), &HelpHandler{bot: bot})

	return hm
}

// SettingsHandler handles interactions related to the settings menu.

type SettingsHandler struct {
	bot *Bot
}

// Handle processes interactions to show the settings menu.

func (b *SettingsHandler) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID
	b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuSetting), lang)
}

// BackHandler handles interactions related to going back in the menu.

type BackHandler struct {
	bot *Bot
}

// Handle processes interactions to go back to the previous menu.

func (b *BackHandler) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID
	previousMenu := popState(int(userID), chatID)
	b.bot.MenuManager.menuInteraction(int(userID), chatID, previousMenu, lang)
}

// SettingLanguageHandler handles interactions related to changing the bot's language settings.

type SettingLanguageHandler struct {
	bot *Bot
}

// Handle processes interactions to show the language settings menu.

func (b *SettingLanguageHandler) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID
	b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuSettingLanguage), lang)
}

// SettingSelectLanguageHandler handles interactions related to selecting a new language for the bot.

type SettingSelectLanguageHandler struct {
	bot *Bot
}

// Handle processes interactions to save the selected language and show the language settings menu.

func (b *SettingSelectLanguageHandler) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	var msg string
	var err error
	var botLang key.Language
	selectLang := callback.Data
	userID := callback.From.ID

	// Save the selected language to storage
	if selectLang == string(key.LangEN) {
		err = storange.SaveBotLanguage(chatID, int(userID), selectLang)
		if err != nil {
			log.Printf("error saving bot language: %v", err)
		}
		msg = sendChangLangMessage(chatID, int(userID), err)
		message := tgbotapi.NewMessage(chatID, msg)
		b.bot.API.Send(message)
	}

	if selectLang == string(key.LangFA) {
		err = storange.SaveBotLanguage(chatID, int(userID), selectLang)
		if err != nil {
			log.Printf("error saving bot language: %v", err)
		}
		msg = sendChangLangMessage(chatID, int(userID), err)
		message := tgbotapi.NewMessage(chatID, msg)
		b.bot.API.Send(message)
	}

	// Retrieve the updated language setting and show the language settings menu
	botLang, err = setting.BotLanguage(int(userID), chatID)
	if err != nil {
		log.Println(err)
	}
	b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuSettingLanguage), botLang)
}

// sendChangeLangMessage generates a message indicating whether the language change was successful or failed.

func sendChangLangMessage(chatID int64, userID int, err error) string {

	var msg string
	lang, er := setting.BotLanguage(userID, chatID)
	if er != nil {
		log.Println(er)
	}
	if err == nil {
		msg = key.GetMenuMessage(lang, key.ChangeLanguageMessage)
	} else {
		msg = key.GetMenuMessage(lang, key.FailedChangeLanguageMessage)
	}
	return msg
}

type TranslationHandler struct {
	bot *Bot
}

func (h *TranslationHandler) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID
	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuTranslation), lang)
}

type TranslationSentMessagesHandler struct {
	bot *Bot
}

func (h *TranslationSentMessagesHandler) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	text := callback.Data
	userID := callback.From.ID

	if text == string(key.KeyTranslateSentMessage) {
		err := storange.SavetTranslateMessageSetting(int(userID), true)
		if err != nil {
			log.Println(err)
			return
		}
	}

	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuTranslationLanguagePairs), lang)
}

// type TranslationSentLanguagePairs struct {
// 	bot *Bot
// }

// func (h *TranslationSentLanguagePairs) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

// 	userID := callback.From.ID
// 	langPairs := callback.Data

// 	err := storange.SaveLanguagePairs(int(userID), langPairs)
// 	if err != nil {
// 		log.Printf("error handler sent language pairs: %v", err)
// 	}

// 	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuFinishTranslateSetup), lang)
// }

type TranslationFinishSetup struct {
	bot *Bot
}

func (h *TranslationFinishSetup) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID

	if err := storange.ActiveTranslation(int(userID), true); err != nil {
		log.Printf("finish translate setup: %v", err)
	}

	message := tgbotapi.NewMessage(chatID, key.GetMenuMessage(lang, key.TranslateFinishSetupMessage))
	h.bot.API.Send(message)

	for i := 0; i < 2; i++ {
		popState(int(userID), chatID)
	}

	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuTranslation), lang)
}

type TranslateResetSetting struct {
	bot *Bot
}

func (h *TranslateResetSetting) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID

	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuResetTranslate), lang)
}

type TranslationResetSettingYes struct {
	bot *Bot
}

func (h *TranslationResetSettingYes) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID
	var msg string

	if err := storange.ResetTranslationSettings(int(userID)); err != nil {
		log.Printf("error translate setting reset: %v", err)
		msg = "Something is wrong, Please try again"
	} else {
		msg = key.GetMenuMessage(lang, key.FinishResetTranslateSettingMessage)
	}

	message := tgbotapi.NewMessage(chatID, msg)
	h.bot.API.Send(message)

	popState(int(userID), chatID)

	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuTranslation), lang)
}

type HelpHandler struct {
	bot *Bot
}

func (h *HelpHandler) Handle(chatID int64, callback *tgbotapi.CallbackQuery, lang key.Language) {

	userID := callback.From.ID
	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuHelp), lang)
	h.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuMain), lang)
}
