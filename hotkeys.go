package main

import (
	"context"
	"fmt"
	"io"
	"github.com/joho/godotenv"
	"os"
	"net/http"
	"strings"
	"net/url"
	"os/exec"
	"runtime"
	hook "github.com/robotn/gohook"
  //"github.com/wailsapp/wails/v2/pkg/runtime"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
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
						wailsRuntime.WindowShow(a.ctx)
						
						// Force focus by temporarily setting always on top
						wailsRuntime.WindowSetAlwaysOnTop(a.ctx, true)
						/////this code is for testing not working properly
						// Short delay to ensure window is visible and focused
						time.Sleep(100 * time.Millisecond)
						
						// Return to normal state
						wailsRuntime.WindowSetAlwaysOnTop(a.ctx, false)
					}
				}
			}
		}
	}()
}

// HideWindow hides the Wails window (to be called from frontend)
func (a *App) HideWindow() {
	if a.ctx != nil {
		wailsRuntime.WindowHide(a.ctx)
	}
}

// ProcessInput handles the input from the frontend
func (a *App) ProcessInput(input string) string {
	if input == "" {
		return "Please enter something!"
	}
	
	// Check if this is a search query with prefix at the END (e.g., "search term !g")
	if strings.Contains(input, " !") {
		parts := strings.Split(input, " !")
		if len(parts) == 2 {
			query := strings.TrimSpace(parts[0])
			prefix := "!" + strings.TrimSpace(parts[1])
			
			return a.handleSearchQuery(prefix, query)
		}
	}
	
	// Check if this is a search query with prefix at the BEGINNING (e.g., "!g search term")
	if strings.HasPrefix(input, "!") {
		parts := strings.SplitN(input, " ", 2)
		if len(parts) == 2 && parts[1] != "" {
			prefix := parts[0]
			query := parts[1]
			
			return a.handleSearchQuery(prefix, query)
		}
	}
	
	// Not a search query, handle as normal note creation
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

// handleSearchQuery processes search queries with prefixes
func (a *App) handleSearchQuery(prefix string, query string) string {
	if query == "" {
		return "Please enter a search term"
	}
	
	encodedQuery := url.QueryEscape(query)
	var searchURL string
	var searchEngine string
	
	switch prefix {
	case "!g":
		searchURL = "https://www.google.com/search?q=" + encodedQuery
		searchEngine = "Google"
	case "!yt":
		searchURL = "https://www.youtube.com/results?search_query=" + encodedQuery
		searchEngine = "YouTube"
	case "!gh":
		searchURL = "https://github.com/search?q=" + encodedQuery
		searchEngine = "GitHub"
	case "!c":
		searchURL = "https://chat.openai.com/?q=" + encodedQuery
		searchEngine = "ChatGPT"
	case "!grok":
		searchURL = "https://grok.com/?q=" + encodedQuery
		searchEngine = "Grok"
	default:
		return "Unknown search prefix: " + prefix
	}
	
	// Open the URL in the default browser
	if err := a.openBrowser(searchURL); err != nil {
		return "Error opening browser: " + err.Error()
	}
	
	return fmt.Sprintf("Opening %s search for: %s", searchEngine, query)
}

// openBrowser opens the specified URL in the default browser
func (a *App) openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
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




