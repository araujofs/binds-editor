package files

import (
	"fmt"
	"testing"
)

func TestBindsReading(t *testing.T) {
	binds, err := readBindsFile("/home/arthur/.local/share/omarchy/default/hypr/bindings/tiling.conf")	

	if err != nil {
		t.Fatal(err)
	}

	for _, bind := range binds {
		fmt.Printf("%+v\n\n", bind)
	}
}