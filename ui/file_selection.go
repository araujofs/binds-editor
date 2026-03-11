package ui

import (
	config "github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	keys "github.com/araujofs/binds-editor/help"
	msgs "github.com/araujofs/binds-editor/messages"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type FileSelection struct {
	*consts.GlobalState
	list list.Model
	help help.Model
	msgs.InfoModel
}

func InitFileSelection(globalState *consts.GlobalState) *FileSelection {
	if globalState == nil {
		globalState = &consts.GlobalState{
			Configuration: config.GetConfigData(),
		}
	}

	fileList := list.New([]list.Item{}, list.NewDefaultDelegate(), 8, 8)
	fileList.SetShowTitle(false)
	fileList.SetShowHelp(false)
	fileList.SetShowFilter(false)
	fileList.SetFilteringEnabled(false)

	model := &FileSelection{
		GlobalState: globalState,
		list:        fileList,
		help:        help.New(),
		InfoModel:   msgs.GetDefaultInfoModel(),
	}

	items := model.filesToItems()
	model.list.SetItems(items)

	if consts.WindowSize.Height != 0 {
		model.setListSize()
	}

	return model
}

func (m FileSelection) Init() tea.Cmd {
	return tea.SetWindowTitle("Binds Editor")
}

func (m FileSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgs.ErrorMsg:
		m.Message = nil
		m.Error = &msg.Message

	case msgs.MessageMsg:
		m.Error = nil
		m.Message = &msg.Message

	case tea.WindowSizeMsg:
		consts.WindowSize = msg
		m.setListSize()

	case tea.KeyMsg:
		m.Message = nil
		m.Error = nil

		if msg.String() == "ctrl+c" {
			m.GlobalState.Configuration.SaveConfiguration()
			return m, tea.Quit
		}

		return updateList(msg, &m)
	}

	return m, nil
}

func (m FileSelection) View() string {
	title := "Binds Editor | Select File"
	title += m.GetStyledMessage()

	return consts.FullScreenStyle.Render(title + m.list.View() + "\n\n" + m.help.FullHelpView(keys.FileSelectionKeys.FullHelp()))
}

func (m *FileSelection) filesToItems() []list.Item {
	files := m.GlobalState.Configuration.Files

	if len(files) == 0 {
		return make([]list.Item, 0)
	}

	items := make([]list.Item, len(files))

	for i, file := range files {
		items[i] = list.DefaultItem(file)
	}

	return items
}

func (m *FileSelection) setListSize() {
	top, right, bottom, left := consts.FullScreenStyle.GetMargin()
	msg := consts.WindowSize

	m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-6)
}

func updateList(msg tea.KeyMsg, m *FileSelection) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keys.CoreKeys.Close):
		return m, tea.Quit

	// the list interprets "esc" as Quit
	case msg.String() == "esc":
		return m, cmd

	case key.Matches(msg, keys.FileSelectionKeys.Select):
		selectedItem := m.list.SelectedItem()
		if selectedItem == nil {
			return m, nil
		}

		selectedFileIdx := config.SearchSlice(m.GlobalState.Configuration.Files, "Name", selectedItem.FilterValue())
		if selectedFileIdx == -1 {
			return m, msgs.SendErrorMsg("selected file doesn't exist in your configuration")
		}

		selectedFile := m.GlobalState.Configuration.Files[selectedFileIdx]
		if selectedFile == nil {
			return m, nil
		}

		m.GlobalState.SelectedFile = selectedFile

		return InitTable(m.GlobalState)

	case key.Matches(msg, keys.FileSelectionKeys.Save):
		err := m.GlobalState.Configuration.SaveConfiguration()
		if err != nil {
			cmd = msgs.SendErrorMsg(err.Error())
		}
		return m, cmd

	case key.Matches(msg, keys.FileSelectionKeys.Add):
		return InitFileManipulation(nil, m.GlobalState)

	case key.Matches(msg, keys.FileSelectionKeys.Delete):
		selectedItem := m.list.SelectedItem()
		if selectedItem == nil {
			return m, nil
		}

		err := m.GlobalState.Configuration.RemoveFile(m.list.SelectedItem().FilterValue())
		if err != nil {
			cmd = msgs.SendErrorMsg(err.Error())
			return m, cmd
		}

		m.list.SetItems(m.filesToItems())

		return m, cmd

	case key.Matches(msg, keys.FileSelectionKeys.Edit):
		selectedItem := m.list.SelectedItem()
		if selectedItem == nil {
			return m, nil
		}

		selectedFileIdx := config.SearchSlice(m.GlobalState.Configuration.Files, "Name", selectedItem.FilterValue())
		if selectedFileIdx == -1 {
			return m, msgs.SendErrorMsg("selected file doesn't exist in your configuration")
		}

		selectedFile := m.GlobalState.Configuration.Files[selectedFileIdx]
		if selectedFile == nil {
			return m, nil
		}

		return InitFileManipulation(selectedFile, m.GlobalState)

	case key.Matches(msg, keys.FileSelectionKeys.Help):
		return m, nil
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}
