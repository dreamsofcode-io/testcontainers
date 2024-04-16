package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(uri string) (*sql.DB, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return db, nil
}

func Migrate(uri string) error {
	path, exists := os.LookupEnv("MIGRATIONS_PATH")
	if !exists {
		path = "file://migrations"
	}

	m, err := migrate.New(path, uri)
	if err != nil {
		return fmt.Errorf("failed to connect migrator: %w", err)
	}

	// Migrate all the way up ...
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate up: %w", err)
	}
	return nil
}
