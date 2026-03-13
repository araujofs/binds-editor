package binds

import "strings"

type Unbind struct {
	Shortcut
	GeneralInfo
}

func NewUnbind(shortcut Shortcut, lineNumber int, rawLine string, file *File, commented bool) (*Unbind, error) {
	unbind := Unbind{
		Shortcut: shortcut,
		GeneralInfo: GeneralInfo{
			LineNumber: lineNumber,
			RawLine:    rawLine,
			File:       file,
			Commented:  commented,
		},
	}

	isValid, missingProperties := unbind.IsValid()
	if !isValid {
		return nil, missingPropertiesToError(missingProperties)
	}

	return &unbind, nil
}

func (u *Unbind) GetLine() string {
	var bind strings.Builder
	var bindDefinition strings.Builder
	var bindContent strings.Builder

	if u.Commented {
		bindDefinition.WriteString("# ")
	}

	bindDefinition.WriteString("unbind")

	bindContent.WriteString(strings.Join(u.Shortcut.ModKeys, " "))
	bindContent.WriteString(", " + u.Shortcut.Key)

	bind.WriteString(bindDefinition.String())
	bind.WriteString(" = " + bindContent.String())

	line := bind.String()

	return line
}

func (u *Unbind) GetLineNumber() int {
	return u.LineNumber
}

func (u Unbind) GetRow() []string {
	var shortcut = ""

	if joinedModKeys := strings.Join(u.Shortcut.ModKeys, "+"); joinedModKeys != "" {
		shortcut = strings.Join([]string{joinedModKeys, u.Shortcut.Key}, "+")
	} else {
		shortcut = u.Shortcut.Key
	}

	bindTypeString := "Unbind"

	if u.Commented {
		bindTypeString = "Commented"
	}

	return []string{
		shortcut,
		bindTypeString,
		"",
		"",
		"",
	}
}

func (b Unbind) IsValid() (isValid bool, missingProperties []string) {
	_, missingProperties = b.Shortcut.isValid()
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
