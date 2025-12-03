package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Shortcut struct {
	ModKeys []ModKey
	Key []string
}

type Bind struct {
	BindCore
	LineNumber int
	RawLine int
}

type BindCore struct {
	Shortcut Shortcut
	ActionType string
	Action string
	Description *string
}

type ModKey string

const (
	SHIFT ModKey = "SHIFT"
	CAPS ModKey = "CAPS"	
	CTRL ModKey = "CTRL"	
	CONTROL ModKey = "CONTROL"	
	ALT ModKey = "ALT"	
	MOD2 ModKey = "MOD2"	
	MOD3 ModKey = "MOD3"	
	SUPER ModKey = "SUPER"	
	WIN ModKey = "WIN"	
	LOGO ModKey = "LOGO"	
	MOD4 ModKey = "MOD4"	
	MOD5 ModKey = "MOD5"	
)

func readBindsFile(path string) ([]Bind, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bindings := make([]Bind, 0, 50)

	scanner := bufio.NewScanner(f)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()
		bindings = append(bindings, parseBind(line, lineNumber))
	}

	return nil, nil
}

func parseBind(rawBindLine string, bindLineNumber int) (*Bind, error) {
	if (rawBindLine[0] == '#') {
		return nil, nil
	}
	errorMsg := fmt.Errorf("Invalid bind format on line %d! Raw line: (%s)", bindLineNumber, rawBindLine)

	bindType, bindContent, found := strings.Cut(rawBindLine, "=")

	if !found {
		return nil, errorMsg
	}

	bindType = strings.TrimSpace(bindType)

	if !(strings.Contains(bindType, "bind")) || !(len(bindType) <= 5) {
		return nil, errorMsg
	}

	if len(bindType) == 4 {
		parseBindContent(bindContent)
	}

	switch bindType[4] {
	case 'd':
		parseBindContentWithDescription(bindContent)
	default:
		return nil, errorMsg
	}
}

func parseBindContent(bindContent string) (*BindCore, error) {
	parts := strings.Split(bindContent, ",")

	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid bind format!")
	}
	
	modKeys := strings.TrimSpace(parts[0])
	keys := strings.TrimSpace(parts[1])
	actionType := strings.TrimSpace(parts[2])
	action := strings.TrimSpace(parts[3])

	
}

func parseBindContentWithDescription(bindContent string) (*BindCore, error) {
	parts := strings.Split(bindContent, ",")

	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid bind format!")
	}

}