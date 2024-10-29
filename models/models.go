package models

import (
	"time"
)

type Game struct {
	ID       int64     // Unique identifier
	ChatId   int64     // Chat ID of the game
	Date     time.Time // Date of the game
	Location string    // Location of the game
	Opponent string    // Opponent for the game
	Players  []Player
}

type Player struct {
	ID       int64  // Unique identifier
	PlayerId int64  // Telegram ID of the player
	Name     string // Name of the player
}
