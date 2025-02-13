package postgre

import (
	"context"
	"database/sql"
	"eshop/internal/infrastructure/constants"
	"fmt"
)

// Optionally could be created models (for gorm e.g)
// TODO: NewDB to pkg, as argument take cfg, validate data in constructor

type DBinteraction interface {
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type DBconfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	SSLmode  string `yaml:"ssl_mode"`
}

func NewPostgres(cfg DBconfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=public",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLmode)

	db, err := sql.Open(constants.PostgresDriver, connStr)
	if err != nil {
		return nil, fmt.Errorf("NewPostgres: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("NewPostgres: %w", err)
	}

	return db, nil
}
