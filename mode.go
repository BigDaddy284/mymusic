package main

// MyTube has two sides, Music and Video, which the mark at the head of the
// sidebar switches between: clicked, it spins over a few times, like a
// coin with the Music mark on one face and the Video mark on the other,
// and stops on the other side, as the name beside it changes. Each side
// has its own menu, pages and player bar, and goes back to the page it
// showed last.

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/brand"
)

// videoMark is the Video mark, the back of the Music mark.
var videoMark = mustSVG(brand.VideoMark)

const (
	// modeSpin is how far the mark turns, in degrees, as the side
	// switches: three and a half turns, which end on the other face.
	modeSpin = 1260
	// modeSpinTime is how long it turns for, slowing to a stop.
	modeSpinTime = 2400 * time.Millisecond
)

// appMode is the side the window shows, and how far the mark turned.
type appMode struct {
	video bool
	turn  float32
	// lastMusic and lastVideo are the pages each side showed last.
	lastMusic, lastVideo string
}

// toggleMode switches between Music and Video, spinning the mark, and goes
// to the page the other side showed last.
func (s *shell) toggleMode() {
	m := &s.mode
	path := s.router.Path()
	if m.video {
		m.lastVideo = path
	} else {
		m.lastMusic = path
	}
	m.video = !m.video
	m.turn += modeSpin
	to := m.lastMusic
	if m.video {
		to = m.lastVideo
	}
	// The watch page left goes on with the video playing now, which may
	// have moved on to the next while the window showed Music.
	if cur := s.videos.player.current.ID; m.video && cur != "" && strings.HasPrefix(to, "/watch/") {
		to = watchPath(cur)
	}
	if to == "" || isVideoPath(to) != m.video {
		to = "/"
		if m.video {
			to = "/videos"
		}
	}
	s.navigate(to)
}

// isVideoPath reports whether path is a page of the Video side.
func isVideoPath(path string) bool {
	return path == "/videos" || strings.HasPrefix(path, "/videos/") || strings.HasPrefix(path, "/videos?") ||
		strings.HasPrefix(path, "/watch/")
}

// isSharedPath reports whether path is a page both sides show: the
// settings and the stats.
func isSharedPath(path string) bool {
	first, _, _ := strings.Cut(strings.TrimPrefix(path, "/"), "/")
	first, _, _ = strings.Cut(first, "?")
	return first == "settings" || first == "stats" || first == "plugin"
}

// followMode puts the window on the side of the page path, when it is one
// side's alone, as one the palette or a link opened: the mark turns over
// once, quickly.
func (s *shell) followMode(path string) {
	if isSharedPath(path) {
		return
	}
	if video := isVideoPath(path); video != s.mode.video {
		s.mode.video = video
		s.mode.turn += 180
	}
}

// railTitle is the name beside the mark: the side shown.
func (s *shell) railTitle() string {
	if s.mode.video {
		return "Video"
	}
	return "Music"
}

// modeMark is the mark, size wide, which switches the side when clicked.
// Its two faces turn about its vertical axis, each hiding while it faces
// away.
func (s *shell) modeMark(c *ui.Context, size float32) {
	m := &s.mode
	label := "Switch to Video"
	if m.video {
		label = "Switch to Music"
	}
	b := s.buttonBase(c.Key("mode-mark")).Size(size, size).Shrink(0).Label(label).Tooltip(label).Cursor(ui.CursorPointer)
	if b.Clicked() {
		s.toggleMode()
	}
	turn := b.AnimateWith("turn", m.turn, modeSpinTime, ui.EaseOut)
	b.Children(func() {
		ui.Image(c.Key("mark-music"), brandMark).Size(size, size).Absolute().Left(0).Top(0).
			RotateY(turn).Perspective(6 * size).BackfaceHidden()
		ui.Image(c.Key("mark-video"), videoMark).Size(size, size).Absolute().Left(0).Top(0).
			RotateY(turn + 180).Perspective(6 * size).BackfaceHidden()
	})
}

// videoNav is the Video side's menu.
var videoNav = []navItem{
	{"vhome", "Home", "home", "/videos"},
	{"vsearch", "Search", "search", "/videos/search"},
	{"vtrending", "Trending", "compass", "/videos/trending"},
	{"vlibrary", "Library", "library", "/videos/library"},
	{"stats", "Stats", "stats", "/stats"},
	{"settings", "Settings", "settings", "/settings"},
}

// navItems is the menu of the side shown, and the entry of path in it.
func (s *shell) navItems(path string) ([]navItem, string) {
	if !s.mode.video {
		return nav, section(path)
	}
	p, _, _ := strings.Cut(path, "?")
	switch {
	case p == "/videos" || strings.HasPrefix(p, "/watch/"):
		return videoNav, "vhome"
	case strings.HasPrefix(p, "/videos/search"):
		return videoNav, "vsearch"
	case strings.HasPrefix(p, "/videos/trending"):
		return videoNav, "vtrending"
	case strings.HasPrefix(p, "/videos/library"), strings.HasPrefix(p, "/videos/playlist/"):
		return videoNav, "vlibrary"
	}
	return videoNav, section(path)
}

// videoLists are the Video side's playlists in the sidebar, where the
// Music side has its mixes: Watch later and Liked videos, then the
// account's own.
func (s *shell) videoLists(c *ui.Context, path string) {
	k := &s.lk
	ui.Column(c).Grow(1).MinHeight(mixesLeast).Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Margin(0, 0, 8).Children(func() {
			ui.Text(c, k.words("Your playlists")).FontSize(k.bodySize()).TextColor(k.accent)
		})
		if !s.svc.Videos.SignedIn() {
			ui.Text(c, "Sign in to YouTube Music to see your playlists here.").FontSize(k.smallSize()).TextColor(k.muted)
			return
		}
		lists := s.videoPlaylistsLoad()
		ui.Scroll(c).Grow(1).MinHeight(0).Children(func() {
			if !lists.done || lists.err != nil {
				s.sideLink(c, "vl-WL", k.words("Watch later"), -1, videoPlaylistPath("WL"), path)
				s.sideLink(c, "vl-LL", k.words("Liked videos"), -1, videoPlaylistPath("LL"), path)
				return
			}
			for _, p := range lists.value {
				s.sideLink(c, "vl-"+p.ID, p.Title, -1, videoPlaylistPath(p.ID), path)
			}
		})
	})
}
