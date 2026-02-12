package help

import (
	"github.com/charmbracelet/bubbles/key"
)

type FileSelectionKeyMap struct {
	CommonKeyMap
	Find key.Binding
}

func (k FileSelectionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Find, k.Close}
}

func (k FileSelectionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Find, k.Close}}
}

var FileSelectionKeys = FileSelectionKeyMap{
	CommonKeys,
	key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "find binding file"),
	),
}
