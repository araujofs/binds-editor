package help

import (
	"github.com/charmbracelet/bubbles/key"
)

type TableKeyMap struct {
	CommonKeyMap
	Create, Delete, Edit, Unbind key.Binding
}

func (k TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Create, k.Delete, k.Edit, k.Unbind, k.Close}
}

func (k TableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Create, k.Delete, k.Edit, k.Unbind, k.Close}}
}


var TableKeys = TableKeyMap{
	CommonKeys,
	key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create bind"),
	),
	key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete bind"),
	),
	key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit bind"),
	),
	key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unbind"),
	),
}

