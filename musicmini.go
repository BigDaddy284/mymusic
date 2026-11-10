package main

// The music's own small player on the Video side, at the top right of the
// page, so the music that plays can be controlled while the user looks
// through videos: the Video side's bar at the foot of the window plays the
// videos.

import (
	"github.com/egoist/mygo/ui"
)

// musicMiniWidth is how wide the music's small player is.
const musicMiniWidth = 400

// musicMiniHeight is how tall it is, the artwork's side.
const musicMiniHeight = 92

// showsMusicMini reports whether the Video side shows the music's small
// player: while a track is the queue's, playing or paused for a video, to
// go back to, but for the watch page, which the video fills, and a track
// the user closed it for.
func (s *shell) showsMusicMini() bool {
	if !s.mode.video {
		return false
	}
	cur, ok := s.st.Current()
	if !ok || cur.Track.ID == s.musicMiniClosed {
		return false
	}
	vp := &s.videos.player
	if vp.filling() {
		return false
	}
	return vp.current.ID == "" || s.router.Path() != watchPath(vp.current.ID)
}

// musicMini is the music's small player: the artwork, the title and the
// artists with the time, the transport, and the progress along its foot,
// which seeks, with a button that closes it until the next track.
func (s *shell) musicMini(c *ui.Context) {
	if !s.showsMusicMini() {
		return
	}
	k := &s.lk
	cur, _ := s.st.Current()
	ui.Box(c.Key("music-mini")).Absolute().Top(8).Right(20).Size(musicMiniWidth, musicMiniHeight).Radius(k.radiusSm).
		Background(k.surface2).Clip().Shadow(0, 12, 32, -8, ui.RGBA(0, 0, 0, 0.55)).Children(func() {
		ui.Row(c).Size(musicMiniWidth, musicMiniHeight).Children(func() {
			s.miniArtwork(c, musicMiniHeight, musicMiniHeight)
			ui.Column(c).Grow(1).MinWidth(0).Padding(6, 6, miniSeekRoom).PaddingStart(12).Gap(2).Justify(ui.Center).Children(func() {
				ui.Row(c).Gap(4).AlignItems(ui.Center).MinWidth(0).Children(func() {
					ui.Text(c, cur.Track.Title).FontSize(k.bodySize()).FontWeight(k.titleWeight()).SingleLine().Grow(1).MinWidth(0)
					if s.miniButton(c, "music-mini-close", "close", "Close the music player", false, 24, 14).Clicked() {
						s.musicMiniClosed = cur.Track.ID
					}
				})
				ui.Row(c).Gap(6).AlignItems(ui.Center).MinWidth(0).Children(func() {
					s.sourceMark(c, cur.Track.Source, false)
					ui.Text(c, names(cur.Track.Artists)).FontSize(k.smallSize()).TextColor(k.muted).SingleLine().Grow(1).MinWidth(0)
					if !s.st.Live {
						ui.Text(c, clock(s.st.Position)+" / "+clock(s.miniDuration())).FontSize(k.captionSize()).
							TextColor(k.muted).FontFeatures("tnum").Shrink(0).SingleLine()
					}
				})
				s.miniControls(c, musicMiniWidth-musicMiniHeight-18)
			})
		})
		if !s.st.Live {
			s.miniSeek(c, musicMiniWidth)
		}
	})
}
