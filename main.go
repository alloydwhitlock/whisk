package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Italic(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)
)

type state int

const (
	stateProjectName state = iota
	stateGitRepo
	stateProjectType
	stateConfirm
	stateCreating
	stateDone
)

// Custom message type for completion
type doneMsg struct{}

type model struct {
	state       state
	projectName textinput.Model
	gitRepo     textinput.Model
	typeList    list.Model
	err         error
}

type projectType struct {
	name        string
	description string
}

func (p projectType) Title() string       { return p.name }
func (p projectType) Description() string { return p.description }
func (p projectType) FilterValue() string { return p.name }

// Initialize the model
func initialModel() model {
	// Project name input
	pn := textinput.New()
	pn.Placeholder = "my-awesome-project"
	pn.Focus()
	pn.CharLimit = 50
	pn.Width = 30

	// Git repo input
	gr := textinput.New()
	gr.Placeholder = "github.com/username/my-awesome-project"
	gr.CharLimit = 100
	gr.Width = 40

	// Project type list
	items := []list.Item{
		projectType{name: "cli", description: "Command-line application"},
		projectType{name: "server", description: "HTTP/API server"},
		projectType{name: "library", description: "Reusable package/library"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#DDDDDD"))

	l := list.New(items, delegate, 30, 10)
	l.Title = "Select project type"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	return model{
		state:       stateProjectName,
		projectName: pn,
		gitRepo:     gr,
		typeList:    l,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			switch m.state {
			case stateProjectName:
				if m.projectName.Value() != "" {
					m.state = stateGitRepo
					m.gitRepo.Focus()
					return m, textinput.Blink
				}

			case stateGitRepo:
				if m.gitRepo.Value() != "" {
					m.state = stateProjectType
					return m, nil
				}

			case stateProjectType:
				m.state = stateConfirm
				return m, nil

			case stateConfirm:
				m.state = stateCreating
				return m, tea.Sequence(
					func() tea.Msg {
						err := createProject(m)
						if err != nil {
							return err
						}
						return doneMsg{}
					},
				)

			case stateDone:
				return m, tea.Quit
			}
		}
	case error:
		m.err = msg
		m.state = stateDone
		return m, nil

	case doneMsg:
		m.state = stateDone
		return m, nil
	}

	// Handle input updates based on current state
	switch m.state {
	case stateProjectName:
		m.projectName, cmd = m.projectName.Update(msg)
		return m, cmd

	case stateGitRepo:
		m.gitRepo, cmd = m.gitRepo.Update(msg)
		return m, cmd

	case stateProjectType:
		m.typeList, cmd = m.typeList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	var s string

	switch m.state {
	case stateProjectName:
		s = titleStyle.Render("Whisk Go Project Creator") + "\n\n"
		s += "Enter project name:\n\n"
		s += m.projectName.View() + "\n\n"
		s += infoStyle.Render("Press Enter to continue, Esc to quit")

	case stateGitRepo:
		s = titleStyle.Render("Whisk Go Project Creator") + "\n\n"
		s += fmt.Sprintf("Project name: %s\n\n", m.projectName.Value())
		s += "Enter Git repository path:\n\n"
		s += m.gitRepo.View() + "\n\n"
		s += infoStyle.Render("Press Enter to continue, Esc to quit")

	case stateProjectType:
		s = titleStyle.Render("Whisk Go Project Creator") + "\n\n"
		s += fmt.Sprintf("Project name: %s\n", m.projectName.Value())
		s += fmt.Sprintf("Git repository: %s\n\n", m.gitRepo.Value())
		s += "Select project type:\n\n"
		s += m.typeList.View()

	case stateConfirm:
		s = titleStyle.Render("Whisk Go Project Creator") + "\n\n"
		s += "Please confirm your choices:\n\n"
		s += fmt.Sprintf("Project name: %s\n", m.projectName.Value())
		s += fmt.Sprintf("Git repository: %s\n", m.gitRepo.Value())
		
		item, ok := m.typeList.SelectedItem().(projectType)
		if ok {
			s += fmt.Sprintf("Project type: %s (%s)\n\n", item.name, item.description)
		}
		
		s += infoStyle.Render("Press Enter to create project, Esc to quit")

	case stateCreating:
		s = titleStyle.Render("Whisk Go Project Creator") + "\n\n"
		s += "Creating project...\n"

	case stateDone:
		s = titleStyle.Render("Whisk Go Project Creator") + "\n\n"
		if m.err != nil {
			s += fmt.Sprintf("Error creating project: %v\n\n", m.err)
		} else {
			s += successStyle.Render("✓ Project created successfully!") + "\n\n"
			s += fmt.Sprintf("Project created at: %s\n\n", m.projectName.Value())
			s += "To get started:\n\n"
			s += fmt.Sprintf("  cd %s\n", m.projectName.Value())
			s += "  go mod tidy\n"
			s += "  go run .\n\n"
		}
		s += infoStyle.Render("Press Enter to exit")
	}

	return s
}

func createProject(m model) error {
	projectName := m.projectName.Value()
	gitRepo := m.gitRepo.Value()
	
	item, ok := m.typeList.SelectedItem().(projectType)
	if !ok {
		return fmt.Errorf("invalid project type selection")
	}
	projectType := item.name

	// Create project directory
	err := os.MkdirAll(projectName, 0755)
	if err != nil {
		return err
	}

	// Create go.mod file
	goModContent := fmt.Sprintf("module %s\n\ngo 1.21\n", gitRepo)
	err = os.WriteFile(filepath.Join(projectName, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		return err
	}

	// Create main.go file based on project type
	var mainContent string
	
	switch projectType {
	case "cli":
		mainContent = generateCliMainFile(gitRepo)
	case "server":
		mainContent = generateServerMainFile(gitRepo)
	default:
		mainContent = generateLibraryMainFile(gitRepo)
	}

	err = os.WriteFile(filepath.Join(projectName, "main.go"), []byte(mainContent), 0644)
	if err != nil {
		return err
	}

	// Create README.md
	readmeContent := fmt.Sprintf("# %s\n\nA Go %s project.\n", projectName, projectType)
	err = os.WriteFile(filepath.Join(projectName, "README.md"), []byte(readmeContent), 0644)
	if err != nil {
		return err
	}

	// Create .gitignore
	gitignoreContent := "# Binaries for programs and plugins\n*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n\n# Test binary, built with `go test -c`\n*.test\n\n# Output of the go coverage tool, specifically when used with LiteIDE\n*.out\n\n# Dependency directories (remove the comment below to include it)\n# vendor/\n"
	err = os.WriteFile(filepath.Join(projectName, ".gitignore"), []byte(gitignoreContent), 0644)
	if err != nil {
		return err
	}

	// Create additional structure based on project type
	switch projectType {
	case "cli":
		// Create cmd directory for CLI commands
		os.MkdirAll(filepath.Join(projectName, "cmd"), 0755)
		os.MkdirAll(filepath.Join(projectName, "internal"), 0755)
		
	case "server":
		// Create directories for server project
		os.MkdirAll(filepath.Join(projectName, "api"), 0755)
		os.MkdirAll(filepath.Join(projectName, "internal/handlers"), 0755)
		os.MkdirAll(filepath.Join(projectName, "internal/middleware"), 0755)
		
		// Create simple handler
		handlerContent := fmt.Sprintf(`package handlers

import (
	"net/http"
)

// HealthHandler returns a simple health check response
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\\"status\\": \\"ok\\"}"))
}
`)
		os.WriteFile(filepath.Join(projectName, "internal/handlers/health.go"), []byte(handlerContent), 0644)
	}

	return nil
}

func generateCliMainFile(gitRepo string) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
		Padding(1, 2)
	
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Bold(true)
)

type keyMap struct {
	Help key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	keys keyMap
	help help.Model
}

func initialModel() model {
	return model{
		keys: keys,
		help: help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	
	return m, nil
}

func (m model) View() string {
	s := titleStyle.Render("CLI Application") + "\n\n"
	s += "Welcome to your new CLI application!\n\n"
	s += m.help.View(m.keys)
	
	return appStyle.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %%v\n", err)
		os.Exit(1)
	}
}
`)
}

func generateServerMainFile(gitRepo string) string {
	return fmt.Sprintf(`package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"%s/internal/handlers"
)

func main() {
	// Setup logger
	logger := log.New(os.Stdout, "", log.LstdFlags)
	
	// Create router (using standard library for simplicity)
	mux := http.NewServeMux()
	
	// Register routes
	mux.HandleFunc("/health", handlers.HealthHandler)
	
	// Create server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		logger.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %%v", err)
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	// Shutdown gracefully
	logger.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %%v", err)
	}
	
	logger.Println("Server exited properly")
}
`, gitRepo)
}

func generateLibraryMainFile(gitRepo string) string {
	return `// Package main is a placeholder for the library
// In a real library, you would have multiple packages organized in subdirectories
package main

import "fmt"

func main() {
	fmt.Println("This is a placeholder for your library.")
	fmt.Println("In a real library, this file would be replaced with proper packages.")
}
`
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
