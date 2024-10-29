package services

import (
	"fmt"
	"log"
	"strings"
	"tg-sunday-league/models"
	"tg-sunday-league/repositories"
	"time"
)

type GameService struct {
	GameRepository *repositories.GameRepository
}

func (g *GameService) CreateNewGame(chatId int64, gameData []string) (string, error) {

	dateStr, timeStr := gameData[0], gameData[1]
	dateTimeStr := fmt.Sprintf("%s %s", dateStr, timeStr)
	dateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return "", fmt.Errorf("Invalid date or time format. Please use YYYY-MM-DD HH:MM.")
	}

	location := gameData[2]
	opponent := strings.Join(gameData[3:], " ")
	game := &models.Game{
		Date:     dateTime,
		Location: location,
		Opponent: opponent,
		ChatId:   chatId,
	}
	err = g.GameRepository.CreateGame(game)
	if err != nil {
		return "", fmt.Errorf("Invalid date or time format. Please use YYYY-MM-DD HH:MM.")
	}
	return fmt.Sprintf("Game created!\nDate: %s\nLocation: %s\nTime: %s\nOpponent: %s",
		game.Date.Format("2006-01-02"), game.Location, game.Date.Format("15:04"), game.Opponent), nil
}

func (g *GameService) RegisterPlayer(chatID *int64, playerId *int64, playerName *string) (string, error) {

	game, err := g.GameRepository.GetLatestGameByChatID(*chatID)
	if err != nil {
		log.Printf("Could not find the latest game: %v", err)
		return "", fmt.Errorf("Could not find the latest game, please try again.")
	}

	player := &models.Player{
		PlayerId: *playerId,
		Name:     *playerName,
	}

	idFound, _ := g.GameRepository.GetPlayerById(playerId)

	if idFound == 0 {
		*playerId, err = g.GameRepository.CreatePlayer(player)
		if err != nil {
			log.Printf("Error creating player: %v", err)
			return "", fmt.Errorf("Could not create player, please try again.")
		}
	} else {
		playerId = &idFound
	}

	playerForGameId, _ := g.GameRepository.GetPlayerForGame(*playerId, game.ID)

	if playerForGameId != 0 {
		log.Printf("Player already registered to the game")
		return "", fmt.Errorf("Player already registered to the game.")
	}

	_, err = g.GameRepository.RegisterPlayerToGame(game, player)
	if err != nil {
		log.Printf("Error registering player to game: %v", err)
		return "", fmt.Errorf("Could not register player, please try again.")
	}
	return fmt.Sprintf("Player %d registered to the game on %s", playerName, game.Date.Format("2006-01-02 15:04")), nil
}
