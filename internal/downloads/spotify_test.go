package downloads

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// TestForgetSpotify reads the downloads of MyTube 0.x, with Spotify
// tracks and a Spotify album: a Spotify track downloaded on its own
// becomes the recording it downloaded, the album and the tracks only it
// kept go, and so does the audio no track needs then.
func TestForgetSpotify(t *testing.T) {
	dir := t.TempDir()
	old := `{"version": 1, "seq": 4,
  "tracks": {
    "sp:own": {"track": {"id": "sp:own", "source": "spotify"}, "recording": {"id": "yt:aaaaaaaaaaa", "source": "youtube", "title": "A"}, "quality": "best", "pins": ["track"], "seq": 1},
    "sp:album": {"track": {"id": "sp:album", "source": "spotify"}, "recording": {"id": "yt:bbbbbbbbbbb", "source": "youtube"}, "quality": "best", "pins": ["album:sp:x"], "seq": 2, "art": "art/b.jpg"},
    "sp:both": {"track": {"id": "sp:both", "source": "spotify"}, "recording": {"id": "yt:ccccccccccc", "source": "youtube"}, "quality": "best", "pins": ["album:sp:x", "playlist:p"], "seq": 3},
    "yt:ccccccccccc": {"track": {"id": "yt:ccccccccccc", "source": "youtube"}, "quality": "best", "pins": ["track"], "seq": 4}
  },
  "files": {
    "aaaaaaaaaaa": {"path": "audio/aaaaaaaaaaa.webm"},
    "bbbbbbbbbbb": {"path": "audio/bbbbbbbbbbb.webm"},
    "ccccccccccc": {"path": "audio/ccccccccccc.webm"}
  },
  "collections": {
    "album:sp:x": {"id": "album:sp:x", "kind": "album", "source": "spotify", "tracks": ["sp:album", "sp:both"], "art": "art/x.jpg"},
    "playlist:p": {"id": "playlist:p", "kind": "playlist", "source": "local", "tracks": ["sp:both"]}
  },
  "art": {"https://img.test/b": "art/b.jpg", "https://img.test/x": "art/x.jpg"}
}`
	if err := os.WriteFile(filepath.Join(dir, "index.json"), []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	x, ok, err := readIndex(dir)
	if err != nil || !ok {
		t.Fatal(ok, err)
	}
	var ids []string
	for id := range x.Tracks {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	if !slices.Equal(ids, []string{"yt:aaaaaaaaaaa", "yt:ccccccccccc"}) {
		t.Fatalf("tracks %v", ids)
	}
	if e := x.Tracks["yt:aaaaaaaaaaa"]; e.Track.Title != "A" || e.Recording != nil || !slices.Equal(e.Pins, []string{"track"}) {
		t.Errorf("the track downloaded on its own: %+v", e)
	}
	if e := x.Tracks["yt:ccccccccccc"]; !slices.Equal(e.Pins, []string{"track", "playlist:p"}) {
		t.Errorf("the recording kept twice is pinned by %v", e.Pins)
	}
	if _, ok := x.Collections["album:sp:x"]; ok || len(x.Collections) != 1 {
		t.Errorf("collections %v", x.Collections)
	}
	if c := x.Collections["playlist:p"]; !slices.Equal(c.Tracks, []string{"yt:ccccccccccc"}) {
		t.Errorf("the playlist holds %v", c.Tracks)
	}
	if _, ok := x.Files["bbbbbbbbbbb"]; ok || len(x.Files) != 2 {
		t.Errorf("files %v", x.Files)
	}
	slices.Sort(x.Pending)
	if want := []string{"art/b.jpg", "art/x.jpg", "audio/bbbbbbbbbbb.webm"}; !slices.Equal(x.Pending, want) {
		t.Errorf("to delete %v, want %v", x.Pending, want)
	}
	if len(x.Art) != 0 {
		t.Errorf("covers %v", x.Art)
	}
}
