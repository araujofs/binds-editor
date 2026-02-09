package messages

type ErrorMsg struct {
	error
}

func NewErrorMsg(e error) ErrorMsg {
	return ErrorMsg{e}
}
