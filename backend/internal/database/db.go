package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

// Open connects to Turso (if TURSO_DATABASE_URL is set) or falls back to local SQLite.
func Open(dbPath string) (*sql.DB, error) {
	tursoURL := os.Getenv("TURSO_DATABASE_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	var db *sql.DB
	var err error

	if tursoURL != "" {
		// Turso remote database
		dsn := tursoURL
		if tursoToken != "" {
			dsn += "?authToken=" + tursoToken
		}
		db, err = sql.Open("libsql", dsn)
		if err != nil {
			return nil, fmt.Errorf("open turso db: %w", err)
		}
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		fmt.Println("Database: Turso (remote)")
	} else {
		// Local SQLite fallback
		dsn := fmt.Sprintf("%s?_pragma=journal_mode(wal)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(on)", dbPath)
		db, err = sql.Open("sqlite", dsn)
		if err != nil {
			return nil, fmt.Errorf("open local db: %w", err)
		}
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		fmt.Println("Database: Local SQLite")
	}

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
