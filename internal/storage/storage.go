package storage

//TODO: integrate migrate function with init

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func InitDB(url string) (*Storage, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed opening a db connection: %w", err)
	}

	if err = dbMigrate(db); err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func dbMigrate(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed setting dialect: %w", err)
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed migrations executing: %w", err)
	}

	return nil
}
