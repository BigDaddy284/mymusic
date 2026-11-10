package app

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestVideoStats counts the time each video played, day by day, a view
// once it played thirty seconds, and keeps them from one start to the
// next.
func TestVideoStats(t *testing.T) {
	path := filepath.Join(t.TempDir(), "video-stats.json")
	vs, err := openVideoStats(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.Local)
	vs.now = func() time.Time { return now }
	a := youtube.Video{ID: "aaaaaaaaaaa", Title: "A", Channel: "Chan"}
	b := youtube.Video{ID: "bbbbbbbbbbb", Title: "B", Channel: "Chan"}
	for range 100 { // 50 seconds
		vs.watched(a, 0.5)
	}
	for range 20 { // 10 seconds: no view
		vs.watched(b, 0.5)
	}
	vs.started()
	now = now.AddDate(0, 0, -10)
	for range 4 * 60 {
		vs.watched(b, 0.5)
	}
	vs.flush()

	again, err := openVideoStats(path)
	if err != nil {
		t.Fatal(err)
	}
	again.now = func() time.Time { return time.Date(2026, 11, 10, 18, 0, 0, 0, time.Local) }
	week := again.stats(7)
	if week.Seconds != 60 || week.Views != 1 {
		t.Errorf("this week: %v seconds, %d views", week.Seconds, week.Views)
	}
	if len(week.Daily) != 7 || week.Daily[6].Date != "2026-11-10" || week.Daily[6].Minutes != 1 {
		t.Errorf("this week's days: %+v", week.Daily)
	}
	all := again.stats(0)
	if all.Seconds != 180 || all.Views != 2 {
		t.Errorf("all time: %v seconds, %d views", all.Seconds, all.Views)
	}
	if len(all.TopVideos) != 2 || all.TopVideos[0].Video.ID != b.ID || all.TopVideos[0].Seconds != 130 || all.TopVideos[0].Video.Title != "B" {
		t.Errorf("top videos: %+v", all.TopVideos)
	}
	if len(all.TopChannels) != 1 || all.TopChannels[0].Seconds != 180 || all.TopChannels[0].Views != 2 {
		t.Errorf("top channels: %+v", all.TopChannels)
	}
}

// TestVideoPositions keeps where each video was left, from one start to
// the next, and forgets one at 0.
func TestVideoPositions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "video-stats.json")
	vs, err := openVideoStats(path)
	if err != nil {
		t.Fatal(err)
	}
	vs.leftAt("a", 95)
	vs.leftAt("b", 40)
	vs.leftAt("b", 0)
	vs.flush()
	again, err := openVideoStats(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := again.position("a"); got != 95 {
		t.Errorf("a was left at %v", got)
	}
	if got := again.position("b"); got != 0 {
		t.Errorf("b, forgotten, was left at %v", got)
	}
}
