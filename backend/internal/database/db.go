package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Open creates a SQLite connection with WAL mode and busy timeout.
func Open(dbPath string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?_pragma=journal_mode(wal)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(on)", dbPath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Set connection pool limits for SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

// Migrate runs the schema migrations.
func Migrate(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
