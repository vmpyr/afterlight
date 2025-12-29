package store

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/schema.sql
var schemaSQL string

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	// mode=rwc: Read/Write/Create
	// _journal_mode=WAL: Write-Ahead Logging
	// _busy_timeout=5000: Wait 5s before failing if DB is locked
	dsn := fmt.Sprintf("file:%s?mode=rwc&_journal_mode=WAL&_busy_timeout=5000", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, fmt.Errorf("failed to apply schema: %w", err)
	}

	log.Println("Database connected and schema applied")
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) DB() *sql.DB {
	return s.db
}
