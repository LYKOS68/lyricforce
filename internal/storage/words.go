package storage

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// Error values returned to the UI. Their messages are stable and mapped to
// user-facing text in the frontend.
var (
	ErrEmptyWord     = errors.New("word is empty")
	ErrInvalidVowel  = errors.New("invalid vowel")
	ErrDuplicateWord = errors.New("word already exists")
	ErrWordNotFound  = errors.New("word not found")
)

// Vowels are the vowel classes of the archive (MVP: simple vowels only).
var Vowels = []string{"a", "e", "i", "o", "u"}

// Word is one entry of the rhyme archive.
type Word struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Vowel     string `json:"vowel"`
	Rhyme     string `json:"rhyme"`
	CreatedAt string `json:"createdAt"`
}

// WordFilter narrows ListWords. Empty fields match everything.
type WordFilter struct {
	Query string `json:"query"` // substring of word or rhyme sound
	Vowel string `json:"vowel"`
}

// AddWord validates and stores a new word.
func (db *DB) AddWord(text, vowel, rhyme string) (Word, error) {
	text = strings.TrimSpace(text)
	vowel = strings.ToLower(strings.TrimSpace(vowel))
	rhyme = strings.ToLower(strings.TrimSpace(rhyme))
	if text == "" {
		return Word{}, ErrEmptyWord
	}
	if !slices.Contains(Vowels, vowel) {
		return Word{}, ErrInvalidVowel
	}

	var w Word
	err := db.conn.QueryRow(
		`INSERT INTO words (text, key, vowel, rhyme) VALUES (?, ?, ?, ?)
		 RETURNING id, text, vowel, rhyme, created_at`,
		text, strings.ToLower(text), vowel, rhyme,
	).Scan(&w.ID, &w.Text, &w.Vowel, &w.Rhyme, &w.CreatedAt)
	if err != nil {
		var se *sqlite.Error
		if errors.As(err, &se) && se.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return Word{}, ErrDuplicateWord
		}
		return Word{}, fmt.Errorf("insert word: %w", err)
	}
	return w, nil
}

// ListWords returns matching words, grouped by rhyme sound (words without one last).
func (db *DB) ListWords(f WordFilter) ([]Word, error) {
	where := []string{"1=1"}
	var args []any
	if v := strings.ToLower(strings.TrimSpace(f.Vowel)); v != "" {
		where = append(where, "vowel = ?")
		args = append(args, v)
	}
	if q := strings.ToLower(strings.TrimSpace(f.Query)); q != "" {
		pattern := "%" + escapeLike(q) + "%"
		where = append(where, `(key LIKE ? ESCAPE '\' OR rhyme LIKE ? ESCAPE '\')`)
		args = append(args, pattern, pattern)
	}

	rows, err := db.conn.Query(
		`SELECT id, text, vowel, rhyme, created_at FROM words
		 WHERE `+strings.Join(where, " AND ")+`
		 ORDER BY rhyme = '', rhyme, key`, args...)
	if err != nil {
		return nil, fmt.Errorf("list words: %w", err)
	}
	defer rows.Close()

	words := []Word{}
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Text, &w.Vowel, &w.Rhyme, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan word: %w", err)
		}
		words = append(words, w)
	}
	return words, rows.Err()
}

// CountWordsByVowel returns the number of words per vowel class (all classes present).
func (db *DB) CountWordsByVowel() (map[string]int, error) {
	counts := make(map[string]int, len(Vowels))
	for _, v := range Vowels {
		counts[v] = 0
	}
	rows, err := db.conn.Query(`SELECT vowel, COUNT(*) FROM words GROUP BY vowel`)
	if err != nil {
		return nil, fmt.Errorf("count words: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var v string
		var n int
		if err := rows.Scan(&v, &n); err != nil {
			return nil, fmt.Errorf("scan count: %w", err)
		}
		counts[v] = n
	}
	return counts, rows.Err()
}

// DeleteWord removes a word by id.
func (db *DB) DeleteWord(id int64) error {
	res, err := db.conn.Exec(`DELETE FROM words WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete word: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrWordNotFound
	}
	return nil
}

func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
