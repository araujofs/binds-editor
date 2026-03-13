package binds

import (
	"fmt"
	"slices"
	"strings"
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

type BindCore struct {
	Shortcut
	Dispatcher  string
	Action      string
	Description string
	Flags       []*Flag
}

type Bind struct {
	BindCore
	GeneralInfo
}

func NewBind(bindCore BindCore, lineNumber int, rawLine string, file *File, commented bool) (*Bind, error) {
	bind := Bind{
		BindCore: bindCore,
		GeneralInfo: GeneralInfo{
			LineNumber: lineNumber,
			RawLine:    rawLine,
			File:       file,
			Commented:  commented,
		},
	}

	isValid, missingProperties := bind.IsValid()
	if !isValid {
		return nil, missingPropertiesToError(missingProperties)
	}

	return &bind, nil
}

func (b *Bind) GetLine() string {
	var bind strings.Builder
	var bindDefinition strings.Builder
	var bindContent strings.Builder

	if b.Commented {
		bindDefinition.WriteString("# ")
	}

	bindDefinition.WriteString("bind")

	for _, flag := range b.Flags {
		bindDefinition.WriteString(flag.Flag)
	}

	bindContent.WriteString(strings.Join(b.Shortcut.ModKeys, " "))
	bindContent.WriteString(", " + b.Shortcut.Key)

	if slices.Contains(b.Flags, DefaultFlags["d"]) {
		bindContent.WriteString(", " + b.Description)
	}

	bindContent.WriteString(", " + b.Dispatcher)
	bindContent.WriteString(", " + b.Action)

	bind.WriteString(bindDefinition.String())
	bind.WriteString(" = " + bindContent.String())

	line := bind.String()

	return line
}

func (b *Bind) GetLineNumber() int {
	return b.LineNumber
}

func (b *Bind) ToggleComment() {
	b.Commented = !b.Commented
}

func (b *Bind) Unbind() error {
	var output *File

	output = b.File
	if b.File.IsReadonly() {
		output = b.File.Output
	}

	unbind, err := NewUnbind(b.Shortcut, b.GetLineNumber(), b.RawLine, output, false)
	if err != nil {
		return err
	}

	output.AddLine(unbind)
	return nil
}

func (b Bind) GetRow() []string {
	var shortcut = ""

	if joinedModKeys := strings.Join(b.Shortcut.ModKeys, "+"); joinedModKeys != "" {
		shortcut = strings.Join([]string{joinedModKeys, b.Shortcut.Key}, "+")
	} else {
		shortcut = b.Shortcut.Key
	}

	bindTypeString := "Normal"

	if b.Commented {
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

func (bc BindCore) isValid() (bool, []string) {
	_, missingProperties := bc.Shortcut.isValid()

	if bc.Dispatcher == "" {
		missingProperties = append(missingProperties, "dispatcher")
	}

	if bc.Action == "" {
		missingProperties = append(missingProperties, "action")
	}

	if flags := bc.Flags; len(flags) > 0 && slices.Contains(flags, DefaultFlags["d"]) && bc.Description == "" {
		missingProperties = append(missingProperties, "description")
	}

	if len(missingProperties) != 0 {
		return false, missingProperties
	}

	return true, missingProperties
}

func (b Bind) IsValid() (isValid bool, missingProperties []string) {
	_, missingProperties = b.BindCore.isValid()
	isGeneralInfoValid, giMissingProperties := b.GeneralInfo.isValid()

	if !isGeneralInfoValid {
		missingProperties = append(missingProperties, giMissingProperties...)
	}

	if len(missingProperties) != 0 {
		return
	}

	isValid = true
	return
}

func missingPropertiesToError(missingProperties []string) error {
	missing := strings.Join(missingProperties, ", ")
	err := fmt.Errorf("invalid bind, missing: %s", missing)

	return err
}
