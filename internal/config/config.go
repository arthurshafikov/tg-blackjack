package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseConfig    `mapstructure:",squash"`
	TelegramBotConfig `mapstructure:",squash"`
	Messages
}

type TelegramBotConfig struct {
	APIKey string `mapstructure:"BOT_API_KEY"`
}

type DatabaseConfig struct {
	Scheme   string `mapstructure:"MONGODB_SCHEME"`
	Host     string `mapstructure:"MONGODB_HOST"`
	Username string `mapstructure:"MONGODB_USER"`
	Password string `mapstructure:"MONGODB_PASSWORD"`
}

type Messages struct {
	Start string
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
	c.DatabaseConfig.Scheme = os.Getenv("MONGODB_SCHEME")
	c.DatabaseConfig.Host = os.Getenv("MONGODB_HOST")
	c.DatabaseConfig.Username = os.Getenv("MONGODB_USER")
	c.DatabaseConfig.Password = os.Getenv("MONGODB_PASSWORD")
	c.TelegramBotConfig.APIKey = os.Getenv("BOT_API_KEY")
}
