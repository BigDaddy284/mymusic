package library

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestForgetSpotify opens a library MyMusic 0.x kept, with Spotify's
// tracks, artists and podcasts beside YouTube Music's: what came from
// Spotify goes, the rest stays, and the file as it was stays beside it.
func TestForgetSpotify(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "library.json")
	old := `{
  "liked": [{"id": "sp:a", "source": "spotify", "title": "A"}, {"id": "yt:b", "source": "youtube", "title": "B"}],
  "disliked": [{"id": "sp:c", "source": "spotify", "title": "C"}],
  "playlists": [{"id": "mx:p", "title": "P", "tracks": [{"id": "yt:d", "source": "youtube"}, {"id": "sp:e", "source": "spotify"}]}],
  "history": [{"track": {"id": "sp:a", "source": "spotify"}, "playedAt": "2026-01-01T00:00:00Z"}],
  "searches": ["spotify"],
  "artists": [{"id": "sp:f", "source": "spotify", "name": "F"}, {"id": "yt:g", "source": "youtube", "name": "G"}],
  "podcasts": [{"id": "sp:h", "source": "spotify", "title": "H"}]
}`
	if err := os.WriteFile(path, []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	l, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	d := l.file.Get()
	if len(d.Liked) != 1 || d.Liked[0].ID != "yt:b" || len(d.Disliked) != 0 || len(d.History) != 0 {
		t.Errorf("liked %v, disliked %v, history %v", d.Liked, d.Disliked, d.History)
	}
	if len(d.Playlists) != 1 || len(d.Playlists[0].Tracks) != 1 || d.Playlists[0].Tracks[0].ID != "yt:d" {
		t.Errorf("playlists %+v", d.Playlists)
	}
	if len(d.Artists) != 1 || d.Artists[0].ID != "yt:g" || len(d.Podcasts) != 0 || len(d.Searches) != 1 {
		t.Errorf("artists %v, podcasts %v, searches %v", d.Artists, d.Podcasts, d.Searches)
	}
	if kept, err := os.ReadFile(path + ".spotify"); err != nil || string(kept) != old {
		t.Errorf("the library as it was: %q, %v", kept, err)
	}
	if data, _ := os.ReadFile(path); strings.Contains(string(data), `"sp:`) {
		t.Errorf("library.json still holds Spotify's: %s", data)
	}
	// Opened again, with nothing of Spotify, it stays as it is.
	if err := os.Remove(path + ".spotify"); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".spotify"); err == nil {
		t.Error("a library without Spotify's was kept aside")
	}
}
