package app

import (
	"testing"
	"time"

	"github.com/bigdaddy284/mymusic/internal/audio"
	"github.com/bigdaddy284/mymusic/internal/media"
	"github.com/bigdaddy284/mymusic/internal/nowplaying"
	"github.com/bigdaddy284/mymusic/internal/settings"
)

// TestSpeedRanges has the engine, the settings and the system's controls
// offer the same speeds.
func TestSpeedRanges(t *testing.T) {
	if audio.MinSpeed != settings.MinSpeed || audio.MaxSpeed != settings.MaxSpeed ||
		nowplaying.MinRate != settings.MinSpeed || nowplaying.MaxRate != settings.MaxSpeed {
		t.Errorf("the engine plays from %v to %v, the settings from %v to %v, the system's controls from %v to %v",
			audio.MinSpeed, audio.MaxSpeed, settings.MinSpeed, settings.MaxSpeed, nowplaying.MinRate, nowplaying.MaxRate)
	}
}

// TestControlsRate tells the system's controls the speed.
func TestControlsRate(t *testing.T) {
	np := NowPlaying{Track: &media.Track{ID: "yt:1", Title: "Song"}, Playing: true, Speed: 1.5}
	if info := controlsInfo(np, time.Now()); info.Rate != 1.5 {
		t.Errorf("rate %v", info.Rate)
	}
	np.Speed = 0
	if info := controlsInfo(np, time.Now()); info.Rate != 1 {
		t.Errorf("rate %v without a speed", info.Rate)
	}
}
