package binds

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

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

func (b Bind) ReplaceInFile(path string) error {
	// ler o arquivo
	content, err := os.ReadFile(path)
	if err != nil && err != io.EOF {
		return err
	}

	// separar por linhas
	lines := strings.Split(string(content), "\n")

	// transformar bind em linha
	newBindLine, err := b.KeybindToLine()
	if err != nil {
		return err
	}

	// substituir linha antiga por nova linha
	lines[b.LineNumber] = *newBindLine

	// transformar linhas em bytes
	content = []byte(strings.Join(lines, "\n"))

	// sobrescrever arquivo
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func (b Bind) KeybindToRow() []string {
	var shortcut = ""

	if joinedModKeys := strings.Join(b.Shortcut.ModKeys, "+"); joinedModKeys != "" {
		shortcut = strings.Join([]string{joinedModKeys, b.Shortcut.Key}, "+")
	} else {
		shortcut = b.Shortcut.Key
	}

	var bindTypeString string
	switch b.Type {
	case Normal:
		bindTypeString = "Normal"

	case Unbind:
		bindTypeString = "Unbind"

	case Comment:
		bindTypeString = "Commented"
	}

	var flags []string
	for _, flag := range b.Flags {
		flags = append(flags, flag.Flag)
	}

	return []string{
		shortcut,
		bindTypeString,
		b.Description,
		b.Dispatcher,
		strings.Join(flags, ", "),
	}
}

func (b Bind) KeybindToLine() (*string, error) {
	if valid, err := b.isBindValid(); !valid {
		return nil, err
	}

	var bind strings.Builder
	var bindDefinition strings.Builder
	var bindContent strings.Builder

	switch b.Type {
	case Normal:
		bindDefinition.WriteString("bind")
	case Comment:
		bindDefinition.WriteString("# bind")
	case Unbind:
		bindDefinition.WriteString("unbind")
	}

	if b.Type != Unbind {
		for _, flag := range b.Flags {
			bindDefinition.WriteString(flag.Flag)
		}
	}

	bindContent.WriteString(strings.Join(b.Shortcut.ModKeys, " "))
	bindContent.WriteString(", " + b.Shortcut.Key)

	if b.Type != Unbind {
		if slices.Contains(b.Flags, DefaultFlags["d"]) {
			bindContent.WriteString(", " + b.Description)
		}

		bindContent.WriteString(", " + b.Dispatcher)
		bindContent.WriteString(", " + b.Action)
	}

	bind.WriteString(bindDefinition.String())
	bind.WriteString(" = " + bindContent.String())

	bindLine := bind.String()

	return &bindLine, nil
}

func (b Bind) isBindValid() (bool, error) {
	missingProperties := []string{}

	if b.Dispatcher == "" {
		missingProperties = append(missingProperties, "dispatcher")
	}

	if b.Action == "" {
		missingProperties = append(missingProperties, "action")
	}

	if b.Type == 0 {
		missingProperties = append(missingProperties, "type")
	}

	if b.Shortcut.Key == "" {
		missingProperties = append(missingProperties, "key")
	}

	if modKeys := b.Shortcut.ModKeys; modKeys == nil || len(modKeys) == 0 {
		missingProperties = append(missingProperties, "mod keys")
	}

	if flags := b.Flags; flags != nil && len(flags) > 0 && slices.Contains(flags, DefaultFlags["d"]) && b.Description == "" {
		missingProperties = append(missingProperties, "description")
	}

	if b.LineNumber == 0 {
		missingProperties = append(missingProperties, "line number")
	}

	if b.RawLine == "" {
		missingProperties = append(missingProperties, "raw line")
	}

	if len(missingProperties) != 0 {
		missing := strings.Join(missingProperties, ", ")
		errorMsg := fmt.Errorf("invalid bind, missing: %s", missing)

		return false, errorMsg
	}

	return true, nil
}
