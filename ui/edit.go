package ui

import (
	"github.com/araujofs/binds-editor/binds"
	tea "github.com/charmbracelet/bubbletea"
)

type Edit struct {
	bind binds.Bind
}

func InitEdit(bind *binds.Bind) tea.Model
func (m Edit) Init() tea.Cmd
func (m Edit) Update() (tea.Model, tea.Cmd)
func (m Edit) View() string
