package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner              spinner.Model
	paused               bool
	resting              bool
	rounds               int
	currentRound         int
	countdownSeconds     int
	restDurationSeconds  int
	roundDurationSeconds int
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line
	return model{
		spinner:              s,
		paused:               false,
		resting:              false,
		rounds:               1,
		restDurationSeconds:  3,
		roundDurationSeconds: 5,
	}
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if !m.paused {
			m.countdownSeconds--
		}
		if m.countdownSeconds <= 0 {
			m.resting = !m.resting
			if m.resting {
				m.countdownSeconds = m.restDurationSeconds
			} else {
				m.currentRound++
				if m.currentRound <= m.rounds {
					m.countdownSeconds = m.roundDurationSeconds
				}
			}
		}
		return m, doTick()
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			m.paused = !m.paused
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().Align(lipgloss.Center)
	if m.currentRound >= m.rounds {
		return style.Render("Finished")
	}
	blocks := []string{
		style.Render("BJJ Timer by Shivan"),
		style.Render(
			m.spinner.View(),
			"Countdown:",
			fmt.Sprintf("%vs", strconv.Itoa(m.countdownSeconds)),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			style.Render("Remaining Rounds:", strconv.Itoa(m.rounds-m.currentRound)),
			style.Render("Current Round:", strconv.Itoa(m.currentRound+1)),
		),
	}
	if m.resting {
		blocks = append(blocks, style.Render("Resting"))
	} else {
		blocks = append(blocks, style.Render("Active"))
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		blocks...,
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
