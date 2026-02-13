package help

import (
	"github.com/charmbracelet/bubbles/key"
)

type FileSelectionKeyMap struct {
	CommonKeyMap
	Add, Edit, Delete, Save key.Binding
}

var FileSelectionKeys = FileSelectionKeyMap{
	CommonKeys,
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add file"),
	),
	key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit file"),
	),
	key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete file"),
	),
	key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save configuration list"),
	),
}

func (k FileSelectionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Add, k.Edit, k.Delete, k.Save, k.Close}
}

func (k FileSelectionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Save}, {k.Edit, k.Delete, k.Add}, {k.Close}}
}

type FileSelectionInputKeyMap struct {
	CoreKeyMap
	Enter, Back key.Binding
}

var FileSelectionInputKeys = FileSelectionInputKeyMap{
	CoreKeys,
	key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

func (k FileSelectionInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Help, k.Close}
}

func (k FileSelectionInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Enter, k.Back}, {k.Help, k.Close}}
}
