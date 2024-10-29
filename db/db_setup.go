package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SetupDatabase initializes the SQLite database and creates the required tables
func SetupDatabase(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY,
			chat_id INTEGER64 NOT NULL,
			opponent TEXT NOT NULL,
			location TEXT NOT NULL,
			date TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_active BOOLEAN DEFAULT TRUE
		);`,
		`CREATE TABLE IF NOT EXISTS players (
			id TEXT PRIMARY KEY, -- UUID stored as TEXT
			player_id INTEGER64 NOT NULL,
			name TEXT NOT NULL,
			is_active BOOLEAN DEFAULT TRUE
		);`,
		`CREATE TABLE IF NOT EXISTS game_attendance (
			game_id INTEGER,
			player_id TEXT,
			has_paid BOOLEAN DEFAULT FALSE,
			FOREIGN KEY (game_id) REFERENCES games(id),
			FOREIGN KEY (player_id) REFERENCES players(id),
			PRIMARY KEY (game_id, player_id)
		);`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Printf("Error creating table: %v", err)
			return err
		}
	}

	log.Println("Database setup completed successfully.")
	return nil
}
