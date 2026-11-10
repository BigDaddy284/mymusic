package main

// Casting videos: the Video side's cast button lists the cast devices,
// and the one chosen plays the video on its own YouTube app, with the
// player bar as its remote. The video waits, paused, in the window.

import (
	"context"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/app"
)

// videoCasting is a cast of videos: the device, and where the video is on
// it, as MyTube last told it, which goes on while it plays.
type videoCasting struct {
	device  app.CastDevice
	playing bool
	at      float64
	since   time.Time
	// starting is set until the device plays.
	starting bool
}

// now is where the video is on the device.
func (vc *videoCasting) now() float64 {
	if !vc.playing || vc.starting {
		return vc.at
	}
	return vc.at + time.Since(vc.since).Seconds()
}

// set has the video be at at, playing or not, from now.
func (vc *videoCasting) set(at float64, playing bool) {
	vc.at, vc.playing, vc.since = at, playing, time.Now()
}

// videoCastButton lists the cast devices, with this computer first, and
// casts to the one chosen.
func (s *shell) videoCastButton(c *ui.Context, has bool) {
	vp := &s.videos.player
	label := "Cast to a device"
	on := vp.casting != nil
	if on {
		label = "Playing on " + vp.casting.device.Name
	}
	b := s.iconButton(c, "v-cast", "cast", label, on, 18).Disabled(!has)
	b.Menu(func(m *ui.Menu) {
		if m.Item("This computer").Checked(!on).Chosen() && on {
			s.stopVideoCast()
		}
		devices := s.svc.Videos.CastDevices()
		if len(devices) == 0 {
			m.Item("Looking for cast devices…").Disabled(true)
		}
		for _, d := range devices {
			name := d.Name
			if d.Model != "" {
				name += ", " + d.Model
			}
			chosen := on && vp.casting.device.ID == d.ID
			if m.Item(name).Checked(chosen).Chosen() && !chosen {
				s.castVideoTo(d)
			}
		}
	})
}

// castVideoTo has the device d play the video, from where it is.
func (s *shell) castVideoTo(d app.CastDevice) {
	vp := &s.videos.player
	if vp.current.ID == "" {
		return
	}
	at := vp.state.Time
	if vp.casting != nil {
		at = vp.casting.now()
	}
	s.localVideoCommand("pause", nil)
	vp.casting = &videoCasting{device: d, starting: true}
	vp.casting.set(at, true)
	s.castCurrent(at)
}

// castCurrent has the device play the video that plays now, from at.
func (s *shell) castCurrent(at float64) {
	vp := &s.videos.player
	vc := vp.casting
	if vc == nil {
		return
	}
	d, id := vc.device, vp.current.ID
	vc.starting = true
	vc.set(at, true)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
		defer cancel()
		err := s.svc.Videos.Cast(ctx, d.ID, id, at)
		s.update(func() {
			if vp.casting != vc {
				return
			}
			if err != nil {
				vp.casting = nil
				s.notes = append(s.notes, noticeOf(err))
				return
			}
			vc.starting = false
			vc.set(at, true)
		})
	}()
}

// stopVideoCast brings the video back to the window, from where it was on
// the device.
func (s *shell) stopVideoCast() {
	vp := &s.videos.player
	vc := vp.casting
	if vc == nil {
		return
	}
	at, playing := vc.now(), vc.playing
	vp.casting = nil
	go s.svc.Videos.StopVideoCast(context.Background())
	s.localVideoCommand("seek", at)
	if playing {
		s.localVideoCommand("play", nil)
	}
	vp.state.Time, vp.state.Playing = at, playing
}

// castCommand runs a command of the player on the device the video is
// cast to, and tells whether it took it.
func (s *shell) castCommand(cmd string, arg any) bool {
	vp := &s.videos.player
	vc := vp.casting
	if vc == nil {
		return false
	}
	videos := s.svc.Videos
	run := func(fn func(ctx context.Context) error) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			if err := fn(ctx); err != nil {
				s.update(func() { s.notes = append(s.notes, noticeOf(err)) })
			}
		}()
	}
	switch cmd {
	case "play":
		vc.set(vc.now(), true)
		run(videos.CastPlay)
	case "pause":
		vc.set(vc.now(), false)
		run(videos.CastPause)
	case "seek":
		at, _ := arg.(float64)
		vc.set(at, vc.playing)
		run(func(ctx context.Context) error { return videos.CastSeek(ctx, at) })
	case "volume":
		v, _ := arg.(float64)
		run(func(ctx context.Context) error { return videos.CastVolume(ctx, v) })
	case "mute":
		if muted, _ := arg.(bool); muted {
			run(func(ctx context.Context) error { return videos.CastVolume(ctx, 0) })
		} else {
			v := s.set.VideoVolume
			run(func(ctx context.Context) error { return videos.CastVolume(ctx, v) })
		}
	default:
		// Captions and speed are the device's own: its remote sets them.
		return true
	}
	return true
}

// followVideoCast keeps the window's video paused while the device plays it,
// and shows where it is there.
func (s *shell) followVideoCast(st *videoState) {
	vc := s.videos.player.casting
	if vc == nil {
		return
	}
	if st.Playing {
		s.localVideoCommand("pause", nil)
	}
	st.Playing, st.Time, st.Ended = vc.playing, vc.now(), false
	if st.Duration > 0 && vc.playing && !vc.starting && st.Time >= st.Duration {
		st.Time, st.Ended = st.Duration, true
	}
}
