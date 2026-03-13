package help

import (
	"github.com/charmbracelet/bubbles/key"
)

type TableKeyMap struct {
	CommonKeyMap
	Create, Delete, Edit, Unbind, Comment, Details, GoBack key.Binding
}

func (k TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Create, k.Delete, k.Edit, k.Unbind, k.Close, k.GoBack}
}

func (k TableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Create}, {k.Delete, k.Edit, k.Unbind}, {k.Comment, k.Details}, {k.Close, k.GoBack}}
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
	key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "toggle comment"),
	),
	key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "bind details"),
	),
	key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go to previous page"),
	),
}

type ReadonlyTableKeyMap struct {
	CommonKeyMap
	Unbind, Details, GoBack key.Binding
}

func (k ReadonlyTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Unbind, k.Close, k.GoBack}
}

func (k ReadonlyTableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Unbind}, {k.Details, k.Close, k.GoBack}}
}

var ReadonlyTableKeys = ReadonlyTableKeyMap{
	CommonKeyMap: CommonKeys,
	Unbind: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unbind"),
	),
	Details: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "bind details"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go to previous page"),
	),
}
