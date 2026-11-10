package app

import (
	"math"
	"testing"

	"github.com/bigdaddy284/mymusic/internal/audio"
	"github.com/bigdaddy284/mymusic/internal/player"
	"github.com/bigdaddy284/mymusic/internal/settings"
	"github.com/bigdaddy284/mymusic/internal/volume"
)

// TestVolumeCurve plays the volume on the curve of the settings in
// MyMusic's own engine (#383), as the curve changes, and leaves the
// volume's position where it was.
func TestVolumeCurve(t *testing.T) {
	a := newQuietApp(t, nil)
	// An engine of its own, which plays into no device, in place of the
	// one of the computer's output.
	e := audio.NewEngine(a.openNative, startSilent, engineEvents{a})
	t.Cleanup(e.Close)
	a.nativeOnce.Do(func() { a.native = e })
	a.applyEngine(a.settings.Get().Clean())
	p := a.GoPlayer()
	p.SetVolume(0.5)
	waitPlayer(t, p, "the volume at 0.5", func(s player.State) bool { return s.Volume == 0.5 })
	check := func(c volume.Curve) {
		t.Helper()
		if a.nativeEngine() != nil {
			if got := e.VolumeCurve(); got != c {
				t.Errorf("the engine plays on %q, not %q", got, c)
			}
		}
		if v := p.State().Volume; v != 0.5 {
			t.Errorf("on %q, the volume moved to %v", c, v)
		}
	}
	check(volume.Natural)
	for _, c := range []volume.Curve{volume.Linear, volume.Natural} {
		if _, err := (&Preferences{a}).Update(func(s *settings.Settings) { s.VolumeCurve = c }); err != nil {
			t.Fatal(err)
		}
		check(c)
	}
	// What settings.json keeps of the volume is its position, whatever the
	// curve.
	if _, err := (&Preferences{a}).Update(func(s *settings.Settings) { s.VolumeCurve = volume.Linear }); err != nil {
		t.Fatal(err)
	}
	p.Flush()
	if v := a.settings.Get().Volume; math.Abs(v-volume.Saved(0.5)) > 1e-9 {
		t.Errorf("on Linear, settings.json keeps the volume as %v, not %v", v, volume.Saved(0.5))
	}
	// A volume key steps the position, which plays on the curve.
	a.command(PlayerCommand{Command: CommandVolumeBy, Value: 0.05})
	waitPlayer(t, p, "a step up", func(s player.State) bool { return math.Abs(s.Volume-0.55) < 1e-9 })
}
