package help

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Create key.Binding
	Delete key.Binding
	Edit   key.Binding
	Unbind key.Binding
	Close  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Create, k.Delete, k.Edit, k.Unbind, k.Close}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Create, k.Delete, k.Edit, k.Unbind, k.Close}}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create bind"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete bind"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit bind"),
	),
	Unbind: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unbind"),
	),
	Close: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
