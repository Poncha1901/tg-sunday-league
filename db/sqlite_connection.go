package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Sqllite struct {
	DB      *sql.DB
	Connect func() (*sql.DB, error)
}

func Connect(db_path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}
	return db, nil
}
