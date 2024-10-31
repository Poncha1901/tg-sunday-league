package repositories

import (
	"database/sql"
	"log"
	"tg-sunday-league/models"
	"time"

	"github.com/google/uuid"
)

type GameRepository struct {
	Db *sql.DB
}

// InsertGame inserts a new game into the database
func (r *GameRepository) InsertGame(game *models.Game) (*models.Game, error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return nil, err
	}

	// Prepare the SQL statement
	stmt, err := tx.Prepare(`
		INSERT INTO games (
			id, 
			chat_id,
			opponent, 
			location, 
			date, 
			price, 
			created_at, 
			created_by, 
			is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(&game.Id, &game.ChatId, &game.Opponent,
		&game.Location, &game.Date, &game.Price, time.Now(), &game.CreatedBy, true)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	log.Println("Game created and committed successfully.")
	return game, nil
}

func (r *GameRepository) GetLatestGameByChatID(chatID int64) (*models.Game, error) {
	stmt, err := r.Db.Prepare(
		`SELECT 
			id, 
			chat_id, 
			opponent, 
			location, 
			price,
			date FROM games 
		WHERE chat_id = ? 
		AND is_active = 1 
		ORDER BY date DESC LIMIT 1`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(chatID)

	game := &models.Game{}
	err = row.Scan(&game.Id, &game.ChatId, &game.Opponent, &game.Location, &game.Price, &game.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return game, nil
}

func (r *GameRepository) InsertUser(user *models.User) (int64, error) {

	tx, err := r.Db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO users (
			id, user_id, name)
		VALUES (?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(&user.Id, &user.UserId, &user.Name)
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

func (r *GameRepository) InsertGamePlayer(game *models.Game, player *models.User) (string, error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return "", err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO game_players (
				game_id, 
				user_id,
				status,
				has_paid
				) VALUES 
				(?, ?, ?, ?)
		`)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&game.Id, &player.Id, "ATTENDING", false)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return "", err
	}

	log.Println("User registered to game and committed successfully.")
	return player.Name, nil
}

func (r *GameRepository) GetUserById(playerId *int64) (*models.User, error) {
	stmt, err := r.Db.Prepare(
		`SELECT id 
		FROM players 
		WHERE player_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(playerId)

	var user models.User
	err = row.Scan(&user.Id, &user.UserId, &user.Name)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *GameRepository) GetUserByUserID(userId int64) (*models.User, error) {
	stmt, err := r.Db.Prepare(
		`SELECT 
			id,
			user_id,
			name
		FROM users 
		WHERE user_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userId)

	var user models.User
	err = row.Scan(&user.Id, &user.UserId, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *GameRepository) GetPlayerForGame(playerId uuid.UUID, gameId uuid.UUID) (*uuid.UUID, error) {
	stmt, err := r.Db.Prepare("SELECT user_id FROM game_players WHERE user_id = ? AND game_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id uuid.UUID
	row := stmt.QueryRow(playerId.String(), gameId.String())
	err = row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No player found for game_id %s and player_id %s", gameId.String(), playerId.String())
			return nil, nil
		}
		return nil, err
	}

	return &id, nil
}

func (r *GameRepository) GetGamePlayers(gameId uuid.UUID) ([]models.User, error) {
	stmt, err := r.Db.Prepare(
		`SELECT 
			u.id, 
			u.user_id, 
			u.name,
			gp.status,
			gp.has_paid
		FROM users u 
		JOIN game_players gp 
		ON u.id = gp.user_id 
		WHERE gp.game_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(gameId.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.User
	for rows.Next() {
		var player models.User
		err := rows.Scan(&player.Id, &player.UserId, &player.Name, &player.Status, &player.HasPaid)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (r *GameRepository) UpdatePlayerPayment(gameId uuid.UUID, playerId uuid.UUID) error {
	stmt, err := r.Db.Prepare(
		`UPDATE game_players 
		SET has_paid = 1 
		WHERE game_id = ? 
		AND user_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(gameId.String(), playerId.String())
	if err != nil {
		return err
	}

	return nil

}
