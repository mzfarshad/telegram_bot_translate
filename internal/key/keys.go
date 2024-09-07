package key

// Define custom types for various uses in the project
type Language string
type TextButton string
type TextMessage string
type MenuState string
type HandlerName string

// Constants defining language codes and button keys
const (
	LangEN Language = "en"
	LangFA Language = "fa"

	// Button key
	KeyTranslaion              TextButton = "translation"
	KeyHelp                    TextButton = "help"
	KeyContactUs               TextButton = "contactUs"
	KeySettings                TextButton = "settings"
	KeyBack                    TextButton = "back"
	KeySettingsLanguage        TextButton = "settingsLanguage"
	KeyTranslateSentMessage    TextButton = "translateSentMessage"
	KeyResetTranslationSetting TextButton = "resetTranslationSetting"
	KeyFinishSetup             TextButton = "finishSetup"
	KeyResetTranslateYes       TextButton = "resetTranslateYes"

	// Message keys
	MainMessage                        TextMessage = "mainMessage"
	StartMessage                       TextMessage = "startMessage"
	SettingMessage                     TextMessage = "settingsMessage"
	SettingLanguageMessage             TextMessage = "settingLanguageMessage"
	ChangeLanguageMessage              TextMessage = "changLanguageMessage"
	FailedChangeLanguageMessage        TextMessage = "failedChangeLanguageMessage"
	TranslationMenuMessage             TextMessage = "translationMenuMessage"
	SelectLanguagePairsMessage         TextMessage = "selectLanguagePairsMessage"
	TranslateFinishMessage             TextMessage = "translateFinishMessage"
	TranslateFinishSetupMessage        TextMessage = "translateFinishSetupMessage"
	ResetTranslateSettingMessage       TextMessage = "resetTranslateMessage"
	FinishResetTranslateSettingMessage TextMessage = "finishResetTranslateSettingMessage"

	// Menu states
	MenuMain                     MenuState = "main"
	MenuSetting                  MenuState = "setting"
	MenuSettingLanguage          MenuState = "settingLanguage"
	MenuTranslation              MenuState = "translationMenu"
	MenuTranslationLanguagePairs MenuState = "translationLanguagePairs"
	MenuFinishTranslateSetup     MenuState = "finishTranslateSetup"
	MenuResetTranslate           MenuState = "resetTranslate"
	MenuHelp                     MenuState = "help"

	// Handler names
	StartHandler HandlerName = "start"

	//Language pairs
	// EN_FA TextButton = "en-fa"
	// FA_EN TextButton = "fa-en"
	// FA_AR TextButton = "fa-ar"
	// AR_FA TextButton = "ar-fa"
	// EN_AR TextButton = "en-ar"
	// AR_EN TextButton = "ar-en"
)

// Map of button texts for different languages
var Buttons = map[Language]map[TextButton]string{
	LangEN: {
		KeyTranslaion:              "Translation",
		KeyContactUs:               "Contact us",
		KeyHelp:                    "Help",
		KeySettings:                "Settings",
		KeyBack:                    "Back",
		KeySettingsLanguage:        "Language",
		KeyTranslateSentMessage:    "Translate Sent Message",
		KeyResetTranslationSetting: "Reset Translation Settings",
		KeyFinishSetup:             "Finish Setup",
		// EN_AR:                      "English to Arabic",
		// EN_FA:                      "English to Persian",
		// FA_AR:                      "Persian to Arabic",
		// FA_EN:                      "Persian to English",
		// AR_EN:                      "Arabic to English",
		// AR_FA:                      "Arabic to Persian",
		KeyResetTranslateYes: "Yes",
	},
	LangFA: {
		KeyTranslaion:              "ترجمه",
		KeyHelp:                    "راهنما",
		KeyContactUs:               "ارتباط با ما",
		KeySettings:                "تنظیمات",
		KeyBack:                    "بازگشت",
		KeySettingsLanguage:        "زبان",
		KeyTranslateSentMessage:    "ترجمه پیام های ارسالی",
		KeyResetTranslationSetting: "بازنشانی تنظیمات ترجمه",
		KeyFinishSetup:             "اتمام تنظیمات",
		// EN_AR:                      "انگلیسی به عربی",
		// EN_FA:                      "انگلیسی به فارسی",
		// FA_AR:                      "فارسی به عربی",
		// FA_EN:                      "فارسی به انگلیسی",
		// AR_EN:                      "عربی به انگلیسی",
		// AR_FA:                      "عربی به فارسی",
		KeyResetTranslateYes: "بله",
	},
}

// Get the text for a button based on the language and button key
func GetKey(lang Language, button TextButton) string {
	return Buttons[lang][button]
}

// Map of messages for different languages
var Messages = map[Language]map[TextMessage]string{
	LangEN: {
		MainMessage:                        "Main Menu",
		StartMessage:                       "Welcome",
		SettingMessage:                     "Settings Menu",
		SettingLanguageMessage:             "Select bot language",
		ChangeLanguageMessage:              "Your language has been changed to English",
		FailedChangeLanguageMessage:        "Sorry, there is a problem changing the language",
		TranslationMenuMessage:             "Please select one of the keys:",
		SelectLanguagePairsMessage:         "Please separate the languages with the ( - ) symbol without a space",
		TranslateFinishMessage:             "Please save the translation settings:",
		TranslateFinishSetupMessage:        "The Translation settings are saved and activated",
		ResetTranslateSettingMessage:       "Do you want to reset the translation settings for sending messages?",
		FinishResetTranslateSettingMessage: "Settings reset successfully",
	},
	LangFA: {
		MainMessage:                        "منو اصلی",
		StartMessage:                       "خوش امدید",
		SettingMessage:                     "منو تنظیمات",
		SettingLanguageMessage:             "زبات بات را انتخاب کنید",
		ChangeLanguageMessage:              "زبان شما به فارسی تغییر یافت",
		FailedChangeLanguageMessage:        "متاسفانه  مشکلی  برای تغییر زبان برای وجود دارد",
		TranslationMenuMessage:             "لطفا یکی از کلیدها را انتخاب کنید:",
		SelectLanguagePairsMessage:         "لطفا زبان های مبدا و مقصد  را با نماد ` - ` بدون فاصله از یکدیگر جدا کنید",
		TranslateFinishMessage:             "لطفا تنظیمات ترجمه را ذخیره کنید:",
		TranslateFinishSetupMessage:        "تنظیمات ترجمه ذخیره و فعال شده است",
		ResetTranslateSettingMessage:       "ایا می خواهید تنظیمات ترجمه برای ارسال پیام را بازنشانی کنید؟",
		FinishResetTranslateSettingMessage: "تنظیمات با موفقیت بازنشانی شد",
	},
}

// Get the menu message based on the language and message key
func GetMenuMessage(lang Language, message TextMessage) string {
	return Messages[lang][message]
}
