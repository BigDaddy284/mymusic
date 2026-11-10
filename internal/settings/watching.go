package settings

import "strings"

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
	// SponsorBlockOff stops the player skipping the parts of videos that
	// SponsorBlock's users marked; Segments is what it does with each
	// kind, by its category, where it doesn't do the kind's default.
	SponsorBlockOff bool              `json:"sponsorBlockOff,omitempty"`
	Segments        map[string]string `json:"segments,omitempty"`
	// BlockedChannels and BlockedWords keep videos off every page of the
	// Video side: those of the channels, by name, and those whose titles
	// hold a word or phrase, in any case.
	BlockedChannels []string `json:"blockedChannels,omitempty"`
	BlockedWords    []string `json:"blockedWords,omitempty"`
	// ShowShorts shows YouTube's Shorts among the videos; off, they're left
	// out of every page.
	ShowShorts bool `json:"showShorts,omitempty"`
}

// The actions of a kind of segment.
const (
	SegmentSkip = "skip"
	SegmentMute = "mute"
	SegmentShow = "show"
	SegmentOff  = "off"
)

// A SegmentKind is a category of SponsorBlock's segments.
type SegmentKind struct {
	ID, Name, About string
	// Default is what the player does with it in a new install; Point is
	// set for a moment rather than a stretch, which shows alone.
	Default string
	Point   bool
}

// SegmentKinds are SponsorBlock's categories, in the order Settings shows
// them.
var SegmentKinds = []SegmentKind{
	{"sponsor", "Sponsor", "Paid promotions and adverts read in the video", SegmentSkip, false},
	{"selfpromo", "Self-promotion", "Merch, Patreon, and the creator's other channels", SegmentSkip, false},
	{"interaction", "Like and subscribe", "Asking you to like, subscribe or follow", SegmentSkip, false},
	{"intro", "Intro", "Intro animations and stretches with nothing in them", SegmentShow, false},
	{"outro", "Endcards and credits", "The end screen and credits", SegmentShow, false},
	{"preview", "Preview and recap", "Clips of what's coming, or a recap of before", SegmentShow, false},
	{"filler", "Filler and tangents", "Jokes and asides the video doesn't need", SegmentOff, false},
	{"music_offtopic", "Non-music in music videos", "The talking and skits of a music video", SegmentSkip, false},
	{"poi_highlight", "Highlight", "The part most people came for, marked on the seek bar", SegmentShow, true},
}

// SegmentAction is what the player does with segments of the category id.
func (w Watching) SegmentAction(id string) string {
	if a, ok := w.Segments[id]; ok {
		return a
	}
	for _, k := range SegmentKinds {
		if k.ID == id {
			return k.Default
		}
	}
	return SegmentOff
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
	segs := map[string]string{}
	for _, k := range SegmentKinds {
		a, ok := w.Segments[k.ID]
		valid := a == SegmentSkip || a == SegmentShow || a == SegmentOff || a == SegmentMute && !k.Point
		if ok && valid && a != k.Default {
			segs[k.ID] = a
		}
	}
	w.Segments = nil
	if len(segs) > 0 {
		w.Segments = segs
	}
	w.BlockedChannels = cleanList(w.BlockedChannels)
	w.BlockedWords = cleanList(w.BlockedWords)
	return w
}

// cleanList is list without its blanks and its repeats, in any case, each
// without the spaces around it.
func cleanList(list []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, x := range list {
		x = strings.TrimSpace(x)
		if x == "" || seen[strings.ToLower(x)] {
			continue
		}
		seen[strings.ToLower(x)] = true
		out = append(out, x)
	}
	return out
}
