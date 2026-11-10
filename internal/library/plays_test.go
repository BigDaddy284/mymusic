package library

import (
	"maps"
	"testing"

	"github.com/bigdaddy284/mymusic/internal/media"
)

func TestPlayCounts(t *testing.T) {
	l := open(t)
	if got := l.PlayCounts(); len(got) != 0 {
		t.Errorf("an empty history counts %v", got)
	}
	a, b := tr("yt:a", media.SourceYouTube, "A"), tr("file:b", media.SourceFile, "B")
	for _, p := range []media.Track{a, b, a, a} {
		if err := l.RecordPlay(p); err != nil {
			t.Fatal(err)
		}
	}
	if got, want := l.PlayCounts(), map[string]int{"yt:a": 3, "file:b": 1}; !maps.Equal(got, want) {
		t.Errorf("counts = %v, want %v", got, want)
	}
	if err := l.ClearHistory(); err != nil {
		t.Fatal(err)
	}
	if got := l.PlayCounts(); len(got) != 0 {
		t.Errorf("the history cleared still counts %v", got)
	}
}
