package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"tg-sunday-league/models"
	"tg-sunday-league/repositories"
	"time"

	"github.com/google/uuid"
)

type GameService struct {
	GameRepository *repositories.GameRepository
}

func (g *GameService) CreateNewGame(chatId int64, userId int64, userName string, gameData []string) (string, error) {

	userFound, err := g.GameRepository.GetUserByUserID(userId)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return "", fmt.Errorf("Could not retrieve user, please try again.")
	}

	if userFound == nil {
		newUser := &models.User{
			Id:     uuid.New(),
			UserId: userId,
			Name:   userName,
		}
		_, err = g.GameRepository.InsertUser(newUser)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return "", fmt.Errorf("Could not create user, please try again.")
		}
	}

	prev_game, err := g.GameRepository.GetLatestGameByChatID(chatId)
	if err != nil {
		log.Printf("Could not find the latest game: %v", err)
		return "", fmt.Errorf("Could not find the latest game, please try again.")
	}

	if prev_game != nil && prev_game.Date.After(time.Now()) {
		return "", fmt.Errorf("There is already a game scheduled on %s against %s",
			prev_game.Date.Format("2006-01-02 15:04"),
			prev_game.Opponent)
	}

	dateStr, timeStr := gameData[0], gameData[1]
	dateTimeStr := fmt.Sprintf("%s %s", dateStr, timeStr)
	dateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return "", fmt.Errorf("Invalid date or time format. Please use YYYY-MM-DD HH:MM.")
	}

	location := gameData[2]
	opponent := strings.Join(gameData[3:], " ")
	game := &models.Game{
		Id:        uuid.New(),
		Date:      dateTime,
		Location:  location,
		Opponent:  opponent,
		ChatId:    chatId,
		CreatedBy: userFound.Id,
	}
	if len(gameData) > 4 {
		priceStr := gameData[4]
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return "", fmt.Errorf("Invalid price format. Please provide a valid number.")
		}
		game.Price = price
	}

	_, err = g.GameRepository.InsertGame(game)
	if err != nil {
		return "", fmt.Errorf("Could not create game, please try again. %v", err)
	}
	return fmt.Sprintf("Game created!\nDate: %s\nLocation: %s\nTime: %s\nOpponent: %s",
		game.Date.Format("2006-01-02"), game.Location, game.Date.Format("15:04"), game.Opponent), nil
}

func (g *GameService) RegisterPlayer(chatID *int64, userId *int64, userName *string) (string, error) {

	game, err := g.GameRepository.GetLatestGameByChatID(*chatID)
	if err != nil {
		log.Printf("Could not find the latest game: %v", err)
		return "", fmt.Errorf("Could not find the latest game, please try again.")
	}

	idFound, _ := g.GameRepository.GetUserByUserID(*userId)

	var player *models.User
	if idFound == nil {
		player := &models.User{
			Id:     uuid.New(),
			UserId: *userId,
			Name:   *userName,
		}
		*userId, err = g.GameRepository.InsertUser(player)
		if err != nil {
			log.Printf("Error creating player: %v", err)
			return "", fmt.Errorf("Could not create player, please try again.")
		}
	} else {
		player = idFound
	}

	playerForGameId, _ := g.GameRepository.GetPlayerForGame(player.Id, game.Id)

	if playerForGameId != nil {
		log.Printf("Player already registered to the game")
		return "", fmt.Errorf("%s already registered for the game.", player.Name)
	}

	_, err = g.GameRepository.InsertGamePlayer(game, player)
	if err != nil {
		log.Printf("Error registering player to game: %v", err)
		return "", fmt.Errorf("Could not register player, please try again.")
	}
	return fmt.Sprintf("%s registered to the game against %s on %s", player.Name, game.Opponent, game.Date.Format("2006-01-02 15:04")), nil
}

func (g *GameService) GetGameDetails(chatID int64) (string, error) {
	game, err := g.GameRepository.GetLatestGameByChatID(chatID)
	if err != nil {
		log.Printf("Error retrieving game details: %v", err)
		return "", fmt.Errorf("Could not retrieve game details, please try again.")
	}
	if game == nil {
		return "No game scheduled.", nil
	}

	players, err := g.GameRepository.GetGamePlayers(game.Id)

	if err != nil {
		log.Printf("Error retrieving game players: %v", err)
		return "", fmt.Errorf("Could not retrieve game players, please try again.")
	}

	playerList := ""
	for i, player := range players {
		playerList += fmt.Sprintf("%d. %s", i+1, player.Name)
		if player.HasPaid {
			playerList += " ✅"
		}
		playerList += "\n"
	}

	return fmt.Sprintf("Game against %s on %s at %s\nLocation: %s\nPrice: %.2f\nPlayers:\n%s",
		game.Opponent, game.Date.Format("2006-01-02"), game.Date.Format("15:04"), game.Location, game.Price, playerList), nil
}

func (g *GameService) RepayGame(chatID *int64, userId *int64, userName *string) (string, error) {
	game, err := g.GameRepository.GetLatestGameByChatID(*chatID)
	if err != nil {
		log.Printf("Could not find the latest game: %v", err)
		return "", fmt.Errorf("Could not find the latest game, please try again.")
	}

	player, err := g.GameRepository.GetUserByUserID(*userId)
	if err != nil {
		log.Printf("Could not find the player: %v", err)
		return "", fmt.Errorf("Could not find the player, please try again.")
	}

	playerForGameId, err := g.GameRepository.GetPlayerForGame(player.Id, game.Id)
	if err != nil {
		log.Printf("Could not find the player for the game: %v", err)
		return "", fmt.Errorf("Could not find the player for the game, please try again.")
	}

	if playerForGameId == nil {
		log.Printf("Player not registered for the game")
		return "", fmt.Errorf("%s is not registered for the game.", player.Name)
	}

	err = g.GameRepository.UpdatePlayerPayment(game.Id, player.Id)
	if err != nil {
		log.Printf("Could not update player payment: %v", err)
		return "", fmt.Errorf("Could not update player payment, please try again.")
	}

	return fmt.Sprintf("%s has paid for the game against %s on %s", player.Name, game.Opponent, game.Date.Format("2006-01-02 15:04")), nil
}
