package tui

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case CopySuccessMsg:
		m.view.statusMessage = "✓ Copied to clipboard"
		return m, nil

	case CopyErrorMsg:
		m.view.statusMessage = "✗ Failed to copy to clipboard"
		return m, nil

	case BrowserOpenMsg:
		if msg.success {
			m.view.statusMessage = "✓ Opened in browser"
		} else {
			m.view.statusMessage = "✗ Failed to open browser"
		}
		return m, nil

	case tea.KeyMsg:
		if m.view.showDetail {
			switch msg.String() {
			case "t", "esc":
				m.view.showDetail = false
				m.view.statusMessage = ""
				m.view.detailReady = false
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			case "y":
				if m.list.SelectedItem() != nil {
					return m, copyToClipboard(m.list.SelectedItem().(Item).OptionName)
				}
				return m, nil
			case "o":
				if m.list.SelectedItem() != nil {
					item := m.list.SelectedItem().(Item)
					url := getDocURL(item)
					if url != "" {
						return m, openInBrowser(url)
					}
				}
				return m, nil
			case "O":
				if m.list.SelectedItem() != nil {
					item := m.list.SelectedItem().(Item)
					url := getSourceURL(item)
					if url != "" {
						return m, openInBrowser(url)
					}
				}
				return m, nil
			case "r":
				if m.list.SelectedItem() != nil {
					item := m.list.SelectedItem().(Item)
					m.view.showDetail = false
					m.view.detailReady = false
					return m, m.filterByPrefix(item.OptionName)
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.view.detailViewport, cmd = m.view.detailViewport.Update(msg)
				return m, cmd
			}
		}

		switch {
		case msg.String() == "ctrl+c", msg.String() == "q":
			return m, tea.Quit

		case key.Matches(msg, m.keys.toggle):
			m.view.showDetail = !m.view.showDetail
			if m.view.showDetail {
				m.initDetailView()
			} else {
				m.view.detailReady = false
			}
			return m, nil

		case msg.String() == "enter":
			m.view.showDetail = true
			m.initDetailView()
			return m, nil

		case msg.String() == "n":
			m.sourceFilter.NixOS = !m.sourceFilter.NixOS
			return m, m.rebuildList()

		case msg.String() == "h":
			m.sourceFilter.HomeManager = !m.sourceFilter.HomeManager
			return m, m.rebuildList()

		case msg.String() == "d":
			m.sourceFilter.Darwin = !m.sourceFilter.Darwin
			return m, m.rebuildList()
		}

	case tea.WindowSizeMsg:
		m.display.width = msg.Width
		m.display.height = msg.Height
		h, v := m.docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-2)

		if m.view.showDetail && !m.view.detailReady {
			m.view.detailViewport = viewport.New(msg.Width, msg.Height-4)
			m.view.detailViewport.YPosition = 0
			m.view.detailReady = true

			selectedItem := m.list.SelectedItem().(Item)
			content := renderDetailedView(selectedItem)
			m.view.detailViewport.SetContent(content)
		} else if m.view.showDetail && m.view.detailReady {
			m.view.detailViewport.Width = msg.Width
			m.view.detailViewport.Height = msg.Height - 4
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case DoneMsg:
		m.allItems = msg.List
		m.totalOptions = len(msg.List)
		filteredItems := m.applySourceFilter(msg.List)
		cmds := []tea.Cmd{}
		for _, newItem := range filteredItems {
			insCmd := m.list.InsertItem(0, newItem)
			cmds = append(cmds, insCmd)
		}
		return m, tea.Batch(cmds...)
	}

	if !m.view.showDetail {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}


func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(text)
		if err != nil {
			return CopyErrorMsg{err: err}
		}
		return CopySuccessMsg{text: text}
	}
}


func (m *Model) initDetailView() {
	if m.display.width == 0 || m.display.height == 0 {
		m.view.showDetail = false
		return
	}

	if m.list.SelectedItem() == nil {
		m.view.showDetail = false
		return
	}

	m.view.detailViewport = viewport.New(m.display.width, m.display.height-4)
	m.view.detailViewport.YPosition = 0
	m.view.detailReady = true

	selectedItem := m.list.SelectedItem().(Item)
	content := renderDetailedView(selectedItem)
	m.view.detailViewport.SetContent(content)
}

func (m Model) rebuildList() tea.Cmd {
	filteredItems := m.applySourceFilter(m.allItems)

	return func() tea.Msg {
		return tea.Batch(
			m.list.SetItems(filteredItems),
		)()
	}
}

