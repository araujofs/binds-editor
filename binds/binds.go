package binds

type BindType int

const (
	Normal BindType = iota + 1
	Unbind
	Comment
)

var ModKeys = []string{
	"SHIFT",
	"CAPS",
	"CTRL",
	"CONTROL",
	"ALT",
	"MOD2",
	"MOD3",
	"SUPER",
	"WIN",
	"LOGO",
	"MOD4",
	"MOD5",
}

type Flag struct {
	Flag        string
	Name        string
	Description string
}

type Flags map[string]*Flag

var DefaultFlags = Flags{
	"l": &Flag{Flag: "l", Name: "locked", Description: "Will also work when an input inhibitor (e.g. a lockscreen) is active."},
	"r": &Flag{Flag: "r", Name: "release", Description: "Will trigger on release of a key."},
	"c": &Flag{Flag: "c", Name: "click", Description: "Will trigger on release of a key or button as long as the mouse cursor stays inside binds:drag_threshold."},
	"g": &Flag{Flag: "g", Name: "drag", Description: "Will trigger on release of a key or button as long as the mouse cursor moves outside binds:drag_threshold."},
	"o": &Flag{Flag: "o", Name: "long press", Description: "Will trigger on long press of a key."},
	"e": &Flag{Flag: "e", Name: "repeat", Description: "Will repeat when held."},
	"n": &Flag{Flag: "n", Name: "non", Description: "consuming	Key/mouse events will be passed to the active window in addition to triggering the dispatcher."},
	"m": &Flag{Flag: "m", Name: "mouse", Description: "See the dedicated Mouse Binds section."},
	"t": &Flag{Flag: "t", Name: "transparent", Description: "Cannot be shadowed by other binds."},
	"i": &Flag{Flag: "i", Name: "ignore mods", Description: "Will ignore modifiers."},
	"d": &Flag{Flag: "d", Name: "has Description", Description: "Will allow you to write a description for your bind."},
	"p": &Flag{Flag: "p", Name: "bypass", Description: "Bypasses the app’s requests to inhibit keybinds."},
	"u": &Flag{Flag: "u", Name: "submap universal", Description: "Will be active no matter the submap."},
}

type Shortcut struct {
	ModKeys []string
	Key     string
}

type BindCore struct {
	Shortcut    Shortcut
	Dispatcher  string
	Action      string
	Description string
	Flags       []*Flag
	Type        BindType
}

type Bind struct {
	BindCore
	LineNumber int
	RawLine    string
}

func (b *Bind) Edit(newBind Bind) {
	if flags := newBind.Flags; flags != nil {
		b.Flags = flags
	}
}
