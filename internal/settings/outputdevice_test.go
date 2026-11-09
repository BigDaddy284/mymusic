package settings

import "testing"

// TestOutputDevice keeps the output device chosen, by its ID and its name,
// trimmed, and the system's default, with no name, for a device of no ID,
// as in settings from before the choice.
func TestOutputDevice(t *testing.T) {
	for _, c := range []struct{ in, want OutputDevice }{
		{OutputDevice{ID: "{0.0.0.00000000}.{a1b2}", Name: " Speakers (Realtek(R) Audio) "}, OutputDevice{ID: "{0.0.0.00000000}.{a1b2}", Name: "Speakers (Realtek(R) Audio)"}},
		{OutputDevice{ID: "{0.0.0.00000000}.{c3d4}"}, OutputDevice{ID: "{0.0.0.00000000}.{c3d4}"}},
		{OutputDevice{ID: "  ", Name: "Speakers"}, OutputDevice{}},
		{OutputDevice{}, OutputDevice{}},
	} {
		s := Default()
		s.OutputDevice = c.in
		if got := s.Clean().OutputDevice; got != c.want {
			t.Errorf("%+v cleans to %+v, want %+v", c.in, got, c.want)
		}
	}
	if d := Default().OutputDevice; d != (OutputDevice{}) {
		t.Errorf("a new install plays on %+v", d)
	}
}
