package ui

import (
	"fmt"

	consts "github.com/araujofs/binds-editor/constants"
	"github.com/araujofs/binds-editor/files"
	keys "github.com/araujofs/binds-editor/help"
	msgs "github.com/araujofs/binds-editor/messages"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	paddingStyle = lipgloss.NewStyle().Padding(0, consts.DefaultPadding)
	tableStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderBottomForeground(lipgloss.Color("240"))
)

type Table struct {
	binds  []*files.Bind
	cursor int
	table  table.Model
	help   help.Model
	keys   keys.TableKeyMap
}

func InitTable() (*Table, tea.Msg) {
	binds, err := files.ReadBindsFile("/home/arthur/.config/hypr/bindings.conf")

	if err != nil {
		fmt.Printf("%s", err)
		return nil, func() tea.Msg { return msgs.NewErrorMsg(err) }
	}

	columns := []table.Column{
		{Title: "Shortcut", Width: 20},
		{Title: "Description", Width: 20},
		{Title: "Type", Width: 20},
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
		table.WithHeight(getTableHeight(len(rows))),
		table.WithKeyMap(table.DefaultKeyMap()),
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
		binds:  binds,
		cursor: 0,
		table:  t,
		help:   help.New(),
		keys:   keys.TableKeys,
	}, nil
}

func (m Table) Init() tea.Cmd {
	return tea.SetWindowTitle("Binds Editor")
}

func (m Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		consts.WindowSize = msg
		m.table.SetHeight(getTableHeight(len(m.table.Rows())))
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
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m Table) View() string {

	width := consts.WindowSize.Width - consts.DefaultPadding

	columns := []table.Column{
		{Title: "Shortcut", Width: width / 4},
		{Title: "Description", Width: width / 4},
		{Title: "Type", Width: width / 4},
		{Title: "Flags", Width: width / 4},
	}

	m.table.SetColumns(columns)

	helpView := "\n" + m.help.View(m.keys)
	return paddingStyle.Render("Binds Editor | Binds: " + fmt.Sprint(len(m.binds)) + "\n" + tableStyle.Render(m.table.View()) + helpView + "\n")
}

func getTableHeight(rowsSize int) int {
	var tableHeight int

	tableHeight = min(consts.WindowSize.Height-consts.DefaultPadding-1, rowsSize)

	return tableHeight
}
