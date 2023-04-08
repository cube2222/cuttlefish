package main

import (
	"context"
	"database/sql"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	_ "github.com/mattn/go-sqlite3"

	"gptui/database"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

//go:embed database/schema.sql
var ddl string

func main() {
	ctx := context.Background()

	// db, err := sql.Open("sqlite", "file:gptui.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	// db, err := sql.Open("sqlite3", "file:gptui.db?cache=shared&mode=rwc&_foreign_keys=1&_journal_mode=WAL&_busy_timeout=5000&_loc=auto")
	db, err := sql.Open("sqlite3", "file:gptui.db?mode=rwc&_foreign_keys=1&_busy_timeout=5000&_loc=auto")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: Add proper migrations
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	// Create an instance of the app structure
	app := NewApp(ctx, queries)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "wails-events",
		Width:  1024,
		Height: 768,
		// Frameless:        true,
		// AlwaysOnTop:      true,
		// WindowStartState: options.Minimised,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
