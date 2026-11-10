package main

// The keyboard's media keys, the system's media controls and the tray
// play, pause, skip and seek the video while one plays, as they do the
// music otherwise.

import (
	"github.com/bigdaddy284/mymusic/internal/app"
)

// videoCommands are the player's commands the video takes.
var videoCommands = map[app.Command]bool{
	app.CommandToggle: true, app.CommandPlay: true, app.CommandPause: true,
	app.CommandNext: true, app.CommandPrevious: true, app.CommandSeek: true,
}

// videoTakes takes a command of the player for the video, while one is
// loaded: on any goroutine, which the window's next update runs it on.
func (s *shell) videoTakes(cmd app.PlayerCommand) bool {
	if !s.videoKeys.Load() || !videoCommands[cmd.Command] {
		return false
	}
	s.update(func() { s.videoCommandOf(cmd) })
	return true
}

// videoCommandOf runs a command of the player on the video.
func (s *shell) videoCommandOf(cmd app.PlayerCommand) {
	vp := &s.videos.player
	if vp.current.ID == "" {
		return
	}
	switch cmd.Command {
	case app.CommandToggle:
		s.videoPlayPause()
	case app.CommandPlay:
		if !vp.state.Playing {
			s.videoPlayPause()
		}
	case app.CommandPause:
		if vp.state.Playing {
			s.videoPlayPause()
		}
	case app.CommandNext:
		s.videoNext(false)
	case app.CommandPrevious:
		s.videoPrevious()
	case app.CommandSeek:
		at := cmd.Value / 1000
		s.videoCommand("seek", at)
		vp.state.Time = at
	}
}
