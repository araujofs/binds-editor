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
		return InitFileSelection(globalState), msgs.SendErrorMsg(err.Error())
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
	switch msg := msg.(type) {
	case msgs.UpdateTableBindsMsg:
		binds, err := binds.ParseBindsFile(m.GlobalState.SelectedFile.Path)
		if err != nil {
			return InitFileSelection(m.GlobalState), msgs.SendErrorMsg(err.Error())
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

		return m.handleKeyPress(msg)
	}

	return m, nil
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

	var fullHelp [][]key.Binding
	if m.GlobalState.SelectedFile.Readonly {
		fullHelp = keys.ReadonlyTableKeys.FullHelp()
	} else {
		fullHelp = keys.TableKeys.FullHelp()
	}

	helpView := "\n" + m.help.FullHelpView(fullHelp)
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

func (m Table) handleKeyPress(pressedKey tea.KeyMsg) (tea.Model, tea.Cmd) {
	readonly := m.GlobalState.SelectedFile.Readonly

	switch {
	case key.Matches(pressedKey, keys.TableKeys.Close):
		return m, tea.Quit

	case key.Matches(pressedKey, keys.TableKeys.Up):
		break

	case key.Matches(pressedKey, keys.TableKeys.Down):
		break

	case key.Matches(pressedKey, keys.TableKeys.Create) && !readonly:
		return InitCreate(m.GlobalState)

	case key.Matches(pressedKey, keys.TableKeys.Edit) && !readonly:
		selectedBind := m.getSelectedBind()
		if selectedBind == nil {
			return m, nil
		}

		return InitEdit(selectedBind, m.GlobalState)

	case key.Matches(pressedKey, keys.TableKeys.Delete) && !readonly:
		selectedBind := m.getSelectedBind()
		if selectedBind == nil {
			return m, nil
		}

		err := selectedBind.Delete()
		if err != nil {
			return m, msgs.SendErrorMsg(err.Error())
		}

		return m, msgs.SendUpdateTableBindsMsg()

	case key.Matches(pressedKey, keys.TableKeys.Unbind):
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

	case key.Matches(pressedKey, keys.TableKeys.Comment) && !readonly:
		selectedBind := m.getSelectedBind()
		if selectedBind == nil {
			return m, nil
		}

		err := selectedBind.Comment()
		if err != nil {
			return m, msgs.SendErrorMsg(err.Error())
		}

		return m, msgs.SendUpdateTableBindsMsg()

	case key.Matches(pressedKey, keys.TableKeys.Details):
		return m, nil

	case key.Matches(pressedKey, keys.TableKeys.GoBack):
		return InitFileSelection(m.GlobalState), nil
	}

	updatedTable, cmd := m.table.Update(pressedKey)
	m.table = updatedTable

	return m, cmd
}
