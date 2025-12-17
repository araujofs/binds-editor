package main

import (
	"fmt"
	"os"

	"github.com/araujofs/binds-editor/files"
	myKeys "github.com/araujofs/binds-editor/help"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type model struct {
	binds    []*files.Bind
	cursor   int
	table table.Model
	help help.Model
	keys myKeys.KeyMap
	width int
}

func initialModel() model {
	terminalWidth, _, err := term.GetSize(0)

	if err != nil {
		terminalWidth = 160
	}

	return model{
		binds:    []*files.Bind{},
		cursor: 0,
		help: help.New(),
		keys: myKeys.Keys,
		width: terminalWidth,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Binds Editor")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.binds)-1 {
				m.cursor++
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m model) View() string {
	baseStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8f5aff")).Align(lipgloss.Center).Width(m.width)
	tableStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#fafffa")).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false).Align(lipgloss.Center)

	paddedWidth, spaces := m.width - 5, 16

	columns := []table.Column{
		{Title: "Shortcut", Width: (paddedWidth / spaces) * 2},
		{Title: "Type", Width: (paddedWidth / spaces) * 1},
		{Title: "Action", Width: (paddedWidth / spaces) * 10},
		{Title: "Description", Width: (paddedWidth / spaces) * 2},
		{Title: "Flags", Width: (paddedWidth / spaces) * 1},
	}

	m.table.SetColumns(columns)

	helpView := "\n" + m.help.View(m.keys)
	return baseStyle.Render("Binds Editor\n" + tableStyle.Render(m.table.View()) + helpView + "\n")
}

func main() {
	m := initialModel()

	columns := []table.Column{
		{Title: "Shortcut", Width: 20},
		{Title: "Type", Width: 20},
		{Title: "Action", Width: 20},
		{Title: "Description", Width: 20},
		{Title: "Flags", Width: 20},
	}

	rows := []table.Row{}

	binds, err := files.ReadBindsFile("/home/arthur/.config/hypr/bindings.conf")

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	for _, bind := range  binds {
		rows = append(rows, bind.KeybindToRow())
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(lipgloss.Color("245")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderTop(true).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)



	m.table = t


	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
