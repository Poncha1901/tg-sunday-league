package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"tg-sunday-league/models"
	"time"
)

type GameRepository struct {
	Db *sql.DB
}

// CreateGame inserts a new game into the database
func (r *GameRepository) CreateGame(game *models.Game) error {
	if r.Db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// Begin a new transaction
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	stmt, err := tx.Prepare("INSERT INTO games (chat_id, opponent, location, date, created_at, is_active) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(game.ChatId, game.Opponent, game.Location, game.Date, time.Now(), true)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	log.Println("Game created and committed successfully.")
	return nil
}

func (r *GameRepository) GetLatestGameByChatID(chatID int64) (*models.Game, error) {
	if r.Db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	// Prepare the SQL statement
	stmt, err := r.Db.Prepare("SELECT id, chat_id, opponent, location, date FROM games WHERE chat_id = ? AND is_active = 1 ORDER BY date DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	row := stmt.QueryRow(chatID)

	// Scan the result into a Game struct
	game := &models.Game{}
	err = row.Scan(&game.ID, &game.ChatId, &game.Opponent, &game.Location, &game.Date)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *GameRepository) CreatePlayer(player *models.Player) (int64, error) {
	if r.Db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}

	tx, err := r.Db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO players (player_id, name) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(player.PlayerId, player.Name)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	playerID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return playerID, nil
}

func (r *GameRepository) RegisterPlayerToGame(game *models.Game, player *models.Player) (int64, error) {
	if r.Db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}

	// Begin a new transaction
	tx, err := r.Db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare("INSERT INTO game_attendances (game_id, player_id) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(game.ID, player.ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	log.Println("Player registered to game and committed successfully.")
	return player.ID, nil
}

func (r *GameRepository) GetPlayerById(playerId *int64) (int64, error) {
	if r.Db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}
	stmt, err := r.Db.Prepare("SELECT id FROM players WHERE player_id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(playerId)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *GameRepository) GetPlayerForGame(playerId int64, gameId int64) (int64, error) {
	if r.Db == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}

	stmt, err := r.Db.Prepare("SELECT player_id FROM game_players WHERE player_id = ? AND game_id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(playerId, gameId)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil

}
