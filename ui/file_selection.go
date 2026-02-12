package ui

import (
	"fmt"
	"strings"

	config "github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	keys "github.com/araujofs/binds-editor/help"
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
	mode             mode
	selectedFilePath string
	selectedFileName string
	message          string
}

func InitFileSelection(path *string) *FileSelection {
	config := config.GetConfigData()

	items := filesToItems(config.Files)

	fileList := list.New(items, list.NewDefaultDelegate(), 8, 8)
	fileList.SetShowTitle(false)
	fileList.AdditionalFullHelpKeys = keys.FileSelectionKeys.ShortHelp
	fileList.Help.ShowAll = true

	input := textinput.New()
	input.Placeholder = "General config"
	input.CharLimit = 50
	input.Width = 20

	model := &FileSelection{
		config:           config,
		list:             fileList,
		input:            input,
		mode:             navigating,
		selectedFilePath: "",
		selectedFileName: "",
		message:          "",
	}

	if path != nil {
		model.selectedFilePath = *path
		model.input.Focus()
		model.mode = adding
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		consts.WindowSize = msg
		m.setListSize()
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.config.SaveConfiguration()
			return m, tea.Quit
		}

		if m.input.Focused() && m.mode != navigating {
			switch {
			case key.Matches(msg, keys.FileSelectionInputKeys.Enter):
				configFileName := m.input.Value()
				m.input.SetValue("")
				m.input.Blur()
				m.setListSize()

				if strings.Trim(configFileName, " ") == "" {
					m.selectedFileName = ""
					m.selectedFilePath = ""
					m.mode = navigating
					return m, cmd
				}

				var err error
				if m.mode == adding {
					err = m.config.AddFile(m.selectedFilePath, configFileName)
				}

				if m.mode == editing {
					err = m.config.EditFile(m.selectedFileName, configFileName)
				}

				if err != nil {
					m.message = err.Error()
				}

				m.selectedFileName = ""
				m.selectedFilePath = ""
				m.mode = navigating
				m.list.SetItems(filesToItems(m.config.Files))
			case key.Matches(msg, keys.FileSelectionInputKeys.Back):
				m.input.SetValue("")
				m.input.Blur()
			}

			m.input, cmd = m.input.Update(msg)
		} else {
			if key.Matches(msg, keys.CoreKeys.Close) {
				m.config.SaveConfiguration()
				return m, tea.Quit
			}

			switch {
			case key.Matches(msg, keys.FileSelectionKeys.Save):
				err := m.config.SaveConfiguration()

				if err != nil {
					m.message = err.Error()
				}

			case key.Matches(msg, keys.FileSelectionKeys.Add):
				return InitFileSearch()

			case key.Matches(msg, keys.FileSelectionKeys.Delete):
				selectedItem := m.list.SelectedItem()
				if selectedItem == nil {
					return m, cmd
				}

				err := m.config.RemoveFile(m.list.SelectedItem().FilterValue())
				if err != nil {
					m.message = err.Error()
				}

				m.list.SetItems(filesToItems(m.config.Files))

			case key.Matches(msg, keys.FileSelectionKeys.Edit):
				selectedItem := m.list.SelectedItem()
				if selectedItem == nil {
					return m, cmd
				}

				selectedFileIdx := config.SearchSlice(m.config.Files, "Name", selectedItem.FilterValue())
				if selectedFileIdx == -1 {
					return m, cmd
				}

				selectedFile := m.config.Files[selectedFileIdx]
				if selectedFile == nil {
					return m, cmd
				}

				m.input.Focus()
				m.input.SetValue(selectedFile.Name)
				m.selectedFileName = selectedFile.Name
				m.mode = editing

				return m, cmd
			case key.Matches(msg, keys.FileSelectionKeys.Help):
				return m, cmd
			}

			m.list, cmd = m.list.Update(msg)
		}
	}

	return m, cmd
}

func (m FileSelection) View() string {
	title := fmt.Sprintf("Binds Editor | File Selection%s\n", m.message)

	if m.input.Focused() {
		return consts.FullScreenStyle.Render(title + m.list.View() + "\n\n" + m.input.View())
	}

	return consts.FullScreenStyle.Render(title + m.list.View())
}

func filesToItems(files []*config.File) []list.Item {
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
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-3)
		return
	}

	m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
}
