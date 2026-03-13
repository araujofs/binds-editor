package ui

import (
	"fmt"
	"os"
	"strings"

	config "github.com/araujofs/binds-editor/configuration"
	consts "github.com/araujofs/binds-editor/constants"
	msgs "github.com/araujofs/binds-editor/messages"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type mode int

const (
	create mode = iota
	edit
)

type FileManipulation struct {
	form         *huh.Form
	originalFile *config.File
	newFile      *config.File
	mode         mode
	formKeys     *huh.KeyMap
	*consts.GlobalState
	msgs.InfoModel
}

func InitFileManipulation(originalFile *config.File, globalState *consts.GlobalState) (tea.Model, tea.Cmd) {
	var mode mode

	var newFile config.File
	if originalFile != nil {
		newFile = *originalFile
		mode = edit
	}

	model := &FileManipulation{
		originalFile: originalFile,
		newFile:      &newFile,
		mode:         mode,
		formKeys:     huh.NewDefaultKeyMap(),
		GlobalState:  globalState,
	}

	model.createForm()

	return model, model.form.Init()
}

func (m FileManipulation) Init() tea.Cmd {
	return nil
}

func (m FileManipulation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		consts.WindowSize = msg

		top, _, bottom, _ := consts.FullScreenStyle.GetMargin()

		m.form = m.form.WithHeight(msg.Height - top - bottom - 3)

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
		case key.Matches(msg, m.formKeys.Quit):
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State != huh.StateCompleted {
		return m, tea.Batch(cmds...)
	}

	if m.mode == edit {
		err := m.Configuration.EditFile(m.originalFile.Name, *m.newFile)
		if err != nil {
			cmds = append(cmds, msgs.SendErrorMsg(err.Error()))
			return InitFileSelection(m.GlobalState), tea.Batch(cmds...)
		}

		cmds = append(cmds, msgs.SendMessageMsg(fmt.Sprintf("file \"%s\" edited", m.originalFile.Name)))
		return InitFileSelection(m.GlobalState), tea.Batch(cmds...)
	}

	err := m.Configuration.AddFile(m.newFile)
	if err != nil {
		cmds = append(cmds, msgs.SendErrorMsg(err.Error()))
		return InitFileSelection(m.GlobalState), tea.Batch(cmds...)
	}

	cmds = append(cmds, msgs.SendMessageMsg(fmt.Sprintf("file \"%s\" created", m.newFile.Name)))
	return InitFileSelection(m.GlobalState), tea.Batch(cmds...)
}

func (m FileManipulation) View() string {
	title := "Binds Editor | Edit Configuration File"
	title += m.GetStyledMessage()

	return consts.FullScreenStyle.Render(title + m.form.View())
}

func (fe *FileManipulation) createForm() {
	top, _, bottom, _ := consts.FullScreenStyle.GetMargin()

	dir := "."

	if homeDir, err := os.UserHomeDir(); err == nil {
		dir = homeDir
	}

	fe.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("File name").Value(&fe.newFile.Name).CharLimit(30),
			huh.NewFilePicker().Validate(isPathEmpty).Title("File path").Value(&fe.newFile.Path).Description("Change the path of the configuration file").AllowedTypes([]string{".conf"}).CurrentDirectory(dir).ShowHidden(true),
		),
		huh.NewGroup(
			huh.NewConfirm().Title("The file should be readonly?").Affirmative("Yes!").Negative("No!").Value(&fe.newFile.Readonly),
		),
		huh.NewGroup(
			huh.NewFilePicker().Validate(isPathEmpty).Title("File output").Value(&fe.newFile.Output).Description("Change the path of the output file").AllowedTypes([]string{".conf"}).CurrentDirectory(dir).ShowHidden(true),
		).WithHideFunc(func() bool {
			return !fe.newFile.Readonly
		}),
	).WithHeight(consts.WindowSize.Height - top - bottom - 3)
}

func isPathEmpty(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path can't be empty")
	}

	return nil
}
