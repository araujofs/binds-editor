package files

import "testing"

func TestBindsReading(t *testing.T) {
	_, err := readBindingsFile("/home/arthur/.config/hypr/bindings.conf")	
	if err != nil {
		t.Fatal(err)
	}
}