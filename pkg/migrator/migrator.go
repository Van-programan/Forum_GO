package migrator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
)

type Migrator struct {
	migrate *migrate.Migrate
	logger  zerolog.Logger
}

func NewMigrator(dbURL, migrationsPath string, logger zerolog.Logger) *Migrator {
	log := logger.With().Str("component", "migrator").Logger()

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		log.Error().
			Str("path", migrationsPath).
			Msg("Migrations directory does not exist")
		return nil
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", migrationsPath).
			Msg("Failed to get absolute path to migrations")
		return nil
	}

	log.Info().
		Str("path", absPath).
		Msg("Initializing migrator")

	sourceURL := fmt.Sprintf("file://%s", absPath)
	fullDBURL := fmt.Sprintf("%s?sslmode=disable", dbURL)

	var m *migrate.Migrate
	attempts := 3
	for i := 0; i < attempts; i++ {
		m, err = migrate.New(sourceURL, fullDBURL)
		if err == nil {
			break
		}

		log.Warn().
			Err(err).
			Int("attempt", i+1).
			Int("max_attempts", attempts).
			Msg("Migration connection attempt failed")

		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Error().
			Err(err).
			Str("db_url", dbURL).
			Msg("Failed to initialize migrator")
		return nil
	}

	log.Info().Msg("Migrator initialized successfully")
	return &Migrator{
		migrate: m,
		logger:  log,
	}
}

func (m *Migrator) Up() {
	log := m.logger.With().Str("operation", "up").Logger()

	log.Info().Msg("Applying migrations...")

	if err := m.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("Database is up to date - no migrations applied")
			return
		}

		log.Error().
			Err(err).
			Msg("Failed to apply migrations")
		return
	}

	log.Info().Msg("Migrations applied successfully")
}

func (m *Migrator) Down() {
	log := m.logger.With().Str("operation", "down").Logger()

	log.Info().Msg("Rolling back migrations...")

	if err := m.migrate.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("No migrations to rollback")
			return
		}

		log.Error().
			Err(err).
			Msg("Failed to rollback migrations")
		return
	}

	log.Info().Msg("Migrations rolled back successfully")
}

func (m *Migrator) Close() {
	if m.migrate != nil {
		log := m.logger.With().Str("operation", "close").Logger()

		sourceErr, dbErr := m.migrate.Close()
		if sourceErr != nil || dbErr != nil {
			log.Error().
				Err(sourceErr).
				Err(dbErr).
				Msg("Error closing migrator")
		} else {
			log.Debug().Msg("Migrator closed successfully")
		}
	}
}
