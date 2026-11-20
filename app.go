package main

import (
	"context"
	"errors"
	"log"
	"os"

	"lyricforce/internal/appdir"
	"lyricforce/internal/modules"
	"lyricforce/internal/storage"
)

// Version of the application shown in the UI.
const Version = "0.1.0"

var errNoDatabase = errors.New("database not available")

// App is the bridge between the Go application logic and the frontend.
// Every exported method is callable from JavaScript via window.go.main.App.
type App struct {
	ctx     context.Context
	db      *storage.DB
	dbErr   error
	logFile *os.File
}

// AppInfo is general information about the running build.
type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	LogPath string `json:"logPath"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logFile = openLogFile()
	path, err := appdir.DatabasePath()
	if err != nil {
		a.dbErr = err
		log.Printf("database path: %v", err)
		return
	}
	a.db, a.dbErr = storage.Open(path)
	if a.dbErr != nil {
		log.Printf("database open: %v", a.dbErr)
		return
	}
	log.Printf("database ready: %s", path)
}

func (a *App) shutdown(_ context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("database close: %v", err)
		}
	}
	if a.logFile != nil {
		log.Printf("shutdown")
		a.logFile.Close()
	}
}

func (a *App) store() (*storage.DB, error) {
	if a.db == nil {
		return nil, errNoDatabase
	}
	return a.db, nil
}

// GetAppInfo returns name, version and log location.
func (a *App) GetAppInfo() AppInfo {
	logPath, _ := appdir.LogPath()
	return AppInfo{Name: "LYKOS REIMARCHIV", Version: Version, LogPath: logPath}
}

// GetModules returns the navigation targets.
func (a *App) GetModules() []modules.Module {
	return modules.All()
}

// GetDatabaseStatus returns the state of the local SQLite database.
func (a *App) GetDatabaseStatus() storage.Status {
	if a.db == nil {
		s := storage.Status{}
		if a.dbErr != nil {
			s.Error = a.dbErr.Error()
		}
		return s
	}
	return a.db.Status()
}

// ListWords returns archive words matching the filter.
func (a *App) ListWords(f storage.WordFilter) ([]storage.Word, error) {
	db, err := a.store()
	if err != nil {
		return nil, err
	}
	return db.ListWords(f)
}

// CountWordsByVowel returns the number of words per vowel class.
func (a *App) CountWordsByVowel() (map[string]int, error) {
	db, err := a.store()
	if err != nil {
		return nil, err
	}
	return db.CountWordsByVowel()
}

// AddWord stores a new word in the archive.
func (a *App) AddWord(text, vowel, rhyme string) (storage.Word, error) {
	db, err := a.store()
	if err != nil {
		return storage.Word{}, err
	}
	w, err := db.AddWord(text, vowel, rhyme)
	if err == nil {
		log.Printf("word added: id=%d %q", w.ID, w.Text)
	}
	return w, err
}

// GetUIState returns a persisted UI state document ("" if none).
func (a *App) GetUIState(key string) (string, error) {
	db, err := a.store()
	if err != nil {
		return "", err
	}
	return db.GetUIState(key)
}

// SetUIState persists a UI state document (JSON).
func (a *App) SetUIState(key, value string) error {
	db, err := a.store()
	if err != nil {
		return err
	}
	return db.SetUIState(key, value)
}

// DeleteWord removes a word from the archive.
func (a *App) DeleteWord(id int64) error {
	db, err := a.store()
	if err != nil {
		return err
	}
	if err := db.DeleteWord(id); err != nil {
		return err
	}
	log.Printf("word deleted: id=%d", id)
	return nil
}
