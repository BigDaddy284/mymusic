package main

import (
	"slices"
	"testing"

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
