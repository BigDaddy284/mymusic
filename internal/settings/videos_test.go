package settings

import (
	"encoding/json"
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
		if got := s.Clean().Watching; got != want {
			t.Errorf("%s: %+v, not %+v", file, got, want)
		}
	}
}
