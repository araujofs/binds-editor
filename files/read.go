package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var variables = map[string]string{}

type Shortcut struct {
	ModKeys []string
	Key     string
}

type Bind struct {
	BindCore
	LineNumber int
	RawLine    string
	Flags      []string
}

type BindCore struct {
	Shortcut    Shortcut
	ActionType  string
	Action      string
	Description string
}

func ReadBindsFile(path string) ([]*Bind, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bindings := make([]*Bind, 0, 50)

	scanner := bufio.NewScanner(f)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()
		bind, err := parseBind(line, lineNumber)

		if err == nil && bind != nil {
			bindings = append(bindings, bind)
		}
	}

	return bindings, nil
}

func parseBind(rawLine string, bindLineNumber int) (*Bind, error) {
	if len(rawLine) == 0 || rawLine[0] == '#' {
		return nil, nil
	}

	errorMsg := fmt.Errorf("invalid bind format on line %d! raw line: (%s)", bindLineNumber, rawLine)

	bindType, bindContent, found := strings.Cut(strings.TrimSpace(rawLine), "=")

	if !found {
		return nil, errorMsg
	}

	bindContent = strings.TrimSpace(bindContent)
	bindType = strings.TrimSpace(bindType)
	bindTypeLen := len(bindType)

	if !(strings.Contains(bindType, "bind")) {
		return nil, errorMsg
	}

	bindFlags := ""

	if bindTypeLen != 4 {
		bindFlags = bindType[4:]
	}

	if strings.Contains(bindFlags, "s") {
		return nil, fmt.Errorf("for now binds with the 's' flag are not parsed! Line: %d, Raw line: %s", bindLineNumber, rawLine)
	}

	bind, err := parseBindContent(bindContent)

	if err != nil {
		return nil, err
	}

	return &Bind{
		BindCore:   *bind,
		LineNumber: bindLineNumber,
		RawLine:    rawLine,
		Flags:      strings.Split(bindFlags, ""),
	}, nil
}

func parseBindContent(bindContent string) (*BindCore, error) {
	parts := strings.Split(bindContent, ",")
	partsLen := len(parts)

	if partsLen > 5 || partsLen < 4 {
		return nil, fmt.Errorf("invalid bind format")
	}

	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}

	modKeys := parseModKeys(parts[0])
	key := parts[1]
	action := parts[3]
	actionType := parts[2]
	var description string

	if partsLen == 5 {
		actionType = parts[3]
		action = parts[4]
		description = parts[2]
	}

	return &BindCore{
		Shortcut: Shortcut{
			ModKeys: modKeys,
			Key:     key,
		},
		ActionType:  actionType,
		Action:      action,
		Description: description,
	}, nil

}

func parseModKeys(modKeys string) []string {
	if modKeys == "" {
		return nil
	}

	s := strings.ToUpper(modKeys)
	separatedModKeys := []string{}

	// it needs to be like that because hyprland doesnt define an especific separator
	if strings.Contains(s, "SHIFT") {
		separatedModKeys = append(separatedModKeys, "SHIFT")
	}

	if strings.Contains(s, "CAPS") {
		separatedModKeys = append(separatedModKeys, "CAPS")
	}

	if strings.Contains(s, "CTRL") {
		separatedModKeys = append(separatedModKeys, "CTRL")
	}

	if strings.Contains(s, "CONTROL") {
		separatedModKeys = append(separatedModKeys, "CONTROL")
	}

	if strings.Contains(s, "ALT") {
		separatedModKeys = append(separatedModKeys, "ALT")
	}

	if strings.Contains(s, "MOD2") {
		separatedModKeys = append(separatedModKeys, "MOD2")
	}

	if strings.Contains(s, "MOD3") {
		separatedModKeys = append(separatedModKeys, "MOD3")
	}

	if strings.Contains(s, "SUPER") {
		separatedModKeys = append(separatedModKeys, "SUPER")
	}

	if strings.Contains(s, "WIN") {
		separatedModKeys = append(separatedModKeys, "WIN")
	}

	if strings.Contains(s, "LOGO") {
		separatedModKeys = append(separatedModKeys, "LOGO")
	}

	if strings.Contains(s, "MOD4") {
		separatedModKeys = append(separatedModKeys, "MOD4")
	}

	if strings.Contains(s, "MOD5") {
		separatedModKeys = append(separatedModKeys, "MOD5")
	}

	return separatedModKeys
}

func (b Bind) KeybindToRow() []string {
	var shortcut = ""

	if joinedModKeys := strings.Join(b.Shortcut.ModKeys, "+"); joinedModKeys != "" {
		shortcut = strings.Join([]string{joinedModKeys, b.Shortcut.Key}, "+")
	} else {
		shortcut = b.Shortcut.Key
	}

	return []string{
		shortcut,
		b.Description,
		b.ActionType,
		strings.Join(b.Flags, ", "),
	}
}
