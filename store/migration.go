package store

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // to support file:// scheme
)

const (
	migrationFileLocation = "file://store/schema"
	MIGRATION_CMD_UP      = "up"
)

type ErrInvalidMigrationCMD struct {
	InvalidCMD string
}

func (e *ErrInvalidMigrationCMD) Error() string {
	return fmt.Sprintf("invalid migration command: %s", e.InvalidCMD)
}

func Migrate(db *sql.DB, cmd string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to open db driver: %w", err)
	}
	migration, err := migrate.NewWithDatabaseInstance(migrationFileLocation, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	switch cmd {
	case MIGRATION_CMD_UP:
		err = migration.Up()
		if err != nil {
			if err == migrate.ErrNoChange {
				return nil
			}
			return fmt.Errorf("failed to migration up: %w", err)
		}
		return nil
	default:
		return &ErrInvalidMigrationCMD{InvalidCMD: cmd}
	}
}
