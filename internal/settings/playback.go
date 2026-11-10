package settings

import (
	"math"

	"github.com/bigdaddy284/mymusic/internal/volume"
)

// Quality chooses between the best audio and smaller streams.
type Quality string

const (
	QualityBest  Quality = "best"
	QualitySaver Quality = "saver"
)

// Playback is how MyMusic plays: the sound, the speed,
// the music videos and the equalizer of Settings › Playback, and the
// volume.
type Playback struct {
	// Quality chooses between the best audio of YouTube Music and smaller
	// streams.
	Quality Quality `json:"quality"`
	// Crossfade overlaps tracks by this many seconds, from 0 to 12.
	Crossfade int `json:"crossfade"`
	// Normalize runs Volume levelling, a compressor that squeezes loud and
	// quiet moments closer together as they play.
	Normalize bool `json:"normalize"`
	// Loudness turns tracks louder than -14 LKFS down to it, from YouTube's
	// measure of each: Even loudness.
	Loudness bool `json:"loudness"`
	// Surround plays a track's 5.1 surround audio, when YouTube has one,
	// in place of its stereo audio: on six speakers where the output device
	// has them, else mixed down to stereo as the track's stream asks.
	Surround bool `json:"surround"`
	// OutputDevice is the device MyMusic's own player plays on: the
	// system's default when its ID is empty.
	OutputDevice OutputDevice `json:"outputDevice"`
	// Autoplay continues with similar tracks when the queue ends.
	Autoplay bool `json:"autoplay"`
	// SongSpeed plays songs faster or slower, at their own pitch, from
	// MinSpeed to MaxSpeed by SpeedStep: 1 is their own speed.
	SongSpeed float64 `json:"songSpeed"`
	// PodcastSpeed plays podcasts' episodes faster or slower, as SongSpeed
	// plays songs, from MinSpeed to MaxSpeed by SpeedStep: 1 is their own
	// speed.
	PodcastSpeed float64 `json:"podcastSpeed"`
	// SkipDisliked passes over the tracks the user disliked as the queue
	// moves on, and leaves them out of radios.
	SkipDisliked bool `json:"skipDisliked"`
	// SkipAI passes over the tracks YouTube labels as made with AI, as the
	// player moves on by itself.
	SkipAI bool `json:"skipAi"`
	// MusicVideos shows the picture of music videos on the now playing
	// page, in place of their cover. It is off at first: it fetches more
	// than the sound, and shows a web page there.
	MusicVideos bool `json:"musicVideos"`
	// Equalizer shapes the sound of MyMusic's own player.
	Equalizer Equalizer `json:"equalizer"`
	// VolumeCurve is how the positions of the volume slider sound, in
	// MyMusic's player (#383): Natural at first.
	VolumeCurve volume.Curve `json:"volumeCurve"`
	// Volume is from 0 to 1.
	Volume float64 `json:"volume"`
}

// defaultPlayback is how a new install plays.
func defaultPlayback() Playback {
	var p Playback
	p.Quality = QualityBest
	p.Loudness = true
	p.Autoplay = true
	p.SongSpeed = 1
	p.PodcastSpeed = 1
	p.SkipDisliked = true
	p.Equalizer = defaultEqualizer()
	p.VolumeCurve = volume.Natural
	p.Volume = 0.8
	return p
}

// clean is p with every value in its range, and d's in place of values
// that are missing or invalid.
func (p Playback) clean(d Playback) Playback {
	if p.Quality != QualitySaver {
		p.Quality = QualityBest
	}
	p.Crossfade = min(max(p.Crossfade, 0), 12)
	p.OutputDevice = p.OutputDevice.clean()
	p.SongSpeed = CleanSpeed(p.SongSpeed)
	p.PodcastSpeed = CleanSpeed(p.PodcastSpeed)
	p.Volume = min(max(p.Volume, 0), 1)
	p.VolumeCurve = cleanVolumeCurve(p.VolumeCurve, d.VolumeCurve)
	p.Equalizer = p.Equalizer.clean()
	return p
}

// The speeds MyMusic's own player plays at: from half a track's own speed
// to three times it, by twentieths.
const (
	MinSpeed  = 0.5
	MaxSpeed  = 3
	SpeedStep = 0.05
)

// CleanSpeed returns the speed of the steps nearest v, from MinSpeed to
// MaxSpeed: 1, a track's own speed, for none, as in settings from before
// the speed, and for one that isn't a number.
func CleanSpeed(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return 1
	}
	v = min(max(v, MinSpeed), MaxSpeed)
	// Dividing a whole number of steps gives the speed as the settings
	// write it, 1.1 rather than 1.1000000000000001.
	return math.Round(v*(1/SpeedStep)) / (1 / SpeedStep)
}
