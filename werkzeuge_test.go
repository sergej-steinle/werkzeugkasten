package werkzeugkasten

import "testing"

func TestWerkzeuge_RandomString(t *testing.T) {
	var werkzeug Werkzeug
	length := 10
	s := werkzeug.RandomString(length)

	if len(s) != length {
		t.Errorf("Want length of %d but got length of %d\n", length, len(s))
	}
}
