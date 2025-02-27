# Whisk

A CLI tool that creates skeleton Go projects with proper structure and Charm UI components.

## Features

- Interactive TUI using [Charm libraries](https://github.com/charmbracelet)
- Customizable project types:
  - CLI applications with BubbleTea
  - HTTP/API servers
  - Reusable libraries
- Sets up proper Go module structure
- Creates essential files and directories
- Automatically configures Git repository paths

## Development

### With Flox

[Flox](https://flox.dev) provides a convenient development environment for this project.

```bash
# Activate the development environment
flox activate

# Install dependencies
go mod tidy

# Run the project
go run .
```

This will set up all necessary dependencies in an isolated environment without affecting your system.

## Installation

### Prerequisites

- Go 1.18 or higher

### Building from source

```bash
# Clone the repository
git clone https://github.com/alloydwhitlock/whisk.git
cd whisk

# Install dependencies
go mod tidy

# Build the binary
go build -o whisk .

# Optional: Move to your PATH
mv whisk /usr/local/bin/
```

## Usage

```bash
whisk
```

The interactive UI will guide you through:

1. Entering a project name
2. Specifying Git repository path (e.g., github.com/username/project)
3. Selecting project type (CLI, server, or library)
4. Confirming your choices
5. Creating the project structure

## Project Structures

### CLI Projects

```
my-cli-project/
├── cmd/
├── internal/
├── main.go       # BubbleTea application
├── go.mod
├── README.md
└── .gitignore
```

### Server Projects

```
my-server-project/
├── api/
├── internal/
│   ├── handlers/
│   │   └── health.go
│   └── middleware/
├── main.go       # HTTP server
├── go.mod
├── README.md
└── .gitignore
```

### Library Projects

```
my-library-project/
├── main.go       # Placeholder
├── go.mod
├── README.md
└── .gitignore
```

## Credits

Built with [Charm](https://charm.sh) libraries:
- BubbleTea: Terminal UI framework
- Bubbles: TUI components
- Lipgloss: Style definitions

## License

MIT
