package help

import (
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
)

type FilePickerKeyMap struct {
	CommonKeyMap
	Open, Back, GoToTop, GoToLast /* PageUp, PageDown,  */, Select, Close, Help key.Binding
}

func (k FilePickerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Open, k.Back, k.Select, k.Close, k.Help}
}

func (k FilePickerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Open}, {k.Back, k.GoToTop, k.GoToLast} /*{  k.PageUp, k.PageDown}, */, {k.Select, k.Close, k.Help}}
}

var FilePickerKeys = FilePickerKeyMap{
	CommonKeyMap: CommonKeys,
	GoToTop:      key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "first")),
	GoToLast:     key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "last")),
	// PageUp:       key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("K/pgup", "page up")),
	// PageDown:     key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("J/pgdown", "page down")),
	Back:   key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("←/h", "back")),
	Open:   key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("→/l", "open")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
}

func GetFilepickerKeyMap() filepicker.KeyMap {
	return filepicker.KeyMap{
		GoToTop:  FilePickerKeys.GoToTop,
		GoToLast: FilePickerKeys.GoToLast,
		Down:     FilePickerKeys.Down,
		Up:       FilePickerKeys.Up,
		// PageUp:   FilePickerKeys.PageUp,
		// PageDown: FilePickerKeys.PageDown,
		Back:   FilePickerKeys.Back,
		Open:   FilePickerKeys.Open,
		Select: FilePickerKeys.Select,
	}
}
