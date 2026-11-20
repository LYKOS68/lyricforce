package storage

import (
	"errors"
	"path/filepath"
	"testing"
)

func openTemp(t *testing.T) *DB {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "sub", "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestOpenMigratesAndReopens(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if s := db.Status(); !s.Connected || s.SchemaVersion != SchemaVersion {
		t.Fatalf("unexpected status: %+v", s)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// Reopening an existing database must not fail or re-run migrations.
	db, err = Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer db.Close()
	if got := db.Status().SchemaVersion; got != SchemaVersion {
		t.Fatalf("schema version after reopen = %d, want %d", got, SchemaVersion)
	}
}

func TestUIState(t *testing.T) {
	db := openTemp(t)

	if v, err := db.GetUIState("canvas"); err != nil || v != "" {
		t.Fatalf("missing key: %q, %v", v, err)
	}
	if err := db.SetUIState("canvas", `{"zoom":1}`); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := db.SetUIState("canvas", `{"zoom":2}`); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	if v, _ := db.GetUIState("canvas"); v != `{"zoom":2}` {
		t.Errorf("get = %q", v)
	}
	if err := db.SetUIState("canvas", `{broken`); !errors.Is(err, ErrInvalidUIState) {
		t.Errorf("invalid json: got %v", err)
	}
	if err := db.SetUIState("", `{}`); !errors.Is(err, ErrInvalidUIState) {
		t.Errorf("empty key: got %v", err)
	}
}

func TestAddWordValidation(t *testing.T) {
	db := openTemp(t)

	if _, err := db.AddWord("  ", "a", ""); !errors.Is(err, ErrEmptyWord) {
		t.Errorf("empty word: got %v", err)
	}
	if _, err := db.AddWord("Nacht", "x", ""); !errors.Is(err, ErrInvalidVowel) {
		t.Errorf("invalid vowel: got %v", err)
	}
	w, err := db.AddWord(" Nacht ", "A", " ACHT ")
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if w.Text != "Nacht" || w.Vowel != "a" || w.Rhyme != "acht" || w.ID == 0 {
		t.Errorf("normalised word = %+v", w)
	}
	// Duplicate detection is case-insensitive, including umlauts.
	if _, err := db.AddWord("nacht", "a", ""); !errors.Is(err, ErrDuplicateWord) {
		t.Errorf("duplicate: got %v", err)
	}
	if _, err := db.AddWord("Äste", "e", ""); err != nil {
		t.Fatalf("add umlaut: %v", err)
	}
	if _, err := db.AddWord("äste", "e", ""); !errors.Is(err, ErrDuplicateWord) {
		t.Errorf("umlaut duplicate: got %v", err)
	}
}

func TestListCountDelete(t *testing.T) {
	db := openTemp(t)
	for _, in := range [][3]string{
		{"Nacht", "a", "acht"}, {"Macht", "a", "acht"}, {"Licht", "i", "icht"}, {"Haus", "a", ""}, {"100%", "u", ""},
	} {
		if _, err := db.AddWord(in[0], in[1], in[2]); err != nil {
			t.Fatalf("add %s: %v", in[0], err)
		}
	}

	all, err := db.ListWords(WordFilter{})
	if err != nil || len(all) != 5 {
		t.Fatalf("list all: %d, %v", len(all), err)
	}
	// Grouped by rhyme sound, words without rhyme sound last.
	if all[0].Text != "Macht" || all[1].Text != "Nacht" || all[2].Text != "Licht" {
		t.Errorf("order = %v", all)
	}

	byVowel, _ := db.ListWords(WordFilter{Vowel: "a"})
	if len(byVowel) != 3 {
		t.Errorf("vowel a: got %d", len(byVowel))
	}
	byRhyme, _ := db.ListWords(WordFilter{Query: "acht"})
	if len(byRhyme) != 2 {
		t.Errorf("query acht: got %d", len(byRhyme))
	}
	// LIKE wildcards in the query are matched literally.
	literal, _ := db.ListWords(WordFilter{Query: "%"})
	if len(literal) != 1 || literal[0].Text != "100%" {
		t.Errorf("literal %%: got %v", literal)
	}

	counts, err := db.CountWordsByVowel()
	if err != nil || counts["a"] != 3 || counts["i"] != 1 || counts["e"] != 0 {
		t.Errorf("counts = %v, %v", counts, err)
	}

	if err := db.DeleteWord(all[0].ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := db.DeleteWord(all[0].ID); !errors.Is(err, ErrWordNotFound) {
		t.Errorf("delete missing: got %v", err)
	}
}
