package files

import (
	"fmt"
	"testing"
)

func TestBindsReading(t *testing.T) {
	binds, err := ReadBindsFile("../keybinds.test.conf")

	if err != nil {
		t.Fatal(err)
	}

	for _, bind := range binds {
		fmt.Printf("%+v\n\n", bind)
	}

	fmt.Printf("%+v\n", variables)
}
