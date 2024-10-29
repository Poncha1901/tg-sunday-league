package main

import (
	"log"
	"tg-sunday-league/bot"
	"tg-sunday-league/config"
	"tg-sunday-league/db"
	"tg-sunday-league/repositories"
	"tg-sunday-league/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Connect to the SQLite database
	dbInstance, dbErr := db.Connect()
	if dbErr != nil {
		log.Fatalf("Could not connect to database: %v", dbErr)
	}
	defer dbInstance.Close() // Ensure the DB connection closes when main exits

	// Setup database tables
	if err := db.SetupDatabase(dbInstance); err != nil {
		log.Fatalf("Could not setup database: %v", err)
	}

	// Initialize GameRepository
	gameRepo := &repositories.GameRepository{Db: dbInstance}
	gameService := &services.GameService{GameRepository: gameRepo}

	// Start the bot with repository dependency
	b, err := bot.NewBot(cfg.BotToken, gameService)
	if err != nil {
		log.Fatalf("Could not create bot: %v", err)
	}

	log.Println("Bot is running...")
	b.TelegramBot.Start()
}
