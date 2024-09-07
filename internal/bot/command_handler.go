package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/key"
	"github.com/mzfarshad/tlg_bot/internal/storange"
)

// StartHandler handles the /start command.

type StartHandler struct {
	bot *Bot
}

// Handle processes the /start command by clearing the menu state and showing the main menu.

func (b *StartHandler) Handle(chatID int64, msg *tgbotapi.Message, lang key.Language) {

	// Clear previous menu state, if user
	userID := msg.From.ID
	if err := storange.ClearMenuState(chatID, int(userID)); err != nil {
		log.Println(err)
		return
	}
	// Show the main menu
	b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuMain), lang)
}

type SelectLanguagePairs struct {
	bot *Bot
}

var languagePirs = [6]string{"fa", "en", "fr", "ar", "de", "es"}

func (b *SelectLanguagePairs) Handle(chatID int64, msg *tgbotapi.Message, lang key.Language) {

	var mssg string
	userID := msg.From.ID
	langPairs := msg.Text

	langPairs = strings.TrimPrefix(langPairs, "/")
	langPairs = strings.ToLower(langPairs)

	log.Println("<<<<<<<<<<", langPairs, ">>>>>>>>>>>")

	pairs, err := split(langPairs)
	if err != nil {
		mssg = key.GetMenuMessage(lang, key.SelectLanguagePairsMessage)
		message := tgbotapi.NewMessage(chatID, mssg)
		b.bot.API.Send(message)
		b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuTranslationLanguagePairs), lang)
		return
	}

	suorceLang, targetLang := pairs[0], pairs[1]

	if !contain(suorceLang, languagePirs) {
		mssg = wrongSelectLnagMessage(lang, suorceLang)
		message := tgbotapi.NewMessage(chatID, mssg)
		b.bot.API.Send(message)
		b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuTranslationLanguagePairs), lang)
		return
	}

	if !contain(targetLang, languagePirs) {
		mssg = wrongSelectLnagMessage(lang, targetLang)
		message := tgbotapi.NewMessage(chatID, mssg)
		b.bot.API.Send(message)
		b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuTranslationLanguagePairs), lang)
		return
	}

	if err := storange.SaveLanguagePairs(int(userID), suorceLang, targetLang); err != nil {
		log.Println(err)
	}

	b.bot.MenuManager.menuInteraction(int(userID), chatID, string(key.MenuFinishTranslateSetup), lang)

}

func contain(langSymbol string, arr [6]string) bool {
	for _, v := range arr {
		if v == langSymbol {
			return true
		}
	}
	return false
}

func split(pairs string) ([]string, error) {

	parts := strings.Split(pairs, "-")

	if len(parts) != 2 {

		return nil, errors.New("error entering language pairs")
	}

	return parts, nil
}

func wrongSelectLnagMessage(lang key.Language, s string) string {
	if lang == key.LangEN {
		return fmt.Sprintf("The selected language %s is not available in the list. Get help from the message below", s)
	} else {
		return fmt.Sprintf("زبان انتخابی  %s در لیست موجود نیست. از پیام پایین کمک بگیرید", s)
	}
}
