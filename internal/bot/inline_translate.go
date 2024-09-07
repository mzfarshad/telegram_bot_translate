package bot

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/storange"
	"github.com/mzfarshad/tlg_bot/internal/translation"
)

func (b *Bot) inlineQueryHandle(userID int, inlineQuery *tgbotapi.InlineQuery) {

	queryID := inlineQuery.ID
	queryText := inlineQuery.Query

	log.Println("Query Text:", queryText)

	setting, err := storange.GetTranslationSetting(userID)
	if err != nil {
		log.Printf("inline query: %v, UserID: %d", err, userID)
		return
	}

	log.Println("Source Language:", setting.SourceLanguage, "Target Language:", setting.TargetLanguage, "UserID: ", userID)

	if setting.ActiveTranslation {

		translateText, err := translation.TranslateText(queryText, setting.SourceLanguage, setting.TargetLanguage)

		log.Printf("Translate Text: %s, UserID: %d", translateText, userID)

		if err != nil {
			log.Printf("error in translate inline query from api translate: %v, UserID: %d", err, userID)
			// translateText = "Translation error"
		}

		if translateText == "" {
			translateText = "No translation avialable"
		}

		result := tgbotapi.NewInlineQueryResultArticle(
			generateUniqueID(userID),
			"Translate",
			translateText,
		)

		results := []interface{}{result}

		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: queryID,
			Results:       results,
			CacheTime:     10,
			IsPersonal:    true,
		}
		_, err = b.API.Request(inlineConf)
		if err != nil {
			log.Println("Error sending inline query response:", err, "UserID: ", userID)

		}
	}
}

func generateUniqueID(userID int) string {
	timeStamp := time.Now().UnixNano()
	return strconv.Itoa(userID) + "-" + strconv.FormatInt(timeStamp, 10)
}
