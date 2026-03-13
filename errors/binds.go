package errors

import "fmt"

type UnsupportedBindFlagError struct {
	lineNumber int
	rawLine    string
	flag       string
}

func NewUnsupportedBindFlagError(lineNumber int, rawLine string, flag string) UnsupportedBindFlagError {
	return UnsupportedBindFlagError{
		lineNumber: lineNumber,
		rawLine:    rawLine,
		flag:       flag,
	}
}

func IsUnsupportedBindFlagError(err error) bool {
	switch err.(type) {
	case UnsupportedBindFlagError:
		return true
	default:
		return false
	}
}

func (u UnsupportedBindFlagError) Error() string {
	return fmt.Sprintf("binds with the '%s' flag are not supported! Line: %d, Raw line: %s", u.flag, u.lineNumber, u.rawLine)
}

type NotABindError struct{}

func (nab NotABindError) Error() string {
	return "this line is not a bind"
}

func IsNotABindError(err error) bool {
	switch err.(type) {
	case NotABindError:
		return true
	default:
		return false
	}
}
