package ui

import (
	config "github.com/araujofs/binds-editor/configuration"
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
	// tableStyle   = lipgloss.NewStyle().
	// 		BorderStyle(lipgloss.NormalBorder()).
	// 		BorderBottom(true).
	// 		BorderBottomForeground(lipgloss.Color("240"))
)

type Table struct {
	binds  []*files.Bind
	cursor int
	table  table.Model
	help   help.Model
	keys   keys.TableKeyMap
	config *config.Configuration
	msgs.InfoModel
}

func InitTable(selectedFile *config.File, configuration *config.Configuration) (tea.Model, tea.Cmd) {
	binds, err := files.ReadBindsFile(selectedFile.Path)

	if err != nil {
		return InitFileSelection(nil, configuration), msgs.SendErrorMsg(err.Error())
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

	model := &Table{
		binds:     binds,
		cursor:    0,
		table:     t,
		help:      help.New(),
		keys:      keys.TableKeys,
		config:    configuration,
		InfoModel: msgs.GetDefaultInfoModel(),
	}

	model.setTableHeight()

	return model, nil
}

func (m Table) Init() tea.Cmd {
	return tea.SetWindowTitle("Binds Editor")
}

func (m Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case msgs.ErrorMsg:
		m.Message = nil
		m.Error = &msg.Message

	case msgs.MessageMsg:
		m.Error = nil
		m.Message = &msg.Message

	case tea.WindowSizeMsg:
		consts.WindowSize = msg
		m.setTableHeight()

	case tea.KeyMsg:
		m.Message = nil
		m.Error = nil

		switch {
		case key.Matches(msg, m.keys.Close):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			testing := "TESTANDO"
			m.Error = &testing
			if m.cursor < len(m.binds)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.GoBack):
			return InitFileSelection(nil, m.config), nil
		}
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m Table) View() string {
	title := "Binds Editor"

	title += m.GetStyledMessage()

	width := consts.WindowSize.Width - consts.DefaultPadding

	columns := []table.Column{
		{Title: "Shortcut", Width: width / 4},
		{Title: "Description", Width: width / 4},
		{Title: "Type", Width: width / 4},
		{Title: "Flags", Width: width / 4},
	}

	m.table.SetColumns(columns)

	helpView := "\n" + m.help.FullHelpView(m.keys.FullHelp())
	return paddingStyle.Render(title + m.table.View() + helpView)
}

func (m *Table) setTableHeight() {
	fixedWindow := consts.WindowSize.Height - consts.DefaultPadding - 2

	m.table.SetHeight(fixedWindow)
}
