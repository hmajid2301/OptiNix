package nix

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type ProgressBarMessenger struct {
	program *tea.Program
	model   *progressBarModel
	wg      sync.WaitGroup
}

type progressBarModel struct {
	progress progress.Model
	message  string
	finished bool
	err      error
	percent  float64
}

type progressMsg struct {
	message string
	percent float64
}

func (m progressBarModel) Init() tea.Cmd {
	return nil
}

func (m progressBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil

	case progressMsg:
		m.message = msg.message
		m.percent = msg.percent
		return m, nil

	case string:
		m.message = msg
		return m, nil

	case finishMsg:
		m.finished = true
		m.message = string(msg)
		m.percent = 1.0
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

func (m progressBarModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s %s\n", crossMark, errorStyle.Render(m.message))
	}
	if m.finished {
		return fmt.Sprintf("%s %s\n", checkMark, successStyle.Render(m.message))
	}

	pad := strings.Repeat(" ", 2)
	return "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n" +
		pad + infoStyle.Render(m.message) + "\n"
}

func NewProgressBarMessenger() *ProgressBarMessenger {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
	)

	model := &progressBarModel{
		progress: prog,
		message:  "Starting...",
		finished: false,
		percent:  0.0,
	}

	return &ProgressBarMessenger{
		model: model,
	}
}

func (p *ProgressBarMessenger) Start() {
	p.program = tea.NewProgram(p.model)
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.program.Run()
	}()
}

func (p *ProgressBarMessenger) Send(msg string) {
	if p.program != nil {
		p.program.Send(msg)
	}
}

func (p *ProgressBarMessenger) SendWithProgress(msg string, percent float64) {
	if p.program != nil {
		p.program.Send(progressMsg{message: msg, percent: percent})
	}
}

func (p *ProgressBarMessenger) Finish(msg string) {
	if p.program != nil {
		p.program.Send(finishMsg(msg))
	}
}

func (p *ProgressBarMessenger) Error(msg string, err error) {
	if p.program != nil {
		p.program.Send(errorMsg{message: msg, err: err})
	}
}

func (p *ProgressBarMessenger) Stop() {
	if p.program != nil {
		p.program.Quit()
		p.wg.Wait()
	}
}
