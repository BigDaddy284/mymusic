package main

import (
	"slices"
	"testing"
	"time"

	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestVideoCardMenu opens a card's menu from its ⋮ button, without playing
// the video, adds the video to Up next, and takes a video YouTube is told
// not to recommend off the cards.
func TestVideoCardMenu(t *testing.T) {
	s := newTestShell(t)
	s.loads["videos:search:frogs"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{
		{ID: "frogs000001", Title: "Frogs singing", Channel: "Pond"},
		{ID: "toads000001", Title: "Toads singing", Channel: "Pond", Feedback: youtube.VideoFeedback{NotInterested: "T"}},
	}}}
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	vp.current = youtube.Video{ID: "playing0001", Title: "Playing"}
	vp.list = []youtube.Video{vp.current, {ID: "after000001"}}
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, videoSearchPath("frogs"))
	if err := tt.Click("More actions"); err != nil {
		t.Fatal(err)
	}
	if vp.current.ID != "playing0001" {
		t.Fatalf("the ⋮ button played %q", vp.current.ID)
	}
	items := tt.Menu()
	for _, want := range []string{"Add to queue", "Save to Watch Later", "Save to playlist", "Share: copy the link", "Not interested", "Don't recommend channel", "Report on YouTube"} {
		if !slices.Contains(items, want) {
			t.Errorf("the menu has no %q: %q", want, items)
		}
	}
	if err := tt.ChooseMenuItem("Add to queue"); err != nil {
		t.Fatal(err)
	}
	ids := []string{}
	for _, v := range vp.list {
		ids = append(ids, v.ID)
	}
	if !slices.Equal(ids, []string{"playing0001", "frogs000001", "after000001"}) {
		t.Errorf("Up next is %q", ids)
	}
	s.videos.hidden = map[string]bool{"toads000001": true}
	tt.Frame()
	if showsText(tt.Texts(), "Toads singing") {
		t.Error("a video YouTube was told not to recommend is still shown")
	}
}

// TestVideoGridFillsThePage stretches the cards to fill the page's width,
// with the side panel open and closed.
func TestVideoGridFillsThePage(t *testing.T) {
	s := newTestShell(t)
	var vids []youtube.Video
	for i := range 12 {
		vids = append(vids, youtube.Video{ID: "video" + count(i), Title: "Video " + count(i)})
	}
	s.loads["videos:search:frogs"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: vids}}
	tt := testerOf(t, s, 1600, 900)
	for _, panel := range []string{"queue", ""} {
		s.panel = panel
		show(s, tt, videoSearchPath("frogs"))
		tt.Advance(50 * time.Millisecond)
		first, _ := tt.Access("Video 0")
		right := first.Bounds.X
		for i := range vids {
			r, _ := tt.Access("Video " + count(i))
			if r.Bounds.Y != first.Bounds.Y {
				break
			}
			right = r.Bounds.X + r.Bounds.W
		}
		// The cards end a gutter short of the panel, or of the window.
		end := float32(1600)
		if panel != "" {
			end -= 340
		}
		if gap := end - right; gap > s.lk.gutter+30 {
			t.Errorf("with the panel %q, the cards end at %.0f, %.0fpx short", panel, right, gap)
		}
	}
}
