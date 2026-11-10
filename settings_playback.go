package main

import (
	"github.com/egoist/mygo/ui"

	"github.com/bigdaddy284/mymusic/internal/settings"
)

// The groups of Settings › Playback: the sound here, the others in their
// features' files.

// soundRows are the rows of Settings › Playback › Sound, which register
// from their features' files.
const soundRows = "playback › Sound"

func init() { addSettings("playback", 20, (*shell).soundSettings) }

// soundSettings show the rows of the sound.
func (s *shell) soundSettings(c *ui.Context) {
	s.group(c, "Sound", func() { s.settingsOf(c, soundRows) })
}

func init() { addSettings(soundRows, 10, (*shell).qualityRow) }

// qualityRow chooses between the best audio and smaller streams.
func (s *shell) qualityRow(c *ui.Context) {
	s.row(c, "Quality", "Data saver streams smaller files of YouTube Music.", func() {
		saver := 0
		if s.set.Quality == settings.QualitySaver {
			saver = 1
		}
		if s.segmented(c, "quality", "Quality", &saver, "best", "data saver") {
			s.set.Quality = settings.QualityBest
			if saver == 1 {
				s.set.Quality = settings.QualitySaver
			}
			s.save()
		}
	}, "best bitrate")
}

func init() { addSettings(soundRows, 60, (*shell).autoplayRow) }

// autoplayRow carries on with a radio as the queue runs out.
func (s *shell) autoplayRow(c *ui.Context) {
	s.row(c, "Autoplay", "When the queue runs out, a radio of what you played carries on.", func() {
		if s.switchOf(c.Key("autoplay"), &s.set.Autoplay).Label("Autoplay").Changed() {
			s.save()
		}
	})
}
