package migrations

import (
	"database/sql"
	"embed"
	"eshop/internal/infrastructure/constants"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"io/fs"
)

//go:embed postgres-migs
var EmbedFS embed.FS

func PostgreMigrate(db *sql.DB, migFiles fs.FS) error {
	goose.SetBaseFS(migFiles)

	if err := goose.SetDialect(constants.PostgresDriver); err != nil {
		return fmt.Errorf("driver's fault: %w", err)
	}

	if err := goose.Up(db, "postgres-migs"); err != nil {
		return fmt.Errorf("failed migrations executing: %w", err)
	}

	return nil
}
