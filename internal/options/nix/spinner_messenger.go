package nix

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	checkMark    = successStyle.Render("✓")
	crossMark    = errorStyle.Render("✗")
)

type SpinnerMessenger struct {
	program *tea.Program
	model   *spinnerModel
	wg      sync.WaitGroup
}

type spinnerModel struct {
	spinner  spinner.Model
	message  string
	finished bool
	err      error
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case string:
		m.message = msg
		return m, nil

	case finishMsg:
		m.finished = true
		m.message = string(msg)
		return m, tea.Quit

	case errorMsg:
		m.finished = true
		m.err = msg.err
		m.message = msg.message
		return m, tea.Quit

	default:
		return m, nil
	}
}

func (m spinnerModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s %s\n", crossMark, errorStyle.Render(m.message))
	}
	if m.finished {
		return fmt.Sprintf("%s %s\n", checkMark, successStyle.Render(m.message))
	}
	return fmt.Sprintf("%s %s\n", spinnerStyle.Render(m.spinner.View()), infoStyle.Render(m.message))
}

type finishMsg string
type errorMsg struct {
	message string
	err     error
}

func NewSpinnerMessenger() *SpinnerMessenger {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	model := &spinnerModel{
		spinner:  s,
		message:  "Starting...",
		finished: false,
	}

	return &SpinnerMessenger{
		model: model,
	}
}

func (s *SpinnerMessenger) Start() {
	s.program = tea.NewProgram(s.model)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.program.Run()
	}()
}

func (s *SpinnerMessenger) Send(msg string) {
	if s.program != nil {
		s.program.Send(msg)
	}
}

func (s *SpinnerMessenger) Finish(msg string) {
	if s.program != nil {
		s.program.Send(finishMsg(msg))
	}
}

func (s *SpinnerMessenger) Error(msg string, err error) {
	if s.program != nil {
		s.program.Send(errorMsg{message: msg, err: err})
	}
}

func (s *SpinnerMessenger) Stop() {
	if s.program != nil {
		s.program.Quit()
		s.wg.Wait()
	}
}
