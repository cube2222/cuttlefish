package main

import (
	"context"
	"database/sql"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	_ "modernc.org/sqlite"

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

	db, err := sql.Open("sqlite", "file:gptui.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create tables
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
