package settings

// Watching is how the Video side's player plays, kept from one start to
// the next: its volume, its speed, its captions and the order of up next.
type Watching struct {
	// VideoVolume is from 0 to 100, as YouTube's player takes it.
	VideoVolume float64 `json:"videoVolume"`
	VideoMuted  bool    `json:"videoMuted,omitempty"`
	// VideoRate is the speed, 1 for normal.
	VideoRate float64 `json:"videoRate"`
	// VideoCaptions shows the captions of the videos that have them; off
	// in new installs.
	VideoCaptions bool `json:"videoCaptions,omitempty"`
	VideoShuffle  bool `json:"videoShuffle,omitempty"`
	// VideoRepeat is 0 for off, 1 for all of up next, 2 for one video.
	VideoRepeat int `json:"videoRepeat,omitempty"`
	// VideoOutput is the device the videos' sound plays on, apart from
	// the music's: none for the system's default.
	VideoOutput OutputDevice `json:"videoOutput"`
}

func defaultWatching() Watching { return Watching{VideoVolume: 100, VideoRate: 1} }

// VideoRates are the speeds of YouTube's player.
var VideoRates = []float64{0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2}

func (w Watching) clean() Watching {
	w.VideoVolume = min(max(w.VideoVolume, 0), 100)
	known := false
	for _, r := range VideoRates {
		known = known || r == w.VideoRate
	}
	if !known {
		w.VideoRate = 1
	}
	w.VideoOutput = w.VideoOutput.clean()
	if w.VideoRepeat < 0 || w.VideoRepeat > 2 {
		w.VideoRepeat = 0
	}
	return w
}
