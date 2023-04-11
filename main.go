package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	_ "modernc.org/sqlite"

	"cuttlefish/database"
	"cuttlefish/database/migrate"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	ctx := context.Background()

	if err := os.MkdirAll("~/.cuttlefish", 0755); err != nil {
		panic(err)
	}

	// Open somewhere in `~/` instead of current directory.
	db, err := sql.Open("sqlite", "file:~/.cuttlefish/data.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	// db, err := sql.Open("sqlite", "file:gptui.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	// db, err := sql.Open("sqlite3", "file:gptui.db?cache=shared&mode=rwc&_foreign_keys=1&_journal_mode=WAL&_busy_timeout=5000&_loc=auto")
	// db, err := sql.Open("sqlite3", "file:gptui.db?mode=rwc&_foreign_keys=1&_busy_timeout=5000&_loc=auto")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrate.Migrate(db); err != nil {
		log.Fatal("could not run migrations:", err)
	}

	queries := database.New(db)

	// Create an instance of the app structure
	app := NewApp(ctx, queries)

	// Create application with options
	err = wails.Run(&options.App{
		Title:            "Cuttlefish",
		Width:            1024,
		Height:           768,
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
