package main

import (
	"github.com/charmbracelet/bubbles/list"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

const (
	projectView sessionState = iota
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	//these values are overwritten at start
	ListTitleColor = lipgloss.Color("#05b4ff")
	ListItemNameColor = lipgloss.Color("#4287f5")
	ListItemDescColor = lipgloss.Color("#2d579c")
)
type (
	sessionState int

	item struct {
		title, desc string
	}
	model struct {
//		State sessionState
		list list.Model
	
	}
)


func (i item) Title()       string { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }


func (m model) Init() tea.Cmd {
	return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
/*	case ffplayStruct.Status:
		switch msg {
		case "start":
			
		case "stop":
		}*/
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}


func (m model) View() string {
/*	switch m.State {
	case playerView:
		return docStyle.Render(m.
	case libraryView:*/
		return docStyle.Render(m.list.View())
}


func startUI(items []list.Item) {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
						Foreground(ListItemNameColor).
						BorderLeftForeground(ListItemDescColor)

	d.Styles.SelectedDesc = d.Styles.SelectedTitle.
						Foreground(ListItemDescColor).
						BorderLeftForeground(ListItemDescColor)

	m := model{list: list.New(items, d, 0, 0)}
	m.list.Title = "radio"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		ferr(err)
	}
}
