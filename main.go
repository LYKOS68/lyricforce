package main

import (
	"embed"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"lyricforce/internal/appdir"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:            "LYKOS REIMARCHIV",
		Width:            1280,
		Height:           820,
		MinWidth:         900,
		MinHeight:        600,
		BackgroundColour: &options.RGBA{R: 14, G: 15, B: 16, A: 255},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind:             []any{app},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// openLogFile sends the standard logger to %APPDATA%\LykosReimarchiv\logs\app.log
// (and stderr). Logging falls back to stderr only if the file cannot be opened.
// Called from App.startup so that build-time runs (binding generation) don't log.
func openLogFile() *os.File {
	path, err := appdir.LogPath()
	if err == nil {
		err = os.MkdirAll(filepath.Dir(path), 0o755)
	}
	if err != nil {
		log.Printf("log dir: %v", err)
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("log file: %v", err)
		return nil
	}
	log.SetOutput(io.MultiWriter(os.Stderr, f))
	log.Printf("start v%s", Version)
	return f
}
