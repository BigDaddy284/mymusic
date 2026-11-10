package app

// Casting the Video side's videos to a device's own YouTube app, through
// YouTube's lounge (internal/cast/lounge.go): the device plays the video
// itself, and MyTube is its remote.

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/egoist/mygo/netproxy"

	"github.com/bigdaddy284/mymusic/internal/cast"
)

var loungeClient = &http.Client{Transport: netproxy.Transport(), Timeout: 15 * time.Second}

// videoCast is the cast of videos under way.
type videoCast struct {
	mu     sync.Mutex
	device CastDevice
	lounge *cast.Lounge
	// screens are the screen IDs the devices told, by device.
	screens map[string]string
}

// CastDevices returns the cast devices on the network, as the music's
// cast button lists them, and looks for them again unless it did a
// moment ago.
func (v *Videos) CastDevices() []CastDevice {
	p := &Player{a: v.a}
	p.FindCastDevices()
	return p.CastDevices()
}

// Cast has the device id play the video videoID from at seconds in, on its
// own YouTube app.
func (v *Videos) Cast(ctx context.Context, id, videoID string, at float64) error {
	a := v.a
	a.cast.mu.Lock()
	f, ok := a.cast.found[id]
	a.cast.mu.Unlock()
	if !ok {
		return errors.New("That cast device is no longer on the network")
	}
	vc := &a.videoCast
	vc.mu.Lock()
	defer vc.mu.Unlock()
	if vc.lounge == nil || vc.device.ID != id {
		screen := vc.screens[id]
		if screen == "" {
			s, err := cast.YouTubeScreen(ctx, f.device)
			if err != nil {
				return errors.New(f.device.Name + " couldn't start YouTube: " + err.Error())
			}
			screen = s
			if vc.screens == nil {
				vc.screens = map[string]string{}
			}
			vc.screens[id] = screen
		}
		l, err := cast.BindLounge(ctx, loungeClient, screen)
		if err != nil {
			delete(vc.screens, id)
			return err
		}
		vc.lounge, vc.device = l, castDeviceOf(f.device)
	}
	return vc.lounge.PlayVideo(ctx, videoID, at)
}

// VideoCastingTo returns the device videos play on, and whether they do.
func (v *Videos) VideoCastingTo() (CastDevice, bool) {
	vc := &v.a.videoCast
	vc.mu.Lock()
	defer vc.mu.Unlock()
	return vc.device, vc.lounge != nil
}

// castDo runs fn on the lounge of the cast under way.
func (v *Videos) castDo(fn func(l *cast.Lounge) error) error {
	vc := &v.a.videoCast
	vc.mu.Lock()
	l := vc.lounge
	vc.mu.Unlock()
	if l == nil {
		return errors.New("no video is being cast")
	}
	return fn(l)
}

// CastPlay, CastPause, CastSeek and CastVolume control the video cast.
func (v *Videos) CastPlay(ctx context.Context) error {
	return v.castDo(func(l *cast.Lounge) error { return l.Play(ctx) })
}

func (v *Videos) CastPause(ctx context.Context) error {
	return v.castDo(func(l *cast.Lounge) error { return l.Pause(ctx) })
}

func (v *Videos) CastSeek(ctx context.Context, at float64) error {
	return v.castDo(func(l *cast.Lounge) error { return l.Seek(ctx, at) })
}

func (v *Videos) CastVolume(ctx context.Context, volume float64) error {
	return v.castDo(func(l *cast.Lounge) error { return l.SetVolume(ctx, volume) })
}

// StopVideoCast pauses the video on the device and lets it go.
func (v *Videos) StopVideoCast(ctx context.Context) {
	vc := &v.a.videoCast
	vc.mu.Lock()
	l := vc.lounge
	vc.lounge, vc.device = nil, CastDevice{}
	vc.mu.Unlock()
	if l != nil {
		_ = l.Pause(ctx)
	}
}
