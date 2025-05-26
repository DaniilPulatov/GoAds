package migrations

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func New(fileDir, dsn string) error {
	migration, err := migrate.New(fileDir, dsn)
	if err != nil {
		log.Println("Error creating new migration:", err)
		return fmt.Errorf("new migration: %w", err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println("Error up new migration:", err)
		return fmt.Errorf("up migration: %w", err)
	}

	log.Println("Migration completed successfully")

	return nil
}
