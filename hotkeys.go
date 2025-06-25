package main

import (
	"context"
	"fmt"
	"io"
	"github.com/joho/godotenv"
	"os"
	"net/http"
	"strings"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
	"crypto/tls"
	"path/filepath"
	"encoding/json"
)

// App struct
type App struct {
	ctx context.Context
}

// Configuration struct
type Config struct {
	ObsidianAPIKey string `json:"obsidian_api_key"`
	ObsidianURL    string `json:"obsidian_url"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the app starts
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	
	// Try multiple locations for .env file
	envPaths := []string{
		".env",                           // Current directory
		filepath.Join(os.Getenv("HOME"), ".config", "pop-up", ".env"), // User config dir
		filepath.Join(filepath.Dir(os.Args[0]), ".env"), // App directory
	}
	
	envLoaded := false
	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				fmt.Println("Loaded .env from:", path)
				envLoaded = true
				break
			}
		}
	}
	
	if !envLoaded {
		fmt.Println("Warning: Could not load .env file")
	}
	
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
						/////this code is for testing not working properly
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
	
	// Simple format: First line is title, rest is content
	parts := strings.SplitN(input, "\n", 2)
	title := parts[0]
	
	var content string
	if len(parts) > 1 {
		content = parts[1]
	} else {
		// If no content provided, use the title as content
		content = title
	}
	
	// Format the title for a file name (replace spaces with underscores)
	fileName := strings.ReplaceAll(title, " ", "_") + ".md"
	
	// Create the note in Obsidian
	return a.CreateObsidianNote(fileName, content)
}

// CreateObsidianNote creates a new note in Obsidian via the Local REST API
func (a *App) CreateObsidianNote(title string, content string) string {
	// First try to get the API key from the config file
	config, err := LoadConfig()
	var apiKey string
	
	if err == nil && config.ObsidianAPIKey != "" {
		apiKey = config.ObsidianAPIKey
	} else {
		// Fall back to environment variable
		apiKey = os.Getenv("OBSIDIAN_API_KEY")
	}
	
	// Obsidian Local REST API URL
	url := fmt.Sprintf("https://127.0.0.1:27124/vault/%s", title)
	
	// Create a custom transport that skips certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	// Create a new HTTP client with the custom transport
	client := &http.Client{Transport: tr}
	
	// Create the request
	req, err := http.NewRequest("PUT", url, strings.NewReader(content))
	if err != nil {
		return fmt.Sprintf("Error creating request: %s", err.Error())
	}
	
	// Add authorization header with API key
	if apiKey != "" {
		req.Header.Add("Authorization", "Bearer " + apiKey)
	} else {
		return "Error: No API key found. Please add your Obsidian API key in settings."
	}
	
	req.Header.Add("Content-Type", "text/plain")
	
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error making request: %s", err.Error())
	}
	defer resp.Body.Close()
	
	// Check the response
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return fmt.Sprintf("Note '%s' created successfully", title)
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		responseText := string(bodyBytes)
		return fmt.Sprintf("Error: %s - %s", resp.Status, responseText)
	}
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	configDir := filepath.Join(home, ".config", "pop-up")
	os.MkdirAll(configDir, 0755) // Create directory if it doesn't exist
	
	configPath := filepath.Join(configDir, "config.json")
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := &Config{
			ObsidianAPIKey: os.Getenv("OBSIDIAN_API_KEY"), // Try to get from env first
			ObsidianURL:    "https://127.0.0.1:27124/vault/",
		}
		
		// Save default config
		configBytes, _ := json.MarshalIndent(defaultConfig, "", "  ")
		os.WriteFile(configPath, configBytes, 0644)
		
		return defaultConfig, nil
	}
	
	// Read existing config
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	var config Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}
	
	// Override with environment variables if available
	if envKey := os.Getenv("OBSIDIAN_API_KEY"); envKey != "" {
		config.ObsidianAPIKey = envKey
	}
	
	return &config, nil
}

// SaveAPIKey saves the API key to the config file
func (a *App) SaveAPIKey(apiKey string) string {
	config, err := LoadConfig()
	if err != nil {
		return "Error loading config"
	}
	
	config.ObsidianAPIKey = apiKey
	
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "pop-up", "config.json")
	
	configBytes, _ := json.MarshalIndent(config, "", "  ")
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		return "Error saving config"
	}
	
	return "API Key saved successfully"
}




