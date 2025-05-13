package migrator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Van-programan/Forum_GO/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	migrate *migrate.Migrate
	logger  logger.Interface
}

func NewMigrator(dbURL, migrationsPath string, logger logger.Interface) *Migrator {
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		logger.Error("migrations directory does not exist: %s", migrationsPath)
		return nil
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		logger.Error("failed to get migrations path: %w", err)
		return nil
	}

	logger.Info("Initializing migrator", "path", absPath)

	sourceURL := fmt.Sprintf("file://%s", absPath)
	fullDBURL := fmt.Sprintf("%s?sslmode=disable", dbURL)

	var m *migrate.Migrate
	attempts := 3
	for i := 0; i < attempts; i++ {
		m, err = migrate.New(sourceURL, fullDBURL)
		if err == nil {
			break
		}

		logger.Warn(fmt.Sprintf("Migration connection attempt failed (attempt %d/%d): %v", i+1, attempts, err))
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Error(fmt.Errorf("failed to initialize migrator: %w", err))
		return nil
	}

	logger.Info("Migrator initialized successfully")
	return &Migrator{migrate: m, logger: logger}
}

func (m *Migrator) Up() {
	if err := m.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("No new migrations to apply")
		}

		m.logger.Error(fmt.Errorf("failed to apply migrations: %w", err))
	}

	m.logger.Info("Migrations applied successfully")
}

func (m *Migrator) Down() {
	if err := m.migrate.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("No migrations to rollback")
		}

		m.logger.Error("Failed to rollback migrations")
	}

	m.logger.Info("Migrations rolled back successfully")
}

func (m *Migrator) Close() {
	if m.migrate != nil {
		if sourceErr, dbErr := m.migrate.Close(); sourceErr != nil || dbErr != nil {
			m.logger.Error("source error: %v, database error: %v", sourceErr, dbErr)
		}
	}
}
