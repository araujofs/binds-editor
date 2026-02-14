package messages

import (
	"fmt"

	consts "github.com/araujofs/binds-editor/constants"
	tea "github.com/charmbracelet/bubbletea"
)

type InfoModel struct {
	Message *string
	Error   *string
}

func GetDefaultInfoModel() InfoModel {
	return InfoModel{
		Message: nil,
		Error:   nil,
	}
}

type ErrorMsg struct {
	Message string
}

type MessageMsg struct {
	Message string
}

func (model *InfoModel) GetStyledMessage() string {
	message := ""
	if model.Message != nil {
		message = " |" + consts.MessageStyle.Render(fmt.Sprintf(" (%s)", *model.Message))
	}

	if model.Error != nil {
		message = " >" + consts.ErrorStyle.Render(fmt.Sprintf(" (%s)", *model.Error))
	}

	return message + "\n\n"

}

func NewErrorMsg(message string) ErrorMsg {
	return ErrorMsg{message}
}

func SendErrorMsg(message string) tea.Cmd {
	return func() tea.Msg {
		return NewErrorMsg(message)
	}
}

func NewMessageMsg(message string) MessageMsg {
	return MessageMsg{message}
}

func SendMessageMsg(message string) tea.Cmd {
	return func() tea.Msg {
		return NewErrorMsg(message)
	}
}
