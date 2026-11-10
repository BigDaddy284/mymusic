package app

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/bigdaddy284/mymusic/internal/library"
	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// mondayEpisode is a YouTube Music podcast's episode.
var mondayEpisode = media.Track{ID: "yt:e1", Source: media.SourceYouTube, Title: "Monday", DurationMs: 1830000,
	Episode: &media.Episode{Show: media.Ref{ID: "yt:s1", Name: "The Daily Thing"}}}

// TestKeepPosition keeps where the episode that plays is, as the player
// reports it, but not while it loads, nor of a song.
func TestKeepPosition(t *testing.T) {
	pos, err := library.OpenPositions(filepath.Join(t.TempDir(), "positions.json"))
	if err != nil {
		t.Fatal(err)
	}
	// The places are written 10 s after a change, by when the folder may
	// be gone: at once, before it goes, as App.Close does.
	t.Cleanup(func() { pos.Flush() })
	a := &App{positions: pos}
	ep := mondayEpisode
	a.keepPosition(NowPlaying{Track: &ep, Playing: true, PositionMs: 120000, DurationMs: 1830000})
	if got := (playerServices{a}).ResumeAt(ep); got != 115*time.Second {
		t.Errorf("ResumeAt = %v", got)
	}
	a.keepPosition(NowPlaying{Track: &ep, Loading: true, PositionMs: 300000, DurationMs: 1830000})
	if p, _ := pos.Get(ep.ID); p.PositionMs != 120000 {
		t.Errorf("a load moved the place to %d", p.PositionMs)
	}
	song := media.Track{ID: "yt:s", DurationMs: 200000}
	a.keepPosition(NowPlaying{Track: &song, Playing: true, PositionMs: 60000})
	if _, ok := pos.Get(song.ID); ok {
		t.Error("a song's place was kept")
	}
	(&Collection{a}).SetPlayed(ep, true)
	if p := (&Collection{a}).InProgress(0); len(p) != 0 {
		t.Errorf("in progress: %+v", p)
	}
}

// TestNoScrobbleOfEpisodes tells Last.fm of no episode.
func TestNoScrobbleOfEpisodes(t *testing.T) {
	ep := mondayEpisode
	if id, _ := lastfmTrack(NowPlaying{Track: &ep, Playing: true}); id != "" {
		t.Errorf("an episode is observed as %q", id)
	}
}

// TestSearchPodcastsFails tells why the search found nothing when YouTube
// Music fails.
func TestSearchPodcastsFails(t *testing.T) {
	yt := youtube.New()
	yt.BaseURL = "http://127.0.0.1:9"
	a := &App{youtube: yt}
	res, err := (&Catalog{a}).SearchPodcasts(context.Background(), "daily")
	if err == nil || len(res.Errors) != 1 || res.YouTube.Podcasts == nil {
		t.Errorf("SearchPodcasts = %+v, %v", res, err)
	}
	if res, err := (&Catalog{a}).SearchPodcasts(context.Background(), "  "); err != nil || len(res.Errors) != 0 {
		t.Errorf("an empty search = %+v, %v", res, err)
	}
}
