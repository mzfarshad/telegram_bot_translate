package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mzfarshad/tlg_bot/internal/bot"
	"github.com/mzfarshad/tlg_bot/internal/config"
)

// init function is executed before the main function.
// It loads environment variables from the .env file.
// If the file is not found, the program will terminate with a fatal error.
// func init() {
// 	if err := godotenv.Load("/app/.env"); err != nil {
// 		log.Fatal("failed to load .env file")
// 	}
// }

func main() {

	// Retrieve the bot token from the environment variables.
	// If the token is not available, the program will terminate with a panic.
	token, err := config.TokenFromENV()
	if err != nil {
		log.Panic(err)
	}

	// Create a new instance of the bot using the token.
	// If the bot cannot be initialized, the program will terminate with a panic
	bot, err := bot.NewBot(token)
	if err != nil {
		log.Panic(err)
	}

	// Enable debug mode for the bot's API.
	bot.API.Debug = true

	log.Printf("Authorized on account %s", bot.API.Self.UserName)

	router := gin.Default()

	router.POST("/webhook", func(ctx *gin.Context) {
		var update tgbotapi.Update

		if err := ctx.ShouldBindJSON(&update); err != nil {
			log.Printf("error parsing update: %v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		}

		bot.HandleUpdate(update)

		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if err := router.Run(":7171"); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
