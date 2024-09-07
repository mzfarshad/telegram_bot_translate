package config

import (
	"errors"
	"os"
)

// TokenFromENV retrieves the Telegram bot token from the environment variables.

func TokenFromENV() (string, error) {
	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		return "", errors.New("please set token at .env file")
	}

	return token, nil
}
