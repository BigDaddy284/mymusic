package main

import (
	"testing"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestVideoInfoBesidePlay lines the video's title up with the top of the
// play button, with its views under the channel.
func TestVideoInfoBesidePlay(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	vp.current = youtube.Video{ID: "a", Title: "A video title", Channel: "A channel", Views: "12K views"}
	vp.list = []youtube.Video{vp.current}
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, "/videos/trending")
	views, ok := tt.Find("12K views")
	if !ok {
		t.Fatalf("no views: %q", tt.Texts())
	}
	play, ok := tt.Access("Play")
	if !ok {
		t.Fatal("no Play button")
	}
	like, ok := tt.Access("Like")
	if !ok {
		t.Fatal("no Like button")
	}
	seek, _ := tt.Access("Seek")
	t.Logf("views %+v like %+v play %+v seek %+v", views, like.Bounds, play.Bounds, seek.Bounds)
	// The title and the channel sit above the views, so the column starts
	// 36px or so above them.
	if top := views.Y - 38; top < play.Bounds.Y-4 || top > play.Bounds.Y+4 {
		t.Errorf("the info starts at %.0f, Play at %.0f", top, play.Bounds.Y)
	}
	if bottom := like.Bounds.Y + like.Bounds.H; bottom > seek.Bounds.Y+seek.Bounds.H+2 {
		t.Errorf("the like button runs out of the bar, to %.0f", bottom)
	}
}
