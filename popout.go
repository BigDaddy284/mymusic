package main

// The pop-out player: the video that plays, alone in a small window of its
// own over the others, which the user moves and sizes where they like, so
// they watch while they work in other apps. The window's web view takes
// the place of the main window's, from where the video was: the player
// bar, the keys and Up next control it as they did, and closing the window
// brings the video back to the main window.

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// popOut is the pop-out player's own state.
type popOut struct {
	// win is the window, while it is open; opening is set while it opens.
	win     *mygo.Window
	opening bool
	// full is set while the window fills the screen for the player, and
	// title is the title it was given.
	full  bool
	title string
}

// popped reports whether the video plays in the pop-out window.
func (vp *videoPlayer) popped() bool { return vp.pop.win != nil }

// popOutVideo moves the video that plays to a window of its own.
func (s *shell) popOutVideo() {
	vp := &s.videos.player
	if vp.popped() || vp.pop.opening || vp.current.ID == "" || vp.casting != nil {
		return
	}
	vp.fullApp = false
	vp.pop.opening = true
	id := vp.current.ID
	var win *mygo.Window
	view := func(c *ui.Context) {
		w, h := c.Size()
		root := ui.Box(c).Fill().Background(ui.RGBA(0, 0, 0, 1))
		root.Children(func() {
			if web := vp.web; web != nil && vp.pop.win == win && win != nil {
				ui.WebView(c.Key("pop-web"), web).Size(w, h).Label("Video: " + vp.current.Title)
			}
		})
	}
	win = mygo.NewWindow(mygo.WindowOptions{
		Title:           vp.current.Title,
		Width:           480,
		Height:          270,
		MinWidth:        240,
		MinHeight:       135,
		AlwaysOnTop:     true,
		BackgroundColor: "#000000",
		StateKey:        "video-pop-out",
		Content:         ui.View(view),
	})
	captions := s.set.VideoCaptions
	go func() {
		web, err := win.NewWebView(mygo.WebViewOptions{URL: watchPageURL(id, captions), BackgroundColor: "#000000"})
		s.update(func() {
			vp.pop.opening = false
			if err != nil || vp.current.ID != id {
				if err != nil {
					log.Printf("mytube: the pop-out player can't open: %v", err)
					s.notes = append(s.notes, notice("The video couldn't pop out"))
				}
				win.Close()
				return
			}
			// The main window's player stops where the pop-out's goes on.
			s.keepVideoPlace()
			s.swapVideoWeb(web, id)
			vp.pop.win, vp.pop.title = win, vp.current.Title
			win.OnClosed(func() { s.update(func() { s.popBackIn(win) }) })
			win.Invalidate()
		})
	}()
}

// swapVideoWeb has web play the video id in place of the player's web view,
// which closes: it settles again, from where the video was left.
func (s *shell) swapVideoWeb(web *mygo.WebView, id string) {
	vp := &s.videos.player
	if vp.stop != nil {
		close(vp.stop)
	}
	if old := vp.web; old != nil {
		go old.Destroy()
	}
	vp.web, vp.loaded = web, id
	vp.settledFor = ""
	vp.state.Playing = false
	vp.stop = make(chan struct{})
	go s.pollVideo(web, vp.stop)
}

// popBackIn brings the video back to the main window as the pop-out
// window win closes, from where it was.
func (s *shell) popBackIn(win *mygo.Window) {
	vp := &s.videos.player
	if vp.pop.win != win {
		return
	}
	vp.pop.win, vp.pop.full = nil, false
	if vp.stop != nil {
		close(vp.stop)
		vp.stop = nil
	}
	// The window's web view closed with it.
	vp.web, vp.loaded, vp.settledFor = nil, "", ""
	if vp.current.ID != "" {
		s.loadVideo()
	}
}

// closePopOut closes the pop-out window, which brings the video back.
func (s *shell) closePopOut() {
	if win := s.videos.player.pop.win; win != nil {
		win.Close()
	}
}

// poppedOutFrame stands for the video in the main window while it plays in
// the pop-out window, w × h, with a button that brings it back.
func (s *shell) poppedOutFrame(c *ui.Context, key string, w, h float32) {
	white := ui.RGBA(255, 255, 255, 1)
	ui.Column(c).Size(w, h).AlignItems(ui.Center).Justify(ui.Center).Gap(12).Padding(16).Children(func() {
		ui.Text(c, "The video plays in its own window.").TextColor(white).TextAlign(ui.Center)
		if s.button(c, key+"-pop-in", "bring it back", true).Clicked() {
			s.closePopOut()
		}
	})
}

// followPopScreen has the pop-out window fill the screen while the
// player's own full screen button asks it to.
func (s *shell) followPopScreen() {
	vp := &s.videos.player
	win := vp.pop.win
	if win == nil {
		return
	}
	if vp.pop.title != vp.current.Title {
		vp.pop.title = vp.current.Title
		title := vp.pop.title
		go win.SetTitle(title)
	}
	if vp.state.FS == vp.pop.full {
		return
	}
	vp.pop.full = vp.state.FS
	full := vp.pop.full
	go win.SetFullScreen(full)
}

// videoMoreButton is the ⋮ button at the end of the Video player bar: its
// menu jumps to a chapter, pops the video out, or plays it as audio, and,
// with fullApp set, fills the window with it.
func (s *shell) videoMoreButton(c *ui.Context, has, fullApp bool) {
	vp := &s.videos.player
	s.iconButton(c, "v-more", "moreVertical", "More", false, 18).Disabled(!has).Menu(func(m *ui.Menu) {
		s.chaptersMenu(m)
		if fullApp && m.Item("Full App: fill the window with the video").Disabled(vp.popped()).Chosen() {
			vp.fullApp = true
		}
		pop := "Pop out: watch in a window over your other apps"
		if vp.popped() {
			pop = "Bring the video back to this window"
		}
		if m.Item(pop).Disabled(vp.casting != nil).Chosen() {
			if vp.popped() {
				s.closePopOut()
			} else {
				s.popOutVideo()
			}
		}
		if m.Item("Play as audio on the music side").Chosen() {
			s.playVideoAudio(vp.current)
		}
		if m.Item("Copy the link").Chosen() {
			copyText(youtubeWatchURL(vp.current.ID))
			s.notes = append(s.notes, notice("Copied the link to "+vp.current.Title))
		}
	})
}
