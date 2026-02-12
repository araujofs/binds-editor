package help

import "github.com/charmbracelet/bubbles/key"

type CoreKeyMap struct {
	Close, Help key.Binding
}

type CommonKeyMap struct {
	CoreKeyMap
	Up, Down key.Binding
}

var CoreKeys = CoreKeyMap{
	key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

var CommonKeys = CommonKeyMap{
	CoreKeys,
	key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
}
