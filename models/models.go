package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	Id        uuid.UUID // Unique identifier
	ChatId    int64     // Chat ID of the game
	Price     float64   // Price of the game
	Date      time.Time // Date of the game
	Location  string    // Location of the game
	Opponent  string    // Opponent for the game
	Players   []User
	CreatedBy uuid.UUID
}

type User struct {
	Id      uuid.UUID // Unique identifier
	UserId  int64     // Telegram ID of the player
	Name    string    // Name of the player
	Status  string    // Status of the player (Attending, Not Attending, Paid)
	HasPaid bool      // Whether the player has paid
}
