package main

import (
	"slices"
	"testing"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestDockVideo docks the video that plays from its watch page with the
// player bar's button: it keeps playing small, and the page shown before
// the video started shows again; the button then opens its page again.
func TestDockVideo(t *testing.T) {
	s := newTestShell(t)
	frogs := youtube.Video{ID: "frogs000001", Title: "Frogs singing", Channel: "Pond"}
	s.loads["videos:search:frogs"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{frogs}}}
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, videoSearchPath("frogs"))
	from := s.router.Location()
	if !slices.Contains(tt.Texts(), "Dock: keep it playing small and go back to the page before") {
		t.Errorf("no dock button: %q", tt.Texts())
	}
	s.watch(frogs, nil)
	s.drain()
	tt.Frame()
	if !isWatchPath(s.router.Path()) {
		t.Fatalf("the video didn't open its page: %q", s.router.Path())
	}
	// No web view plays the video in tests, so the button stays off: what
	// it does is called.
	s.dockVideo()
	s.drain()
	tt.Frame()
	if got := s.router.Location(); got != from {
		t.Errorf("docked, the page is %q, not %q", got, from)
	}
	if vp.current.ID != frogs.ID || !s.showsMini() {
		t.Errorf("docked, the video plays %q, small: %v", vp.current.ID, s.showsMini())
	}
	s.dockVideo()
	s.drain()
	tt.Frame()
	if s.router.Path() != watchPath(frogs.ID) {
		t.Errorf("undocked, the page is %q", s.router.Path())
	}
}

// TestDockVideoWithNoPageBefore docks a video that started with no page
// before it to Home.
func TestDockVideoWithNoPageBefore(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	vp.current = youtube.Video{ID: "frogs000001"}
	s.mode.video = true
	s.navigate(watchPath(vp.current.ID))
	s.dockVideo()
	if got := s.router.Path(); got != "/videos" {
		t.Errorf("docked to %q", got)
	}
}

// TestInWatchLater tells a video is in Watch Later as its list has it,
// unless it was saved or taken off here since.
func TestInWatchLater(t *testing.T) {
	s := newTestShell(t)
	s.loads[watchLaterKey] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{{ID: "saved000001"}}}}
	if !s.inWatchLater("saved000001") || s.inWatchLater("other000001") {
		t.Fatal("Watch Later's list isn't followed")
	}
	s.markWatchLater("saved000001", false)
	s.markWatchLater("other000001", true)
	if s.inWatchLater("saved000001") || !s.inWatchLater("other000001") {
		t.Error("what was saved or taken off here isn't followed")
	}
	// A video taken off leaves Watch later's list at once.
	s.loads[watchLaterKey].(*loaded[youtube.VideoPage]).value.Videos[0].Title = "Frogs singing"
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, videoPlaylistPath(youtube.WatchLater))
	if showsText(tt.Texts(), "Frogs singing") {
		t.Error("a video taken off Watch Later still shows in it")
	}
	delete(s.videos.watchLater, "saved000001")
	tt.Frame()
	if !showsText(tt.Texts(), "Frogs singing") {
		t.Errorf("Watch later doesn't show its video: %q", tt.Texts())
	}
}

// TestVideoHomeLists puts the account's lists, once the Library's tabs,
// beside the home feed.
func TestVideoHomeLists(t *testing.T) {
	want := []string{"All", "Subscriptions", "History", "Playlists", "Watch later", "Liked videos"}
	if !slices.Equal(videoHomeChips, want) {
		t.Errorf("Home's chips are %q", videoHomeChips)
	}
	for _, n := range videoNav {
		if n.label == "Library" {
			t.Error("the Video side still has a Library")
		}
	}
}
