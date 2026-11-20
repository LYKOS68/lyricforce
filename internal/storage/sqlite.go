// Package storage owns the local SQLite database: opening, schema migrations
// and the queries of the rhyme archive.
package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // pure-Go SQLite driver, no CGO required
)

// migrations are applied in order; index i upgrades the schema to version i+1.
// Never edit an existing entry — append a new one.
var migrations = []string{
	// v1: metadata
	`CREATE TABLE IF NOT EXISTS app_meta (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	INSERT OR IGNORE INTO app_meta (key, value) VALUES ('created_at', datetime('now'));`,

	// v2: rhyme words. Categorisation is phonetic and assigned by the user:
	// vowel = vowel class of the stressed vowel, rhyme = sound sequence from the stressed vowel.
	`CREATE TABLE words (
		id         INTEGER PRIMARY KEY,
		text       TEXT NOT NULL,
		key        TEXT NOT NULL UNIQUE,
		vowel      TEXT NOT NULL CHECK (vowel IN ('a','e','i','o','u')),
		rhyme      TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE INDEX words_vowel ON words (vowel);
	CREATE INDEX words_rhyme ON words (rhyme);`,

	// v3: persisted UI state (e.g. canvas layout) as JSON documents per key.
	`CREATE TABLE ui_state (
		key        TEXT PRIMARY KEY,
		value      TEXT NOT NULL,
		updated_at TEXT NOT NULL DEFAULT (datetime('now'))
	);`,
}

// SchemaVersion is the schema version this build expects.
var SchemaVersion = len(migrations)

// DB wraps the SQLite connection and its file location.
type DB struct {
	conn *sql.DB
	path string
}

// Status describes the current database state for display in the UI.
type Status struct {
	Connected     bool   `json:"connected"`
	Path          string `json:"path"`
	SizeBytes     int64  `json:"sizeBytes"`
	SchemaVersion int    `json:"schemaVersion"`
	Error         string `json:"error,omitempty"`
}

// Open creates the parent directory if needed, opens the database and applies migrations.
func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	conn, err := sql.Open("sqlite", "file:"+filepath.ToSlash(path)+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	conn.SetMaxOpenConns(1)
	db := &DB{conn: conn, path: path}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, err
	}
	return db, nil
}

// migrate applies all pending migrations, each in its own transaction.
func (db *DB) migrate() error {
	var current int
	if err := db.conn.QueryRow("PRAGMA user_version").Scan(&current); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}
	for v := current; v < len(migrations); v++ {
		tx, err := db.conn.Begin()
		if err != nil {
			return fmt.Errorf("begin migration v%d: %w", v+1, err)
		}
		if _, err := tx.Exec(migrations[v]); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrate to v%d: %w", v+1, err)
		}
		if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", v+1)); err != nil {
			tx.Rollback()
			return fmt.Errorf("set schema version v%d: %w", v+1, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration v%d: %w", v+1, err)
		}
	}
	return nil
}

// Status reports connection, file and schema information.
func (db *DB) Status() Status {
	s := Status{Path: db.path}
	if err := db.conn.Ping(); err != nil {
		s.Error = err.Error()
		return s
	}
	s.Connected = true
	_ = db.conn.QueryRow("PRAGMA user_version").Scan(&s.SchemaVersion)
	if fi, err := os.Stat(db.path); err == nil {
		s.SizeBytes = fi.Size()
	}
	return s
}

// Close closes the underlying connection.
func (db *DB) Close() error {
	return db.conn.Close()
}
