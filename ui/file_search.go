package ui

import (
	"os"

	"github.com/araujofs/binds-editor/binds"
	"github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	keys "github.com/araujofs/binds-editor/help"
	msgs "github.com/araujofs/binds-editor/messages"
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
	emptyDirectory bool
	msgs.InfoModel
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

	m := FileSearch{
		filePicker:     fp,
		help:           h,
		config:         configuration,
		emptyDirectory: false,
		InfoModel:      msgs.GetDefaultInfoModel(),
	}

	if consts.WindowSize.Height != 0 {
		m.setFilePickerHeight()
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
		m.setFilePickerHeight()
		return m, nil

	case msgs.ErrorMsg:
		m.Message = nil
		m.Error = &msg.Message

	case msgs.MessageMsg:
		m.Error = nil
		m.Message = &msg.Message

	case tea.KeyMsg:
		m.Message = nil
		m.Error = nil

		switch {
		case key.Matches(msg, keys.CoreKeys.Close):
			return m, tea.Quit

		case key.Matches(msg, keys.FilePickerKeys.Select):
			m.filePicker, cmd = m.filePicker.Update(msg)

			if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
				if _, err := binds.ReadBindsFile(path); err != nil {
					return m, msgs.SendErrorMsg(err.Error())
				}

				selection := InitFileSelection(&path, m.config)

				return selection, nil
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

			return selection, nil

		case key.Matches(msg, keys.CoreKeys.Help):
			return m, nil
		}
	}

	m.filePicker, cmd = m.filePicker.Update(msg)

	return m, cmd

}

func (m FileSearch) View() string {
	title := "Binds Editor | Search File"

	title += m.GetStyledMessage()

	if m.emptyDirectory {
		return consts.FullScreenFileSearchStyle.Render(title + m.filePicker.View() + lipgloss.NewStyle().MarginTop(1).Render(m.help.FullHelpView(keys.FilePickerKeys.FullHelp())))
	}

	return consts.FullScreenFileSearchStyle.Render((title + m.filePicker.View() + m.help.FullHelpView(keys.FilePickerKeys.FullHelp())))
}

func (m *FileSearch) setFilePickerHeight() {
	top, _, bottom, _ := consts.FullScreenStyle.GetMargin()
	windowHeight := consts.WindowSize.Height

	helpHeight := len(keys.FilePickerKeys.FullHelp()[0])
	m.filePicker.SetHeight(windowHeight - top - bottom - 2 - helpHeight)
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
