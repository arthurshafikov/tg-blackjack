package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	envFileContent = []byte(
		`MONGO_INITDB_DATABASE=homestead
		MONGODB_HOST=localhost
		MONGODB_USER=homestead
		MONGODB_PASSWORD=secret
		MONGODB_SCHEME=mongodb
		BOT_API_KEY=someKey`,
	)
	envFolderPath     = "./"
	configFolderPath  = "./fakeconfigs"
	configFileContent = []byte(`Messages:
  ChatAlreadyRegistered: "ChatAlreadyRegistered"
  ChatNotExists: "ChatNotExists"`)

	expectedDatabaseConfig = Database{
		Host:     "localhost",
		Scheme:   "mongodb",
		Username: "homestead",
		Password: "secret",
	}
	expectedTelegramBotConfig = TelegramBot{
		APIKey: "someKey",
	}
	expectedMessages = Messages{
		ChatAlreadyRegistered: "ChatAlreadyRegistered",
		ChatNotExists:         "ChatNotExists",
	}
)

func TestNewConfigFromEnvFile(t *testing.T) {
	createFakeConfigFileAndFolder(t)
	defer deleteFakeConfigFileAndFolder(t)
	createFakeEnvFile(t)
	defer deleteFakeEnvFile(t)

	config := NewConfig(envFolderPath, configFolderPath)

	require.Equal(t, expectedDatabaseConfig, config.Database)
	require.Equal(t, expectedTelegramBotConfig, config.TelegramBot)
	require.Equal(t, expectedMessages, config.Messages)
}

func TestNewConfigFromSystem(t *testing.T) {
	createFakeConfigFileAndFolder(t)
	defer deleteFakeConfigFileAndFolder(t)

	require.NoError(t, os.Setenv("MONGODB_SCHEME", "mongodb"))
	require.NoError(t, os.Setenv("MONGODB_HOST", "localhost"))
	require.NoError(t, os.Setenv("MONGODB_USER", "homestead"))
	require.NoError(t, os.Setenv("MONGODB_PASSWORD", "secret"))
	require.NoError(t, os.Setenv("BOT_API_KEY", "someKey"))

	config := NewConfig("", configFolderPath)

	require.Equal(t, expectedDatabaseConfig, config.Database)
	require.Equal(t, expectedTelegramBotConfig, config.TelegramBot)
	require.Equal(t, expectedMessages, config.Messages)
}

func createFakeEnvFile(t *testing.T) {
	t.Helper()
	if err := os.WriteFile("app.env", envFileContent, 0600); err != nil { //nolint:gofumpt
		t.Fatal(err)
	}
}

func deleteFakeEnvFile(t *testing.T) {
	t.Helper()
	if err := os.Remove("app.env"); err != nil {
		t.Fatal(err)
	}
}

func createFakeConfigFileAndFolder(t *testing.T) {
	t.Helper()
	if err := os.Mkdir(configFolderPath, 0777); err != nil { //nolint:gofumpt
		t.Fatal(err)
	}
	err := os.WriteFile(fmt.Sprintf("%s/main.yml", configFolderPath), configFileContent, 0600) //nolint:gofumpt
	if err != nil {
		t.Fatal(err)
	}
}

func deleteFakeConfigFileAndFolder(t *testing.T) {
	t.Helper()
	if err := os.RemoveAll(configFolderPath); err != nil {
		t.Fatal(err)
	}
}
