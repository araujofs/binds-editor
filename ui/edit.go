package ui

import (
	"github.com/araujofs/binds-editor/binds"
	consts "github.com/araujofs/binds-editor/constants"
	msgs "github.com/araujofs/binds-editor/messages"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type Edit struct {
	form         *huh.Form
	originalBind *binds.Bind
	newBind      binds.Bind
	formKeys     *huh.KeyMap
	confirmed    bool
	msgs.InfoModel
}

func InitEdit(bind *binds.Bind) (tea.Model, tea.Cmd) {
	model := &Edit{
		originalBind: bind,
		newBind:      *bind,
		formKeys:     huh.NewDefaultKeyMap(),
		confirmed:    false,
	}

	model.createForm()

	return model, model.form.Init()
}
func (m Edit) Init() tea.Cmd {
	return nil
}

func (m Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	if m.form.State == huh.StateCompleted && m.confirmed {
		cmds = append(cmds, msgs.SendMessageMsg("Nao confirmado"))
	}

	if m.form.State == huh.StateCompleted && !m.confirmed {
		m.createForm()
		cmds = append(cmds, m.form.Init())
	}

	return m, tea.Batch(cmds...)
}

func (m Edit) View() string {
	title := "Binds Editor | Select File"
	title += m.GetStyledMessage()

	return consts.FullScreenStyle.Render(title + m.form.View())
}

func getFlagOptions() []huh.Option[*binds.Flag] {
	var options []huh.Option[*binds.Flag]

	for _, flag := range binds.DefaultFlags {
		options = append(options, huh.NewOption(
			flag.Flag+" - "+flag.Name, flag,
		))
	}

	return options
}

func getModKeysOptions() []huh.Option[string] {
	var options []huh.Option[string]

	for _, key := range binds.ModKeys {
		options = append(options, huh.NewOption(
			key, key,
		))
	}

	return options
}

func (e *Edit) createForm() {
	top, _, bottom, _ := consts.FullScreenStyle.GetMargin()

	e.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[*binds.Flag]().Title("Flags").Value(&e.newBind.Flags).Options(
				getFlagOptions()...,
			),
			huh.NewInput().Title("Dispatcher").Value(&e.newBind.Dispatcher),
			huh.NewInput().Title("Action").Value(&e.newBind.Action),
			huh.NewInput().Title("Description").Value(&e.originalBind.Description),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("Mod keys").Value(&e.newBind.Shortcut.ModKeys).Options(
				getModKeysOptions()...,
			),
			huh.NewInput().Title("Main key").Value(&e.newBind.Shortcut.Key),
		),
		huh.NewGroup(
			huh.NewConfirm().Key("confirm").Title("Confirm edit?").Affirmative("Yes!").Negative("No.").Value(&e.confirmed),
		),
	).WithHeight(consts.WindowSize.Height - top - bottom - 3)
}
