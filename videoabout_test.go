package main

import (
	"slices"
	"testing"

	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestVideoAboutPanel opens the Description panel from the Video side's
// player bar: it shows the views, date and description of the video that
// plays, and its links work. A web page opens in the browser, another
// video plays here, and a channel opens its page.
func TestVideoAboutPanel(t *testing.T) {
	var opened []string
	was := openWebPage
	openWebPage = func(_ *shell, u string) { opened = append(opened, u) }
	t.Cleanup(func() { openWebPage = was })
	s := newTestShell(t)
	frogs := youtube.Video{ID: "frogs000001", Title: "Frogs singing", Channel: "Pond"}
	s.loads["videos:details:"+frogs.ID] = &loaded[youtube.VideoDetails]{done: true, value: youtube.VideoDetails{
		Views: "541 views",
		About: youtube.VideoAbout{Views: "541 views", Date: "11 Nov 2026", Description: [][]media.TextRun{
			{{Text: "Resources: "}, {Text: "https://example.com/r", Link: "https://example.com/r"}},
			{{Text: "Featured — "}, {Text: "• Toads at dusk...", Link: "https://www.youtube.com/watch?v=toads000001"}},
			{{Text: "By "}, {Text: "@pond", Link: "https://www.youtube.com/channel/UCpond"}},
		}},
	}}
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	tt := testerOf(t, s, 1400, 900)
	s.mode.video = true
	s.watch(frogs, nil)
	s.drain()
	tt.Frame()
	click(t, s, tt, "Description")
	if s.panel != "about" || !tt.HasText("Close description") {
		t.Fatalf("the Description panel didn't open: %q", tt.Texts())
	}
	if !showsText(tt.Texts(), "541 views") || !showsText(tt.Texts(), "11 Nov 2026") || !showsText(tt.Texts(), "Resources: ") {
		t.Errorf("the panel shows %q", tt.Texts())
	}
	click(t, s, tt, "https://example.com/r")
	if !slices.Equal(opened, []string{"https://example.com/r"}) {
		t.Errorf("opened in the browser %q", opened)
	}
	click(t, s, tt, "@pond")
	if got := s.router.Path(); got != videoChannelPath("UCpond") {
		t.Errorf("the channel's link went to %q", got)
	}
	s.navigate(watchPath(frogs.ID))
	s.drain()
	tt.Frame()
	click(t, s, tt, "• Toads at dusk...")
	if vp.current.ID != "toads000001" || vp.current.Title != "Toads at dusk..." || s.router.Path() != watchPath("toads000001") {
		t.Errorf("the video's link plays %+v, on %q", vp.current, s.router.Path())
	}
	if len(opened) != 1 {
		t.Errorf("a YouTube link opened in the browser: %q", opened)
	}
}

// TestYoutubeLinkOf reads what links to YouTube open here.
func TestYoutubeLinkOf(t *testing.T) {
	for raw, want := range map[string]youtubeLink{
		"https://www.youtube.com/watch?v=abc&t=90s": {video: "abc", at: 90},
		"https://youtu.be/abc?t=1m5s":               {video: "abc", at: 65},
		"https://m.youtube.com/shorts/abc":          {video: "abc"},
		"https://www.youtube.com/channel/UCx":       {channel: "UCx"},
		"https://www.youtube.com/playlist?list=PL1": {playlist: "PL1"},
		"https://www.youtube.com/hashtag/frogs":     {search: "#frogs"},
		"https://music.youtube.com/watch?v=abc":     {video: "abc"},
	} {
		if got, ok := youtubeLinkOf(raw); !ok || got != want {
			t.Errorf("%s reads %+v, %v", raw, got, ok)
		}
	}
	for _, raw := range []string{"https://example.com/watch?v=abc", "https://www.youtube.com/@pond", "mailto:a@b.c"} {
		if _, ok := youtubeLinkOf(raw); ok {
			t.Errorf("%s opens here", raw)
		}
	}
	if linkTitle("  • Toads at dusk") != "Toads at dusk" || linkTitle("https://youtu.be/x") != "" {
		t.Error("the titles of links")
	}
}
