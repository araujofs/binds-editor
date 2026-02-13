package main

import (
	"fmt"
	"os"

	"github.com/araujofs/binds-editor/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := ui.InitFileSelection(nil, nil)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
