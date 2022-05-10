package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App
	TelegramBot `mapstructure:",squash"`
	Database    `mapstructure:",squash"`
	Messages
}

type App struct {
	NumOfDecks int
}

type TelegramBot struct {
	APIKey string `mapstructure:"BOT_API_KEY"`
}

type Database struct {
	Scheme   string `mapstructure:"MONGODB_SCHEME"`
	Host     string `mapstructure:"MONGODB_HOST"`
	Username string `mapstructure:"MONGODB_USER"`
	Password string `mapstructure:"MONGODB_PASSWORD"`
}

type Messages struct {
	ChatAlreadyRegistered   string
	ChatNotExists           string
	ChatCreatedSuccessfully string
	ChatHasActiveGame       string
	ChatHasNoActiveGame     string

	Blackjack       string
	GameOver        string
	Win             string
	Lose            string
	Push            string
	BlackjackResult string

	PlayerCantDraw       string
	PlayerHand           string
	PlayerHandBusted     string
	PlayerAlreadyStopped string
	PlayerAlreadyBusted  string
	StoppedDrawing       string

	DealerHand      string
	DealerBlackjack string
	GameEnterHint   string

	TopPlayers string
}

func NewConfig(envPath, configFolder string) *Config {
	var config Config

	// Read from yml
	viper.AddConfigPath(configFolder)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	if envPath == "" {
		config.readEnvVarsFromSystem()
	} else {
		config.readEnvVarsFromFile(envPath)
	}

	return &config
}

func (c *Config) readEnvVarsFromFile(envPath string) {
	viper.AddConfigPath(envPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	if err := viper.Unmarshal(c); err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) readEnvVarsFromSystem() {
	c.Database.Scheme = os.Getenv("MONGODB_SCHEME")
	c.Database.Host = os.Getenv("MONGODB_HOST")
	c.Database.Username = os.Getenv("MONGODB_USER")
	c.Database.Password = os.Getenv("MONGODB_PASSWORD")
	c.TelegramBot.APIKey = os.Getenv("BOT_API_KEY")
}
