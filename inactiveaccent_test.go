package main

import (
	"testing"

	"github.com/bigdaddy284/mymusic/internal/settings"
)

// TestAccentWhileBehind keeps the accent of the sliders and switches in a
// window another is in front of, as MyMusic's own parts keep theirs, in
// every layout: MyGo 0.8 fades its widgets' to gray
// there unless the theme keeps it.
func TestAccentWhileBehind(t *testing.T) {
	for _, layout := range settings.Layouts {
		if k := lookOf(layout, settings.DensityComfortable); k.theme.InactiveAccent != k.theme.Accent {
			t.Errorf("%s: the accent fades to %v behind, from %v", layout, k.theme.InactiveAccent, k.theme.Accent)
		}
	}
}
