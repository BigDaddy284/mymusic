package downloads

import (
	"path/filepath"
	"testing"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// TestVideoFile tells the file of a video downloaded.
func TestVideoFile(t *testing.T) {
	dir := t.TempDir()
	video := yt("fffffffffff", "Teardrop")
	m := open(t, dir, Options{Fetcher: &fetcher{}})
	m.Start()
	m.Add([]media.Track{video}, "best")
	until(t, "the download", func() bool { return kind(m, video.ID) == Done })
	if p, ok := m.Snapshot().VideoFile("fffffffffff"); !ok || p != filepath.Join(dir, "audio", "fffffffffff.webm") {
		t.Errorf("the file %q, %v", p, ok)
	}
	if p, ok := m.Snapshot().VideoFile("ggggggggggg"); ok {
		t.Errorf("a video not downloaded: %q", p)
	}
	var none *Snapshot
	if _, ok := none.VideoFile("fffffffffff"); ok {
		t.Error("a nil snapshot holds a file")
	}
}
