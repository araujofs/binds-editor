package ui

import (
	"fmt"
	"os"

	"github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	keys "github.com/araujofs/binds-editor/help"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileSearch struct {
	filePicker     filepicker.Model
	help           help.Model
	config         *configuration.Configuration
	message        string
	emptyDirectory bool
}

func InitFileSearch(configuration *configuration.Configuration) (*FileSearch, tea.Cmd) {
	fp := filepicker.New()

	fp.AllowedTypes = []string{".conf"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.ShowHidden = true
	fp.KeyMap = keys.GetFilepickerKeyMap()
	fp.ShowSize = true
	fp.Styles.EmptyDirectory = fp.Styles.EmptyDirectory.PaddingLeft(0)
	h := help.New()
	h.ShowAll = true

	m := FileSearch{
		filePicker:     fp,
		help:           h,
		config:         configuration,
		message:        "",
		emptyDirectory: false,
	}

	if consts.WindowSize.Height != 0 {
		setFilePickerHeight(&m)
	}

	return &m, fp.Init()
}

func (m FileSearch) Init() tea.Cmd {
	return nil
}

func (m FileSearch) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		consts.WindowSize = msg
		setFilePickerHeight(&m)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.CoreKeys.Close):
			return m, tea.Quit

		case key.Matches(msg, keys.FilePickerKeys.Select):
			m.filePicker, cmd = m.filePicker.Update(msg)

			if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
				selection := InitFileSelection(&path, m.config)

				return selection, cmd
			}

			return m, cmd

		case key.Matches(msg, keys.FilePickerKeys.Open):
			m.filePicker, cmd = m.filePicker.Update(msg)

			if isDirEmpty(m.filePicker.CurrentDirectory) {
				m.emptyDirectory = true
			}

			return m, cmd

		case key.Matches(msg, keys.FilePickerKeys.Back):
			m.emptyDirectory = false

		case key.Matches(msg, keys.FilePickerKeys.GoBack):
			selection := InitFileSelection(nil, m.config)

			return selection, cmd

		case key.Matches(msg, keys.CoreKeys.Help):
			return m, cmd
		}
	}

	m.filePicker, cmd = m.filePicker.Update(msg)

	return m, cmd

}

func (m FileSearch) View() string {
	title := fmt.Sprintf("Binds Editor | File Searching%s\n", m.message)

	if m.emptyDirectory {
		return consts.FullScreenFileSearchStyle.Render(title + m.filePicker.View() + lipgloss.NewStyle().MarginTop(1).Render(m.help.View(keys.FilePickerKeys)))
	}

	return consts.FullScreenFileSearchStyle.Render((title + m.filePicker.View() + m.help.View(keys.FilePickerKeys)))
}

func setFilePickerHeight(model *FileSearch) {
	top, _, bottom, _ := consts.FullScreenStyle.GetMargin()
	windowHeight := consts.WindowSize.Height

	if !model.help.ShowAll {
		model.filePicker.SetHeight(windowHeight - top - bottom - 2)
		return
	}

	helpHeight := len(keys.FilePickerKeys.FullHelp()[0])
	model.filePicker.SetHeight(windowHeight - top - bottom - 1 - helpHeight)
}

func isDirEmpty(path string) bool {
	file, _ := os.Open(path)
	defer file.Close()

	_, err := file.ReadDir(1)
	if err != nil {
		return true
	}

	if os.IsNotExist(err) {
		return true
	}

	return false
}
