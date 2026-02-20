package constants

import "github.com/araujofs/binds-editor/configuration"

type GlobalState struct {
	SelectedFile  *configuration.File
	Configuration *configuration.Configuration
}
