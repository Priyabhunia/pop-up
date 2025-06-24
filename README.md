# README

## About

This is the official Wails Vanilla template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

# Pop-up

A minimalist, global hotkey-activated input box for quick commands and searches with Obsidian integration.

## Features

- **Global Hotkey**: Press `Ctrl+Space` anywhere to activate the input box
- **Always On Top**: Input box appears above all other windows
- **Escape to Hide**: Press `Esc` to hide the input box
- **Quick Obsidian Notes**: Type anything to create notes directly in Obsidian

## Screenshots

![Pop-up Screenshot](screenshots/screenshot1.png)

## Installation & Setup

1. **Install Wails** (if not already installed):
   ```
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

2. **Clone the repository**:
   ```
   git clone https://github.com/Priyabhunia/pop-up.git
   cd pop-up
   ```

3. **Setup Obsidian Integration**:
   - Install the "Local REST API" plugin in Obsidian
   - Enable the plugin and copy your API key
   - Create a `.env` file in the project root with:
     ```
     OBSIDIAN_API_KEY=your_api_key_here
     ```

4. **Build and Run**:
   ```
   wails build
   ```
   Then run the executable in the `build/bin` directory.

## Development

- **Development Mode**:
  ```
  wails dev
  ```

- **Frontend Customization**: Edit files directly in `frontend/dist/`

## Usage

1. Launch the application
2. Make sure Obsidian is running with the Local REST API plugin enabled
3. Press `Ctrl+Space` anywhere to show the input box
4. Type your note content and press `Enter` to save it to Obsidian
5. Press `Esc` to hide the input box

## Obsidian Integration

This app integrates with Obsidian through its Local REST API:

- **Requirements**: 
  - Obsidian must be running
  - The "Local REST API" plugin must be installed and enabled
  - Your API key must be set in the `.env` file

- **How it works**:
  - Text entered in the pop-up is sent to Obsidian as a new note
  - The first line becomes the note title (and filename)
  - The content is formatted as Markdown

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests for new features or improvements.

## License

[MIT License](LICENSE)

## Acknowledgements

- Built with [Wails](https://wails.io)
- Global hotkey using [robotn/gohook](https://github.com/robotn/gohook)
- Integration with [Obsidian](https://obsidian.md) via Local REST API
