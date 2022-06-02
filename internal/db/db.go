package db

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"

	"webChat/internal/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var (
	//go:embed migrations/*.sql
	fs embed.FS
	//go:embed migrations/*migratable*.sql
	migrateDownFS embed.FS
)

func NewConnection(cfg *config.Database) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	if err = migrateUp(cfg.URL); err != nil {
		return nil, fmt.Errorf("failed to migrate db up: %w", err)
	}

	if cfg.MigrateDown {
		if err = migrateDown(cfg.URL); err != nil {
			return nil, fmt.Errorf("failed to migrate db down: %w", err)
		}
	}

	return db, nil
}

func migrateUp(url string) error {
	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", sourceInstance, url)
	if err != nil {
		return fmt.Errorf("failed to create new migrate instance: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations up: %w", err)
	}

	return nil
}

func migrateDown(url string) error {
	sourceInstance, err := iofs.New(migrateDownFS, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", sourceInstance, url)
	if err != nil {
		return fmt.Errorf("failed to create new migrate instance: %w", err)
	}

	if err = m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations down: %w", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations up: %w", err)
	}

	return nil
}
