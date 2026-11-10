package main

import (
	"testing"
	"time"

	"github.com/bigdaddy284/mymusic/internal/app"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// TestModeSwitch switches to Video with the mark, which turns over and
// names the side, shows the Video menu and its home, and back to Music.
func TestModeSwitch(t *testing.T) {
	s := newTestShell(t)
	tt := testerOf(t, s, 1280, 820)
	if !tt.HasText("Music") || tt.HasText("trending") {
		t.Fatalf("the window doesn't start on Music: %q", tt.Texts())
	}
	click(t, s, tt, "Switch to Video")
	s.drain()
	tt.Frame()
	tt.Advance(3 * time.Second)
	if !s.mode.video || s.mode.turn != modeSpin {
		t.Fatalf("the mark didn't turn to Video: %+v", s.mode)
	}
	for _, want := range []string{"Video", "trending", "your playlists", "Search YouTube videos", "Switch to Music"} {
		if !tt.HasText(want) {
			t.Errorf("Video doesn't show %q: %q", want, tt.Texts())
		}
	}
	if got := s.router.Path(); got != "/videos" {
		t.Errorf("Video opens %q, not its home", got)
	}
	click(t, s, tt, "Switch to Music")
	s.drain()
	tt.Frame()
	tt.Advance(3 * time.Second)
	if s.mode.video || s.router.Path() != "/" || tt.HasText("trending") {
		t.Errorf("the window didn't go back to Music: %q at %q", tt.Texts(), s.router.Path())
	}
}

// TestModeFollowsPath puts the window on the side of a page opened by its
// path, and keeps the side on the pages both share.
func TestModeFollowsPath(t *testing.T) {
	s := newTestShell(t)
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, "/videos/trending")
	if !s.mode.video {
		t.Fatal("a Video page leaves the window on Music")
	}
	show(s, tt, "/settings")
	if !s.mode.video {
		t.Error("the settings switched the side")
	}
	show(s, tt, "/library")
	if s.mode.video {
		t.Error("a Music page leaves the window on Video")
	}
}

// TestVideoSearchShowsCards shows the videos a search finds as cards.
func TestVideoSearchShowsCards(t *testing.T) {
	s := newTestShell(t)
	s.loads["videos:search:frogs"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{
		{ID: "abcdefghijk", Title: "Frogs of the caatinga", Channel: "Nature Films", Duration: "14:30", Views: "48K views", Published: "4 years ago"},
	}}}
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, videoSearchPath("frogs"))
	for _, want := range []string{"Frogs of the caatinga", "Nature Films", "14:30", "48K views · 4 years ago"} {
		if !showsText(tt.Texts(), want) {
			t.Errorf("the search doesn't show %q: %q", want, tt.Texts())
		}
	}
}

// TestRandomFallsBackToAWord searches a word of the list when the title of
// the random article finds nothing.
func TestRandomFallsBackToAWord(t *testing.T) {
	s := newTestShell(t)
	randomWord = func() string { return "volcano" }
	t.Cleanup(func() { randomWord = app.RandomWord })
	s.loads["videos:search:Obscure title"] = &loaded[youtube.VideoPage]{done: true}
	s.loads["videos:search:volcano"] = &loaded[youtube.VideoPage]{done: true}
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, videoSearchPath("Obscure title")+"&random=wikipedia")
	tt.Frame()
	if s.videos.query != "volcano" || !tt.HasText("Random picked a word from its own list.") {
		t.Errorf("an empty random search doesn't go on to a word: %q, %q", s.videos.query, tt.Texts())
	}
}

// TestVideoNext plays up next in order, around with repeat all, and the
// same video again with repeat one.
func TestVideoNext(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	list := []youtube.Video{{ID: "a"}, {ID: "b"}, {ID: "c"}}
	vp.list, vp.current = list, list[1]
	s.videoNext(false)
	if vp.current.ID != "c" {
		t.Fatalf("next after b plays %q", vp.current.ID)
	}
	s.videoNext(true)
	if vp.current.ID != "c" {
		t.Errorf("the end of up next without repeat plays %q", vp.current.ID)
	}
	s.set.VideoRepeat = videoRepeatAll
	s.videoNext(true)
	if vp.current.ID != "a" {
		t.Errorf("repeat all goes on to %q, not the first", vp.current.ID)
	}
	s.set.VideoRepeat = videoRepeatOne
	s.videoNext(true)
	if vp.current.ID != "a" {
		t.Errorf("repeat one goes on to %q", vp.current.ID)
	}
}

// TestResizerClamps keeps the dragged widths within their bounds.
func TestResizerClamps(t *testing.T) {
	s := newTestShell(t)
	s.set.SidebarWidth = 9999
	if got := s.railWidth(); got != railMost {
		t.Errorf("a sidebar dragged too wide is %v wide", got)
	}
	s.set.SidebarWidth = 10
	if got := s.railWidth(); got != railLeast {
		t.Errorf("a sidebar dragged too narrow is %v wide", got)
	}
}

// TestVideoOnMusicSide keeps the video playing as the Music side shows,
// until music starts to play, which pauses it, and keeps it there.
func TestVideoOnMusicSide(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	vp.current = youtube.Video{ID: "a"}
	vp.state.Playing = true
	s.mode.video = false
	s.followWatch()
	if vp.current.ID != "a" {
		t.Fatal("the video stopped as the Music side showed")
	}
	if !s.showsMini() && s.router.Path() != "/" {
		t.Error("the video doesn't show on the Music side")
	}
	s.st.Playing = true
	s.followWatch()
	if vp.state.Playing || vp.current.ID != "a" {
		t.Errorf("as music starts, the video plays %v, or went: %q", vp.state.Playing, vp.current.ID)
	}
}

// TestVideoTakesMediaKeys gives the media keys to the video while one is
// loaded, and to the music otherwise.
func TestVideoTakesMediaKeys(t *testing.T) {
	s := newTestShell(t)
	if s.videoTakes(app.PlayerCommand{Command: app.CommandToggle}) {
		t.Error("the video took Play / Pause with no video")
	}
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	vp.list = []youtube.Video{{ID: "a"}, {ID: "b"}}
	vp.current = vp.list[0]
	s.followWatch()
	if !s.videoTakes(app.PlayerCommand{Command: app.CommandNext}) {
		t.Fatal("the video didn't take Next")
	}
	if s.videoTakes(app.PlayerCommand{Command: app.CommandLike}) {
		t.Error("the video took Like")
	}
	s.videoCommandOf(app.PlayerCommand{Command: app.CommandNext})
	if vp.current.ID != "b" {
		t.Errorf("Next plays %q", vp.current.ID)
	}
}

// TestVideoCastCommands sends the player bar's commands to the device the
// video is cast to, and follows where the video is there.
func TestVideoCastCommands(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	vp.current = youtube.Video{ID: "a"}
	vp.casting = &videoCasting{device: app.CastDevice{ID: "tv", Name: "TV"}}
	vp.casting.set(30, true)
	s.videoCommand("pause", nil)
	if vp.casting.playing {
		t.Error("Pause didn't pause the cast")
	}
	s.videoCommand("seek", 90.0)
	st := videoState{Duration: 600, Playing: true}
	s.followVideoCast(&st)
	if st.Playing || st.Time != 90 {
		t.Errorf("the cast's state: %+v", st)
	}
	s.stopWatching()
	if vp.casting != nil {
		t.Error("the cast goes on with no video")
	}
}

// TestWatchedShare fills the red bar from where a video was left here,
// and from what YouTube tells of one that wasn't.
func TestWatchedShare(t *testing.T) {
	s := newTestShell(t)
	if got := s.watchedShare(youtube.Video{ID: "aaaaaaaaaaa", Watched: 40}); got != 0.4 {
		t.Errorf("YouTube's 40%% fills %v", got)
	}
	s.svc.Videos.LeftAt("aaaaaaaaaaa", 30, 120)
	if got := s.watchedShare(youtube.Video{ID: "aaaaaaaaaaa", Watched: 40}); got != 0.25 {
		t.Errorf("left at 30 s of 120 fills %v", got)
	}
	if got := s.watchedShare(youtube.Video{ID: "bbbbbbbbbbb"}); got != 0 {
		t.Errorf("a video never watched fills %v", got)
	}
}

// TestRandomSkipsShownVideos leaves out of a random search the videos an
// earlier random search showed.
func TestRandomSkipsShownVideos(t *testing.T) {
	s := newTestShell(t)
	s.loads["videos:search:Frogs"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{
		{ID: "popular0001", Title: "A popular video"}, {ID: "frogs000001", Title: "Frogs singing"},
	}}}
	s.loads["videos:search:Toads"] = &loaded[youtube.VideoPage]{done: true, value: youtube.VideoPage{Videos: []youtube.Video{
		{ID: "popular0001", Title: "A popular video"}, {ID: "toads000001", Title: "Toads singing"},
	}}}
	tt := testerOf(t, s, 1280, 820)
	show(s, tt, videoSearchPath("Frogs")+"&random=wikipedia")
	if !showsText(tt.Texts(), "A popular video") {
		t.Fatalf("the first search doesn't show its videos: %q", tt.Texts())
	}
	show(s, tt, videoSearchPath("Toads")+"&random=wikipedia")
	if showsText(tt.Texts(), "A popular video") || !showsText(tt.Texts(), "Toads singing") {
		t.Errorf("the second search shows %q", tt.Texts())
	}
}

// TestVideoSideBackToThePlayingVideo goes back to the watch page of the
// video playing now, though the video moved on while Music showed.
func TestVideoSideBackToThePlayingVideo(t *testing.T) {
	s := newTestShell(t)
	vp := &s.videos.player
	vp.problem = "no web view in tests"
	list := []youtube.Video{{ID: "first000001", Title: "First"}, {ID: "second00001", Title: "Second"}}
	tt := testerOf(t, s, 1280, 820)
	s.watch(list[0], list)
	tt.Frame()
	s.toggleMode()
	tt.Frame()
	s.videoNext(true)
	s.toggleMode()
	tt.Frame()
	if vp.current.ID != "second00001" || s.router.Path() != watchPath("second00001") {
		t.Errorf("back on Video, %q plays at %q", vp.current.ID, s.router.Path())
	}
}
