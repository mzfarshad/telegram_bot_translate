package storange

import (
	"errors"
	"fmt"
	"log"
)

type TranslationSetting struct {
	UserID            int
	SourceLanguage    string
	TargetLanguage    string
	SentMessage       bool
	ActiveTranslation bool
}

func SavetTranslateMessageSetting(userID int, sentMessage bool) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO translation 
					   (user_id, sent_message, sent_source_language, sent_target_language)
						VALUES (?, ?, NULL, NULL )`, userID, sentMessage)
	if err != nil {
		return fmt.Errorf("failed to save message type in db: %v", err)
	}
	return nil
}

func SaveLanguagePairs(userID int, sourceLang, targetLang string) error {

	log.Println("Source Language: ", sourceLang)
	log.Println("Target Language: ", targetLang)

	_, err := db.Exec(`UPDATE translation
					   SET sent_source_language = ?, sent_target_language = ?
					   WHERE user_id = ?`, sourceLang, targetLang, userID)

	if err != nil {
		return errors.New("failed to save sent language pairs")
	}

	return nil
}

func ActiveTranslation(userID int, active bool) error {
	_, err := db.Exec("UPDATE translation SET active_translation = ? WHERE user_id = ?",
		active, userID)

	if err != nil {
		return fmt.Errorf("failed active translation in db: %v", err)
	}

	return nil
}

func ResetTranslationSettings(userID int) error {
	res, err := db.Exec(`UPDATE translation SET 
					   sent_message = FALSE, 
					   sent_source_language = NULL,
					   sent_target_language = NULL,
					   active_translation = FALSE
					   WHERE user_id = ?`, userID)

	if err != nil {
		return fmt.Errorf("failed to reset translation settings: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no settings found to reset for user_id: %d", userID)
	}

	return nil
}

func GetTranslationSetting(userID int) (*TranslationSetting, error) {

	setting := &TranslationSetting{}
	err := db.QueryRow(`SELECT user_id, sent_source_language,
	                    sent_target_language,sent_message, active_translation 
	                    FROM translation
					    WHERE user_id = ?`, userID).
		Scan(&setting.UserID, &setting.SourceLanguage,
			&setting.TargetLanguage, &setting.SentMessage,
			&setting.ActiveTranslation)
	if err != nil {
		return nil, fmt.Errorf("failed to fetching translation setting in db: %v", err)
	}

	return setting, nil

}
