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

type PlayerStatus string

const (
	ATTENDING PlayerStatus = "ATTENDING"
	OUT       PlayerStatus = "OUT"
)

type IGameService interface {
	CreateNewGame(chatId int64, userId int64, userName string, gameData []string) (*models.Game, *[]models.User, *[]models.User, error)
	RegisterPlayer(chatId *int64, userId *int64, userName *string, status PlayerStatus) (*models.Game, *[]models.User, *[]models.User, error)
	GetGameDetails(chatId int64) (*models.Game, *[]models.User, *[]models.User, error)
	RepayGame(chatId *int64, userId *int64, userName *string) (*models.Game, *[]models.User, *[]models.User, error)
}

type GameService struct {
	GameRepository repositories.IGameRepository
}

func (g *GameService) CreateNewGame(chatId int64, userId int64, userName string, gameData []string) (*models.Game, *[]models.User, *[]models.User, error) {

	userFound, err := g.GameRepository.GetUserByUserID(userId)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not retrieve user, please try again.")
	}
	log.Printf("User found: %v", userFound)
	if userFound == nil {
		newUser := &models.User{
			Id:     uuid.New(),
			UserId: userId,
			Name:   userName,
		}
		_, err = g.GameRepository.InsertUser(newUser)
		userFound = newUser
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return nil, nil, nil, fmt.Errorf("Could not create user, please try again.")
		}
	}

	prev_game, err := g.GameRepository.GetLatestGameByChatID(chatId)
	if err != nil {
		log.Printf("Could not find the latest game: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not find the latest game, please try again.")
	}

	if prev_game != nil && prev_game.Date.After(time.Now()) {
		return nil, nil, nil, fmt.Errorf("There is already a game scheduled on %s against %s",
			prev_game.Date.Format("2006-01-02 15:04"),
			prev_game.Opponent)
	}

	dateStr, timeStr := gameData[0], gameData[1]
	dateTimeStr := fmt.Sprintf("%s %s", dateStr, timeStr)
	dateTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Invalid date or time format. Please use YYYY-MM-DD HH:MM.")
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
			return nil, nil, nil, fmt.Errorf("Invalid price format. Please provide a valid number.")
		}
		game.Price = price
	}

	_, err = g.GameRepository.InsertGame(game)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not create game, please try again. %v", err)
	}

	var players *[]models.User
	var absentees *[]models.User

	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not retrieve game details, please try again.")
	}
	game, players, absentees, err = g.GetGameDetails(chatId)
	log.Printf("Game created successfully: %v", game)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not retrieve game details, please try again.")
	}
	return game, players, absentees, nil

}

func (g *GameService) RegisterPlayer(chatID *int64, userId *int64, userName *string, status PlayerStatus) (*models.Game, *[]models.User, *[]models.User, error) {

	game, err := g.GameRepository.GetLatestGameByChatID(*chatID)
	if err != nil {
		log.Printf("Could not find the latest game: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not find the latest game, please try again.")
	}

	idFound, _ := g.GameRepository.GetUserByUserID(*userId)

	var player *models.User
	if idFound == nil {
		player := &models.User{
			Id:     uuid.New(),
			UserId: *userId,
			Name:   *userName,
			Status: string(status),
		}
		*userId, err = g.GameRepository.InsertUser(player)
		if err != nil {
			log.Printf("Error creating player: %v", err)
			return nil, nil, nil, fmt.Errorf("Could not create player, please try again.")
		}
	} else {
		player = idFound
		player.Status = string(status)
	}

	playerForGameId, _ := g.GameRepository.GetPlayerForGame(player.Id, game.Id)

	if playerForGameId != nil {
		g.GameRepository.UpdatePlayerGameStatus(game.Id, player.Id, string(status))
		game, players, absentees, err := g.GetGameDetails(*chatID)
		if err != nil {
			log.Printf("Error retrieving game details: %v", err)
			return nil, nil, nil, fmt.Errorf("Could not retrieve game details, please try again.")
		}
		return game, players, absentees, nil
	}

	_, err = g.GameRepository.InsertGamePlayer(game, player)
	if err != nil {
		log.Printf("Error registering player to game: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not register player, please try again.")
	}
	game, players, absentees, err := g.GetGameDetails(*chatID)
	if err != nil {
		log.Printf("Error retrieving game details: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not retrieve game details, please try again.")
	}
	return game, players, absentees, nil
}

func (g *GameService) GetGameDetails(chatID int64) (*models.Game, *[]models.User, *[]models.User, error) {
	game, err := g.GameRepository.GetLatestGameByChatID(chatID)
	if err != nil || game == nil {
		log.Printf("Error retrieving game details: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not retrieve game details, please try again.")
	}

	allPlayers, err := g.GameRepository.GetGamePlayers(game.Id)
	var players []models.User
	var absentees []models.User
	for _, player := range allPlayers {
		if player.Status == "OUT" {
			absentees = append(absentees, player)
		}
		if player.Status == "ATTENDING" {
			players = append(players, player)
		}
	}

	if err != nil {
		log.Printf("Error retrieving game players: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not retrieve game players, please try again.")
	}
	return game, &players, &absentees, nil
}

func (g *GameService) RepayGame(chatID *int64, userId *int64, userName *string) (*models.Game, *[]models.User, *[]models.User, error) {
	game, err := g.GameRepository.GetLatestGameByChatID(*chatID)
	if err != nil {
		log.Printf("Could not find the latest game: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not find the latest game, please try again.")
	}

	player, err := g.GameRepository.GetUserByUserID(*userId)
	if err != nil {
		log.Printf("Could not find the player: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not find the player, please try again.")
	}

	playerForGameId, err := g.GameRepository.GetPlayerForGame(player.Id, game.Id)
	if err != nil {
		log.Printf("Could not find the player for the game: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not find the player for the game, please try again.")
	}

	if playerForGameId == nil {
		log.Printf("Player not registered for the game")
		return nil, nil, nil, fmt.Errorf("%s is not registered for the game.", player.Name)
	}

	err = g.GameRepository.UpdatePlayerPayment(game.Id, player.Id)
	if err != nil {
		log.Printf("Could not update player payment: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not update player payment, please try again.")
	}

	game, players, absentees, err := g.GetGameDetails(*chatID)
	if err != nil {
		log.Printf("Error retrieving game details: %v", err)
		return nil, nil, nil, fmt.Errorf("Could not retrieve game details, please try again.")
	}
	return game, players, absentees, nil
}
