package main

import (
	"testing"
	"time"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestMiddleButtonScrolls scrolls the page while the middle button is held
// below where it went down, and stops when it lets go; a click of it
// keeps scrolling until the next click.
func TestMiddleButtonScrolls(t *testing.T) {
	s := newTestShell(t)
	var vids []youtube.Video
	for i := range 40 {
		vids = append(vids, youtube.Video{ID: "video" + count(i), Title: "Video " + count(i)})
	}
	s.loads["videos:search:frogs"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: vids}}
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, videoSearchPath("frogs"))
	tt.PressButton(2, 600, 400)
	tt.Move(600, 600)
	for range 20 {
		tt.Advance(25 * time.Millisecond)
	}
	moved := s.pageY()
	if moved < 100 {
		t.Fatalf("holding the middle button 200px down scrolled %.0fpx", moved)
	}
	tt.ReleaseButton(2, 600, 600)
	tt.Advance(200 * time.Millisecond)
	if s.autoScroll.on || s.pageY()-moved > 60 {
		t.Errorf("the page went on scrolling after the button let go: %.0f to %.0f", moved, s.pageY())
	}
	// A click without moving keeps it on, until the next click.
	tt.PressButton(2, 600, 400)
	tt.ReleaseButton(2, 600, 400)
	if !s.autoScroll.on {
		t.Fatal("a click of the middle button didn't start scrolling")
	}
	tt.Move(600, 300)
	at := s.pageY()
	for range 10 {
		tt.Advance(25 * time.Millisecond)
	}
	if s.pageY() >= at {
		t.Errorf("the pointer above the anchor didn't scroll up: %.0f to %.0f", at, s.pageY())
	}
	tt.Press(600, 300)
	tt.Release(600, 300)
	if s.autoScroll.on || s.videos.player.current.ID != "" {
		t.Errorf("the click that ends it didn't, or played a video: %+v", s.autoScroll)
	}
}

// TestMusicMiniOnVideoSide shows the music that plays at the top right of
// the Video side, whose buttons control it, until it's closed.
func TestMusicMiniOnVideoSide(t *testing.T) {
	s := playing(t)
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, "/videos")
	if !s.mode.video {
		t.Fatal("not on the Video side")
	}
	cur, _ := s.st.Current()
	if !showsText(tt.Texts(), cur.Track.Title) {
		t.Fatalf("the Video side doesn't show the music playing, %q: %q", cur.Track.Title, tt.Texts())
	}
	click(t, s, tt, "Close the music player")
	tt.Frame()
	if s.showsMusicMini() {
		t.Error("the music player didn't close")
	}
	show(s, tt, "/")
	if s.showsMusicMini() {
		t.Error("the Music side shows it too")
	}
}

// TestMusicMiniStaysForAVideo keeps the music's player on the Video side's
// home while a video plays, paused, to go back to.
func TestMusicMiniStaysForAVideo(t *testing.T) {
	s := playing(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, "/videos")
	vp.current = youtube.Video{ID: "a", Title: "A video"}
	vp.state.Playing = true
	tt.Frame()
	if !s.showsMusicMini() {
		t.Error("the music's player went as a video played")
	}
}
