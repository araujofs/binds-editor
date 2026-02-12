package ui

import (
	config "github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	keys "github.com/araujofs/binds-editor/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type FileSelection struct {
	config *config.Configuration
	list   list.Model
}

func InitFileSelection() *FileSelection {
	config := config.GetConfigData()

	items := filesToItems(config.Files)

	l := list.New(items, list.NewDefaultDelegate(), 8, 8)
	l.SetShowTitle(false)
	l.AdditionalShortHelpKeys = keys.FileSelectionKeys.ShortHelp

	m := &FileSelection{
		config,
		l,
	}

	if consts.WindowSize.Height != 0 {
		setListSize(&l, consts.WindowSize)
	}

	return m
}

func (m FileSelection) Init() tea.Cmd {
	return nil
}

func (m FileSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		consts.WindowSize = msg
		setListSize(&m.list, msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.CoreKeys.Close):
			return m, tea.Quit
		case key.Matches(msg, keys.FileSelectionKeys.Find):
			return InitFileSearch()
		}

		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m FileSelection) View() string {
	title := "Binds Editor | File Selection\n"

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

func setListSize(list *list.Model, msg tea.WindowSizeMsg) {
	top, right, bottom, left := consts.FullScreenStyle.GetMargin()
	list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
}
