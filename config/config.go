package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken      string
	SqlliteDbPath string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var botToken string = os.Getenv("API_KEY")

	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}
	var sqlliteDbPath string = os.Getenv("SQL_LITE_DB_PATH")

	if sqlliteDbPath == "" {
		return nil, fmt.Errorf("SQL_LITE_DB_PATH is not set")
	}

	return &Config{
		BotToken:      botToken,
		SqlliteDbPath: sqlliteDbPath,
	}, nil
}
