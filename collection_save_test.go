package main

import "testing"

// The account's own lists can't be saved, which YouTube answers with "not
// found"; any other playlist can.
func TestSavable(t *testing.T) {
	for id, want := range map[string]bool{
		"yt:LM":          false,
		"yt:SE":          false,
		"yt:PLabc":       true,
		"yt:RDCLAK5uy_x": true,
		"mx:LM":          true,
	} {
		if got := savable(id); got != want {
			t.Errorf("savable(%q) = %v, want %v", id, got, want)
		}
	}
}
