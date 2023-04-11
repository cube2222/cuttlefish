package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
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

	userHomedir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	cuttlefishHomedir := filepath.Join(userHomedir, ".cuttlefish")

	if err := os.MkdirAll(cuttlefishHomedir, 0755); err != nil {
		panic(err)
	}

	dbFile := filepath.Join(cuttlefishHomedir, "data.db")

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", dbFile))
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
