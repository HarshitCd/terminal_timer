package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type Model struct {
	width  int
	height int

	now time.Time
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tickMsg:
		m.now = time.Now()
		return m, tickCmd()
	}

	return m, nil
}

func (m Model) View() string {
	text := m.now.Format("15:04")
	quitText := "ctrl+c or q to quit"
	text = figure.NewFigure(text, "banner3", false).String()

	// backgroundStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#45A1E8"))
	centerStyle := lipgloss.NewStyle().Height(m.height).Width(m.width).Align(lipgloss.Center, lipgloss.Center)
	headingStyle := lipgloss.NewStyle().Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6D6D6D"))

	return centerStyle.Render(fmt.Sprintf("%s\n\n%s", headingStyle.Render(text), hintStyle.Render(quitText)))
}

func main() {
	p := tea.NewProgram(Model{now: time.Now()}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("failed to run")
		os.Exit(1)
	}
}
