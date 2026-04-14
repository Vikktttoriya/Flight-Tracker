package migrator

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	migrator      *migrate.Migrate
	dbURL         string
	migrationPath string
}

func NewMigrator(dbURL, migrationPath string) (*Migrator, error) {
	migrator, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{
		migrator:      migrator,
		dbURL:         dbURL,
		migrationPath: migrationPath,
	}, nil
}

func (mg *Migrator) Up() error {
	if err := mg.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	log.Println("migrations applied")
	return nil
}

func (m *Migrator) Down() error {
	if err := m.migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}
	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	return m.migrator.Version()
}

func (m *Migrator) Close() error {
	if m.migrator != nil {
		sourceErr, dbErr := m.migrator.Close()
		if sourceErr != nil {
			return fmt.Errorf("source error: %w", sourceErr)
		}
		if dbErr != nil {
			return fmt.Errorf("database error: %w", dbErr)
		}
	}
	return nil
}
