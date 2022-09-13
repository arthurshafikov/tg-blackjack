package app

import (
	"context"
	"flag"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/logger"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
	"github.com/arthurshafikov/tg-blackjack/internal/repository/mongodb"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers"
	"github.com/arthurshafikov/tg-blackjack/internal/telegram/handlers/commands"
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
	logger := logger.NewLogger()

	mongo, err := mongodb.NewMongoDB(ctx, mongodb.Config{
		Scheme:   config.Database.Scheme,
		Host:     config.Database.Host,
		Username: config.Database.Username,
		Password: config.Database.Password,
	})
	if err != nil {
		logger.Error(err)
		return
	}

	repository := repository.NewRepository(mongo)

	services := services.NewServices(services.Deps{
		Config:     config,
		Repository: repository,
		Logger:     logger,
	})

	botAPI, err := tgbotapi.NewBotAPI(config.TelegramBot.APIKey)
	if err != nil {
		logger.Error(err)
		return
	}
	telegramHelper := handlers.NewHelper(botAPI)
	commandHandler := commands.NewCommandHandler(handlers.HandlerParams{
		Ctx:      ctx,
		Services: services,
		Logger:   logger,
		Messages: config.Messages,
		Helper:   telegramHelper,
	})
	telegramBot := telegram.NewBot(&telegram.Deps{
		Ctx:      ctx,
		Services: services,
		Logger:   logger,
		Messages: config.Messages,

		CommandHandler: commandHandler,

		Helper: telegramHelper,
	})

	if err := telegramBot.Start(); err != nil {
		logger.Error(err)
	}
}
