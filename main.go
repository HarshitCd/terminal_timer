package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type Model struct {
	width  int
	height int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	now := time.Now()
	text := now.Format("15:04")
	quitText := "press any key to quit"
	text = figure.NewFigure(text, "banner3", false).String()

	// backgroundStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#45A1E8"))
	centerStyle := lipgloss.NewStyle().Height(m.height).Width(m.width).Align(lipgloss.Center, lipgloss.Center)
	headingStyle := lipgloss.NewStyle().Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6D6D6D"))

	return centerStyle.Render(fmt.Sprintf("%s\n\n%s", headingStyle.Render(text), hintStyle.Render(quitText)))
}

func main() {
	p := tea.NewProgram(Model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("failed to run")
		os.Exit(1)
	}
}
