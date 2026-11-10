package main

import (
	"testing"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestVideoTabs keeps the video that plays in a tab of the top bar, whose
// click plays it again, and whose × closes it.
func TestVideoTabs(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	vp.current = youtube.Video{ID: "aaaaaaaaaaa", Title: "Bald Eagles Ep.33", Channel: "Cleetus"}
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, "/videos/search")
	if err := tt.Click("Keep the video that plays in a tab"); err != nil {
		t.Fatal(err)
	}
	s.drain()
	tt.Frame()
	if len(s.set.VideoTabs) != 1 || s.set.VideoTabs[0].Title != "Bald Eagles Ep.33" {
		t.Fatalf("tabs %+v", s.set.VideoTabs)
	}
	s.addVideoTab(youtube.Video{ID: "bbbbbbbbbbb", Title: "Second"})
	s.addVideoTab(youtube.Video{ID: "bbbbbbbbbbb", Title: "Second"})
	if len(s.set.VideoTabs) != 2 {
		t.Fatalf("a video kept twice: %+v", s.set.VideoTabs)
	}
	tt.Frame()
	if err := tt.Click("Second"); err != nil {
		t.Fatalf("no tab: %v: %q", err, tt.Texts())
	}
	s.drain()
	tt.Frame()
	if vp.current.ID != "bbbbbbbbbbb" || s.router.Path() != watchPath("bbbbbbbbbbb") {
		t.Errorf("the tab played %q at %q", vp.current.ID, s.router.Path())
	}
	s.closeVideoTab("aaaaaaaaaaa")
	if len(s.set.VideoTabs) != 1 || s.set.VideoTabs[0].ID != "bbbbbbbbbbb" {
		t.Errorf("after closing: %+v", s.set.VideoTabs)
	}
	// The tabs are kept in the settings.
	if got := s.base.VideoTabs; len(got) != 1 {
		t.Errorf("saved tabs %+v", got)
	}
}
