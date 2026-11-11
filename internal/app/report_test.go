package app

import (
	"testing"
	"time"

	"github.com/bigdaddy284/mymusic/internal/media"
)

// TestTrackStartsOnce tells that a track started once it sounds: not while
// it loads, nor as the engine holds it before any audio comes, nor when it
// plays again after a pause; again after nothing played (#145).
func TestTrackStartsOnce(t *testing.T) {
	var a App
	ta, tb := &media.Track{ID: "yt:a"}, &media.Track{ID: "yt:b"}
	for i, c := range []struct {
		np      NowPlaying
		started bool
	}{
		{NowPlaying{Track: ta, Loading: true}, false},
		{NowPlaying{Track: ta, Playing: true, Loading: true}, false},
		{NowPlaying{Track: ta, Playing: true}, true},
		{NowPlaying{Track: ta, Playing: true, PositionMs: 1000}, false},
		{NowPlaying{Track: ta}, false}, // paused
		{NowPlaying{Track: ta, Playing: true}, false},
		{NowPlaying{Track: tb, Loading: true}, false}, // the next loads
		{NowPlaying{Track: tb, Playing: true}, true},
		{NowPlaying{}, false}, // stopped
		{NowPlaying{Track: tb, Playing: true}, true},
	} {
		if got := a.startOf(c.np); got != c.started {
			t.Errorf("report %d, %+v: started %v, want %v", i, c.np, got, c.started)
		}
	}
}

// TestLoadingPlays has the media controls and the title count a track that
// loads as playing, as it plays once loaded; a track that failed, neither
// playing nor loading, is paused.
func TestLoadingPlays(t *testing.T) {
	track := &media.Track{ID: "yt:a", Title: "Song"}
	if info := controlsInfo(NowPlaying{Track: track, Loading: true}, time.Now()); !info.Playing {
		t.Error("the controls show a track that loads as paused")
	}
	if info := controlsInfo(NowPlaying{Track: track}, time.Now()); info.Playing {
		t.Error("the controls show a paused track as playing")
	}
	if got := playingLine(NowPlaying{Track: track, Loading: true}); got != "Song" {
		t.Errorf("a track that loads is %q", got)
	}
}

// TestLargeArtwork picks an https cover of about 512 pixels.
func TestLargeArtwork(t *testing.T) {
	if got := largeArtwork([]media.Image{{URL: "http://plain.example/a.jpg", Width: 512}}); got != "" {
		t.Errorf("plain http artwork used: %q", got)
	}
	got := largeArtwork([]media.Image{{URL: "https://a/60", Width: 60}, {URL: "https://a/226", Width: 226}})
	if got != "https://a/226" {
		t.Errorf("small artwork: %q, want the largest", got)
	}
}
