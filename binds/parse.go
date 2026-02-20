package binds

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

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
		line := scanner.Text()
		bind, err := parseBind(path, line, lineNumber)
		if err != nil {
			return nil, err
		}

		if bind != nil {
			bindings = append(bindings, bind)
		}

		lineNumber++
	}

	return bindings, nil
}

func parseBind(path string, rawLine string, bindLineNumber int) (*Bind, error) {
	if len(rawLine) == 0 {
		return nil, nil
	}

	rawLine = strings.Trim(rawLine, " ")

	errorMsg := fmt.Errorf("invalid bind format on line %d! raw line: (%s)", bindLineNumber, rawLine)
	commented := false

	if strings.HasPrefix(rawLine, "#") {
		commented = true
		rawLine = strings.TrimPrefix(rawLine, "#")
	}

	bindDefinition, bindContent, found := strings.Cut(rawLine, "=")
	if !found && !commented {
		return nil, errorMsg
	}

	if !found && commented {
		return nil, nil
	}

	if !commented && !(strings.Contains(bindDefinition, "bind")) {
		return nil, errorMsg
	}

	bindDefinition = strings.TrimSpace(bindDefinition)
	bindContent = strings.TrimSpace(bindContent)
	bindFlags := []*Flag{}
	bindDefinitionLen := len(bindDefinition)
	bind := Bind{
		BindCore: BindCore{
			Flags: []*Flag{},
			Type:  Normal,
		},
		LineNumber: bindLineNumber,
		RawLine:    rawLine,
		FilePath:   path,
	}

	if strings.HasPrefix(bindDefinition, "un") {
		bindCore, err := parseUnbind(bindContent)
		if err != nil {
			return nil, errorMsg
		}

		bind.BindCore = *bindCore
		bind.Type = Unbind

		return &bind, nil
	}

	if bindDefinitionLen != 4 {
		bindFlagsStrings := strings.Split(bindDefinition[4:], "")
		if slices.Contains(bindFlagsStrings, "s") {
			// return nil, fmt.Errorf("for now binds with the 's' flag are not supported! Line: %d, Raw line: %s", bindLineNumber, rawLine)
			return nil, nil
		}

		for _, bindFlag := range bindFlagsStrings {
			if flag, ok := DefaultFlags[bindFlag]; ok {
				bindFlags = append(bindFlags, flag)
			}
		}
	}

	if commented {
		bindCore, err := parseBindContent(bindContent)

		if err != nil {
			return nil, nil
		}

		bind.BindCore = *bindCore
		bind.Type = Comment
		bind.Flags = bindFlags

		return &bind, nil
	}

	bindCore, err := parseBindContent(bindContent)
	if err != nil {
		return nil, errorMsg
	}

	bind.BindCore = *bindCore
	bind.Flags = bindFlags

	return &bind, nil
}

func parseUnbind(unbind string) (*BindCore, error) {
	parts := strings.Split(unbind, ",")
	partsLen := len(parts)

	if partsLen != 2 {
		return nil, fmt.Errorf("invalid bind format")
	}

	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}

	modKeys := parseModKeys(parts[0])
	key := parts[1]

	return &BindCore{
		Shortcut: Shortcut{
			ModKeys: modKeys,
			Key:     key,
		},
		Dispatcher:  "",
		Action:      "",
		Description: "",
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
	dispatcher := parts[2]
	var description string

	if partsLen == 5 {
		dispatcher = parts[3]
		action = parts[4]
		description = parts[2]
	}

	return &BindCore{
		Shortcut: Shortcut{
			ModKeys: modKeys,
			Key:     key,
		},
		Dispatcher:  dispatcher,
		Action:      action,
		Description: description,
		Type:        Normal,
	}, nil
}

func parseModKeys(modKeys string) []string {
	if modKeys == "" {
		return nil
	}

	s := strings.ToUpper(modKeys)
	separatedModKeys := []string{}

	// it needs to be like that because hyprland doesnt define an especific separator
	for _, modKey := range ModKeys {
		if strings.Contains(s, modKey) {
			separatedModKeys = append(separatedModKeys, modKey)
		}
	}

	return separatedModKeys
}
