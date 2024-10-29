package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Sqllite struct {
	DB      *sql.DB
	Connect func() (*sql.DB, error)
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "sunday_league.sqlite")
	if err != nil {
		return nil, err
	}
	return db, nil
}
