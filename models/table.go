package models

import (
	"fmt"

	"github.com/araujofs/binds-editor/files"
	myKeys "github.com/araujofs/binds-editor/help"
	msgs "github.com/araujofs/binds-editor/messages"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type Table struct {
	binds  []*files.Bind
	cursor int
	table  table.Model
	help   help.Model
	keys   myKeys.KeyMap
	width  int
}

func InitTable() (*Table, tea.Msg) {
	terminalWidth, _, err := term.GetSize(0)

	if err != nil {
		terminalWidth = 160
	}

	binds, err := files.ReadBindsFile("/home/arthur/.config/hypr/bindings.conf")

	if err != nil {
		fmt.Printf("%s", err)
		return nil, func() tea.Msg { return msgs.NewErrorMsg(err) }
	}

	columns := []table.Column{
		{Title: "Shortcut", Width: 20},
		{Title: "Type", Width: 20},
		{Title: "Action", Width: 20},
		{Title: "Description", Width: 20},
		{Title: "Flags", Width: 20},
	}

	rows := []table.Row{}

	for _, bind := range binds {
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
		Foreground(lipgloss.Color("#fff")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderTop(true).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Background(lipgloss.Color("#7c53de")).
		Foreground(lipgloss.Color("#fff")).
		Bold(true)

	t.SetStyles(s)

	return &Table{
		binds:  []*files.Bind{},
		cursor: 0,
		table:  t,
		help:   help.New(),
		keys:   myKeys.Keys,
		width:  terminalWidth,
	}, nil
}

func (m Table) Init() tea.Cmd {
	return tea.SetWindowTitle("Binds Editor")
}

func (m Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Close):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
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

func (m Table) View() string {
	baseStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width)
	tableStyle := lipgloss.NewStyle().
		BorderBottom(true).
		Bold(false).Align(lipgloss.Center)

	paddedWidth, spaces := m.width-5, 16

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
