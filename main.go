package main

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	_ "modernc.org/sqlite"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var db *sql.DB
	{
		//TODO: get better exePath
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}
	
		pathToApp := filepath.Join(filepath.Dir(exePath), "data.db")
		//TODO: cross platform file:
		db, err = sql.Open("sqlite", "file:///"+pathToApp)
		if err != nil {
			log.Fatal(err)
		}
	
		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatal(pingErr)
		}
	}
	
	// habitRepository := repository.CreateHabitRepository(db)

	err := wails.Run(&options.App{
		Title:  "Taufiqraw's Habit Tracker",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			// habitRepository,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
