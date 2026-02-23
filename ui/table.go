package ui

import (
	"github.com/araujofs/binds-editor/binds"
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
	binds  []*binds.Bind
	cursor int
	table  table.Model
	help   help.Model
	*consts.GlobalState
	msgs.InfoModel
}

func InitTable(globalState *consts.GlobalState) (tea.Model, tea.Cmd) {
	binds, err := binds.ParseBindsFile(globalState.SelectedFile.Path)
	if err != nil {
		return InitFileSelection(nil, globalState), msgs.SendErrorMsg(err.Error())
	}

	columns := []table.Column{
		{Title: "Shortcut", Width: 20},
		{Title: "Type", Width: 20},
		{Title: "Description", Width: 20},
		{Title: "Action Type", Width: 20},
		{Title: "Flags", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(bindsToRows(binds)),
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
		binds:       binds,
		cursor:      0,
		table:       t,
		help:        help.New(),
		GlobalState: globalState,
		InfoModel:   msgs.GetDefaultInfoModel(),
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
	case msgs.UpdateTableBindsMsg:
		binds, err := binds.ParseBindsFile(m.GlobalState.SelectedFile.Path)
		if err != nil {
			return InitFileSelection(nil, m.GlobalState), msgs.SendErrorMsg(err.Error())
		}

		m.binds = binds
		m.table.SetRows(bindsToRows(m.binds))

		return m, nil

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
			return InitCreate(m.GlobalState)

		case key.Matches(msg, keys.TableKeys.Edit):
			selectedBind := m.getSelectedBind()
			if selectedBind == nil {
				return m, nil
			}

			return InitEdit(selectedBind, m.GlobalState)

		case key.Matches(msg, keys.TableKeys.Delete):
			selectedBind := m.getSelectedBind()
			if selectedBind == nil {
				return m, nil
			}

			err := selectedBind.Delete()
			if err != nil {
				return m, msgs.SendErrorMsg(err.Error())
			}

			return m, msgs.SendUpdateTableBindsMsg()

		case key.Matches(msg, keys.TableKeys.Unbind):
			selectedBind := m.getSelectedBind()
			if selectedBind == nil {
				return m, nil
			}

			if selectedBind.Type == binds.Unbind {
				return m, msgs.SendMessageMsg("can't unbind an unbind")
			}

			err := selectedBind.Unbind()
			if err != nil {
				return m, msgs.SendErrorMsg(err.Error())
			}

			return m, msgs.SendUpdateTableBindsMsg()

		case key.Matches(msg, keys.TableKeys.Comment):
			selectedBind := m.getSelectedBind()
			if selectedBind == nil {
				return m, nil
			}

			err := selectedBind.Comment()
			if err != nil {
				return m, msgs.SendErrorMsg(err.Error())
			}

			return m, msgs.SendUpdateTableBindsMsg()

		case key.Matches(msg, keys.TableKeys.Details):
			return m, nil

		case key.Matches(msg, keys.TableKeys.GoBack):
			return InitFileSelection(nil, m.GlobalState), nil

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

func (m *Table) getSelectedBind() *binds.Bind {
	if len(m.binds) <= 0 {
		return nil
	}

	selectedBind := m.binds[m.table.Cursor()]

	return selectedBind
}

func bindsToRows(binds []*binds.Bind) []table.Row {
	rows := []table.Row{}

	for _, bind := range binds {
		rows = append(rows, bind.KeybindToRow())
	}

	return rows
}
