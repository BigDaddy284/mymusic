package match

import (
	"testing"

	"github.com/bigdaddy284/mymusic/internal/media"
)

func TestNormalize(t *testing.T) {
	for in, want := range map[string]string{
		"Beyoncé":                       "beyonce",
		"  AC/DC — Back in Black! ":     "ac dc back in black",
		"Song (Remastered 2011)":        "song remastered 2011",
		"Mötley Crüe":                   "motley crue",
		"Hello, World (feat. Somebody)": "hello world feat somebody",
	} {
		if got := normalize(in); got != want {
			t.Errorf("normalize(%q) = %q, want %q", in, got, want)
		}
	}
	if got := CleanTitle("Song (Remastered 2011)"); got != "Song" {
		t.Errorf("cleanTitle = %q", got)
	}
	if got := CleanTitle("Song - 2011 Remaster"); got != "Song" {
		t.Errorf("cleanTitle = %q", got)
	}
}

func items(tracks ...media.Track) []media.Item {
	out := make([]media.Item, len(tracks))
	for i := range tracks {
		out[i] = media.Item{Kind: media.ItemTrack, Track: &tracks[i]}
	}
	return out
}
