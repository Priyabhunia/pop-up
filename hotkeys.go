package main

import (
	"context"
	"fmt"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}


// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}



// OnStartup is called when the app starts
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	a.startHotkeyListener()
}

// startHotkeyListener starts a goroutine to listen for Ctrl+Space globally
func (a *App) startHotkeyListener() {
	go func() {
		events := hook.Start()
		defer hook.End()
		for ev := range events {
			if ev.Kind == hook.KeyDown {
				// Ctrl+Space: ev.Keycode == 57 (space), ev.Mask == 2 (Ctrl)
				if ev.Keycode == 57 && ev.Mask == 2 {
					if a.ctx != nil {
						// Show the window
						runtime.WindowShow(a.ctx)
						
						// Force focus by temporarily setting always on top
						runtime.WindowSetAlwaysOnTop(a.ctx, true)
						
						// Short delay to ensure window is visible and focused
						time.Sleep(100 * time.Millisecond)
						
						// Return to normal state
						runtime.WindowSetAlwaysOnTop(a.ctx, false)
					}
				}
			}
		}
	}()
}

// HideWindow hides the Wails window (to be called from frontend)
func (a *App) HideWindow() {
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
	}
}

// ProcessInput handles the input from the frontend
func (a *App) ProcessInput(input string) string {
	if input == "" {
		return "Please enter something!"
	}
	
	return fmt.Sprintf("You entered: %s", input)
}




