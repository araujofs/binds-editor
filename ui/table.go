package ui

import (
	"github.com/araujofs/binds-editor/binds"
	config "github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
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
	binds        []*binds.Bind
	cursor       int
	table        table.Model
	help         help.Model
	config       *config.Configuration
	selectedFile *config.File
	msgs.InfoModel
}

func InitTable(selectedFile *config.File, configuration *config.Configuration) (tea.Model, tea.Cmd) {
	binds, err := binds.ReadBindsFile(selectedFile.Path)

	if err != nil {
		return InitFileSelection(nil, configuration), msgs.SendErrorMsg(err.Error())
	}

	columns := []table.Column{
		{Title: "Shortcut", Width: 20},
		{Title: "Type", Width: 20},
		{Title: "Description", Width: 20},
		{Title: "Action Type", Width: 20},
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

	s.Cell = s.Cell.AlignHorizontal(lipgloss.Center)

	t.SetStyles(s)

	model := &Table{
		binds:        binds,
		cursor:       0,
		table:        t,
		help:         help.New(),
		config:       configuration,
		selectedFile: selectedFile,
		InfoModel:    msgs.GetDefaultInfoModel(),
	}

	model.setTableHeight()

	return model, nil
}

func (m Table) Init() tea.Cmd {
	return nil
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
		case key.Matches(msg, keys.TableKeys.Close):
			return m, tea.Quit

		case key.Matches(msg, keys.TableKeys.Up):
			break

		case key.Matches(msg, keys.TableKeys.Down):
			break

		case key.Matches(msg, keys.TableKeys.Create):
			return m, nil

		case key.Matches(msg, keys.TableKeys.Edit):
			if len(m.binds) <= 0 {
				return m, nil
			}

			selectedBind := m.binds[m.table.Cursor()]

			return InitEdit(selectedBind, m.selectedFile, m.config)

		case key.Matches(msg, keys.TableKeys.Delete):
			return m, nil

		case key.Matches(msg, keys.TableKeys.Unbind):
			return m, nil

		case key.Matches(msg, keys.TableKeys.Comment):
			return m, nil

		case key.Matches(msg, keys.TableKeys.Details):
			return m, nil

		case key.Matches(msg, keys.TableKeys.GoBack):
			return InitFileSelection(nil, m.config), nil

		}

		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m Table) View() string {
	title := "Binds Editor | Bindings"

	title += m.GetStyledMessage()

	width := consts.WindowSize.Width - consts.DefaultPadding

	columns := []table.Column{
		{Title: "Shortcut", Width: (width / 20) * 7},
		{Title: "Type", Width: (width / 20) * 3},
		{Title: "Description", Width: (width / 20) * 4},
		{Title: "Action Type", Width: (width / 20) * 3},
		{Title: "Flags", Width: (width / 20) * 3},
	}

	m.table.SetColumns(columns)

	helpView := "\n" + m.help.FullHelpView(keys.TableKeys.FullHelp())
	return paddingStyle.Render(title + m.table.View() + helpView)
}

func (m *Table) setTableHeight() {
	fixedWindow := consts.WindowSize.Height - consts.DefaultPadding - 3

	m.table.SetHeight(fixedWindow)
}
