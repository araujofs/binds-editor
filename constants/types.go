package constants

import (
	"github.com/araujofs/binds-editor/binds"
	"github.com/araujofs/binds-editor/configuration"
)

type GlobalState struct {
	SelectedFile  *binds.File
	Configuration *configuration.Configuration
}
