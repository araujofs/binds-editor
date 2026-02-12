package ui

import (
	"fmt"

	config "github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	keys "github.com/araujofs/binds-editor/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type FileSelection struct {
	config      *config.Configuration
	list        list.Model
	input       textinput.Model
	newFilePath string
	message     string
}

func InitFileSelection(path *string) *FileSelection {
	config := config.GetConfigData()

	items := filesToItems(config.Files)

	fileList := list.New(items, list.NewDefaultDelegate(), 8, 8)
	fileList.SetShowTitle(false)
	fileList.AdditionalShortHelpKeys = keys.FileSelectionKeys.ShortHelp

	input := textinput.New()
	input.Placeholder = "General config"
	input.CharLimit = 50
	input.Width = 20

	model := &FileSelection{
		config:      config,
		list:        fileList,
		input:       input,
		newFilePath: "",
		message:     "",
	}

	if path != nil {
		model.newFilePath = *path
		model.input.Focus()
		model.list.AdditionalShortHelpKeys = keys.FileSelectionInputKeys.ShortHelp
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
		if key.Matches(msg, keys.CoreKeys.Close) {
			m.config.SaveConfiguration()
			return m, tea.Quit
		}
		if key.Matches(msg, keys.FileSelectionKeys.Save) {
			err := m.config.SaveConfiguration()

			if err != nil {
				m.message = err.Error()
			}
		}
		if m.input.Focused() {
			switch {
			case key.Matches(msg, keys.FileSelectionInputKeys.Enter):
				configFileName := m.input.Value()
				m.input.SetValue("")
				m.input.Blur()
				m.setListSize()

				err := m.config.AddFile(m.newFilePath, configFileName)

				if err != nil {
					m.message = err.Error()
				}

				m.newFilePath = ""
				m.list.SetItems(filesToItems(m.config.Files))
			case key.Matches(msg, keys.FileSelectionInputKeys.Back):
				m.input.SetValue("")
				m.input.Blur()
			}

			m.input, cmd = m.input.Update(msg)
		} else {
			switch {
			case key.Matches(msg, keys.FileSelectionKeys.Find):
				return InitFileSearch()
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
