package main

// Full screen of the Video side: YouTube's own full screen button, in
// the player, and the app's Full screen button, both fill the screen with
// the video; Escape leaves it, and Full App, even while the player has
// the keys.

import (
	"log"

	"github.com/egoist/mygo"
)

// escapeKey is the key that leaves Full App and full screen. The page of
// YouTube's player, from another site, keeps the keys pressed in it from
// the window, so while either is on, and the window is in front, the
// app takes Escape from the system, as a global shortcut, and gives it
// back after.
const escapeKey = "Escape"

// videoFullScreen tells whether the video fills the screen: the app's
// Full screen button, or the player's own, asked for it.
func (vp *videoPlayer) videoFullScreen() bool {
	return vp.current.ID != "" && !vp.popped() && (vp.fullScreen || vp.state.FS)
}

// filling tells whether the video fills the window: Full App or full
// screen is on.
func (vp *videoPlayer) filling() bool { return vp.fullApp || vp.videoFullScreen() }

// leaveFullApp leaves Full App and full screen, the player's own too.
func (s *shell) leaveFullApp() {
	vp := &s.videos.player
	if vp.state.FS {
		s.videoCommand("exitFullScreen", nil)
		vp.state.FS = false
	}
	vp.fullApp, vp.fullScreen = false, false
	s.followFullScreen()
}

// followFullScreen makes the window full screen while the video asks for
// it, and takes Escape while the video fills the window in front.
func (s *shell) followFullScreen() {
	vp := &s.videos.player
	if s.win == nil {
		return
	}
	full := s.mode.video && vp.videoFullScreen()
	if full != vp.ownsScreen {
		// Only the full screen the video asked for is given back: the
		// window may have been full screen before.
		if full && s.win.IsFullScreen() {
			full = false
		} else {
			vp.ownsScreen = full
			win := s.win
			go win.SetFullScreen(full)
		}
	}
	take := s.mode.video && vp.filling() && !vp.blurred
	if take == vp.takesEscape {
		return
	}
	if take {
		if err := mygo.GlobalShortcut.Register(escapeKey, func() { s.update(s.leaveFullApp) }); err != nil {
			log.Printf("mytube: Escape for Full App: %v", err)
			return
		}
	} else {
		mygo.GlobalShortcut.Unregister(escapeKey)
	}
	vp.takesEscape = take
}

// followFocus gives Escape back to other apps while the window isn't in
// front, and takes it again as it comes back.
func (s *shell) followFocus(focused bool) func() {
	return func() {
		s.videos.player.blurred = !focused
		s.followFullScreen()
	}
}
