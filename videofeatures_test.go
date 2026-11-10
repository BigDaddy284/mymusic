package main

import (
	"slices"
	"testing"

	"github.com/bigdaddy284/mymusic/internal/app"
	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/player"
	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestBlockedVideosLeftOut leaves blocked channels, blocked words and
// Shorts off the pages, and Shorts back on when the settings show them.
func TestBlockedVideosLeftOut(t *testing.T) {
	s := newTestShell(t)
	s.set.BlockedChannels = []string{"prank central"}
	s.set.BlockedWords = []string{"REACTION"}
	list := []youtube.Video{
		{ID: "keep", Title: "A good video", Channel: "Good"},
		{ID: "chan", Title: "Fine", Channel: "Prank Central"},
		{ID: "word", Title: "My reaction to it", Channel: "Good"},
		{ID: "short", Title: "Quick", Channel: "Good", Short: true},
	}
	ids := func(vs []youtube.Video) []string {
		var out []string
		for _, v := range vs {
			out = append(out, v.ID)
		}
		return out
	}
	if got := ids(s.shownVideos(slices.Clone(list))); !slices.Equal(got, []string{"keep"}) {
		t.Errorf("shown %v", got)
	}
	s.set.ShowShorts = true
	if got := ids(s.shownVideos(slices.Clone(list))); !slices.Equal(got, []string{"keep", "short"}) {
		t.Errorf("with Shorts shown %v", got)
	}
	s.blockChannel(youtube.Video{Channel: "Good"})
	if got := s.shownVideos(slices.Clone(list)); len(got) != 0 {
		t.Errorf("after blocking Good: %v", ids(got))
	}
}

// TestSegmentSkipped skips a sponsor's segment once, with an undo, and
// mutes one the settings mute, giving the sound back after it.
func TestSegmentSkipped(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.current = youtube.Video{ID: "v1", Title: "V"}
	vp.settled = 2
	s.loads["videos:segments:v1"] = &loaded[[]app.Segment]{done: true, value: []app.Segment{
		{Start: 10, End: 40, Category: "sponsor", UUID: "s"},
		{Start: 100, End: 120, Category: "selfpromo", UUID: "p"},
	}}
	vp.state = videoState{Time: 12, Duration: 300, Playing: true}
	s.followSegments()
	if vp.state.Time != 40 || len(s.actionToasts) != 1 || s.actionToasts[0].message != "Skipped: Sponsor" {
		t.Fatalf("after the sponsor: at %v, toasts %+v", vp.state.Time, s.actionToasts)
	}
	// Back into it, as the undo does, it plays.
	s.actionToasts[0].run()
	vp.state.Time = 12
	s.followSegments()
	if vp.state.Time != 12 || len(s.actionToasts) != 1 {
		t.Errorf("skipped again: at %v", vp.state.Time)
	}
	s.set.Segments = map[string]string{"selfpromo": settings.SegmentMute}
	vp.state.Time = 105
	s.followSegments()
	if vp.segments.muting != "p" || vp.state.Time != 105 {
		t.Errorf("not muted: %+v at %v", vp.segments, vp.state.Time)
	}
	vp.state.Time = 125
	s.followSegments()
	if vp.segments.muting != "" {
		t.Errorf("still muted after the segment")
	}
	// Off, nothing skips.
	s.set.SponsorBlockOff = true
	vp.segments = segmentState{}
	vp.state.Time = 12
	s.followSegments()
	if vp.state.Time != 12 {
		t.Errorf("skipped with SponsorBlock off")
	}
}

// TestVideoChannelPage shows a channel's name, subscribers, videos and a
// subscribe button, and a card's channel opens it.
func TestVideoChannelPage(t *testing.T) {
	s := newTestShell(t)
	s.loads["videos:channel:UC1"] = &loaded[youtube.VideoChannel]{done: true, value: youtube.VideoChannel{
		ID: "UC1", Name: "Some Channel", Handle: "@some", Subscribers: "1.2M subscribers",
		Videos: youtube.VideoPage{Videos: []youtube.Video{{ID: "v1", Title: "Latest upload", Channel: "Some Channel", ChannelID: "UC1"}}},
	}}
	tt := testerOf(t, s, 1280, 1000)
	show(s, tt, videoChannelPath("UC1"))
	for _, want := range []string{"Some Channel", "@some · 1.2M subscribers", "Latest upload", "subscribe"} {
		if !showsText(tt.Texts(), want) {
			t.Errorf("no %q: %q", want, tt.Texts())
		}
	}
	s.loads["videos:search:x"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{
		{ID: "v2", Title: "Found", Channel: "Other Channel", ChannelID: "UC2"}}}}
	show(s, tt, videoSearchPath("x"))
	if err := tt.Click("Open the channel Other Channel"); err != nil {
		t.Fatalf("no channel link: %v: %q", err, tt.Texts())
	}
	s.drain()
	tt.Frame()
	if p := s.router.Path(); p != videoChannelPath("UC2") {
		t.Errorf("the channel opened %q", p)
	}
}

// TestSearchBothSides shows YouTube Music's songs under a video search's
// Songs chip.
func TestSearchBothSides(t *testing.T) {
	s := searchVideosShell(t)
	s.loads["videos:search:boiler room"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{
		{ID: "v9", Title: "Boiler Room on YouTube", Channel: "BR"}}}}
	s.videos.searchChip = 1
	tt := testerOf(t, s, 1280, 1400)
	show(s, tt, videoSearchPath("boiler room"))
	if !showsText(tt.Texts(), "Boiler Room Anthem") {
		t.Errorf("the Songs chip shows no song: %q", tt.Texts())
	}
	// The Music side's search shows YouTube's videos under All.
	show(s, tt, searchPath("boiler room"))
	if !showsText(tt.Texts(), "On YouTube") || !showsText(tt.Texts(), "Boiler Room on YouTube") {
		t.Errorf("the music search shows no YouTube videos: %q", tt.Texts())
	}
}

// TestPlayVideoAsAudio plays a video's sound on the music side, as a
// track by its channel.
func TestPlayVideoAsAudio(t *testing.T) {
	s := newTestShell(t)
	v := youtube.Video{ID: "abcdefghijk", Title: "A talk", Channel: "Talks", ChannelID: "UC9", Duration: "1:02:03"}
	s.playVideoAudio(v)
	s.settle(t, func(st player.State) bool { _, ok := st.Current(); return ok })
	cur, _ := s.p.State().Current()
	tr := cur.Track
	if tr.ID != media.ID(media.SourceYouTube, v.ID) || tr.Title != "A talk" || len(tr.Artists) != 1 || tr.Artists[0].Name != "Talks" ||
		tr.DurationMs != 3723000 || !tr.Video {
		t.Errorf("the track is %+v", tr)
	}
}
