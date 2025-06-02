package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
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

	timer       time.Duration
	remaining   time.Duration
	isRunning   bool
	audioPlayed bool
}

func initialModel(timer int) Model {
	m := Model{
		timer:       time.Minute * time.Duration(timer),
		remaining:   time.Minute * time.Duration(timer),
		isRunning:   true,
		audioPlayed: false,
	}

	return m
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
		case "r":
			m.remaining = m.timer
			m.isRunning = true
			m.audioPlayed = false

			return m, nil
		case " ":
			m.isRunning = !m.isRunning
		}
	case tickMsg:
		if m.isRunning && !m.audioPlayed {
			m.remaining -= time.Second
			if int(m.remaining.Seconds()) == -1 {
				m.remaining += time.Second
				m.audioPlayed = true
				playAudio()
			}
		}
		return m, tickCmd()
	}

	return m, nil
}

func (m Model) View() string {
	minutes := int(m.remaining.Minutes())
	seconds := int(m.remaining.Seconds()) % 60
	text := fmt.Sprintf("%02d:%02d", minutes, seconds)

	quitText := "[r] to reset, [space] to start/stop, [ctrl+c] or [q] to quit"
	text = figure.NewFigure(text, "banner3", false).String()

	centerStyle := lipgloss.NewStyle().Height(m.height).Width(m.width).Align(lipgloss.Center, lipgloss.Center)
	headingStyle := lipgloss.NewStyle().Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6D6D6D"))

	return centerStyle.Render(fmt.Sprintf("%s\n\n%s", headingStyle.Render(text), hintStyle.Render(quitText)))
}

func playAudio() {
	file, err := os.Open("audio/time_up.mp3")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the audio file
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// Initialize the speaker
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	// Play the audio
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the audio to finish playing
	<-done
}

func main() {
	var (
		timer int
		err   error
	)
	args := os.Args
	if len(args) > 1 {
		timer, err = strconv.Atoi(args[1])
		if err != nil {
			panic("Enter a valid integer")
		}
	}
	p := tea.NewProgram(initialModel(timer), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("failed to run")
		os.Exit(1)
	}
}
