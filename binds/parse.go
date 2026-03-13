package binds

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/araujofs/binds-editor/errors"
)

// TODO: pass file to parsing functions

func (fi *File) Parse() ([]Line, error) {
	f, err := os.Open(fi.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]Line, 0, 50)
	scanner := bufio.NewScanner(f)

	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		newLine, err := parseLine(fi, line, lineNumber)
		if err != nil {
			return nil, err
		}

		lines = append(lines, newLine)

		lineNumber++
	}

	return lines, nil
}

func parseLine(file *File, line string, lineNumber int) (Line, error) {
	err := fmt.Errorf("invalid bind format on line %d! raw line: (%s)", lineNumber, line)

	if len(line) == 0 {
		return &RawLine{}, nil
	}

	if strings.HasPrefix(line, "#") {
		parsedLine, err := parseComment(file, line, lineNumber)
		if err != nil {
			return nil, err
		}

		return parsedLine, nil
	}

	if strings.HasPrefix(line, "unbind") {
		// never return nil, nil inside parseBind
		return parseUnbind(file, line, lineNumber)
	}

	if strings.HasPrefix(line, "bind") {
		bind, err := parseBind(file, line, lineNumber)
		if errors.IsUnsupportedBindFlagError(err) {
			return &RawLine{
				Content:    &line,
				LineNumber: lineNumber,
			}, nil
		}

		return bind, err
	}

	return nil, err

}

func parseComment(file *File, line string, lineNumber int) (Line, error) {
	err := fmt.Errorf("invalid bind format on line %d! raw line: (%s)", lineNumber, line)
	trimmedLine := strings.TrimSpace(strings.TrimPrefix(line, "#"))

	if strings.HasPrefix(trimmedLine, "unbind") {
		unbind, err := parseUnbind(file, line, lineNumber)
		if err != nil {
			return &RawLine{
				Content:    &line,
				LineNumber: lineNumber,
			}, nil
		}

		unbind.Commented = true
		return unbind, nil
	}

	if strings.HasPrefix(trimmedLine, "bind") {
		bind, err := parseBind(file, line, lineNumber)
		if err != nil {
			return &RawLine{
				Content:    &line,
				LineNumber: lineNumber,
			}, nil
		}

		bind.Commented = true
		return bind, nil
	}

	return nil, err
}

func parseBind(file *File, line string, lineNumber int) (*Bind, error) {
	errMsg := fmt.Errorf("invalid bind format on line %d! raw line: (%s)", lineNumber, line)
	bindDefinition, bindContent, found := strings.Cut(line, "=")

	if !found {
		return nil, errMsg
	}

	var hasDescription bool
	var flags []*Flag

	bindDefinition, bindContent = strings.TrimSpace(bindDefinition), strings.TrimSpace(bindContent)
	bindDefinition = strings.TrimSpace(strings.TrimPrefix(bindDefinition, "bind"))

	if len(bindDefinition) != 0 {
		resultFlags, description, err := parseFlags(bindDefinition, line, lineNumber)
		if err != nil {
			return nil, err
		}

		flags = resultFlags
		hasDescription = description
	}

	bindCore, err := parseContent(bindContent, hasDescription, flags)
	if err != nil {
		return nil, err
	}

	bind, err := NewBind(*bindCore, lineNumber, line, file, false)
	if err != nil {
		return nil, err
	}

	return bind, nil
}

func parseUnbind(file *File, line string, lineNumber int) (*Unbind, error) {
	errMsg := fmt.Errorf("invalid bind format on line %d! raw line: (%s)", lineNumber, line)
	bindDefinition, bindContent, found := strings.Cut(line, "=")

	if !found {
		return nil, errMsg
	}

	bindDefinition, bindContent = strings.TrimSpace(bindDefinition), strings.TrimSpace(bindContent)
	bindDefinition = strings.TrimSpace(strings.TrimPrefix(bindDefinition, "unbind"))

	if len(bindDefinition) != 0 {
		return nil, errMsg
	}

	rawShortcut := strings.Split(bindDefinition, ",")
	if len(rawShortcut) != 2 {
		return nil, errMsg
	}

	for idx, part := range rawShortcut {
		rawShortcut[idx] = strings.TrimSpace(part)
	}

	unbind, err := NewUnbind(parseShortcut(rawShortcut), lineNumber, line, file, false)
	if err != nil {
		return nil, err
	}
	return unbind, nil
}

func parseShortcut(shortcut []string) Shortcut {
	var modKeys []string

	if len(shortcut[0]) != 0 {
		modKeys = parseModKeys(shortcut[0])
	}

	return Shortcut{
		ModKeys: modKeys,
		Key:     shortcut[1],
	}
}

func parseFlags(rawFlags, rawLine string, lineNumber int) ([]*Flag, bool, error) {
	if strings.Contains(rawFlags, "s") {
		return nil, false, errors.NewUnsupportedBindFlagError(lineNumber, rawLine, "s")
	}

	if len(rawFlags) > 14 {
		return nil, false, fmt.Errorf("invalid bind flags! line: %d", lineNumber)
	}

	separatedFlags := strings.Split(rawFlags, "")
	parsedFlags := map[string]*Flag{}
	resultFlags := make([]*Flag, 14)

	for idx, rawFlag := range separatedFlags {
		flag, ok := DefaultFlags[rawFlag]
		if !ok {
			return nil, false, fmt.Errorf("invalid flag! line: %d", lineNumber)
		}

		_, ok = parsedFlags[rawFlag]
		if ok {
			return nil, false, fmt.Errorf("invalid bind flags, flag \"%s\" repeated! line: %d", rawFlag, lineNumber)
		}

		parsedFlags[rawFlag] = flag
		resultFlags[idx] = flag
	}

	_, ok := parsedFlags["d"]
	return resultFlags, ok, nil
}

func parseContent(content string, hasDescription bool, flags []*Flag) (*BindCore, error) {
	parts := strings.Split(content, ",")
	partsLen := len(parts)

	if (hasDescription && partsLen != 5) || (partsLen != 4) {
		return nil, fmt.Errorf("invalid bind format")
	}

	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}

	shortcut := parseShortcut(parts[0:2])

	action := parts[3]
	dispatcher := parts[2]
	var description string

	if hasDescription {
		dispatcher = parts[3]
		action = parts[4]
		description = parts[2]
	}

	return &BindCore{
		Shortcut:    shortcut,
		Dispatcher:  dispatcher,
		Action:      action,
		Description: description,
		Flags:       flags,
	}, nil
}

func parseModKeys(modKeys string) []string {
	s := strings.ToUpper(modKeys)
	var separatedModKeys []string

	for _, modKey := range ModKeys {
		if strings.Contains(s, modKey) {
			separatedModKeys = append(separatedModKeys, modKey)
		}
	}

	return separatedModKeys
}
