package settings

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestMusicVideos keeps the picture of music videos off in new installs
// and in the settings of earlier versions, until the user turns it on.
func TestMusicVideos(t *testing.T) {
	for file, want := range map[string]bool{
		`{}`:                   false,
		`{"autoplay":true}`:    false,
		`{"musicVideos":true}`: true,
	} {
		s := Default()
		if err := json.Unmarshal([]byte(file), &s); err != nil {
			t.Fatal(err)
		}
		if s = s.Clean(); s.MusicVideos != want {
			t.Errorf("%s: music videos %v, not %v", file, s.MusicVideos, want)
		}
	}
}

// TestWatching keeps the video player's volume, speed and captions from
// one start to the next, with captions off in new installs, and cleans
// what it can't take.
func TestWatching(t *testing.T) {
	s := Default().Clean()
	if s.VideoCaptions || s.VideoVolume != 100 || s.VideoRate != 1 {
		t.Errorf("new install: %+v", s.Watching)
	}
	for file, want := range map[string]Watching{
		`{}`: {VideoVolume: 100, VideoRate: 1},
		`{"videoVolume":40,"videoRate":1.5,"videoCaptions":true,"videoRepeat":2}`: {VideoVolume: 40, VideoRate: 1.5, VideoCaptions: true, VideoRepeat: 2},
		`{"videoVolume":0,"videoMuted":true}`:                                     {VideoRate: 1, VideoMuted: true},
		`{"videoVolume":400,"videoRate":3,"videoRepeat":9}`:                       {VideoVolume: 100, VideoRate: 1},
	} {
		s := Default()
		if err := json.Unmarshal([]byte(file), &s); err != nil {
			t.Fatal(err)
		}
		if got := s.Clean().Watching; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %+v, not %+v", file, got, want)
		}
	}
}

// TestSegmentsAndBlocks keeps the actions of SponsorBlock's kinds the user
// changed, and their defaults otherwise, and cleans the block lists.
func TestSegmentsAndBlocks(t *testing.T) {
	s := Default().Clean()
	if s.SponsorBlockOff || s.SegmentAction("sponsor") != SegmentSkip || s.SegmentAction("filler") != SegmentOff || s.ShowShorts {
		t.Fatalf("new install: %+v", s.Watching)
	}
	s.Segments = map[string]string{"sponsor": "mute", "intro": "skip", "filler": "off", "poi_highlight": "mute", "nope": "skip"}
	s.BlockedChannels = []string{" A ", "a", "", "B"}
	s = s.Clean()
	if s.SegmentAction("sponsor") != SegmentMute || s.SegmentAction("intro") != SegmentSkip {
		t.Errorf("changed actions: %v", s.Segments)
	}
	// The default, a highlight's mute, and an unknown kind aren't kept.
	if _, ok := s.Segments["filler"]; ok || len(s.Segments) != 2 {
		t.Errorf("segments kept %v", s.Segments)
	}
	if len(s.BlockedChannels) != 2 || s.BlockedChannels[0] != "A" || s.BlockedChannels[1] != "B" {
		t.Errorf("blocked channels %q", s.BlockedChannels)
	}
}
