package help

import (
	"github.com/charmbracelet/bubbles/key"
)

type FileSelectionKeyMap struct {
	CommonKeyMap
	Find, Save key.Binding
}

var FileSelectionKeys = FileSelectionKeyMap{
	CommonKeys,
	key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "find binding file"),
	),
	key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save configuration list"),
	),
}

func (k FileSelectionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Find, k.Close, k.Save}
}

func (k FileSelectionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Find, k.Close}}
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
