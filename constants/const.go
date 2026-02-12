package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	WindowSize                tea.WindowSizeMsg
	DefaultPadding            = 2
	FullScreenStyle           = lipgloss.NewStyle().Margin(0, 2)
	FullScreenFileSearchStyle = FullScreenStyle.AlignHorizontal(lipgloss.Left).AlignVertical(lipgloss.Center)
	CenterStyle               = lipgloss.NewStyle().Align(lipgloss.Center)
	BackgroundStyle           = lipgloss.NewStyle().Background(lipgloss.Color("#fff"))
)
