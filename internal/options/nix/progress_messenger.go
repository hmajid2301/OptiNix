package nix

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type ProgressMessenger struct {
	useTUI      bool
	spinner     *SpinnerMessenger
	progressBar *ProgressBarMessenger
	steps       []string
	totalSteps  int
}

var (
	plainSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	plainInfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
)

func NewProgressMessenger(useTUI bool) *ProgressMessenger {
	pm := &ProgressMessenger{
		useTUI:     useTUI,
		steps:      []string{},
		totalSteps: 40,
	}

	if useTUI {
		pm.progressBar = NewProgressBarMessenger()
		pm.progressBar.Start()
	}

	return pm
}

func (p *ProgressMessenger) Send(msg string) {
	if p.useTUI && p.progressBar != nil {
		p.steps = append(p.steps, msg)
		currentStep := len(p.steps)
		percent := float64(currentStep) / float64(p.totalSteps)
		if percent > 1.0 {
			percent = 1.0
		}
		p.progressBar.SendWithProgress(msg, percent)
	} else {
		fmt.Println(plainInfoStyle.Render("→ " + msg))
	}
}

func (p *ProgressMessenger) Finish(msg string) {
	if p.useTUI && p.progressBar != nil {
		p.progressBar.Finish(msg)
	} else {
		fmt.Println(plainSuccessStyle.Render("✓ " + msg))
	}
}

func (p *ProgressMessenger) Error(msg string, err error) {
	errorMark := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗")
	if p.useTUI && p.progressBar != nil {
		p.progressBar.Error(msg, err)
	} else {
		fmt.Printf("%s %s: %v\n", errorMark, msg, err)
	}
}

func (p *ProgressMessenger) Stop() {
	if p.useTUI && p.progressBar != nil {
		p.progressBar.Stop()
	}
}
