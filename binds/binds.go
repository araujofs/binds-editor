package binds

type BindType int

const (
	normal BindType = iota + 1
	unbind
	comment
)

type Shortcut struct {
	ModKeys []string
	Key     string
}

type Bind struct {
	BindCore
	LineNumber int
	RawLine    string
	Flags      []string
	Type       BindType
}

type BindCore struct {
	Shortcut    Shortcut
	ActionType  string
	Action      string
	Description string
}

func (b *Bind) Edit(newBind Bind) {
	if flags := newBind.Flags; flags != nil {
		b.Flags = flags
	}
}
