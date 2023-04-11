package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

//go:embed sql_migrations
var sqlMigrations embed.FS

type Migration struct {
	ID      string
	Migrate func(tx *sql.Tx) error
}

func Migrate(db *sql.DB) error {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id TEXT PRIMARY KEY)"); err != nil {
		return fmt.Errorf("could not create migrations table: %w", err)
	}

	var migrations []Migration

	// Those could be supported, but aren't right now, as generating the schema for sqlc is too much work then.
	// That said, as long as they don't include DLL, they're fine, really.
	// migrations = append(migrations, goMigrations...)

	entries, err := sqlMigrations.ReadDir("sql_migrations")
	if err != nil {
		return fmt.Errorf("could not read sql migration list: %w", err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		data, err := sqlMigrations.ReadFile("sql_migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("could not read sql migration file: %w", err)
		}
		migrations = append(migrations, Migration{
			ID: strings.TrimSuffix(entry.Name(), ".sql"),
			Migrate: func(tx *sql.Tx) error {
				if _, err := tx.Exec(string(data)); err != nil {
					return fmt.Errorf("could not execute sql migration: %w", err)
				}
				return nil
			},
		})
	}

	slices.SortFunc(migrations, func(a, b Migration) bool {
		return a.ID < b.ID
	})

	for _, migration := range migrations {
		if _, err := db.Exec("INSERT INTO migrations (id) VALUES (?)", migration.ID); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				log.Printf("Skipping migration %s, already applied.", migration.ID)
				continue
			}
			return fmt.Errorf("could not insert migration id: %w", err)
		}
		log.Printf("Applying migration %s.", migration.ID)
		if err := func() error {
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("could not begin transaction: %w", err)
			}
			defer tx.Rollback()
			if err := migration.Migrate(tx); err != nil {
				return fmt.Errorf("could not execute migration: %w", err)
			}
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("could not commit transaction: %w", err)
			}
			return nil
		}(); err != nil {
			return fmt.Errorf("could not apply migration: %w", err)
		}
	}

	return nil
}
