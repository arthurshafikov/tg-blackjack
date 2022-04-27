package app

import (
	"context"
	"flag"
	"log"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/logger"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
	"github.com/arthurshafikov/tg-blackjack/internal/repository/mongodb"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	envPath          string
	configFolderPath string
)

func init() {
	flag.StringVar(&envPath, "env", "", "Path to .env file folder")
	flag.StringVar(&configFolderPath, "cfgFolder", "", "Path to configs folder")
}

func Run() {
	flag.Parse()

	ctx := context.Background()
	config := config.NewConfig(envPath, configFolderPath)

	botAPI, err := tgbotapi.NewBotAPI(config.TelegramBotConfig.APIKey)
	if err != nil {
		log.Fatal(err)
	}

	mongo, err := mongodb.NewMongoDB(ctx, mongodb.Config{
		Scheme:   config.DatabaseConfig.Scheme,
		Host:     config.DatabaseConfig.Host,
		Username: config.DatabaseConfig.Username,
		Password: config.DatabaseConfig.Password,
	})
	if err != nil {
		log.Fatalln(err)
	}

	repository := repository.NewRepository(mongo)

	logger := logger.NewLogger()
	services := services.NewServices(services.Deps{
		Repository: repository,
		Logger:     logger,
	})

	telegramBot := telegram.NewBot(ctx, botAPI, services, config.Messages)

	if err := telegramBot.Start(); err != nil {
		log.Fatalln(err)
	}
}
