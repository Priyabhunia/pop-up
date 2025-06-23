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

# Quick Input

A minimalist, global hotkey-activated input box for quick commands and searches.

## Features

- **Global Hotkey**: Press `Ctrl+Space` anywhere to activate the input box
- **Always On Top**: Input box appears above all other windows
- **Frameless Design**: Clean, distraction-free interface
- **Escape to Hide**: Press `Esc` to hide the input box
- **Custom Command Processing**: Process any input through a Go backend

## Screenshots

![Quick Input Screenshot](screenshots/screenshot1.png)


## Installation & Setup

1. **Install Wails** (if not already installed):
   ```
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

2. **Clone the repository**:
   ```
   git clone https://github.com/yourusername/quick-input.git
   cd quick-input
   ```

3. **Build and Run**:
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
2. Press `Ctrl+Space` anywhere to show the input box
3. Type your command and press `Enter`
4. Press `Esc` to hide the input box

## License

[MIT License](LICENSE)

## Acknowledgements

- Built with [Wails](https://wails.io)
- Global hotkey using [robotn/gohook](https://github.com/robotn/gohook)
