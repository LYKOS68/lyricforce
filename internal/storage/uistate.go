package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// maxUIStateBytes bounds a single UI state document.
const maxUIStateBytes = 64 << 10

var ErrInvalidUIState = errors.New("invalid ui state")

// GetUIState returns the JSON document stored under key, or "" if none exists.
func (db *DB) GetUIState(key string) (string, error) {
	var v string
	err := db.conn.QueryRow(`SELECT value FROM ui_state WHERE key = ?`, key).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get ui state: %w", err)
	}
	return v, nil
}

// SetUIState stores a JSON document under key, replacing any previous value.
func (db *DB) SetUIState(key, value string) error {
	if key == "" || len(value) > maxUIStateBytes || !json.Valid([]byte(value)) {
		return ErrInvalidUIState
	}
	_, err := db.conn.Exec(
		`INSERT INTO ui_state (key, value) VALUES (?, ?)
		 ON CONFLICT (key) DO UPDATE SET value = excluded.value, updated_at = datetime('now')`,
		key, value)
	if err != nil {
		return fmt.Errorf("set ui state: %w", err)
	}
	return nil
}
