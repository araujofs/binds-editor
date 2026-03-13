package binds

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type File struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Output   *File  `json:"output"`
	IsOutput bool   `json:"isOutput"`
	Lines    []Line
}

type Line interface {
	GetLine() string
	GetLineNumber() int
}

func (f *File) AddLine(line Line) {
	f.Lines = append(f.Lines, line)
}

func (f *File) DeleteLine(line Line) error {
	number := line.GetLineNumber()

	if f.Lines[number] != line {
		return fmt.Errorf("line doesnt exist in file")
	}

	f.Lines = slices.Delete(f.Lines, number, number+1)
	return nil
}

func (f *File) WriteLines() error {
	var content strings.Builder

	for _, line := range f.Lines {
		line := line.GetLine()

		content.WriteString(line)
	}

	return nil
}

func (f File) IsReadonly() bool {
	is := f.Output != nil && !f.IsOutput

	if !is {
		f.Output = nil
	}

	return is
}

func (f File) FilterValue() string {
	return f.Name
}

func (f File) Description() string {
	var description string

	description += f.Path

	if f.IsReadonly() {
		description += " - " + f.Output.Path
	}

	return description
}

func (f File) Title() string {
	readonly := lipgloss.NewStyle().Foreground(lipgloss.Color("#fb1")).Render(" (Readonly)")

	if f.IsReadonly() {
		return f.Name + readonly
	}

	return f.Name
}

func NewFile(name, path string, output *string) *File {
	file := &File{
		Name: name,
		Path: path,
	}

	if output != nil {
		file.Output = &File{
			Path:     *output,
			IsOutput: true,
		}
	}

	return file
}

type RawLine struct {
	Content    *string
	LineNumber int
}

func (rl *RawLine) GetLine() string {
	if rl.Content == nil {
		return "\n"
	}

	return *rl.Content
}

func (rl *RawLine) GetLineNumber() int {
	return rl.LineNumber
}

type Shortcut struct {
	ModKeys []string
	Key     string
}

func (s Shortcut) isValid() (bool, []string) {
	var missingProperties []string

	if s.Key == "" {
		missingProperties = append(missingProperties, "key")
	}

	if modKeys := s.ModKeys; len(modKeys) == 0 {
		missingProperties = append(missingProperties, "mod keys")
	}

	if len(missingProperties) != 0 {
		return false, missingProperties
	}

	return true, missingProperties
}

type GeneralInfo struct {
	LineNumber int
	RawLine    string
	Commented  bool
	File       *File
}

func (gi GeneralInfo) isValid() (bool, []string) {
	var missingProperties []string

	if strings.TrimSpace(gi.RawLine) == "" {
		missingProperties = append(missingProperties, "rawLine")
	}

	if gi.File == nil {
		missingProperties = append(missingProperties, "file path")
	}

	if len(missingProperties) != 0 {
		return false, missingProperties
	}

	return true, missingProperties

}
