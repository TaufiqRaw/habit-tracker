package main

import (
	"embed"
	"habit-tracker/internal/bootstrap"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/repository"
	"habit-tracker/internal/service"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	_ "modernc.org/sqlite"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db := bootstrap.InitDB()
	habitService := service.CreateHabitService(repository.CreateHabitRepository(db))
	trackerService := service.CreateTrackerService(repository.CreateTrackerRepository(db))

	err := wails.Run(&options.App{
		Title:  "Taufiqraw's Habit Tracker",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			habitService,
			trackerService,
		},
		EnumBind: []interface{}{
			domain.AllRestDayMode,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
