package main

// Videos play in YouTube's embedded player, public/watch.html, in a web
// view of the main window: MyGo draws no video. The web view stays while
// the Video side shows: the watch page shows it large, other pages small,
// at the foot of the Up next panel, or in a corner of the page, so the
// video goes on while the user browses; Full App fills the window with it.
// The player bar of the Video side controls it, through YouTube's player
// API (or, on YouTube's own watch page, which plays the videos that can't
// be embedded, its <video>), which the window asks for the player's state
// twice a second.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/url"
	"strconv"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/internal/youtube"
)

// videoPoll is how often the window asks the player for its state.
const videoPoll = 500 * time.Millisecond

// videoStatus is what the player tells of the video it plays.
type videoState struct {
	Time     float64 `json:"t"`
	Duration float64 `json:"d"`
	Playing  bool    `json:"playing"`
	Ended    bool    `json:"ended"`
	Rate     float64 `json:"rate"`
	Volume   float64 `json:"volume"`
	Muted    bool    `json:"muted"`
	Captions bool    `json:"captions"`
	// Esc is set once Escape was pressed in the page.
	Esc bool `json:"esc"`
	// FS is set while the player's own full screen button has the page
	// fill the screen.
	FS bool `json:"fs"`
}

// Repeat modes of the videos.
const (
	videoRepeatOff = iota
	videoRepeatAll
	videoRepeatOne
)

// videoPlayer is the player's state, which belongs to the main thread.
type videoPlayer struct {
	// web is the web view, once made; making is set while it is made, and
	// problem tells why it couldn't be.
	web     *mygo.WebView
	making  bool
	problem string
	// current is the video that plays, and loaded the one the web view
	// was told to load; list is the list it was chosen from: up next.
	current youtube.Video
	loaded  string
	list    []youtube.Video
	// state is what the player told last, and stop ends its polling.
	state videoState
	stop  chan struct{}
	// endedOf is the video whose end was taken in, once.
	endedOf string
	// settled counts the polls since the video loaded that brought the
	// player the user's volume, speed and captions: YouTube's player turns
	// captions on by itself, as the account's own settings ask, after it
	// starts.
	settled    int
	settledFor string
	// fullApp fills the window with the video; fullScreen, the screen,
	// as the player's own button does (fullscreen.go).
	fullApp    bool
	fullScreen bool
	// ownsScreen is set while the window is full screen for the video,
	// takesEscape while Escape is the app's, and blurred while the
	// window isn't in front.
	ownsScreen  bool
	takesEscape bool
	blurred     bool
	// casting is the cast of the videos, nil while they play here
	// (videocast.go).
	casting *videoCasting
	// musicPlayed is whether music played in the last frame: the video
	// stops as music starts.
	musicPlayed bool
	// pausedForMusic is when music starting paused the video last.
	pausedForMusic time.Time
	// seeking and seekTo follow the seek bar while it's dragged; volume
	// follows the volume's slider.
	seeking bool
	seekTo  float64
	volume  float64
	// segments follows SponsorBlock's segments of the video (segments.go).
	segments segmentState
	// pop is the pop-out window the video plays in (popout.go).
	pop popOut
}

// watch plays v, with list as up next, on the watch page.
func (s *shell) watch(v youtube.Video, list []youtube.Video) {
	vs := &s.videos
	if vs.known == nil {
		vs.known = map[string]youtube.Video{}
	}
	for _, x := range list {
		vs.known[x.ID] = x
	}
	vs.known[v.ID] = v
	if len(list) > 0 {
		vs.player.list = list
	}
	s.playVideo(v)
	s.navigate(watchPath(v.ID))
}

// playVideo has the player play v.
func (s *shell) playVideo(v youtube.Video) {
	vp := &s.videos.player
	vp.current = v
	vp.state = videoState{}
	vp.endedOf = ""
	if s.st.Playing {
		s.p.Pause()
	}
	s.loadVideo()
	// A video chosen while the videos are cast plays on the device, from
	// where it was left.
	if vp.casting != nil {
		at := s.svc.Videos.Position(v.ID)
		if at <= videoResumeAfter || s.svc.Videos.WatchedShare(v.ID) > 0.97 {
			at = 0
		}
		s.castCurrent(at)
	}
}

// loadVideo makes the web view, or has it load the video that plays.
func (s *shell) loadVideo() {
	vp := &s.videos.player
	id := vp.current.ID
	if id == "" || vp.making || vp.problem != "" {
		return
	}
	if vp.web == nil {
		if s.win == nil {
			return
		}
		if problem := s.videoProblem(); problem != "" {
			vp.problem = problem
			return
		}
		vp.making = true
		win := s.win
		go func() {
			web, err := win.NewWebView(mygo.WebViewOptions{URL: watchPageURL(id, s.set.VideoCaptions), BackgroundColor: "#000000"})
			s.update(func() {
				vp.making = false
				if err != nil {
					log.Printf("mytube: the video player can't open: %v", err)
					vp.problem = "the video player couldn't open."
					return
				}
				vp.web, vp.loaded = web, id
				vp.stop = make(chan struct{})
				go s.pollVideo(web, vp.stop)
				if vp.current.ID != id {
					s.loadVideo()
				}
			})
		}()
		return
	}
	if vp.loaded == id {
		return
	}
	vp.loaded = id
	web, captions := vp.web, s.set.VideoCaptions
	go func() {
		if err := web.Page().LoadURL(watchPageURL(id, captions)); err != nil {
			log.Printf("mytube: playing the video %s: %v", id, err)
		}
	}()
}

// watchPageURL is the player's page, playing id, with the captions on or
// off: the page keeps them so all through the video, whatever YouTube's
// account settings ask.
func watchPageURL(id string, captions bool) string {
	cc := "0"
	if captions {
		cc = "1"
	}
	return "/watch.html?v=" + url.QueryEscape(id) + "&cc=" + cc
}

// stopWatching closes the player: the video stops.
func (s *shell) stopWatching() {
	vp := &s.videos.player
	if vp.stop != nil {
		close(vp.stop)
		vp.stop = nil
	}
	if vp.casting != nil {
		vp.casting = nil
		go s.svc.Videos.StopVideoCast(context.Background())
	}
	if web := vp.web; web != nil {
		go web.Destroy()
	}
	if win := vp.pop.win; win != nil {
		vp.pop.win = nil
		win.Close()
	}
	vp.web, vp.loaded, vp.current, vp.state, vp.fullApp, vp.fullScreen = nil, "", youtube.Video{}, videoState{}, false, false
	// The watching's last report goes.
	s.svc.Videos.ReportWatching("", 0, false, 1, false)
	s.followFullScreen()
}

// videoStateJS asks the player for its state, the embedded one's or, on
// YouTube's own watch page, its <video>'s.
const videoStateJS = `
const p = window.mtPlayer && window.mtPlayer.getPlayerState ? window.mtPlayer : null;
const esc = !!window.mtEsc;
window.mtEsc = false;
const fs = !!document.fullscreenElement;
if (p) {
  const st = p.getPlayerState();
  let caps = false;
  try { caps = (p.getOptions() || []).indexOf('captions') >= 0 && !!(p.getOption('captions', 'track') || {}).languageCode; } catch (e) {}
  return {t: p.getCurrentTime() || 0, d: p.getDuration() || 0, playing: st === 1 || st === 3, ended: st === 0,
    rate: p.getPlaybackRate() || 1, volume: p.getVolume(), muted: p.isMuted(), captions: caps, esc, fs};
}
const v = document.querySelector('video');
if (!v) return {t: 0, d: 0, playing: false, ended: false, rate: 1, volume: 100, muted: false, captions: false, esc, fs};
const cc = document.querySelector('.ytp-subtitles-button');
return {t: v.currentTime || 0, d: isFinite(v.duration) ? v.duration : 0, playing: !v.paused && !v.ended, ended: v.ended,
  rate: v.playbackRate || 1, volume: Math.round(v.volume * 100), muted: v.muted,
  captions: !!cc && cc.getAttribute('aria-pressed') === 'true', esc, fs};
`

// pollVideo asks the player in web for its state until stop closes.
func (s *shell) pollVideo(web *mygo.WebView, stop chan struct{}) {
	t := time.NewTicker(videoPoll)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
		}
		if web.IsDestroyed() {
			return
		}
		st, err := mygo.EvalAs[videoState](web.Page(), videoStateJS)
		if err != nil {
			continue // the page loads
		}
		s.update(func() {
			vp := &s.videos.player
			if vp.web != web {
				return
			}
			// Watching counts in the stats while the video plays; repeat
			// one, back at its start, counts as another view.
			if st.Time < vp.state.Time-5 && st.Time < 3 {
				s.svc.Videos.Started()
			}
			if st.Playing {
				v := vp.current
				go s.svc.Videos.Watched(v, videoPoll.Seconds())
			}
			s.followVideoCast(&st)
			s.videoStarts(st)
			vp.state = st
			// YouTube learns of the watching from the app: the embedded
			// player plays signed out. A cast's device tells it itself.
			if vp.casting == nil {
				s.svc.Videos.ReportWatching(vp.current.ID, st.Time, st.Playing, st.Rate, st.Duration <= 0)
			}
			s.settleVideo()
			s.followSegments()
			s.keepVideoPlace()
			if st.Esc {
				s.leaveFullApp()
			}
			s.followFullScreen()
			s.followPopScreen()
			if st.Ended && vp.endedOf != vp.current.ID {
				vp.endedOf = vp.current.ID
				s.videoNext(true)
			}
		})
	}
}

// videoStarts pauses the music when the player tells the video started to
// play, as from YouTube's own button in it, as the bar's button does: but
// for the moments after music paused the video, while the player still
// tells it plays.
func (s *shell) videoStarts(st videoState) {
	vp := &s.videos.player
	if st.Playing && !vp.state.Playing && s.st.Playing && time.Since(vp.pausedForMusic) > 2*time.Second {
		s.p.Pause()
		s.st.Playing = false
	}
}

// videoCommand runs a command of the player, as "play", with its argument.
func (s *shell) videoCommand(cmd string, arg any) {
	if s.castCommand(cmd, arg) {
		return
	}
	s.localVideoCommand(cmd, arg)
}

// localVideoCommand runs a command of the window's player, as "play",
// with its argument, whether or not the video is cast.
func (s *shell) localVideoCommand(cmd string, arg any) {
	web := s.videos.player.web
	if web == nil {
		return
	}
	a, _ := json.Marshal(arg)
	code := fmt.Sprintf(videoCommandJS, strconv.Quote(cmd), a)
	go func() {
		if _, err := web.Page().Eval(code); err != nil {
			log.Printf("mytube: the video player's %s: %v", cmd, err)
		}
	}()
}

// videoCommandJS runs a command of the player: %s is its name, %s its
// argument, in JSON.
const videoCommandJS = `
const cmd = %s, arg = %s;
if (cmd === 'exitFullScreen') {
  if (document.fullscreenElement) document.exitFullscreen();
  return true;
}
const p = window.mtPlayer && window.mtPlayer.getPlayerState ? window.mtPlayer : null;
if (p) {
  switch (cmd) {
    case 'play': p.playVideo(); break;
    case 'pause': p.pauseVideo(); break;
    case 'seek': p.seekTo(arg, true); break;
    case 'rate': p.setPlaybackRate(arg); break;
    case 'volume': p.setVolume(arg); if (arg > 0) p.unMute(); break;
    case 'mute': if (arg) p.mute(); else p.unMute(); break;
    case 'captions':
      window.mtCaptions = !!arg;
      if (arg) { p.loadModule('captions'); p.setOption('captions', 'reload', true); }
      else { p.unloadModule('captions'); }
      break;
  }
  return true;
}
const v = document.querySelector('video');
if (!v) return false;
switch (cmd) {
  case 'play': v.play(); break;
  case 'pause': v.pause(); break;
  case 'seek': v.currentTime = arg; break;
  case 'rate': v.playbackRate = arg; break;
  case 'volume': v.volume = arg / 100; if (arg > 0) v.muted = false; break;
  case 'mute': v.muted = !!arg; break;
  case 'captions': { const b = document.querySelector('.ytp-subtitles-button'); if (b && (b.getAttribute('aria-pressed') === 'true') !== !!arg) b.click(); break; }
}
return true;
`

const (
	// videoResumeAfter and videoResumeBefore bound where a video goes on
	// from: not in its first seconds, nor its last.
	videoResumeAfter  = 10
	videoResumeBefore = 15
)

// keepVideoPlace keeps where the video is, to go on from there when it
// plays again; one watched to its end, or gone back to its start, starts
// over. It waits for the player to have gone on from where the video was
// left (settleVideo).
func (s *shell) keepVideoPlace() {
	vp := &s.videos.player
	st := vp.state
	if vp.current.ID == "" || st.Duration <= 0 || vp.settled < 4 {
		return
	}
	// One watched to its end keeps its end, which its red bar fills, and
	// starts over.
	at := st.Time
	switch {
	case st.Ended || at > st.Duration-videoResumeBefore:
		at = st.Duration
	case at < videoResumeAfter:
		at = 0
	}
	id := vp.current.ID
	s.svc.Videos.LeftAt(id, at, st.Duration)
}

// settleVideo brings the player the user's volume, speed and captions,
// kept in the settings: once the video's length is known, and again a
// few polls later, as YouTube's player may turn captions on by itself.
func (s *shell) settleVideo() {
	vp := &s.videos.player
	if vp.settledFor != vp.current.ID {
		vp.settledFor, vp.settled = vp.current.ID, 0
	}
	if vp.state.Duration <= 0 || vp.settled > 6 {
		return
	}
	vp.settled++
	w := s.set.Watching
	switch vp.settled {
	case 1:
		s.localVideoCommand("volume", w.VideoVolume)
		s.localVideoCommand("mute", w.VideoMuted)
		s.localVideoCommand("rate", w.VideoRate)
		vp.state.Volume, vp.state.Muted, vp.state.Rate = w.VideoVolume, w.VideoMuted, w.VideoRate
		// The video goes on from where it was left, as tracks do.
		// From where it was left here, or, when it wasn't, where YouTube
		// tells it was left, on any device.
		at := s.svc.Videos.Position(vp.current.ID)
		if at == 0 && vp.current.Watched > 0 {
			at = vp.state.Duration * float64(vp.current.Watched) / 100
		}
		if at > videoResumeAfter && at < vp.state.Duration-videoResumeBefore {
			s.localVideoCommand("seek", at)
			vp.state.Time = at
		}
		fallthrough
	case 3, 6:
		// The player's sound plays in a process of its own, which starts
		// with the first sound.
		if vp.settled == 3 && s.set.VideoOutput.ID != "" {
			s.applyVideoOutput()
		}
		if !w.VideoCaptions || vp.state.Captions != w.VideoCaptions {
			s.localVideoCommand("captions", w.VideoCaptions)
			vp.state.Captions = w.VideoCaptions
		}
	}
}

// videoPlayPause plays or pauses the video.
func (s *shell) videoPlayPause() {
	vp := &s.videos.player
	if vp.state.Playing {
		s.videoCommand("pause", nil)
		vp.state.Playing = false
		return
	}
	if vp.state.Ended {
		s.videoCommand("seek", 0)
	}
	s.videoCommand("play", nil)
	vp.state.Playing = true
	if s.st.Playing {
		s.p.Pause()
	}
}

// videoNext plays the next video of up next: a random other one with
// shuffle on. auto is set as a video ended, which repeat one plays again;
// at the end of the list, repeat all goes back to its start.
func (s *shell) videoNext(auto bool) {
	vp := &s.videos.player
	if auto && s.set.VideoRepeat == videoRepeatOne {
		vp.endedOf = ""
		s.videoCommand("seek", 0)
		s.videoCommand("play", nil)
		return
	}
	i := videoIndex(vp.list, vp.current.ID)
	next := -1
	switch {
	case len(vp.list) == 0:
	case s.set.VideoShuffle && len(vp.list) > 1:
		for next = i; next == i; {
			next = rand.IntN(len(vp.list))
		}
	case i+1 < len(vp.list):
		next = i + 1
	case s.set.VideoRepeat == videoRepeatAll:
		next = 0
	}
	if next < 0 {
		return
	}
	s.playNext(vp.list[next])
}

// videoPrevious goes back to the start of the video, after its first few
// seconds, or else to the video before it in up next.
func (s *shell) videoPrevious() {
	vp := &s.videos.player
	i := videoIndex(vp.list, vp.current.ID)
	if vp.state.Time > 3 || i <= 0 {
		s.videoCommand("seek", 0)
		return
	}
	s.playNext(vp.list[i-1])
}

// playNext plays v, and follows it with the watch page when that shows.
func (s *shell) playNext(v youtube.Video) {
	onWatch := len(s.router.Path()) > len("/watch/") && s.router.Path()[:len("/watch/")] == "/watch/"
	s.playVideo(v)
	if onWatch {
		s.router.Replace(watchPath(v.ID))
	}
}

// videoIndex is the index of id in list, or -1.
func videoIndex(list []youtube.Video, id string) int {
	for i, v := range list {
		if v.ID == id {
			return i
		}
	}
	return -1
}

// placeVideo shows the player's web view in a frame of w × h, or what
// stands for it while it opens or where it can't.
func (s *shell) placeVideo(c *ui.Context, key string, w, h float32) {
	k := &s.lk
	vp := &s.videos.player
	white := ui.RGBA(255, 255, 255, 1)
	ui.Box(c.Key(key)).Size(w, h).Shrink(0).Radius(k.radiusSm).Clip().Background(ui.RGBA(0, 0, 0, 1)).Children(func() {
		switch {
		case vp.problem != "":
			ui.Column(c).Size(w, h).AlignItems(ui.Center).Justify(ui.Center).Gap(12).Padding(16).Children(func() {
				ui.Text(c, "Videos can't play inside the app here: "+vp.problem).TextColor(white).TextAlign(ui.Center)
				if s.button(c, key+"-on-youtube", "watch on YouTube", true).Clicked() {
					go s.svc.System.OpenExternal(youtubeWatchURL(vp.current.ID))
				}
			})
		case vp.popped():
			s.poppedOutFrame(c, key, w, h)
		case vp.web == nil:
			s.loadVideo()
			ui.Row(c).Size(w, h).AlignItems(ui.Center).Justify(ui.Center).Children(func() { ui.Spinner(c).Label("Loading the video") })
		default:
			ui.WebView(c.Key(key+"-web"), vp.web).Absolute().Left(0).Top(0).Size(w, h).Label("Video: " + vp.current.Title)
		}
	})
}

// watchPage plays the video id as large as the pane lets it, with up
// next under it while the side panel is closed.
func (s *shell) watchPage(c *ui.Context, id string) {
	k := &s.lk
	vs := &s.videos
	vp := &vs.player
	if vp.current.ID != id {
		v, ok := vs.known[id]
		if !ok {
			v = youtube.Video{ID: id}
		}
		s.playVideo(v)
	}
	s.pageScroll(c, func() {
		// The video fills the middle of the window. YouTube's player shows
		// its title and channel, and keeps the picture's shape, with
		// black around it.
		_, winH := c.Size()
		w := s.paneWidth(c) - 2*k.gutter
		h := float32(int(max(240, winH-k.barRoom()-48)))
		ui.Box(c).Margin(16, 0, 32).Children(func() {
			if vp.filling() {
				ui.Box(c).Size(w, h).Radius(k.radiusSm).Background(ui.RGBA(0, 0, 0, 1))
				return
			}
			s.placeVideo(c, "watch-frame", w, h)
		})
		if s.panel == "" {
			var next []youtube.Video
			for _, x := range vp.list {
				if x.ID != id {
					next = append(next, x)
				}
			}
			if len(next) > 0 {
				s.sectionTitle(c, "Up next")
				s.videoGrid(c, "up-next", next)
			}
		}
	})
}

// fullAppView fills the window with the video, under a thin strip that
// leaves it, as Escape does, when Full App is on: it tells whether it
// built the window.
func (s *shell) fullAppView(c *ui.Context) bool {
	vp := &s.videos.player
	if !vp.filling() || !s.mode.video || vp.current.ID == "" {
		vp.fullApp = false
		return false
	}
	k := &s.lk
	exit := s.leaveFullApp
	root := ui.Column(c).Fill().Background(ui.RGBA(0, 0, 0, 1)).TextColor(ui.RGBA(255, 255, 255, 1)).Font(k.body)
	root.OnShortcut(0, ui.KeyEscape, exit)
	root.Children(func() {
		if vp.videoFullScreen() {
			// Full screen, the player's own: the video alone.
			if vp.web != nil {
				w, h := c.Size()
				ui.WebView(c.Key("full-app-web"), vp.web).Size(w, h).Label("Video: " + vp.current.Title)
			}
			return
		}
		bar := c.TitleBar()
		strip := ui.Row(c).Height(max(bar.Height, 32)).Shrink(0).Padding(0, bar.Right+12, 0, bar.Left+12).Gap(12)
		if bar.Height > 0 {
			strip.DragWindow()
		}
		strip.Children(func() {
			ui.Text(c, vp.current.Title).FontSize(k.smallSize()).TextColor(ui.RGBA(255, 255, 255, 0.8)).SingleLine().Grow(1).MinWidth(0)
			b := s.buttonBase(c.Key("exit-full-app")).Height(26).PaddingX(10).Radius(4).Label("Exit Full App").Cursor(ui.CursorPointer).
				Background(ui.RGBA(255, 255, 255, 0.14))
			if b.Clicked() {
				exit()
			}
			b.Children(func() {
				ui.Text(c, "Exit Full App (Esc)").FontSize(k.captionSize()).TextColor(ui.RGBA(255, 255, 255, 1)).SingleLine()
			})
		})
		w, h := c.Size()
		h -= max(bar.Height, 32)
		if vp.web != nil {
			ui.WebView(c.Key("full-app-web"), vp.web).Size(w, h).Label("Video: " + vp.current.Title)
		}
	})
	return true
}

// videoMini shows the video that plays, w wide, while the page shown is
// another than its own, so it goes on as the user browses; its button
// opens its page again.
func (s *shell) videoMini(c *ui.Context, key string, w float32) {
	k := &s.lk
	vp := &s.videos.player
	h := float32(int(w*9/16 + 0.5))
	ui.Column(c.Key(key)).Width(w).Gap(6).Shrink(0).Children(func() {
		ui.Box(c).Size(w, h).Children(func() {
			s.placeVideo(c, key+"-video", w, h)
		})
		ui.Row(c).Gap(8).Children(func() {
			ui.Text(c, vp.current.Title).FontSize(k.smallSize()).SingleLine().Grow(1).MinWidth(0)
			if s.iconButton(c, key+"-open", "external", "Back to the video", false, 16).Clicked() {
				s.navigate(watchPath(vp.current.ID))
			}
			if s.iconButton(c, key+"-close", "close", "Stop the video", false, 16).Clicked() {
				s.stopWatching()
			}
		})
	})
}

// showsMini reports whether the video that plays shows small: on any
// page but its own, and Music's home, which shows it in a tile.
func (s *shell) showsMini() bool {
	vp := &s.videos.player
	p := s.router.Path()
	if vp.current.ID == "" || vp.filling() {
		return false
	}
	if s.mode.video {
		return p != watchPath(vp.current.ID)
	}
	return p != "/" || s.home.editing
}

// floatingMini shows the video that plays in the corner of the page: on
// the Video side when the side panel, which shows it otherwise, is
// closed, and on the Music side, whose side panel is the music's.
func (s *shell) floatingMini(c *ui.Context) {
	if !s.showsMini() || s.mode.video && s.panel == "queue" {
		return
	}
	k := &s.lk
	// On the Music side it goes at the top right, beside the page's
	// title; on the Video side, at the bottom right.
	box := ui.Box(c).Absolute().Right(20).Padding(10).Radius(k.radiusSm).Background(k.surface2).
		Shadow(0, 12, 32, -8, ui.RGBA(0, 0, 0, 0.55))
	w := float32(320)
	if s.mode.video {
		box.Bottom(16)
	} else {
		box.Top(8)
		w = 280
	}
	box.Children(func() { s.videoMini(c, "mini-float", w) })
}

// videoQueuePanel is the Up next panel of the Video side: the video that
// plays, the list it came from, and the video itself, small, while another
// page shows.
func (s *shell) videoQueuePanel(c *ui.Context) {
	k := &s.lk
	vp := &s.videos.player
	s.panelHead(c, "Up next", func() {
		if len(vp.list) > 0 {
			b := s.buttonBase(c.Key("clear-videos")).Label("Clear up next").Cursor(ui.CursorPointer)
			if b.Clicked() {
				vp.list = nil
			}
			b.Children(func() { ui.Text(c, "clear").FontSize(k.smallSize()).TextColor(k.accent) })
		}
	})
	ui.Scroll(c).Grow(1).MinHeight(0).Padding(0, 12, 12).Children(func() {
		if vp.current.ID == "" && len(vp.list) == 0 {
			s.emptyIn(c, "Nothing plays yet. Pick a video, or press Random.")
			return
		}
		if vp.current.ID != "" {
			ui.Text(c, "Now playing").FontSize(k.smallSize()).TextColor(k.accent).Margin(0, 12, 6)
			s.videoRow(c, "now", vp.current, true)
		}
		for i, v := range vp.list {
			if v.ID == vp.current.ID {
				continue
			}
			s.videoRow(c, "next-"+count(i), v, false)
		}
	})
	if s.showsMini() {
		w := s.panelWidth(c) - 32
		ui.Box(c).Padding(12, 16, 16).Shrink(0).BorderWidth(1, 0, 0, 0).BorderColor(k.border).Children(func() {
			s.videoMini(c, "mini-panel", w)
		})
	}
}

// videoRow is a video of up next, which plays when clicked.
func (s *shell) videoRow(c *ui.Context, key string, v youtube.Video, current bool) {
	k := &s.lk
	row := s.buttonBase(c.Key(key)).Gap(12).Padding(8).Radius(k.radiusSm).Label(v.Title).Cursor(ui.CursorPointer)
	switch {
	case current:
		row.Background(k.selected)
	case row.Hovered():
		row.Background(k.hover)
	}
	if row.Clicked() && !current {
		s.playNext(v)
	}
	row.Children(func() {
		ui.Box(c).Size(71, 40).Shrink(0).Radius(4).Clip().Children(func() {
			s.videoThumb(c, key+"-thumb", v, 71, 40)
			s.watchedBar(c, v, 71, 3)
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			ui.Text(c, v.Title).FontSize(k.smallSize()).SingleLine()
			if v.Channel != "" {
				ui.Text(c, v.Channel).FontSize(k.captionSize()).TextColor(k.muted).SingleLine()
			}
		})
		if v.Duration != "" {
			ui.Text(c, v.Duration).FontSize(k.captionSize()).TextColor(k.muted).FontFeatures("tnum")
		}
	})
}

// videoPlayerBar is the Video side's player bar: what plays, the same
// transport and seek bar as the music's, and the same buttons at its side,
// which work on the video.
// videoControlsHeight is how tall the middle of the Video player bar is:
// the play button, 48, the gap, 6, and the seek row, 20.
const videoControlsHeight = 48 + 6 + 20

func (s *shell) videoPlayerBar(c *ui.Context) {
	k := &s.lk
	vp := &s.videos.player
	st := vp.state
	has := vp.current.ID != "" && vp.web != nil
	bar := ui.Row(c).Height(k.player).Shrink(0).Padding(0, 24, 0, 16).Gap(24)
	if k.floating {
		bar.Margin(0, floatMargin, floatMargin).Padding(0, 20, 0, 9).Radius(k.radius).Background(k.surface2)
	} else {
		bar.Background(k.surface).BorderWidth(1, 0, 0, 0).BorderColor(k.border)
	}
	bar.Children(func() {
		// The info starts level with the top of the play button, which
		// sits in the middle of the bar on top of the seek row.
		top := max(0, (k.player-videoControlsHeight)/2)
		ui.Column(c).Grow(1).Basis(0).MinWidth(220).Gap(1).AlignSelf(ui.Start).Margin(top, 0, 0, 0).Children(func() {
			if vp.current.ID == "" {
				ui.Text(c, "Nothing playing").FontSize(k.bodySize()).SingleLine().MinWidth(0)
				ui.Text(c, "Pick a video, or press Random").FontSize(k.smallSize()).TextColor(k.muted).SingleLine().MinWidth(0)
				return
			}
			ui.Text(c, vp.current.Title).FontSize(k.bodySize()).SingleLine().MinWidth(0)
			// The channel opens its page; the chapter that plays follows it.
			ui.Row(c).Gap(0).MinWidth(0).Children(func() {
				s.channelLink(c, "v-channel", vp.current, k.smallSize())
				chapters := s.videoChapters(vp.current.ID)
				if i := videoChapterAt(chapters, st.Time); i >= 0 {
					ui.Text(c, " · "+chapters[i].Title).FontSize(k.smallSize()).TextColor(k.muted).SingleLine().MinWidth(0).Shrink(1)
				}
			})
			s.videoRatingRow(c, vp.current)
		})
		ui.Column(c).Grow(1.6).Basis(0).MinWidth(320).Gap(6).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Gap(12).Children(func() {
				if s.iconButton(c, "v-shuffle", "shuffle", "Shuffle", s.set.VideoShuffle, 18).Clicked() {
					s.set.VideoShuffle = !s.set.VideoShuffle
					s.save()
				}
				if s.roundButton(c, "v-previous", "previous", "Previous", false, !has).Clicked() {
					s.videoPrevious()
				}
				name, label := "play", "Play"
				if st.Playing {
					name, label = "pause", "Pause"
				}
				if s.roundButton(c, "v-toggle", name, label, true, !has).Clicked() {
					s.videoPlayPause()
				}
				if s.roundButton(c, "v-next", "next", "Next", false, !has).Clicked() {
					s.videoNext(false)
				}
				rep, label := "repeat", "Repeat: off"
				switch s.set.VideoRepeat {
				case videoRepeatAll:
					label = "Repeat: all"
				case videoRepeatOne:
					rep, label = "repeatOne", "Repeat: one"
				}
				if s.iconButton(c, "v-repeat", rep, label, s.set.VideoRepeat != videoRepeatOff, 18).Clicked() {
					s.set.VideoRepeat = (s.set.VideoRepeat + 1) % 3
					s.save()
				}
			})
			if !vp.seeking {
				vp.seekTo = st.Time
			}
			ui.Row(c).Gap(8).FillWidth().MaxWidth(640).Children(func() {
				ui.Text(c, clock(time.Duration(vp.seekTo*float64(time.Second)))).FontSize(k.captionSize()).TextColor(k.muted).
					FontFeatures("tnum").Width(44).TextAlign(ui.End)
				// The segments and chapters mark the track, under the bar.
				var seek ui.Element
				track := ui.Box(c.Key("v-seek-track")).Grow(1).MinWidth(0)
				tw := track.Bounds().W
				track.Children(func() {
					s.seekMarks(c, tw, st.Duration)
					seek = s.seekSlider(c, &vp.seekTo, max(st.Duration, 1)).FillWidth().KeepFocus().Label("Seek").Disabled(!has || st.Duration <= 0)
				})
				_, _, held := seek.Dragged()
				if has && !held && (vp.seeking || seek.Changed()) {
					s.videoCommand("seek", vp.seekTo)
					vp.state.Time = vp.seekTo
				}
				vp.seeking = held
				ui.Text(c, clock(time.Duration(st.Duration*float64(time.Second)))).FontSize(k.captionSize()).TextColor(k.muted).FontFeatures("tnum").Width(44)
			})
		})
		ui.Row(c).Grow(1).Basis(0).MinWidth(220).Gap(4).Justify(ui.End).Children(func() { s.videoSide(c, has) })
	})
}

// videoSide is the buttons at the side of the Video side's player bar, as
// the music's: captions in place of lyrics, up next, casting, the speed,
// the device, the volume, and Full App in place of the mini player.
func (s *shell) videoSide(c *ui.Context, has bool) {
	vp := &s.videos.player
	st := vp.state
	// The button shows the user's choice: YouTube's player tells the
	// captions' track for a while after they're turned off.
	on := s.set.VideoCaptions
	if s.iconButton(c, "v-captions", "lyrics", "Captions", on, 18).Disabled(!has).Clicked() {
		s.videoCommand("captions", !on)
		vp.state.Captions = !on
		s.set.VideoCaptions = !on
		s.save()
	}
	s.chaptersButton(c, has)
	if s.iconButton(c, "v-queue", "queue", "Up next", s.panel == "queue", 18).Clicked() {
		s.togglePanel("queue")
	}
	s.videoCastButton(c, has)
	speed := s.iconButton(c, "v-speed", "speed", "Speed: "+videoSpeedLabel(st.Rate), st.Rate != 0 && st.Rate != 1, 18).Disabled(!has)
	speed.Menu(func(m *ui.Menu) {
		for _, r := range settings.VideoRates {
			if m.Item(videoSpeedLabel(r)).Checked(r == st.Rate || st.Rate == 0 && r == 1).Chosen() {
				s.videoCommand("rate", r)
				vp.state.Rate = r
				s.set.VideoRate = r
				s.save()
			}
		}
	})
	s.videoOutputButton(c, has)
	ui.Row(c).Gap(4).Width(140).MinWidth(80).Children(func() {
		name, label := "volume", "Mute"
		if st.Muted || has && st.Volume == 0 {
			name, label = "volumeMute", "Unmute"
		}
		if s.iconButton(c, "v-mute", name, label, false, 18).Disabled(!has).Clicked() {
			s.videoCommand("mute", !st.Muted)
			vp.state.Muted = !st.Muted
			s.set.VideoMuted = !st.Muted
			s.save()
		}
		vp.volume = st.Volume / 100
		if st.Muted {
			vp.volume = 0
		}
		if !has {
			vp.volume = 1
		}
		if ui.Slider(c.Key("v-volume"), &vp.volume, 0, 1).Grow(1).MinWidth(0).KeepFocus().Label("Volume").Disabled(!has).Changed() {
			v := float64(int(vp.volume*100 + 0.5))
			s.videoCommand("volume", v)
			vp.state.Volume, vp.state.Muted = v, false
			s.set.VideoVolume, s.set.VideoMuted = v, false
			s.save()
		}
	})
	pop := "Pop out: watch in a small window over your other apps"
	if vp.popped() {
		pop = "Bring the video back to this window"
	}
	if s.iconButton(c, "v-pop-out", "external", pop, vp.popped(), 18).Disabled(!has || vp.casting != nil).Clicked() {
		if vp.popped() {
			s.closePopOut()
		} else {
			s.popOutVideo()
		}
	}
	if s.iconButton(c, "v-full-app", "mini", "Full App: fill the window with the video", vp.fullApp, 18).Disabled(!has || vp.popped()).Clicked() {
		vp.fullApp = true
	}
}

// videoSpeedLabel writes a speed as YouTube's player does: "Normal", "1.5×".
func videoSpeedLabel(r float64) string {
	if r == 0 || r == 1 {
		return "Normal"
	}
	return strconv.FormatFloat(r, 'f', -1, 64) + "×"
}

// followWatch closes the player as the Music side shows, and pauses the
// video as the music plays.
func (s *shell) followWatch() {
	vp := &s.videos.player
	defer s.followFullScreen()
	defer s.followPreview()
	// The video goes on as the Music side shows, until music starts to
	// play, which pauses it, as a video started while music plays pauses
	// the music (playVideo): it stays where it shows, to go back to.
	playing := s.st.Playing
	started := playing && !vp.musicPlayed
	vp.musicPlayed = playing
	if started && vp.current.ID != "" && vp.state.Playing {
		s.videoCommand("pause", nil)
		vp.state.Playing = false
		vp.pausedForMusic = time.Now()
	}
	s.videoKeys.Store(vp.current.ID != "")
}

// homeVideo is the tile of Music's home that shows the video that plays,
// beside Jump back in, two columns wide and two rows high.
func (s *shell) homeVideo(c *ui.Context, columns int) {
	vp := &s.videos.player
	if vp.current.ID == "" || vp.filling() {
		return
	}
	k := &s.lk
	gap := k.gap / 2
	grid := min(s.paneWidth(c)-2*k.gutter, 1040)
	w := grid
	if columns == 4 {
		w = (grid-3*gap)/2 + gap
	}
	w = float32(int(max(w, 200)))
	ui.Box(c.Key("tile-video")).ColumnSpan(min(2, columns)).RowSpan(2).Radius(k.radius).Clip().Children(func() {
		s.videoMini(c, "home-video", w)
	})
}
