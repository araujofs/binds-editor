package ui

import (
	"strings"

	config "github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	keys "github.com/araujofs/binds-editor/help"
	msgs "github.com/araujofs/binds-editor/messages"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	navigating mode = iota + 1
	adding
	editing
)

type FileSelection struct {
	config           *config.Configuration
	list             list.Model
	input            textinput.Model
	help             help.Model
	mode             mode
	selectedFilePath string
	selectedFileName string
	msgs.InfoModel
}

func InitFileSelection(path *string, configuration *config.Configuration) *FileSelection {
	if configuration == nil {
		configuration = config.GetConfigData()
	}

	fileList := list.New([]list.Item{}, list.NewDefaultDelegate(), 8, 8)
	fileList.SetShowTitle(false)
	fileList.SetShowHelp(false)
	fileList.SetShowFilter(false)
	fileList.SetFilteringEnabled(false)

	input := textinput.New()
	input.Placeholder = "General config"
	input.CharLimit = 50
	input.Width = 20

	model := &FileSelection{
		config:           configuration,
		list:             fileList,
		input:            input,
		help:             help.New(),
		mode:             navigating,
		selectedFilePath: "",
		selectedFileName: "",
		InfoModel:        msgs.GetDefaultInfoModel(),
	}

	items := model.filesToItems()
	model.list.SetItems(items)

	if path != nil {
		model.selectedFilePath = *path
		model.input.Focus()
		model.mode = adding
		model.setListSize()
	}

	if consts.WindowSize.Height != 0 {
		model.setListSize()
	}

	return model
}

func (m FileSelection) Init() tea.Cmd {
	return nil
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
			m.config.SaveConfiguration()
			return m, tea.Quit
		}

		if m.input.Focused() && m.mode != navigating {
			return updateInput(msg, &m)
		}

		return updateList(msg, &m)
	}

	return m, nil
}

func (m FileSelection) View() string {
	title := "Binds Editor | Select File"
	title += m.GetStyledMessage()

	if m.input.Focused() {
		return consts.FullScreenStyle.Render(title + m.list.View() + "\n" + m.help.FullHelpView(keys.FileSelectionKeys.FullHelp()) + "\n" + m.input.View())
	}

	return consts.FullScreenStyle.Render(title + m.list.View() + "\n\n" + m.help.FullHelpView(keys.FileSelectionKeys.FullHelp()))
}

func (m *FileSelection) filesToItems() []list.Item {
	files := m.config.Files

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

	if m.input.Focused() {
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-6)
		return
	}

	m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-6)
}

func updateInput(msg tea.KeyMsg, m *FileSelection) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keys.FileSelectionInputKeys.Enter):
		configFileName := m.input.Value()
		m.input.SetValue("")
		m.input.Blur()
		m.setListSize()

		if strings.Trim(configFileName, " ") == "" {
			m.mode = navigating
		}

		var err error
		if m.mode == adding {
			err = m.config.AddFile(m.selectedFilePath, configFileName)
		}

		if m.mode == editing {
			err = m.config.EditFile(m.selectedFileName, configFileName)
		}

		if err != nil {
			cmd = msgs.SendErrorMsg(err.Error())
		}

		m.selectedFileName = ""
		m.selectedFilePath = ""
		m.mode = navigating
		m.list.SetItems(m.filesToItems())

		return m, cmd

	case key.Matches(msg, keys.FileSelectionInputKeys.Back):
		m.input.SetValue("")
		m.input.Blur()

		m.selectedFileName = ""
		m.selectedFilePath = ""

		m.mode = navigating
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
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

		selectedFileIdx := config.SearchSlice(m.config.Files, "Name", selectedItem.FilterValue())
		if selectedFileIdx == -1 {
			return m, msgs.SendErrorMsg("selected file doesn't exist in your configuration")
		}

		selectedFile := m.config.Files[selectedFileIdx]
		if selectedFile == nil {
			return m, nil
		}

		return InitTable(selectedFile, m.config)

	case key.Matches(msg, keys.FileSelectionKeys.Save):
		err := m.config.SaveConfiguration()
		if err != nil {
			cmd = msgs.SendErrorMsg(err.Error())
		}
		return m, cmd

	case key.Matches(msg, keys.FileSelectionKeys.Add):
		return InitFileSearch(m.config)

	case key.Matches(msg, keys.FileSelectionKeys.Delete):
		selectedItem := m.list.SelectedItem()
		if selectedItem == nil {
			return m, nil
		}

		err := m.config.RemoveFile(m.list.SelectedItem().FilterValue())
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

		selectedFileIdx := config.SearchSlice(m.config.Files, "Name", selectedItem.FilterValue())
		if selectedFileIdx == -1 {
			return m, msgs.SendErrorMsg("selected file doesn't exist in your configuration")
		}

		selectedFile := m.config.Files[selectedFileIdx]
		if selectedFile == nil {
			return m, nil
		}

		m.input.Focus()
		m.input.SetValue(selectedFile.Name)
		m.selectedFileName = selectedFile.Name
		m.mode = editing
		m.setListSize()

		return m, nil

	case key.Matches(msg, keys.FileSelectionKeys.Help):
		return m, nil
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}
