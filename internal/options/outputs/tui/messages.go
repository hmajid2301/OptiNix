package tui

import "github.com/charmbracelet/bubbles/list"

type DoneMsg struct {
	List []list.Item
}

type CopySuccessMsg struct {
	text string
}

type CopyErrorMsg struct {
	err error
}
