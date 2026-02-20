package ui

import (
	"slices"

	"github.com/araujofs/binds-editor/binds"
	consts "github.com/araujofs/binds-editor/constants"
	msgs "github.com/araujofs/binds-editor/messages"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type Create struct {
	form         *huh.Form
	originalBind *binds.Bind
	newBind      *binds.Bind
	formKeys     *huh.KeyMap
	*consts.GlobalState
	msgs.InfoModel
}

func InitCreate(globalState *consts.GlobalState) (tea.Model, tea.Cmd) {
	model := &Create{
		newBind: &binds.Bind{
			FilePath: globalState.SelectedFile.Path,
			BindCore: binds.BindCore{
				Type: binds.Normal,
			},
		},
		formKeys:    huh.NewDefaultKeyMap(),
		GlobalState: globalState,
	}

	model.createForm()

	return model, model.form.Init()
}

func (m Create) Init() tea.Cmd {
	return nil
}

func (m Create) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	if m.form.State == huh.StateCompleted {
		err := m.newBind.AppendToFile()
		if err != nil {
			cmds = append(cmds, msgs.SendErrorMsg(err.Error()))

			m.createForm()

			return m, tea.Batch(cmds...)
		}

		model, cmd := InitTable(m.GlobalState)
		cmds = append(cmds, cmd, msgs.SendMessageMsg("bind created"))

		return model, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m Create) View() string {
	title := "Binds Editor | Create Bind"
	title += m.GetStyledMessage()

	return consts.FullScreenStyle.Render(title + m.form.View())
}

func (c *Create) createForm() {
	top, _, bottom, _ := consts.FullScreenStyle.GetMargin()

	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[*binds.Flag]().Title("Flags").Value(&c.newBind.Flags).Options(
				getFlagOptions()...,
			),
			huh.NewMultiSelect[string]().Title("Mod keys").Value(&c.newBind.Shortcut.ModKeys).Options(
				getModKeysOptions()...,
			),
			huh.NewInput().Title("Main key").Value(&c.newBind.Shortcut.Key).CharLimit(1),
		),
		huh.NewGroup(
			huh.NewInput().Title("Dispatcher").Value(&c.newBind.Dispatcher),
			huh.NewInput().Title("Action").Value(&c.newBind.Action),
		),
		huh.NewGroup(
			huh.NewInput().Title("Description").Value(&c.newBind.Description),
		).WithHideFunc(
			func() bool {
				return !slices.Contains(c.newBind.Flags, binds.DefaultFlags["d"])
			},
		),
	).WithHeight(consts.WindowSize.Height - top - bottom - 3)
}
