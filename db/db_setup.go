package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SetupDatabase initializes the SQLite database and creates the required tables
func SetupDatabase(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			user_id INTEGER NOT NULL,
			name VARCHAR NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS games (
			id VARCHAR(36) PRIMARY KEY,
			chat_id INTEGER,
			opponent VARCHAR,
			location VARCHAR,
			price FLOAT,
			date DATE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			is_active BOOL,
			FOREIGN KEY (created_by) REFERENCES users(id)
	);`,
		`CREATE TABLE IF NOT EXISTS player_status (
			name VARCHAR PRIMARY KEY
	);`,
		`CREATE TABLE IF NOT EXISTS game_players (
			game_id VARCHAR(36),
			user_id VARCHAR(36),
			status VARCHAR,
			has_paid BOOL,
			joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (game_id) REFERENCES games(id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (status) REFERENCES player_status(name),
			PRIMARY KEY (game_id, user_id)
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
