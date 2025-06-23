package main

import (
	//"context"
	"embed"
	//"log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	// "golang.design/x/hotkey"
	// gohook "github.com/robotn/gohook"
	// "github.com/go-vgo/robotgo"
	// "github.com/wailsapp/wails/v2/pkg/runtime"
)


//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		
		Title:         "Quick Input",
		Width:         650,
		Height:        80,
		MinWidth:      650,
		MinHeight:     80,
		MaxWidth:      650,
		MaxHeight:     80,
		Fullscreen:    false,
		StartHidden:   true,
		Frameless:     true,  // This removes the title bar completely
		AlwaysOnTop:   true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 30, G: 30, B: 30, A: 255},
		OnStartup:        app.OnStartup,
		Bind:             []interface{}{app},

		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent: true,

		
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}






