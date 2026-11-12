package settings

import "testing"

// TestPanel keeps the side panels the window shows, the Video side's
// description among them, and the queue for one it doesn't know.
func TestPanel(t *testing.T) {
	for file, want := range map[string]string{
		`{}`:                  "queue",
		`{"panel":"queue"}`:   "queue",
		`{"panel":"lyrics"}`:  "lyrics",
		`{"panel":"about"}`:   "about",
		`{"panel":"none"}`:    "none",
		`{"panel":"unknown"}`: "queue",
	} {
		s, err := Decode([]byte(file))
		if err != nil {
			t.Fatal(err)
		}
		if got := s.Clean().Panel; got != want {
			t.Errorf("%s shows the panel %q, not %q", file, got, want)
		}
	}
}
