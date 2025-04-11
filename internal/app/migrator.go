package app

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Migrator struct {
	db *sql.DB
	fs fs.FS
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db: db,
		fs: os.DirFS(filepath.Join("..", "..", "migrations")),
	}
}

func (m *Migrator) Apply() error {
	entries, err := fs.ReadDir(m.fs, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	var upFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".up.sql") {
			upFiles = append(upFiles, name)
		}
	}

	sort.Slice(upFiles, func(i, j int) bool {
		return getMigrationNumber(upFiles[i]) < getMigrationNumber(upFiles[j])
	})

	for _, file := range upFiles {
		content, err := fs.ReadFile(m.fs, file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		if _, err := m.db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute %s: %w", file, err)
		}
	}

	return nil
}

func (m *Migrator) Rollback() error {
	entries, err := fs.ReadDir(m.fs, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	var downFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".down.sql") {
			downFiles = append(downFiles, name)
		}
	}

	if len(downFiles) == 0 {
		return nil
	}

	sort.Slice(downFiles, func(i, j int) bool {
		return getMigrationNumber(downFiles[i]) > getMigrationNumber(downFiles[j])
	})

	lastDown := downFiles[0]
	content, err := fs.ReadFile(m.fs, lastDown)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", lastDown, err)
	}

	if _, err := m.db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute %s: %w", lastDown, err)
	}

	return nil
}

func getMigrationNumber(filename string) int {
	parts := strings.Split(filename, "_")
	if len(parts) == 0 {
		return 0
	}
	num, _ := strconv.Atoi(parts[0])
	return num
}
